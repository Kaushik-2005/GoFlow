package goflow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// PostgresRepository stores jobs in PostgreSQL through database/sql.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository builds a PostgreSQL-backed job repository.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, job Job) error {
	const query = `
		INSERT INTO jobs (id, job_type, payload, status, attempts, max_attempts)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxAttempts,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			return ErrJobAlreadyExists
		}

		return fmt.Errorf("create job %q: %w", job.ID, err)
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Job, error) {
	const query = `
		SELECT id, job_type, payload, status, attempts, max_attempts
		FROM jobs
		WHERE id = $1
	`

	job, err := scanJob(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrJobNotFound
		}

		return Job{}, fmt.Errorf("get job %q: %w", id, err)
	}

	return job, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Job, error) {
	const query = `
		SELECT id, job_type, payload, status, attempts, max_attempts
		FROM jobs
		ORDER BY created_at ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]Job, 0)
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, fmt.Errorf("list jobs: %w", err)
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	return jobs, nil
}

func (r *PostgresRepository) Update(ctx context.Context, job Job) error {
	const query = `
		UPDATE jobs
		SET job_type = $2, payload = $3, status = $4, attempts = $5, max_attempts = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Attempts,
		job.MaxAttempts,
	)
	if err != nil {
		return fmt.Errorf("update job %q: %w", job.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update job %q: %w", job.ID, err)
	}
	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM jobs
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete job %q: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete job %q: %w", id, err)
	}
	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanJob(s scanner) (Job, error) {
	var job Job
	var status string

	err := s.Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&status,
		&job.Attempts,
		&job.MaxAttempts,
	)
	if err != nil {
		return Job{}, err
	}

	job.Status = JobStatus(status)
	return job, nil
}
