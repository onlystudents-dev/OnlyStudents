package v1

import (
	"context"
	"crypto/rand"
	"errors"

	"fmt"
	"log/slog"
	"math/big"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	opaquepkg "onlystudents/internal/opaque"
	"time"

	"github.com/bytemare/opaque"
	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type loginInitRequest struct {
	User              int32             `json:"user"`
	Role              string            `json:"role"`
	StartLoginRequest opaquepkg.Message `json:"start_login_request"`
}

type forgetPasswordRequest struct {
	User int32  `json:"user"`
	Role string `json:"role"`
}

type loginInitResponse struct {
	LoginResponse opaquepkg.Message `json:"login_response"`
	LoginHandle   string            `json:"login_handle"`
}

type loginFinishRequest struct {
	LoginHandle        string            `json:"login_handle"`
	FinishLoginRequest opaquepkg.Message `json:"finish_login_request"`
}

func LoginInit(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	var req loginInitRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.User <= 0 || req.Role == "" || len(req.StartLoginRequest) == 0 {
		return c.SendStatus(400)
	}

	var record *opaque.ClientRecord

	queries := db_queries.New(pool)
	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, req.Role, req.User, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	switch {
	case err == nil:
		opaque_record, rerr := queries.GetOpaqueRecord(c.Context(), account.ID)
		if rerr != nil {
			record = opaquepkg.FakeRecord
		} else {
			reg, derr := opaque_server.Deserialize.RegistrationRecord(opaque_record.RegistrationRecord)

			if derr != nil {
				slog.Error("corrupt opaque record", "err", derr)
				return c.SendStatus(500)
			}
			record = &opaque.ClientRecord{
				CredentialIdentifier: []byte(account.ID.String()),
				RegistrationRecord:   reg,
			}
		}
	case errors.Is(err, pgx.ErrNoRows):
		record = opaquepkg.FakeRecord
	default:
		slog.Error("login init account lookup", "err", err)
		return c.SendStatus(500)
	}

	ke2, handle, err := opaquepkg.LoginInit(c.Context(), opaque_server, rdb, record, req.StartLoginRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(loginInitResponse{LoginResponse: ke2, LoginHandle: handle})
}

func LoginFinish(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, opaque_server *opaque.Server) error {
	var req loginFinishRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.LoginHandle == "" {
		return c.SendStatus(400)
	}

	sessionToken, err := opaquepkg.LoginFinish(c.Context(), opaque_server, rdb, pool, req.LoginHandle, req.FinishLoginRequest)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "WRONG_CREDENTIALS"})
	}

	duration := time.Duration(helpers.GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   helpers.GetEnvFallback("APP_ENV", "development") == "production",
		SameSite: "Lax",
	})

	return c.SendStatus(200)
}

func Logout(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_token := c.Cookies("session_token", "")

	if session_token == "" {
		return c.SendStatus(401)
	}

	helpers.SessionDelete(c, rdb, session_token)
	c.ClearCookie("session_token")

	return c.SendStatus(200)
}

type ResetPasswordEmailData struct {
	ResetCode string
	Name      string
}

func generateResetCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}

	return string(code), nil
}

func ForgetPassword(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req forgetPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(200)
	}

	if req.User <= 0 || req.Role == "" {
		return c.SendStatus(200)
	}

	queries := db_queries.New(pool)

	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, req.Role, req.User, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	sendMail := err == nil && account.EmailAddress.Valid && account.EmailAddress.String != "" && account.EmailVerified

	code, e := generateResetCode(helpers.GetIntEnvFallback("PASSWORD_RESET_CODE_LEN", 10, 64))
	if e != nil {
		return c.SendStatus(500)
	}
	cacheKey := fmt.Sprintf("pending_password_reset_%s", code)
	sessionData := helpers.SessionData{AccountID: req.User, Role: req.Role}
	dataJSON, e := json.Marshal(sessionData)
	if e != nil {
		return c.SendStatus(500)
	}
	ttl := time.Duration(helpers.GetInt64EnvFallback("PASSWORD_RESET_CODE_TTL", 15, 1440)) * time.Minute
	if e := rdb.Set(c.Context(), cacheKey, string(dataJSON), ttl).Err(); e != nil {
		return c.SendStatus(500)
	}

	if sendMail {
		email := account.EmailAddress.String
		// #nosec G118
		go func() {
			// intentionally detached to avoid timing-based user enumeration
			bg, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			mailer, e := mail.NewFromEnv()
			if e != nil {
				slog.Error("password reset mailer init failed", "err", e)
				return
			}
			if e := mailer.SendTemplateContext(bg, email, "Reset password", "password_reset.html", ResetPasswordEmailData{ResetCode: code, Name: account.Role}); e != nil {
				slog.Error("password reset email failed", "account", email, "err", e)
			}
		}()
	}

	return c.SendStatus(200)
}

type ResetPasswordConfirmType struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PendingPassword string `json:"pending_password"`
}

func ForgetPasswordConfirm(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req ResetPasswordConfirmType

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Password == "" || req.ConfirmPassword == "" || req.PendingPassword == "" {
		return c.SendStatus(400)
	}

	is_password_good, err := helpers.PasswordChecks(c, req.Password, req.ConfirmPassword)

	if !is_password_good {
		return err
	}

	cacheKey := fmt.Sprintf("pending_password_reset_%s", req.PendingPassword)

	dataJSON, err := rdb.GetDel(c.Context(), cacheKey).Result()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "INVALID_PASSWORD_RESET_TOKEN",
		})
	}

	var sessionData helpers.SessionData
	if err = json.Unmarshal([]byte(dataJSON), &sessionData); err != nil {
		slog.Error("json unmarshall error", "err", err)
		return c.SendStatus(500)
	}

	// IMPORTANT TODO: OPAQUE password change
	return c.SendStatus(501)

	// invalidate cached account object which has the old password hash
	helpers.InvalidateCachedAccount(c.Context(), rdb, sessionData.Role, sessionData.AccountID)

	return c.SendStatus(200)
}
