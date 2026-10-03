package studentapi

import (
	"log/slog"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

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
	var req SubmitHomeworkRequest

	if err := c.Bind().Query(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	rows_affected, homework_submit_err := queries.StudentUpsertHomeworkSubmission(c.Context(), db_queries.StudentUpsertHomeworkSubmissionParams{
		HomeworkID: req.HomeworkID,
		StudentID:  account_id,
		Content:    pgtype.Text{String: req.Content, Valid: true},
	})

	if homework_submit_err != nil {
		slog.Error("submit homework db err", "err", homework_submit_err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if rows_affected == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusOK)
}

func SubmitAbsenceReason(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	// TODO: make this work, needs a new sql table
	return c.SendStatus(fiber.StatusNotImplemented)
}
