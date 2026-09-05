package helpers

import (
	"context"
	db_queries "onlystudents/internal/db/store"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CheckPermission(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client, token string, permission string, school_id int64) bool {
	if token == "" {
		return false
	}

	sessionData, err := SessionGet(ctx, rdb, token)
	if err != nil {
		return false
	}

	queries := db_queries.New(pool)

	params := db_queries.CheckPermissionParams{
		Name:      permission,
		TeacherID: sessionData.AccountID,
		SchoolID:  int32(school_id),
	}

	has, err := queries.CheckPermission(ctx, params)

	if err != nil {
		return false
	}

	return has
}
