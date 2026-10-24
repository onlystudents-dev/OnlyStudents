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
	Type       string
	Note       string
	FileType   string
	FileID     int64
}

type SubmitAbsenceReasonRequest struct {
	AbsenceID int32
	Type      string
	Note      string
	FileType  string
	FileID    int64
}

func SubmitHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentModify(c, pool, rdb,
		func(req SubmitHomeworkRequest) bool {
			if req.HomeworkID <= 0 || (req.Type != "note" && req.Type != "file") {
				return true
			}
			if req.Type == "note" {
				return req.Note == "" || req.FileID > 0
			}
			return req.FileType == "" || req.FileID <= 0 || req.Note != ""
		},
		func(ctx context.Context, school_id int32, account_id int32, queries *db_queries.Queries, req SubmitHomeworkRequest) (int64, error) {
			return queries.StudentUpsertHomeworkSubmission(c.Context(), db_queries.StudentUpsertHomeworkSubmissionParams{
				HomeworkID:     req.HomeworkID,
				StudentID:      account_id,
				SchoolID:       school_id,
				SubmissionType: req.Type,
				Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
				FileID:         pgtype.Int8{Int64: req.FileID, Valid: req.FileID > 0},
			})
		})
}

func SubmitAbsenceReason(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentModify(c, pool, rdb,
		func(req SubmitAbsenceReasonRequest) bool {
			if req.AbsenceID <= 0 || (req.Type != "note" && req.Type != "file") {
				return true
			}
			if req.Type == "note" {
				return req.Note == "" || req.FileID > 0
			}
			return req.FileType == "" || req.FileID <= 0
		},
		func(ctx context.Context, school_id int32, account_id int32, queries *db_queries.Queries, req SubmitAbsenceReasonRequest) (int64, error) {
			return queries.StudentUpsertAbsenceReason(c.Context(), db_queries.StudentUpsertAbsenceReasonParams{
				AbsenceID:      req.AbsenceID,
				StudentID:      account_id,
				SchoolID:       school_id,
				SubmissionType: req.Type,
				Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
				FileID:         pgtype.Int8{Int64: req.FileID, Valid: req.FileID > 0},
			})
		})
}
