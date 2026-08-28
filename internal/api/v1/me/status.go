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

func Status(c fiber.Ctx, pool *pgxpool.Pool, session_store *helpers.SessionStore, cache_store *helpers.CacheStore) error {
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

	switch session_data.Role {
	case "student":
		student, err := cache_store.CacheOrGetStudent(context.Background(), *queries, session_data.AccountID, int32(helpers.GetUintEnvFallback("PERSON_CACHE_TTL", 5*60)))

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
		guardian, err := cache_store.CacheOrGetGuardian(context.Background(), *queries, session_data.AccountID, int32(helpers.GetUintEnvFallback("PERSON_CACHE_TTL", 5*60)))

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
		teacher, err := cache_store.CacheOrGetTeacher(context.Background(), *queries, session_data.AccountID, int32(helpers.GetUintEnvFallback("PERSON_CACHE_TTL", 5*60)))

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
