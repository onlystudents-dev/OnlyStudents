package meapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type GradeSummary struct {
	ID          int64       `json:"id"`
	Subject     string      `json:"subject"`
	SubjectCode pgtype.Text `json:"subject_code"`
	Teacher     interface{} `json:"teacher"`
	Term        string      `json:"term"`
	Type        string      `json:"type"`
	Value       int16       `json:"value"`
	Date        pgtype.Date `json:"date"`
	Note        pgtype.Text `json:"note"`
}

type FinalGradeSummary struct {
	ID          int64       `json:"id"`
	Subject     string      `json:"subject"`
	SubjectCode pgtype.Text `json:"subject_code"`
	Teacher     interface{} `json:"teacher"`
	Term        string      `json:"term"`
	Value       int16       `json:"value"`
}

func Grades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var grade_summaries []GradeSummary

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch session_data.Role {
	case "student", "guardian":
		grades, err := helpers.CacheOrGetStudentGrades(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback("GRADES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, grade_row := range grades {
			grade_summaries = append(grade_summaries, GradeSummary{
				ID:          grade_row.ID,
				Subject:     grade_row.Subject,
				SubjectCode: grade_row.SubjectCode,
				Teacher:     grade_row.Teacher,
				Term:        grade_row.Term,
				Type:        grade_row.Type,
				Value:       grade_row.Value,
				Date:        grade_row.Date,
				Note:        grade_row.Note,
			})
		}
	}

	return c.JSON(grade_summaries)
}

func FinalGrades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	queries := db_queries.New(pool)

	var final_grade_summaries []GradeSummary

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch session_data.Role {
	case "student", "guardian":
		grades, err := helpers.CacheOrGetStudentFinalGrades(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback("GRADES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, final_grade_row := range grades {
			final_grade_summaries = append(final_grade_summaries, GradeSummary{
				ID:          final_grade_row.ID,
				Subject:     final_grade_row.Subject,
				SubjectCode: final_grade_row.SubjectCode,
				Teacher:     final_grade_row.Teacher,
				Term:        final_grade_row.Term,
				Value:       final_grade_row.Value,
			})
		}

	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(final_grade_summaries)
}
