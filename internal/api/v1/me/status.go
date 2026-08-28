package meapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusData struct {
	session_data helpers.SessionData
	first_name   string
	last_name    string
}

func Status(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore) error {
	session_token := c.Cookies("session_token", "")
	if session_token == "" {
		return c.SendStatus(401)
	}

	session_data, err := session_store.Get(c, session_token)
	if err != nil {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	var status_data StatusData

	// TODO: add redis caching for this
	switch session_data.Role {
	case "student":
		student, err := queries.GetStudent(context.Background(), session_data.AccountID)

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			session_data: session_data,
			first_name:   student.FirstName,
			last_name:    student.LastName,
		}

	case "guardian":
		guardian, err := queries.GetGuardian(context.Background(), session_data.AccountID)

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			session_data: session_data,
			first_name:   guardian.FirstName,
			last_name:    guardian.LastName,
		}

	case "teacher":
		teacher, err := queries.GetTeacher(context.Background(), session_data.AccountID)

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			session_data: session_data,
			first_name:   teacher.FirstName,
			last_name:    teacher.LastName,
		}

	default:
		return c.SendStatus(400)
	}

	return c.JSON(status_data)
}
