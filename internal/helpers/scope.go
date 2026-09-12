package helpers

import (
	"errors"
	db_queries "onlystudents/internal/db/store"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

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
