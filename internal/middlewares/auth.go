package middlewares

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(c fiber.Ctx, session_store *helpers.SessionStore) error {
	session_token := c.Cookies("session_token", "")

	if session_token == "" {
		return c.SendStatus(401)
	}

	data, err := session_store.Get(c, session_token)

	if err != nil {
		return c.SendStatus(401)
	}

	c.Locals("account", data)

	return c.Next()
}
