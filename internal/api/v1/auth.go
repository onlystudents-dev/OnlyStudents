package v1

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type loginRequest struct {
	User     int32  `json:"user"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var dummyPasswordHash = func() string {
	h, err := helpers.Argon2HashPassword("timing-equalization-dummy")
	if err != nil {
		panic(err)
	}
	return h
}()

func Login(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.User <= 0 || req.Password == "" || req.Role == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)
	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, req.Role, req.User, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	if err != nil {
		// run argon2verify on a dummy hash (if the hash isnt a dummy, you could still do enumeration because "" would fail instantly) to fix timing-based enumeration attacks, if the account does not exist, it would not verify with argon2, which would have a slight latency difference
		helpers.Argon2Verify(req.Password, dummyPasswordHash)
		return c.Status(401).JSON(fiber.Map{
			"error": "WRONG_CREDENTIALS",
		})
	}

	if !helpers.Argon2Verify(req.Password, account.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "WRONG_CREDENTIALS",
		})
	}

	session_token, err := helpers.SessionCreate(c.Context(), rdb, req.User, req.Role)
	if err != nil {
		return c.SendStatus(500)
	}

	duration := time.Duration(helpers.GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   helpers.GetEnvFallback("APP_ENV", "development") == "production",
		SameSite: "Lax",
	})

	return c.JSON(helpers.SessionData{
		Role:      req.Role,
		AccountID: req.User,
	})
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
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(200)
	}

	if req.User <= 0 || req.Role == "" {
		return c.SendStatus(200)
	}

	queries := db_queries.New(pool)

	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, req.Role, req.User, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))
	sendMail := err == nil && account.EmailAddress.Valid && account.EmailAddress.String != ""

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

	hashedPassword, err := helpers.Argon2HashPassword(req.Password)
	if err != nil {
		slog.Error("password hashing error", "err", err)
		return c.SendStatus(500)
	}

	queries := db_queries.New(pool)

	if sessionData.Role == "guardian" {
		err = queries.ResetPasswordGuardian(c.Context(), db_queries.ResetPasswordGuardianParams{PasswordHash: hashedPassword, GuardianID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", "err", err)
			return c.SendStatus(500)
		}
	}

	if sessionData.Role == "student" {
		err = queries.ResetPasswordStudent(c.Context(), db_queries.ResetPasswordStudentParams{PasswordHash: hashedPassword, StudentID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", "err", err)
			return c.SendStatus(500)
		}
	}

	if sessionData.Role == "teacher" {
		err = queries.ResetPasswordTeacher(c.Context(), db_queries.ResetPasswordTeacherParams{PasswordHash: hashedPassword, TeacherID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", "err", err)
			return c.SendStatus(500)
		}
	}

	// invalidate cached account object which has the old password hash
	helpers.InvalidateCachedAccount(c.Context(), rdb, sessionData.Role, sessionData.AccountID)

	return c.SendStatus(200)
}
