package v1

import (
	"context"
	"crypto/rand"
	"math/big"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"onlystudents/internal/mail"
	"time"

	"github.com/gofiber/fiber/v3"
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
	account, err := cache_store.CacheOrGetAccount(context.Background(), *queries, req.Role, req.User, int32(helpers.GetUintEnvFallback("ACCOUNT_CACHE_TTL", 5*60)))

	if err != nil {
		return c.SendStatus(401)
	}

	if !helpers.Argon2Verify(req.Password, account.PasswordHash) {
		return c.SendStatus(401)
	}

	session_token, err := session_store.Create(context.Background(), req.User, req.Role)
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

	err := session_store.Delete(context.Background(), token)

	if err != nil {
		return c.SendStatus(401)
	}

	c.ClearCookie("session_token")

	return c.Redirect().To("/")
}

type ResetPasswordEmailData struct {
	ResetCode string
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
		return c.SendStatus(400)
	}

	if req.User <= 0 || req.Role == "" {
		return c.SendStatus(400)
	}

	queries := db_queries.New(pool)

	account, err := cache_store.CacheOrGetAccount(context.Background(), *queries, req.Role, req.User, int32(helpers.GetUintEnvFallback("ACCOUNT_CACHE_TTL", 5*60)))
	if err != nil {
		return c.SendStatus(401)
	}

	if !account.EmailAddress.Valid || account.EmailAddress.String == "" {
		return c.SendStatus(400)
	}

	code, err := generateResetCode(int(helpers.GetUintEnvFallback("PASSWORD_RESET_CODE_LEN", 8)))
	if err != nil {
		return c.SendStatus(500)
	}

	mailer, err := mail.NewFromEnv()
	if err != nil {
		return c.SendStatus(500)
	}

	if err := mailer.SendTemplate(account.EmailAddress.String, "Reset password", "password_reset.html", ResetPasswordEmailData{ResetCode: code}); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
