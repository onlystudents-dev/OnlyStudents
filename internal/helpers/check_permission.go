package helpers

import (
	"math"

	db_queries "onlystudents/internal/db/store"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CheckPermission(c fiber.Ctx, pool *pgxpool.Pool, rdb *redis.Client, permission string, school_id int64) bool {
	session_data, ok := c.Locals("session").(SessionData)

	if !ok {
		return false
	}

	if school_id < 0 || school_id > math.MaxInt32 {
		return false
	}

	queries := db_queries.New(pool)

	params := db_queries.CheckPermissionParams{
		Name:      permission,
		TeacherID: session_data.AccountID,
		SchoolID:  int32(school_id), // #nosec G115 -- bounds-checked above
	}

	has, err := queries.CheckPermission(c.Context(), params)

	if err != nil {
		return false
	}

	return has
}
