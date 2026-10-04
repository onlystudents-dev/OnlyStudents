package timetable

import (
	"context"
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func TeacherTimeTableSummaryByID[request_type, Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	required_permission string,
	cache_key string,
	request_check_func func(request_type) bool,
	get_id_func func(request_type) int32,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, ID int32, schoolID int32, accountID int32, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	var req request_type

	if err := c.Bind().Query(&req); err != nil {
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

	var summaries []T

	teacher_scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, required_permission)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, get_id_func(req), teacher_scope.SchoolID, teacher_scope.TeacherID, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

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

func TeacherTimeTableSummary[Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	required_permission string,
	cache_key string,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, schoolID int32, accountID int32, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	var summaries []T

	teacher_scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, required_permission)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, teacher_scope.SchoolID, teacher_scope.TeacherID, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

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

func TeacherTimeTableModify[request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	required_permission string,
	request_check_func func(request_type) bool,
	query_func func(context.Context, helpers.TeacherScope, *db_queries.Queries, request_type) (int64, error),
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

	teacher_scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, required_permission)

	if status_code != fiber.StatusOK {
		return helpers.ErrorByStatusCode(c, status_code)
	}

	queries := db_queries.New(pool)

	rows_affected, err := query_func(c.Context(), teacher_scope, queries, req)

	if err != nil {
		slog.Error("teacher modify err", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	if rows_affected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(fiber.StatusOK)
}
