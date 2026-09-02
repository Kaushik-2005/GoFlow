package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"goflow/internal/goflow"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const migrationFile = "migrations/001_create_jobs.sql"

func openPostgresStore(ctx context.Context) (*goflow.PostgresRepository, *sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	if err := applyMigrations(ctx, db); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("apply migrations: %w", err)
	}

	return goflow.NewPostgresRepository(db), db, nil
}

func applyMigrations(ctx context.Context, db *sql.DB) error {
	migrationSQL, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	if _, err := db.ExecContext(ctx, string(migrationSQL)); err != nil {
		return fmt.Errorf("run migration: %w", err)
	}

	return nil
}

func newJobID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate job id: %w", err)
	}

	return "job-" + hex.EncodeToString(bytes), nil
}
