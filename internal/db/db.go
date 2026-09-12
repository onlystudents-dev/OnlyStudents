package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	db_user, user_exists := os.LookupEnv("DB_USER")
	db_password, pass_exists := os.LookupEnv("DB_PASSWORD")
	db_host, host_exists := os.LookupEnv("DB_HOST")
	db_port, port_exists := os.LookupEnv("DB_PORT")
	db_database, db_exists := os.LookupEnv("DB_DATABASE")
	db_sslmode, ssl_exists := os.LookupEnv("DB_SSLMODE")

	if !user_exists || !pass_exists || !host_exists || !port_exists || !db_exists || !ssl_exists {
		fmt.Fprintf(os.Stderr, "Missing required database environment variables: set DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_DATABASE and DB_SSLMODE\n")
		os.Exit(1)
	}

	DATABASE_URL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db_user,
		db_password,
		db_host,
		db_port,
		db_database,
		db_sslmode,
	)

	pool, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}
