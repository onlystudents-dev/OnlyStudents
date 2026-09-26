package manageapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TODO: implement this
func AddExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

// TODO: implement this
func RemoveExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

// TODO: implement this
func EditExam(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}
