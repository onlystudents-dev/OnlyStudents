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
	AdminToken string `json:"admin_token"`
}

func AdminTokenLogin(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	if helpers.GetEnvFallback("ADMIN_TOKEN_ENABLED", "false") != "true" {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	}

	admin_token, exists := os.LookupEnv("ADMIN_TOKEN")

	// no default token for security
	if !exists {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
	}

	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.AdminToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if subtle.ConstantTimeCompare([]byte(req.AdminToken), []byte(admin_token)) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	session_token, err := helpers.SessionCreate(c.Context(), rdb, 999, "00000000-0000-0000-0000-000000000000;", "token", "admin")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
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
