package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
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

func CreateCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req CreateCustomSubjectRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_CUSTOM_SUBJECT")

	if has_permission == false {
		return c.SendStatus(401)
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

func EditCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req EditCustomSubjectRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if req.Name == "" || req.Id == 0 {
		c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_CUSTOM_SUBJECT")

	if has_permission == false {
		return c.SendStatus(401)
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

func DeleteCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req DeleteCustomSubjectRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		c.SendStatus(500)
	}

	if token == "" {
		c.SendStatus(401)
	}

	if req.Id == 0 {
		c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_CUSTOM_SUBJECT")

	if has_permission == false {
		c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteCustomSubjectParams{
		ID:       req.Id,
		SchoolID: int32(school_id),
	}

	err = queries.DeleteCustomSubject(c.Context(), params)

	if err != nil {
		c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		c.SendStatus(500)
	}

	token := c.Cookies("session-token")

	if token == "" {
		c.SendStatus(401)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_CUSTOM_SUBJECT")

	if has_permission == false {
		c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	listofcustomsubject, err := queries.ReadCustomSubject(c.Context(), int32(school_id))

	if err != nil {
		c.SendStatus(500)
	}

	return c.JSON(listofcustomsubject)
}
