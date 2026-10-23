package teacherapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SubjectSummary struct {
	ID               int32  `json:"id"`
	Name             string `json:"name"`
	Custom           bool   `json:"is_custom"`
	TeacherID        int32  `json:"teacher_id"`
	TeacherFirstName string `json:"teacher_first_name"`
	TeacherLastName  string `json:"teacher_last_name"`
	GroupID          int32  `json:"group_id"`
	GroupName        string `json:"group_name"`
}

func convertSubject(row db_queries.GetTeacherSubjectsRow) SubjectSummary {
	return SubjectSummary{
		ID:               row.ID,
		Name:             row.Name,
		Custom:           row.Custom,
		TeacherID:        row.TeacherID,
		TeacherLastName:  row.TeacherLastName,
		TeacherFirstName: row.TeacherFirstName,
		GroupID:          row.GroupID,
		GroupName:        row.GroupName,
	}
}

func Subjects(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return TeacherSummary(c, pool, rdb, "TEACHER_SUBJECTS_CACHE_TTL", func(ctx context.Context, rdb *redis.Client, queries db_queries.Queries, scope helpers.TeacherScope, ttl int32) ([]db_queries.GetTeacherSubjectsRow, error) {
		return helpers.CacheOrGetTeacherSubjects(ctx, rdb, queries, scope, ttl)
	}, convertSubject)
}
