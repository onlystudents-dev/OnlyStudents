package helpers

import (
	"context"
	db_queries "onlystudents/internal/db/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckPermission(ctx context.Context, pool *pgxpool.Pool, store *SessionStore, token string, permission string) bool {

	if token == "" {
		return false
	}

	sessionData, err := store.Get(ctx, token)
	if err != nil {
		return false
	}

	queries := db_queries.New(pool)

	params := db_queries.CheckPermissionParams{
		Name:      permission,
		TeacherID: sessionData.AccountID,
	}

	has, err := queries.CheckPermission(ctx, params)

	if err != nil {
		return false
	}

	return has
}
