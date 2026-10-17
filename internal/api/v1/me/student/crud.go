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

func StudentSummary[Row, T, request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	single bool,
	request_check_func func(request_type) bool,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req request_type, scope helpers.MeScope, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req request_type

	if request_check_func != nil {
		if err := c.Bind().Query(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}

		if request_check_func(req) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}
	}

	queries := db_queries.New(pool)

	summaries := []T{}

	me_scope, err := helpers.ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	switch session_data.Role {
	case "student", "guardian":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, req, me_scope, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if !single {
		return c.JSON(summaries)
	} else {
		return c.JSON(summaries[0])
	}
}

func StudentModify[request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	request_check_func func(request_type) bool,
	query_func func(context.Context, int32, int32, *db_queries.Queries, request_type) (int64, error),
) error {
	var req request_type

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if request_check_func(req) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	me_scope, scope_err := helpers.ResolveMeScope(c, pool, rdb)

	if scope_err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	var rows_affected int64
	var err error

	switch session_data.Role {
	case "student", "guardian":
		rows_affected, err = query_func(c.Context(), me_scope.SchoolID, me_scope.StudentID, queries, req)
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
