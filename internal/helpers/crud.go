package helpers

import (
	"context"
	"log/slog"
	db_queries "onlystudents/internal/db/store"

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
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req request_type, scope MeScope, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var req request_type

	if request_check_func != nil {
		if err := c.Bind().Query(&req); err != nil {
			slog.Debug("StudentSummary bad request check failed")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}

		if request_check_func(req) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}
	}

	queries := db_queries.New(pool)

	summaries := []T{}

	me_scope, err := ResolveMeScope(c, pool, rdb)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED"})
	}

	switch session_data.Role {
	case "student", "guardian":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, req, me_scope, GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			error_code, status_code := MapDBError(err)
			return c.Status(status_code).JSON(fiber.Map{"error": error_code})
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
		slog.Debug("StudentModify bad request check failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	session_data, ok := c.Locals("session").(SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	me_scope, scope_err := ResolveMeScope(c, pool, rdb)

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
		error_code, status_code := MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	if rows_affected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func TeacherSummary[Row, T, request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	required_permission string,
	request_check_func func(request_type) bool,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req request_type, scope TeacherScope, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	queries := db_queries.New(pool)

	summaries := []T{}

	var req request_type

	if request_check_func != nil {
		if err := c.Bind().Query(&req); err != nil {
			slog.Debug("StudentSummary bad request check failed")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}

		if request_check_func(req) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
		}
	}

	var teacher_scope TeacherScope
	status_code := fiber.StatusInternalServerError

	if required_permission != "" {
		teacher_scope, status_code = ResolveTeacherCapabilityScope(c, pool, rdb, required_permission)
	} else {
		teacher_scope, status_code = ResolveTeacherScope(c, pool, rdb)
	}

	if status_code != fiber.StatusOK {
		return ErrorByStatusCode(c, status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, req, teacher_scope, GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			error_code, status_code := MapDBError(err)
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
	required_permission string,
	request_check_func func(request_type) bool,
	query_func func(context.Context, TeacherScope, *db_queries.Queries, request_type) (int64, error),
) error {
	var req request_type

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if request_check_func(req) {
		slog.Debug("TeacherModify bad request check failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	_, ok := c.Locals("session").(SessionData)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "INVALID_SESSION"})
	}

	var teacher_scope TeacherScope
	status_code := fiber.StatusInternalServerError

	if required_permission != "" {
		teacher_scope, status_code = ResolveTeacherCapabilityScope(c, pool, rdb, required_permission)
	} else {
		teacher_scope, status_code = ResolveTeacherScope(c, pool, rdb)
	}

	if status_code != fiber.StatusOK {
		return ErrorByStatusCode(c, status_code)
	}

	queries := db_queries.New(pool)

	rows_affected, err := query_func(c.Context(), teacher_scope, queries, req)

	if err != nil {
		error_code, status_code := MapDBError(err)
		return c.Status(status_code).JSON(fiber.Map{"error": error_code})
	}

	if rows_affected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	return c.SendStatus(fiber.StatusOK)
}
