package middlewares

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(c fiber.Ctx, rdb *redis.Client) error {
	session_token := c.Cookies("session_token", "")

	if session_token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	session_data, err := helpers.SessionGet(c, rdb, session_token)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("session", session_data)

	return c.Next()
}
