package timetable

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CreateCustomSubjectRequest struct {
	SubjectName string `json:"subject_name"`
	SubjectIcon string `json:"subject_icon"`
}

type EditCustomSubjectRequest struct {
	Id          int32  `json:"id"`
	SubjectName string `json:"subject_name"`
	SubjectIcon string `json:"subject_icon"`
}

type DeleteCustomSubjectRequest struct {
	Id int32 `json:"id"`
}

type CustomSubjectSummary struct {
	Id          int32  `json:"id"`
	SubjectName string `json:"subject_name"`
	SubjectIcon string `json:"subject_icon"`
}

func convertCustomSubject(row db_queries.ReadCustomSubjectRow) CustomSubjectSummary {
	return CustomSubjectSummary{
		Id:          row.ID,
		SubjectName: row.SubjectName,
		SubjectIcon: row.Icon,
	}
}

func CreateCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req CreateCustomSubjectRequest) bool {
			return req.SubjectName == "" || (req.SubjectIcon != "" && !slices.Contains(helpers.Icons, req.SubjectIcon))
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req CreateCustomSubjectRequest) (int64, error, func()) {
			row_count, err := queries.CreateCustomSubject(ctx, db_queries.CreateCustomSubjectParams{
				SchoolID:    teacher_scope.SchoolID,
				SubjectName: req.SubjectName,
				Icon:        req.SubjectIcon,
			})

			return row_count, err, func() { helpers.InvalidateCachedCustomSubjects(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func DeleteCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req DeleteCustomSubjectRequest) bool {
			return req.Id <= 0
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req DeleteCustomSubjectRequest) (int64, error, func()) {
			row_count, err := queries.DeleteCustomSubject(ctx, db_queries.DeleteCustomSubjectParams{
				SchoolID: teacher_scope.SchoolID,
				ID:       req.Id,
			})

			return row_count, err, func() { helpers.InvalidateCachedCustomSubjects(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func EditCustomSubject(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherModify(c, pool, rdb, "MANAGE_CUSTOM_SUBJECT",
		func(req EditCustomSubjectRequest) bool {
			return req.Id <= 0 || req.SubjectName == "" || (req.SubjectIcon != "" && !slices.Contains(helpers.Icons, req.SubjectIcon))
		},
		func(ctx context.Context, teacher_scope helpers.TeacherScope, queries *db_queries.Queries, req EditCustomSubjectRequest) (int64, error, func()) {
			row_count, err := queries.EditCustomSubject(ctx, db_queries.EditCustomSubjectParams{
				SubjectName: req.SubjectName,
				Icon:        req.SubjectIcon,
				ID:          req.Id,
				SchoolID:    teacher_scope.SchoolID,
			})

			return row_count, err, func() { helpers.InvalidateCachedCustomSubjects(ctx, rdb, teacher_scope.SchoolID) }
		})
}

func ReadCustomSubjects(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return helpers.TeacherSummary(c, pool, rdb, "CUSTOM_SUBJECT_CACHE_TTL", "MANAGE_CUSTOM_SUBJECT", nil, func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, req bool, scope helpers.TeacherScope, ttl int32) ([]db_queries.ReadCustomSubjectRow, error) {
		return helpers.CacheOrGetTeacherCustomSubjects(ctx, rdb, queries, scope, ttl)
	}, convertCustomSubject)
}
