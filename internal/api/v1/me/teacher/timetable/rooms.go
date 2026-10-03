package timetable

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Name == "" || req.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_ROOMS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.CreateRoomParams{
		SchoolID: scope.SchoolID,
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	err := queries.CreateRoom(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)

}

func UpdateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req EditRoomRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Capacity == 0 || req.Id == 0 || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_ROOMS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.EditRoomParams{
		Name:     req.Name,
		Capacity: req.Capacity,
		ID:       req.Id,
		SchoolID: scope.SchoolID,
	}

	err := queries.EditRoom(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	var req DeleteRoomRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	if req.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "BAD_REQUEST"})
	}

	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_ROOMS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	params := db_queries.DeleteRoomParams{
		SchoolID: scope.SchoolID,
		ID:       req.Id,
	}

	err := queries.DeleteRoom(c.Context(), params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	scope, status_code := helpers.ResolveTeacherCapabilityScope(c, pool, rdb, "MANAGE_ROOMS")

	if status_code != fiber.StatusOK {
		return c.SendStatus(status_code)
	}

	queries := db_queries.New(pool)

	rooms, err := queries.ReadRoom(c.Context(), scope.SchoolID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "SERVER_ERROR"})
	}

	return c.JSON(rooms)
}
