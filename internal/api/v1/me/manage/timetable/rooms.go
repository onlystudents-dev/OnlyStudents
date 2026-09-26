package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateRoomRequest struct {
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

type EditRoomRequest struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

type DeleteRoomRequest struct {
	Id int32 `json:"id"`
}

func CreateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req CreateRoomRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Name == "" || req.Capacity <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if !has_permission {
		return c.SendStatus(fiber.StatusForbidden)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateRoomParams{
		SchoolID: int32(school_id),
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	err = queries.CreateRoom(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)

}

func UpdateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditRoomRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Capacity == 0 || req.Id == 0 || req.Name == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if !has_permission {
		return c.SendStatus(fiber.StatusForbidden)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditRoomParams{
		Name:     req.Name,
		Capacity: req.Capacity,
		ID:       req.Id,
		SchoolID: int32(school_id),
	}

	err = queries.EditRoom(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteRoomRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if req.Id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if !has_permission {
		return c.SendStatus(fiber.StatusForbidden)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteRoomParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteRoom(c.Context(), params)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if school_id == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if !has_permission {
		return c.SendStatus(fiber.StatusForbidden)
	}

	queries := db_queries.New(pool)

	listofrooms, err := queries.ReadRoom(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(listofrooms)
}
