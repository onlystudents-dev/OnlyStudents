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
	SubjectCode string      `json:"subject_code"`
	Teacher     interface{} `json:"teacher"`
	Term        string      `json:"term"`
	Type        string      `json:"type"`
	Value       int16       `json:"value"`
	Date        pgtype.Date `json:"date"`
	Note        pgtype.Text `json:"note"`
}

func Grades(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	session_data, err := helpers.SessionGet(c, rdb, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	var grade_summaries []GradeSummary

	switch session_data.Role {
	case "student":
		grades, err := helpers.CacheOrGetStudentGrades(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("GRADES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(500)
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

	case "guardian":
		grades, err := helpers.CacheOrGetGuardianGrades(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("GRADES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(500)
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

	case "teacher":
		grades, err := helpers.CacheOrGetTeacherGrades(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("GRADES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(500)
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

	default:
		return c.SendStatus(400)
	}

	return c.JSON(grade_summaries)
}
