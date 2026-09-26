package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateCustomSubjectRequest struct {
	Name string `json:"name"`
}

type EditCustomSubjectRequest struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type DeleteCustomSubjectRequest struct {
	Id int32 `json:"id"`
}

func CreateCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateCustomSubjectRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Name == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveManageScope(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateCustomSubjectParams{
		SchoolID:    int32(scope.SchoolID),
		SubjectName: req.Name,
	}

	err := queries.CreateCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func EditCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditCustomSubjectRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Name == "" || req.Id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveManageScope(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditCustomSubjectParams{
		SubjectName: req.Name,
		SchoolID:    int32(scope.SchoolID),
		ID:          req.Id,
	}

	err := queries.EditCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteCustomSubjectRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	scope, status_code := helpers.ResolveManageScope(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteCustomSubjectParams{
		ID:       req.Id,
		SchoolID: scope.SchoolID,
	}

	err := queries.DeleteCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, status_code := helpers.ResolveManageScope(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	custom_subjects, err := queries.ReadCustomSubject(c.Context(), int32(scope.SchoolID))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(custom_subjects)
}
