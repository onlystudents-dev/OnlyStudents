package helpers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckPermission(ctx context.Context, pool *pgxpool.Pool, session_store *SessionStore, token string, permission string) {

}
