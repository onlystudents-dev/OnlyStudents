package adminapi

import (
	"crypto/subtle"
	"onlystudents/internal/helpers"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type loginRequest struct {
	AdminToken string `json:"AdminToken"`
}

func AdminTokenLogin(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	if helpers.GetEnvFallback("ADMIN_TOKEN_ENABLED", "false") != "true" {
		return c.SendStatus(fiber.StatusNotImplemented)
	}

	admin_token, exists := os.LookupEnv("ADMIN_TOKEN")

	// no default token for security
	if !exists {
		return c.SendStatus(fiber.StatusNotImplemented)
	}

	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.AdminToken == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if subtle.ConstantTimeCompare([]byte(req.AdminToken), []byte(admin_token)) != 1 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	session_token, err := helpers.SessionCreate(c.Context(), rdb, 999, "00000000-0000-0000-0000-000000000000;", "token", "admin")
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	duration := time.Duration(helpers.GetInt64EnvFallback("SESSION_TTL", 3600, 2592000)) * time.Second

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session_token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   helpers.GetEnvFallback("APP_ENV", "development") == "production",
		SameSite: "Strict",
	})

	return c.SendStatus(fiber.StatusOK)
}
