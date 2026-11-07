package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

type BellScheduleTypeSummary struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func convertBellScheduleType(row db_queries.ReadBellScheduleTypeRow) BellScheduleTypeSummary {
	return BellScheduleTypeSummary{
		Id:   row.ID,
		Name: row.Name,
	}
}

func CreateBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req CreateBellScheduleTypeRequest) bool {
			return req.Name == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateBellScheduleTypeRequest) (int64, error, func()) {
			row_count, err := queries.CreateBellScheduleType(ctx, db_queries.CreateBellScheduleTypeParams{
				SchoolID: teacher_scope.SchoolID,
				Name:     req.Name,
			})
			return row_count, err, func() { helpers.InvalidateCachedBellScheduleTypes(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func DeleteBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req DeleteBellScheduleTypeRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteBellScheduleTypeRequest) (int64, error, func()) {
			row_count, err := queries.DeleteBellScheduleType(ctx, db_queries.DeleteBellScheduleTypeParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})

			return row_count, err, func() { helpers.InvalidateCachedBellScheduleTypes(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func EditBellScheduleType(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req EditBellScheduleTypeRequest) bool {
			return req.Id <= 0 || req.Name == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditBellScheduleTypeRequest) (int64, error, func()) {
			row_count, err := queries.EditBellScheduleType(ctx, db_queries.EditBellScheduleTypeParams{
				Name:     req.Name,
				ID:       req.Id,
				SchoolID: teacher_scope.SchoolID,
			})

			return row_count, err, func() { helpers.InvalidateCachedBellScheduleTypes(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func ReadBellScheduleTypes(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "BELL_SCHEDULE_TYPE_CACHE_TTL", "MANAGE_BELL_SCHEDULE", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.ReadBellScheduleTypeRow, error) {
		return helpers.CacheOrGetTeacherBellScheduleTypes(ctx, rdb, queries, scope, ttl)
	}, convertBellScheduleType)
}
