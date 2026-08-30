package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBellScheduleTypeRequest struct {
	Name string `json:"name"`
}

type EditBellScheduleTypeRequest struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type DeleteBellScheduleTypeRequest struct {
	Id int32 `json:"id"`
}

func CreateBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req CreateBellScheduleTypeRequest

	token := c.Cookies("session_token")

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

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_BELL_SCHEDULE")

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateBellScheduleTypeParams{
		SchoolID: int32(school_id),
		Name:     req.Name,
	}

	err = queries.CreateBellScheduleType(c.Context(), params)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req DeleteBellScheduleTypeRequest

	token := c.Cookies("session_token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		c.SendStatus(500)
	}

	if req.Id == 0 {
		c.SendStatus(400)
	}

	if token == "" {
		c.SendStatus(401)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_BELL_SCHEDULE")

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteBellScheduleTypeParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteBellScheduleType(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {
	var req EditBellScheduleTypeRequest

	token := c.Cookies("session_token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.Name == "" {
		c.SendStatus(400)
	}

	if token == "" {
		c.SendStatus(401)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_BELL_SCHEDULE")

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditBellScheduleTypeParams{
		Name:     req.Name,
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.EditBellScheduleType(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, store *helpers.SessionStore) error {

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, store, token, "MANAGE_BELL_SCHEDULE")

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	listoftype, err := queries.ReadBellScheduleType(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(listoftype)
}
