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

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Name == "" || req.Capacity <= 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateRoomParams{
		SchoolID: int32(school_id),
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	err = queries.CreateRoom(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)

}

func UpdateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditRoomRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)
	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Capacity == 0 || req.Id == 0 || req.Name == "" {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
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
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func DeleteRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteRoomRequest

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	if req.Id == 0 {
		return c.SendStatus(400)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteRoomParams{
		SchoolID: int32(school_id),
		ID:       req.Id,
	}

	err = queries.DeleteRoom(c.Context(), params)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {

	school_id, err := strconv.ParseInt(c.Get("X-School"), 10, 32)

	if err != nil {
		return c.SendStatus(500)
	}

	if school_id == 0 {
		return c.SendStatus(500)
	}

	has_permission := helpers.CheckPermission(c, pool, rdb, "MANAGE_ROOMS", school_id)

	if has_permission == false {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	listofrooms, err := queries.ReadCustomSubject(c.Context(), int32(school_id))

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(listofrooms)
}
