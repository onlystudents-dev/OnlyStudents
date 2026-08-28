package meapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: implement this and add redis caching
func Homework(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	_, err := session_store.Get(c, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	return c.SendStatus(501)
}
