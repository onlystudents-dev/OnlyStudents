package studentapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SubmitHomeworkRequest struct {
	HomeworkID int32  `json:"homework_id"`
	Type       string `json:"type"`
	Note       string `json:"note"`
	FileType   string `json:"file_type"`
	FileUUID   string `json:"file_uuid"`
}

type SubmitAbsenceReasonRequest struct {
	AbsenceID int32  `json:"absence_id"`
	Type      string `json:"type"`
	Note      string `json:"note"`
	FileType  string `json:"file_type"`
	FileUUID  string `json:"file_uuid"`
}

func SubmitHomework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	return StudentModify(c, pool, rdb,
		func(req SubmitHomeworkRequest) bool {
			if req.HomeworkID <= 0 || (req.Type != "note" && req.Type != "file") {
				return true
			}
			if req.Type == "note" {
				_, err := uuid.Parse(req.FileUUID)
				return req.Note == "" || err == nil
			}

			_, err := uuid.Parse(req.FileUUID)
			return req.FileType == "" || err != nil || req.Note != ""
		},
		func(ctx context.Context, school_id int32, account_id int32, queries *db_queries.Queries, req SubmitHomeworkRequest) (int64, error) {
			file_uuid, err := uuid.Parse(req.FileUUID)
			return queries.StudentUpsertHomeworkSubmission(c.Context(), db_queries.StudentUpsertHomeworkSubmissionParams{
				HomeworkID:     req.HomeworkID,
				StudentID:      account_id,
				SchoolID:       school_id,
				SubmissionType: req.Type,
				Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
				FileUuid:       pgtype.UUID{Bytes: file_uuid, Valid: err == nil},
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
				_, err := uuid.Parse(req.FileUUID)
				return req.Note == "" || err == nil
			}

			_, err := uuid.Parse(req.FileUUID)
			return req.FileType == "" || err != nil || req.Note != ""
		},
		func(ctx context.Context, school_id int32, account_id int32, queries *db_queries.Queries, req SubmitAbsenceReasonRequest) (int64, error) {
			file_uuid, err := uuid.Parse(req.FileUUID)
			return queries.StudentUpsertAbsenceReason(c.Context(), db_queries.StudentUpsertAbsenceReasonParams{
				AbsenceID:      req.AbsenceID,
				StudentID:      account_id,
				SchoolID:       school_id,
				SubmissionType: req.Type,
				Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
				FileUuid:       pgtype.UUID{Bytes: file_uuid, Valid: err == nil},
			})
		})
}
