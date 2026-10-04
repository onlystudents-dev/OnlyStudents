package studentapi

import (
	"context"
	"log/slog"
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

func StudentModify[request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	request_check_func func(request_type) bool,
	query_func func(context.Context, int32, *db_queries.Queries, request_type) (int64, error),
) error {
	var req request_type

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if request_check_func(req) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	account_id, resolve_err := helpers.ResolvePerson(c, *queries, session_data)

	if resolve_err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var rows_affected int64
	var err error

	switch session_data.Role {
	case "student", "guardian":
		rows_affected, err = query_func(c.Context(), account_id, queries, req)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if err != nil {
		slog.Error("student modify err", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if rows_affected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(fiber.StatusOK)
}
