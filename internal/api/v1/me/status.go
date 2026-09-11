package meapi

import (
	db_queries "onlystudents/internal/db/store"
	"onlystudents/internal/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ChildrenData struct {
	AccountID  int32  `json:"account_id"`
	SchoolName string `json:"school_name"`
	SchoolID   string `json:"school_id"`
	ClassID    string `json:"class_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

type StatusData struct {
	Role       string         `json:"role"`
	AccountID  int32          `json:"account_id"`
	SchoolName string         `json:"school_name"`
	SchoolID   string         `json:"school_id"`
	ClassID    string         `json:"class_id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	PfpURL     string         `json:"pfp_url"`
	Children   []ChildrenData `json:"children"`
}

func Status(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client) error {
	session_data, ok := c.Locals("session").(helpers.SessionData)

	if !ok {
		return c.SendStatus(401)
	}

	queries := db_queries.New(pool)

	var status_data StatusData

	account, err := helpers.CacheOrGetAccount(c.Context(), rdb, *queries, session_data.Role, session_data.AccountID, helpers.GetInt32EnvFallback("ACCOUNT_CACHE_TTL", 5*60, 604800))

	if err != nil {
		return c.SendStatus(401)
	}

	switch session_data.Role {
	case "student":
		student, err := helpers.CacheOrGetStudent(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: student.FirstName,
			LastName:  student.LastName,
			PfpURL:    account.PfpUrl,
		}

	case "guardian":
		guardian, err := helpers.CacheOrGetGuardian(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: guardian.FirstName,
			LastName:  guardian.LastName,
			PfpURL:    account.PfpUrl,
		}

	case "teacher":
		teacher, err := helpers.CacheOrGetTeacher(c.Context(), rdb, *queries, session_data.AccountID, helpers.GetInt32EnvFallback("PERSON_CACHE_TTL", 5*60, 604800))

		if err != nil {
			return c.SendStatus(401)
		}

		status_data = StatusData{
			Role:      session_data.Role,
			AccountID: session_data.AccountID,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			PfpURL:    account.PfpUrl}

	default:
		return c.SendStatus(400)
	}

	return c.JSON(status_data)
}
