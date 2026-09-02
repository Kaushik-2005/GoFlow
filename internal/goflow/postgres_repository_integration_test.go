package goflow

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("sql.Open() returned error: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("PingContext() returned error: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			job_type TEXT NOT NULL,
			payload BYTEA,
			status TEXT NOT NULL,
			attempts INTEGER NOT NULL DEFAULT 0,
			max_attempts INTEGER NOT NULL DEFAULT 3,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("create table returned error: %v", err)
	}

	repo := NewPostgresRepository(db)
	jobID := "integration-job"
	_, _ = db.ExecContext(ctx, `DELETE FROM jobs WHERE id = $1`, jobID)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM jobs WHERE id = $1`, jobID)
	})

	job := Job{ID: jobID, Type: "email", Payload: []byte("payload"), Status: StatusPending, Attempts: 0, MaxAttempts: 3}
	if err := repo.Create(ctx, job); err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	got, err := repo.Get(ctx, jobID)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.ID != jobID || got.Type != "email" {
		t.Fatalf("unexpected job after create: %+v", got)
	}

	jobs, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}
	found := false
	for _, listed := range jobs {
		if listed.ID == jobID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected job %q to appear in List()", jobID)
	}

	job.Status = StatusRunning
	job.Attempts = 1
	if err := repo.Update(ctx, job); err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	updated, err := repo.Get(ctx, jobID)
	if err != nil {
		t.Fatalf("Get() after Update() returned error: %v", err)
	}
	if updated.Status != StatusRunning || updated.Attempts != 1 {
		t.Fatalf("unexpected updated job: %+v", updated)
	}

	if err := repo.Delete(ctx, jobID); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	_, err = repo.Get(ctx, jobID)
	if err == nil {
		t.Fatal("expected Get() after Delete() to return an error")
	}
}
