package teacherapi

import (
	"context"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func TeacherSummaryByID[request_type, Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	single bool,
	request_check_func func(request_type) bool,
	get_id_func func(request_type) int32,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, ID int32, scope helpers.TeacherScope, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	var req request_type

	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if request_check_func(req) {
		slog.Debug("TeacherSummaryByID bad request check failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	summaries := []T{}

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, get_id_func(req), teacher_scope, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			error_code, status_code := helpers.MapDBError(err)
			return c.Status(status_code).JSON(fiber.Map{"error": error_code})
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if single {
		if len(summaries) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}
		return c.JSON(summaries[0])
	}
	return c.JSON(summaries)
}

func TeacherSummary[Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope helpers.TeacherScope, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	summaries := []T{}

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, teacher_scope, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			error_code, status_code := helpers.MapDBError(err)
			return c.Status(status_code).JSON(fiber.Map{"error": error_code})
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.JSON(summaries)
}

func TeacherModify[request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	request_check_func func(request_type) bool,
	query_func func(context.Context, helpers.TeacherScope, *db_queries.Queries, request_type) (int64, error),
) error {
	var req request_type

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if request_check_func(req) {
		slog.Debug("TeacherModify bad request check failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	queries := db_queries.New(pool)

	rows_affected, err := query_func(c.Context(), teacher_scope, queries, req)

	if err != nil {
		error_code, status_code := helpers.MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	if rows_affected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(fiber.StatusOK)
}
