package timetable

import (
	"context"
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

type RoomSummary struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

func convertRoom(row db_queries.ReadRoomRow) RoomSummary {
	return RoomSummary{
		ID:       row.ID,
		Name:     row.Name,
		Capacity: row.Capacity,
	}
}

func CreateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_ROOMS",
		func(req CreateRoomRequest) bool {
			return req.Name == "" || req.Capacity <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateRoomRequest) (int64, error) {
			return queries.CreateRoom(ctx, db_queries.CreateRoomParams{
				SchoolID: teacher_scope.SchoolID,
				Name:     req.Name,
				Capacity: req.Capacity,
			})
		})
}

func UpdateRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_ROOMS",
		func(req EditRoomRequest) bool {
			return req.Capacity <= 0 || req.Id <= 0 || req.Name == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditRoomRequest) (int64, error) {
			return queries.EditRoom(ctx, db_queries.EditRoomParams{
				ID:       req.Id,
				SchoolID: teacher_scope.SchoolID,
				Name:     req.Name,
				Capacity: req.Capacity,
			})
		})
}

func DeleteRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_ROOMS",
		func(req DeleteRoomRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteRoomRequest) (int64, error) {
			return queries.DeleteRoom(ctx, db_queries.DeleteRoomParams{
				ID:       req.Id,
				SchoolID: teacher_scope.SchoolID,
			})
		})
}

func ReadRoom(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableSummary(c, pool, rdb, "MANAGE_ROOMS", "ROOMS_CACHE_TTL", helpers.CacheOrGetRooms, convertRoom)
}
