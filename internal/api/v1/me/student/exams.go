package studentapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ExamSummary struct {
	ID          int64       `json:"id"`
	Subject     string      `json:"subject"`
	SubjectCode pgtype.Text `json:"subject_code"`
	Teacher     interface{} `json:"teacher"`
	Title       string      `json:"title"`
	Description pgtype.Text `json:"description"`
	Date        pgtype.Date `json:"date"`
	StartTime   pgtype.Time `json:"start_time"`
	EndTime     pgtype.Time `json:"end_time"`
	Room        pgtype.Text `json:"room"`
}

func Exams(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var exam_summaries []ExamSummary

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch session_data.Role {
	case "student", "guardian":
		exams, err := helpers.CacheOrGetStudentExams(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback("exams_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, exam_row := range exams {
			exam_summaries = append(exam_summaries, ExamSummary{
				ID:          exam_row.ID,
				Subject:     exam_row.Subject,
				SubjectCode: exam_row.SubjectCode,
				Teacher:     exam_row.Teacher,
				Title:       exam_row.Title,
				Description: exam_row.Description,
				Date:        exam_row.Date,
				StartTime:   exam_row.StartTime,
				EndTime:     exam_row.EndTime,
				Room:        exam_row.Room,
			})
		}
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(exam_summaries)
}
