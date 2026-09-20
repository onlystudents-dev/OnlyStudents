package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

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
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	if req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateCustomSubjectParams{
		SchoolID:    int32(school_id),
		SubjectName: req.Name,
	}

	err = queries.CreateCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditCustomSubjectRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(400)
	}

	if req.Name == "" || req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditCustomSubjectParams{
		SubjectName: req.Name,
		SchoolID:    int32(school_id),
		ID:          req.Id,
	}

	err = queries.EditCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteCustomSubjectRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(400)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteCustomSubjectParams{
		ID:       req.Id,
		SchoolID: int32(school_id),
	}

	err = queries.DeleteCustomSubject(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(400)
	}

	if school_id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT", school_id)

	if !has_permission {
		return c.SendStatus(403)
	}

	queries := db_queries.New(pool)

	listofcustomsubject, err := queries.ReadCustomSubject(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(listofcustomsubject)
}
