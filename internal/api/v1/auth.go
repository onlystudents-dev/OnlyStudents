package v1

import (
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
)

type loginRequest struct {
	User     int32  `json:"user"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Login(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore, cache_store *helpers.CacheStore) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.User <= 0 || req.Password == "" || req.Role == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)
	account, err := cache_store.CacheOrGetAccount(c.Context(), *queries, req.Role, req.User, int32(helpers.GetUintEnvFallback("ACCOUNT_CACHE_TTL", 5*60)))

	if err != nil {
		// run argon2verify to fix timing-based enumeration attacks, if the account does not exist, it would not verify with argon2, which would have a slight latency difference
		helpers.Argon2Verify(req.Password, account.PasswordHash)
		return c.SendStatus(401)
	}

	if !helpers.Argon2Verify(req.Password, account.PasswordHash) {
		return c.SendStatus(401)
	}

	session_token, err := session_store.Create(c.Context(), req.User, req.Role)
	if err != nil {
		return c.SendStatus(500)
	}

	duration := time.Duration(helpers.GetUintEnvFallback("SESSION_TTL", 3600)) * time.Second

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

func Logout(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	token := c.Cookies("session_token", "")

	if token == "" {
		return c.SendStatus(401)
	}

	err := session_store.Delete(c.Context(), token)

	if err != nil {
		return c.SendStatus(401)
	}

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

func ForgetPassword(c fiber.Ctx, pool *pgxpool.Pool, cache_store *helpers.CacheStore) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(200)
	}

	if req.User <= 0 || req.Role == "" {
		return c.SendStatus(200)
	}

	queries := db_queries.New(pool)

	account, err := cache_store.CacheOrGetAccount(c.Context(), *queries, req.Role, req.User, int32(helpers.GetUintEnvFallback("ACCOUNT_CACHE_TTL", 5*60)))
	if err != nil {
		return c.SendStatus(200)
	}

	if !account.EmailAddress.Valid || account.EmailAddress.String == "" {
		return c.SendStatus(200)
	}

	code, err := generateResetCode(int(helpers.GetUintEnvFallback("PASSWORD_RESET_CODE_LEN", 10)))
	if err != nil {
		return c.SendStatus(500)
	}

	cacheKey := fmt.Sprintf("pending_password_reset_%s", code)

	sessionData := helpers.SessionData{
		AccountID: req.User,
		Role:      req.Role,
	}

	dataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return c.SendStatus(500)
	}

	ttl := time.Duration(helpers.GetUintEnvFallback("PASSWORD_RESET_CODE_TTL", 15)) * time.Minute
	err = cache_store.RedisDB.Set(c.Context(), cacheKey, string(dataJSON), ttl).Err()
	if err != nil {
		return c.SendStatus(500)
	}

	mailer, err := mail.NewFromEnv()
	if err != nil {
		return c.SendStatus(500)
	}

	if err := mailer.SendTemplate(account.EmailAddress.String, "Reset password", "password_reset.html", ResetPasswordEmailData{ResetCode: code, Name: account.Role}); err != nil {
		slog.Error("password reset email failed", "account", account.EmailAddress.String, "err", err)
		if helpers.GetEnvFallback("APP_ENV", "development") != "production" {
			return c.Status(502).SendString(err.Error())
		}
		return c.SendStatus(502)
	}

	return c.SendStatus(200)
}

type ResetPasswordConfirmType struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PendingPassword string `json:"pending_password"`
}

func ForgetPasswordConfirm(c fiber.Ctx, pool *pgxpool.Pool, cache_store *helpers.CacheStore) error {
	var req ResetPasswordConfirmType

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	if req.Password == "" || req.ConfirmPassword == "" || req.PendingPassword == "" {
		return c.SendStatus(400)
	}

	if req.Password != req.ConfirmPassword {
		return c.SendStatus(401)
	}

	cacheKey := fmt.Sprintf("pending_password_reset_%s", req.PendingPassword)

	dataJSON, err := cache_store.RedisDB.GetDel(c.Context(), cacheKey).Result()
	if err != nil {
		return c.Status(400).SendString("Invalid password reset token!")
	}

	var sessionData helpers.SessionData
	if err = json.Unmarshal([]byte(dataJSON), &sessionData); err != nil {
		slog.Error("json unmarshall error", err)
		return c.SendStatus(500)
	}

	hashedPassword, err := helpers.Argon2HashPassword(req.Password)
	if err != nil {
		slog.Error("password hashing error", err)
		return c.SendStatus(500)
	}

	queries := db_queries.New(pool)

	if sessionData.Role == "guardian" {
		err = queries.ResetPasswordGuardian(c.Context(), db_queries.ResetPasswordGuardianParams{PasswordHash: hashedPassword, GuardianID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", err)
			return c.SendStatus(500)
		}
	}

	if sessionData.Role == "student" {
		err = queries.ResetPasswordStudent(c.Context(), db_queries.ResetPasswordStudentParams{PasswordHash: hashedPassword, StudentID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", err)
			return c.SendStatus(500)
		}
	}

	if sessionData.Role == "teacher" {
		err = queries.ResetPasswordTeacher(c.Context(), db_queries.ResetPasswordTeacherParams{PasswordHash: hashedPassword, TeacherID: pgtype.Int4{Int32: sessionData.AccountID, Valid: true}})

		if err != nil {
			slog.Error("password recovery error", err)
			return c.SendStatus(500)
		}
	}

	// invalidate cached account object which has the old password hash
	cache_store.InvalidateCachedAccount(c.Context(), sessionData.Role, sessionData.AccountID)

	return c.SendStatus(200)
}
