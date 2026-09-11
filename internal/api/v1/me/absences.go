package meapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AbsenceSummary struct {
	ID          int64
	Student     interface{}
	Subject     pgtype.Text
	SubjectCode pgtype.Text
	Teacher     interface{}
	Date        pgtype.Date
	Type        string
	Justified   bool
	Note        pgtype.Text
	VerifiedBy  interface{}
}

func Absences(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	var absence_summaries []AbsenceSummary

	account_id, err := helpers.ResolvePerson(c, *queries, session_data)

	if err != nil {
		return c.SendStatus(401)
	}

	switch session_data.Role {
	case "student", "guardian":
		absences, err := helpers.CacheOrGetStudentAbsences(c.Context(), rdb, *queries, account_id, helpers.GetInt32EnvFallback("ABSENCES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(500)
		}

		for _, absence_row := range absences {
			absence_summaries = append(absence_summaries, AbsenceSummary{
				ID:          absence_row.ID,
				Subject:     absence_row.Subject,
				SubjectCode: absence_row.SubjectCode,
				Teacher:     absence_row.Teacher,
				Date:        absence_row.Date,
				Type:        absence_row.Type,
				Justified:   absence_row.Justified,
				Note:        absence_row.Note,
				VerifiedBy:  absence_row.VerifiedBy,
			})
		}

	case "teacher":
		absences, err := helpers.CacheOrGetTeacherAbsences(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("ABSENCES_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(500)
		}

		for _, absence_row := range absences {
			absence_summaries = append(absence_summaries, AbsenceSummary{
				ID:          absence_row.ID,
				Student:     absence_row.Student,
				Subject:     absence_row.Subject,
				SubjectCode: absence_row.SubjectCode,
				Date:        absence_row.Date,
				Type:        absence_row.Type,
				Justified:   absence_row.Justified,
				Note:        absence_row.Note,
				VerifiedBy:  absence_row.VerifiedBy,
			})
		}

	default:
		return c.SendStatus(400)
	}

	return c.JSON(absence_summaries)
}
