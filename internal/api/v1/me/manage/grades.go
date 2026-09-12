package manageapi

import (
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TODO: implement this
func GetGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func AddGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func RemoveGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func EditGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func GetFinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func AddFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func RemoveFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}

// TODO: implement this
func EditFinalGrade(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}
	return c.SendStatus(501)
}
