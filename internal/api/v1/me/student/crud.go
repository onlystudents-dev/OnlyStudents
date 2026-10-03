package studentapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func StudentSummary[Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, accountID int32, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	var summaries []T

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	switch session_data.Role {
	case "student", "guardian":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(summaries)
}
