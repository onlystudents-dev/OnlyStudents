package teacherapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TODO: implement this
func AddHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

// TODO: implement this
func RemoveHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

// TODO: implement this
func EditHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

// TODO: implement this
func EditHomeworkSubmission(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}
