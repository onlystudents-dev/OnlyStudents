package db

import (
	"context"
	"fmt"

	"onlystudents/internal/helpers"
	"onlystudents/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func seedDemoIfNeeded(ctx context.Context, conn *pgxpool.Pool) error {
	var seeded bool
	if err := conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM school_type WHERE id = 1)").Scan(&seeded); err != nil {
		return fmt.Errorf("check demo seed: %w", err)
	}
	if seeded {
		return nil
	}

	content, err := migrations.SeedsFS.ReadFile("seeds/demo.sql")
	if err != nil {
		return fmt.Errorf("read demo seed: %w", err)
	}
	if _, err := conn.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("apply demo seed: %w", err)
	}
	return nil
}

func RunMigrations(conn *pgxpool.Pool) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(conn)
	if err := goose.UpContext(context.Background(), db, "."); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	if helpers.GetEnvFallback("DEMO_MODE", "false") == "true" {
		if err := seedDemoIfNeeded(context.Background(), conn); err != nil {
			return err
		}
	}

	return nil
}
