package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateCustomSubjectRequest struct {
	SubjectName string `json:"subject_name"`
}

type EditCustomSubjectRequest struct {
	Id          int32  `json:"id"`
	SubjectName string `json:"subject_name"`
}

type DeleteCustomSubjectRequest struct {
	Id int32 `json:"id"`
}

type CustomSubjectSummary struct {
	Id          int32  `json:"id"`
	SubjectName string `json:"subject_name"`
}

func convertCustomSubject(row db_queries.ReadCustomSubjectRow) CustomSubjectSummary {
	return CustomSubjectSummary{
		Id:          row.ID,
		SubjectName: row.SubjectName,
	}
}

func CreateCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req CreateCustomSubjectRequest) bool {
			return req.SubjectName == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateCustomSubjectRequest) (int64, error) {
			return queries.CreateCustomSubject(ctx, db_queries.CreateCustomSubjectParams{
				SchoolID:    teacher_scope.SchoolID,
				SubjectName: req.SubjectName,
			})
		})
}

func DeleteCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req DeleteCustomSubjectRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteCustomSubjectRequest) (int64, error) {
			return queries.DeleteCustomSubject(ctx, db_queries.DeleteCustomSubjectParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})
		})
}

func EditCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req EditCustomSubjectRequest) bool {
			return req.Id <= 0 || req.SubjectName == ""
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditCustomSubjectRequest) (int64, error) {
			return queries.EditCustomSubject(ctx, db_queries.EditCustomSubjectParams{
				SubjectName: req.SubjectName,
				ID:          req.Id,
				SchoolID:    teacher_scope.SchoolID,
			})
		})
}

func ReadCustomSubjects(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherTimeTableSummary(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT", "BELL_SCHEDULE_TYPE_CACHE_TTL", helpers.CacheOrGetCustomSubjects, convertCustomSubject)
}
