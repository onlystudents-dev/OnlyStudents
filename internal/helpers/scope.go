package helpers

import (
	"errors"
	db_queries "onlystudents/internal/db/store"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type readerScope struct {
	SchoolID  int32
	ClassID   int32
	StudentID int32
}

func ResolveReaderScope(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) (readerScope, error) {
	session, ok := c.Locals("session").(SessionData)
	if !ok {
		return readerScope{}, errors.New("invalid scope")
	}

	queries := db_queries.New(pool)
	ttl := GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800)

	var studentID int32
	switch session.Role {
	case "student":
		studentID = session.AccountID
	case "guardian":
		id, err := ResolvePerson(c, *queries, session)
		if err != nil {
			return readerScope{}, errors.New("invalid scope")
		}
		studentID = id
	default:
		return readerScope{}, errors.New("invalid scope")
	}

	student, err := CacheOrGetStudent(c.Context(), rdb, *queries, studentID, ttl)
	if err != nil {
		return readerScope{}, errors.New("invalid scope")
	}

	return readerScope{StudentID: studentID, SchoolID: student.SchoolID, ClassID: student.ClassesID}, nil
}

func ResolvePerson(c fiber.Ctx, queries db_queries.Queries, session_data SessionData) (int32, error) {
	switch session_data.Role {
	case "student":
		return session_data.AccountID, nil
	case "guardian":
		requested_student_id_str := c.Query("student_id")

		requested_student_id, err := strconv.ParseInt(requested_student_id_str, 10, 32)

		if err != nil {
			return 0, err
		}

		can_view_student, err := queries.CanViewStudent(c.Context(), db_queries.CanViewStudentParams{
			GuardianID: session_data.AccountID,
			StudentID:  int32(requested_student_id),
		})

		if err != nil {
			return 0, err
		}

		if can_view_student != 1 {
			return 0, errors.New("No access")
		}

		return int32(requested_student_id), nil

	case "teacher":
		return 0, errors.New("teacher cannot access this")
	default:
		return 0, errors.New("role doesn't exist")
	}
}
