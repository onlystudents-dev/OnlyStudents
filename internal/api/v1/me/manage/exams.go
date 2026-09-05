package manageapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TODO: implement this
func ManageExams(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	_, err := helpers.SessionGet(c, rdb, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	return c.SendStatus(501)
}
