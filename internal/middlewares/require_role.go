package middlewares

import (
	"onlystudents/internal/helpers"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func RequireRoleMiddleware(c fiber.Ctx, rdb *redis.Client, roles []string) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	if !slices.Contains(roles, session_data.Role) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	return c.Next()
}
