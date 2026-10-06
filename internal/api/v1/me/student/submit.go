package studentapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SubmitHomeworkRequest struct {
	HomeworkID int32
	Content    string
}

type SubmitAbsenceRequest struct {
	HomeworkID int32
	Content    string
}

func SubmitHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentModify(c, pool, rdb,
		func(req SubmitHomeworkRequest) bool {
			return req.HomeworkID <= 0 || req.Content == ""
		},
		func(ctx context.Context, school_id int32, account_id int32, queries *db_queries.Queries, req SubmitHomeworkRequest) (int64, error) {
			return queries.StudentUpsertHomeworkSubmission(c.Context(), db_queries.StudentUpsertHomeworkSubmissionParams{
				HomeworkID: req.HomeworkID,
				StudentID:  account_id,
				SchoolID:   school_id,
				Content:    pgtype.Text{String: req.Content, Valid: true},
			})
		})
}

func SubmitAbsenceReason(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	// TODO: make this work, needs a new sql table
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "NOT_IMPLEMENTED"})
}
