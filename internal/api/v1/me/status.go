package meapi

import (
	"context"
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusData struct {
	Role      string `json:"role"`
	AccountID int32  `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: student.FirstName,
			LastName:  student.LastName,
		}

	case "guardian":
		guardian, err := queries.GetGuardian(context.Background(), session_data.AccountID)

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: guardian.FirstName,
			LastName:  guardian.LastName,
		}

	case "teacher":
		teacher, err := queries.GetTeacher(context.Background(), session_data.AccountID)

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
		}

	default:
		return c.SendStatus(400)
	}

	return c.JSON(status_data)
}
