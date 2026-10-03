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
	request_check_func func(request_type) bool,
	get_id_func func(request_type) int32,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, ID int32, schoolID int32, accountID int32, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	var req request_type

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !request_check_func(req) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var summaries []T

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, get_id_func(req), teacher_scope.SchoolID, teacher_scope.TeacherID, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(summaries)
}

func TeacherSummary[Row, T any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	cache_key string,
	cache_or_get_func func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, schoolID int32, accountID int32, ttl int32) ([]Row, error),
	convert func(row Row) T,
) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var summaries []T

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	switch session_data.Role {
	case "teacher":
		data, err := cache_or_get_func(c.Context(), rdb, *queries, teacher_scope.SchoolID, teacher_scope.TeacherID, helpers.GetInt32EnvFallback(cache_key, 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, row := range data {
			summaries = append(summaries, convert(row))
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(summaries)
}

func TeacherModify[request_type any](
	c fiber.Ctx,
	pool *pgxpool.Pool,
	rdb *redis.Client,
	teacher_student_check bool,
	class_subjects_func func(request_type) int32,
	request_check_func func(request_type) bool,
	query_func func(context.Context, helpers.TeacherScope, *db_queries.Queries, request_type) (int64, error),
) error {
	var req request_type

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !request_check_func(req) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	teacher_scope, status_code := helpers.ResolveTeacherScope(c, pool, rdb)

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	if teacher_student_check {
		teacher_accessible, err := queries.TeacherTeachesStudent(c.Context(), db_queries.TeacherTeachesStudentParams{
			SchoolID:  teacher_scope.SchoolID,
			TeacherID: teacher_scope.TeacherID,
			ID:        class_subjects_func(req),
		})

		if err != nil {
			slog.Error("teacher accessible check err", "err", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !teacher_accessible {
			return c.SendStatus(fiber.StatusForbidden)
		}
	}

	rows_affected, err := query_func(c.Context(), teacher_scope, queries, req)

	if err != nil {
		slog.Error("teacher modify err", "err", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if rows_affected == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}
