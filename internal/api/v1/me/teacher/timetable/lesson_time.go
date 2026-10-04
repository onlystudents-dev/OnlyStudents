package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateLessonTimeRequest struct {
	TypeID       int32 `json:"type_id"`
	LessonNumber int32 `json:"lesson_number"`
	Start        int64 `json:"lesson_start"`
	End          int64 `json:"lesson_stop"`
}

type EditLessonTimeRequest struct {
	Id           int32 `json:"id"`
	LessonNumber int32 `json:"lesson_number"`
	Start        int64 `json:"lesson_start"`
	End          int64 `json:"lesson_stop"`
}

type DeleteLessonTimeRequest struct {
	Id int32 `json:"id"`
}

type ReadLessonTimeRequest struct {
	TypeID int32 `json:"type_id" query:"type_id"`
}

type BellScheduleSummary struct {
	ID              int32 `json:"id"`
	SchoolID        int32 `json:"school_id"`
	TypeID          int32 `json:"type_id"`
	HasLessonNumber bool  `json:"has_lesson_number"`
	LessonNumber    int32 `json:"lesson_number"`
	AtStart         int64 `json:"at_start"`
	AtEnd           int64 `json:"at_end"`
}

func convertBellSchedule(row db_queries.BellSchedule) BellScheduleSummary {
	return BellScheduleSummary{
		ID:              row.ID,
		SchoolID:        row.SchoolID,
		TypeID:          row.TypeID,
		HasLessonNumber: row.LessonNumber.Valid,
		LessonNumber:    row.LessonNumber.Int32,
		AtStart:         row.AtStart.Microseconds / 1000000,
		AtEnd:           row.AtEnd.Microseconds / 1000000,
	}
}

func CreateLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req CreateLessonTimeRequest) bool {
			return req.TypeID <= 0 || req.LessonNumber <= 0 || req.Start <= 0 || req.End == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateLessonTimeRequest) (int64, error) {
			params := db_queries.CreateLessonTimeParams{
				SchoolID:     int32(teacher_scope.SchoolID),
				TypeID:       req.TypeID,
				LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
				AtStart:      pgtype.Time{Microseconds: req.Start * 1000000, Valid: true},
				AtEnd:        pgtype.Time{Microseconds: req.End * 1000000, Valid: true},
			}

			return queries.CreateLessonTime(c.Context(), params)
		})
}

func DeleteLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req DeleteLessonTimeRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteLessonTimeRequest) (int64, error) {
			return queries.DeleteLessonTime(c.Context(), db_queries.DeleteLessonTimeParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})
		})
}

func EditLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_BELL_SCHEDULE",
		func(req EditLessonTimeRequest) bool {
			return req.Id <= 0 || req.LessonNumber <= 0 || req.Start <= 0 || req.End == 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditLessonTimeRequest) (int64, error) {
			return queries.EditLessonTime(c.Context(), db_queries.EditLessonTimeParams{
				SchoolID:     int32(teacher_scope.SchoolID),
				ID:           req.Id,
				LessonNumber: pgtype.Int4{Int32: req.LessonNumber, Valid: true},
				AtStart:      pgtype.Time{Microseconds: req.Start * 1000000, Valid: true},
				AtEnd:        pgtype.Time{Microseconds: req.End * 1000000, Valid: true},
			})
		})
}

func ReadLessonTime(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableSummaryByID(c, pool, rdb, "MANAGE_BELL_SCHEDULE", "GROUP_STUDENTS_CACHE_TTL",
		func(req ReadLessonTimeRequest) bool { return req.TypeID <= 0 },
		func(req ReadLessonTimeRequest) int32 { return req.TypeID },
		helpers.CacheOrGetLessonTime, convertBellSchedule)
}
