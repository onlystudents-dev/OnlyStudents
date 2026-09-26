package studentapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HomeworkSummary struct {
	ID          int64              `json:"id"`
	Subject     string             `json:"subject"`
	SubjectCode pgtype.Text        `json:"subject_code"`
	Teacher     interface{}        `json:"teacher"`
	Title       string             `json:"title"`
	Description pgtype.Text        `json:"description"`
	DueDate     pgtype.Date        `json:"due_date"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	SubmittedAt pgtype.Timestamptz `json:"submitted_at"`
	GradedValue pgtype.Int2        `json:"graded_value"`
}

func Homework(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var homework_summaries []HomeworkSummary

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch session_data.Role {
	case "student", "guardian":
		homeworks, err := helpers.CacheOrGetStudentHomework(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback("homeworks_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, homework_row := range homeworks {
			homework_summaries = append(homework_summaries, HomeworkSummary{
				ID:          homework_row.ID,
				Subject:     homework_row.Subject,
				SubjectCode: homework_row.SubjectCode,
				Teacher:     homework_row.Teacher,
				Title:       homework_row.Title,
				Description: homework_row.Description,
				DueDate:     homework_row.DueDate,
				CreatedAt:   homework_row.CreatedAt,
				SubmittedAt: homework_row.SubmittedAt,
				GradedValue: homework_row.GradedValue,
			})
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(homework_summaries)
}
