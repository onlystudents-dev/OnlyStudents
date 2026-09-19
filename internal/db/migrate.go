package db

import (
	"context"
	"fmt"
	"io/fs"
	"onlystudents/internal/helpers"
	constants "onlystudents/internal/helpers"
	"onlystudents/migrations"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func applyMigration(ctx context.Context, conn *pgxpool.Pool, name string) error {
	content, err := migrations.FS.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}
	if _, err := conn.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("apply %s: %w", name, err)
	}
	version := strings.TrimSuffix(name, ".sql")
	if _, err := conn.Exec(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES ($1, now())", version); err != nil {
		return fmt.Errorf("record %s: %w", version, err)
	}
	return nil
}

func RunMigrations(conn *pgxpool.Pool) error {
	conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS schema_migrations (version text PRIMARY KEY, applied_at timestamptz)")

	migration_files, listdir_err := migrations.FS.ReadDir(".")

	if listdir_err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list migrations: %v\n", listdir_err)
		return listdir_err
	}

	migration_files = slices.DeleteFunc(migration_files, func(migration_file fs.DirEntry) bool {
		if !constants.MigrationRegex.MatchString(migration_file.Name()) {
			fmt.Fprintf(os.Stderr, "Found invalid migration file: %v\n", migration_file.Name())
			return true
		}

		return false
	})

	sort.Slice(migration_files, func(i, j int) bool {
		return migration_files[i].Name() < migration_files[j].Name()
	})

	for _, migration_file := range migration_files {
		if migration_file.Name() == "demo.sql" {
			continue
		}

		migration_name, _, _ := strings.Cut(migration_file.Name(), ".sql")
		var version string
		err := conn.QueryRow(context.Background(), "SELECT version FROM schema_migrations WHERE version = $1", migration_name).Scan(&version)
		if err == nil {
			continue // already applied
		}
		if err != pgx.ErrNoRows {
			return fmt.Errorf("check %s: %w", migration_name, err)
		}

		if err := applyMigration(context.Background(), conn, migration_file.Name()); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "applied migration: %s\n", migration_name)
	}

	if helpers.GetEnvFallback("DEMO_MODE", "false") == "true" {
		var version string
		err := conn.QueryRow(context.Background(), "SELECT version FROM schema_migrations WHERE version = $1", "demo").Scan(&version)
		if err != nil {
			if err := applyMigration(context.Background(), conn, "demo.sql"); err != nil {
				return err
			}
		}
	}

	return nil
}
