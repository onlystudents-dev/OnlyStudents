package db

import (
	"context"
	"fmt"
	env "onlystudents/internal/helpers"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	DATABASE_URL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetEnvFallback("DB_USER", "user"),
		env.GetEnvFallback("DB_PASSWORD", "changeme"),
		env.GetEnvFallback("DB_HOST", "localhost"),
		env.GetEnvFallback("DB_PORT", "5432"),
		env.GetEnvFallback("DB_DATABASE", "onlystudents"),
	)

	pool, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}
