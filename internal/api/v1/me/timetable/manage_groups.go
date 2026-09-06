package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateGroupRequest struct {
	BellId int32  `json:"bell_id"`
	Name   string `json:"name"`
}

type DeleteGroupRequest struct {
	Id int32 `json:"id"`
}

type EditGroupRequest struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func CreateGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateGroupRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.BellId == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_GROUPS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateGroupParams{
		SchoolID:  int32(school_id),
		BellID:    req.BellId,
		GroupName: req.Name,
	}

	err = queries.CreateGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteGroupRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_GROUPS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteGroupParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func EditGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditGroupRequest

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_GROUPS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditGroupParams{
		SchoolID:  int32(school_id),
		ID:        req.Id,
		GroupName: req.Name,
	}

	err = queries.EditGroup(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadGroup(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {

	token := c.Cookies("session-token")

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if token == "" {
		return c.SendStatus(401)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c.Context(), pool, rdb, token, "MANAGE_GROUPS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	data, err := queries.ReadGroup(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(data)
}
