package goflow

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

var testAvailableAt = time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)

func newMockRepository(t *testing.T) (*PostgresRepository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() returned error: %v", err)
	}

	t.Cleanup(func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("db.Close() returned error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet SQL expectations: %v", err)
		}
	})

	return NewPostgresRepository(db), mock
}

func TestPostgresRepositoryCreate(t *testing.T) {
	tests := []struct {
		name      string
		job       Job
		setupMock func(sqlmock.Sqlmock)
		checkErr  func(*testing.T, error)
	}{
		{
			name: "success",
			job:  Job{ID: "job-1", Type: "email", Payload: []byte("payload"), Status: StatusPending, Attempts: 0, MaxAttempts: 3, AvailableAt: testAvailableAt},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO jobs (id, job_type, payload, status, attempts, max_attempts, available_at, last_error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).
					WithArgs("job-1", "email", []byte("payload"), StatusPending, 0, 3, testAvailableAt, "").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("Create() returned error: %v", err)
				}
			},
		},
		{
			name: "duplicate key maps to ErrJobAlreadyExists",
			job:  Job{ID: "job-1", Type: "email", Status: StatusPending, Attempts: 0, MaxAttempts: 3, AvailableAt: testAvailableAt},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO jobs (id, job_type, payload, status, attempts, max_attempts, available_at, last_error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).
					WithArgs("job-1", "email", []byte(nil), StatusPending, 0, 3, testAvailableAt, "").
					WillReturnError(&pq.Error{Code: "23505"})
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, ErrJobAlreadyExists) {
					t.Fatalf("expected ErrJobAlreadyExists, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)
			err := repo.Create(context.Background(), tt.job)
			tt.checkErr(t, err)
		})
	}
}

func TestPostgresRepositoryGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(sqlmock.Sqlmock)
		check     func(*testing.T, Job, error)
	}{
		{
			name: "success",
			id:   "job-1",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "job_type", "payload", "status", "attempts", "max_attempts", "available_at", "last_error"}).
					AddRow("job-1", "email", []byte("payload"), "pending", 0, 3, testAvailableAt, "")
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, job_type, payload, status, attempts, max_attempts, available_at, last_error
		FROM jobs
		WHERE id = $1
	`)).WithArgs("job-1").WillReturnRows(rows)
			},
			check: func(t *testing.T, job Job, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("Get() returned error: %v", err)
				}
				if job.ID != "job-1" || job.Type != "email" || job.Status != StatusPending {
					t.Fatalf("unexpected job: %+v", job)
				}
			},
		},
		{
			name: "not found",
			id:   "job-404",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, job_type, payload, status, attempts, max_attempts, available_at, last_error
		FROM jobs
		WHERE id = $1
	`)).WithArgs("job-404").WillReturnError(sql.ErrNoRows)
			},
			check: func(t *testing.T, job Job, err error) {
				t.Helper()
				if !errors.Is(err, ErrJobNotFound) {
					t.Fatalf("expected ErrJobNotFound, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)
			job, err := repo.Get(context.Background(), tt.id)
			tt.check(t, job, err)
		})
	}
}

func TestPostgresRepositoryList(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		check     func(*testing.T, []Job, error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "job_type", "payload", "status", "attempts", "max_attempts", "available_at", "last_error"}).
					AddRow("job-1", "email", []byte("payload-1"), "pending", 0, 3, testAvailableAt, "").
					AddRow("job-2", "report", []byte("payload-2"), "running", 1, 3, testAvailableAt, "")
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, job_type, payload, status, attempts, max_attempts, available_at, last_error
		FROM jobs
		ORDER BY created_at ASC, id ASC
	`)).WillReturnRows(rows)
			},
			check: func(t *testing.T, jobs []Job, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("List() returned error: %v", err)
				}
				if len(jobs) != 2 {
					t.Fatalf("expected 2 jobs, got %d", len(jobs))
				}
				if jobs[1].Status != StatusRunning {
					t.Fatalf("expected second job status %q, got %q", StatusRunning, jobs[1].Status)
				}
			},
		},
		{
			name: "query failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, job_type, payload, status, attempts, max_attempts, available_at, last_error
		FROM jobs
		ORDER BY created_at ASC, id ASC
	`)).WillReturnError(sql.ErrConnDone)
			},
			check: func(t *testing.T, jobs []Job, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected List() to return an error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)
			jobs, err := repo.List(context.Background())
			tt.check(t, jobs, err)
		})
	}
}

func TestPostgresRepositoryClaimPending(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(sqlmock.Sqlmock)
		wantClaimed bool
		wantErr     bool
	}{
		{
			name: "success",
			id:   "job-1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE jobs
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		AND status = $3
	`)).WithArgs("job-1", StatusRunning, StatusPending).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantClaimed: true,
		},
		{
			name: "not pending",
			id:   "job-1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE jobs
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		AND status = $3
	`)).WithArgs("job-1", StatusRunning, StatusPending).WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "exec failure",
			id:   "job-1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE jobs
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		AND status = $3
	`)).WithArgs("job-1", StatusRunning, StatusPending).WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)

			claimed, err := repo.ClaimPending(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected ClaimPending() to return an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ClaimPending() returned error: %v", err)
			}
			if claimed != tt.wantClaimed {
				t.Fatalf("expected claimed %v, got %v", tt.wantClaimed, claimed)
			}
		})
	}
}

func TestPostgresRepositoryUpdate(t *testing.T) {
	tests := []struct {
		name      string
		job       Job
		setupMock func(sqlmock.Sqlmock)
		checkErr  func(*testing.T, error)
	}{
		{
			name: "success",
			job:  Job{ID: "job-1", Type: "email", Status: StatusRunning, Attempts: 1, MaxAttempts: 3, AvailableAt: testAvailableAt},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE jobs
		SET job_type = $2, payload = $3, status = $4, attempts = $5, max_attempts = $6, available_at = $7, last_error = $8, updated_at = NOW()
		WHERE id = $1
	`)).WithArgs("job-1", "email", []byte(nil), StatusRunning, 1, 3, testAvailableAt, "").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("Update() returned error: %v", err)
				}
			},
		},
		{
			name: "not found",
			job:  Job{ID: "job-404", Type: "email", Status: StatusPending, Attempts: 0, MaxAttempts: 3, AvailableAt: testAvailableAt},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE jobs
		SET job_type = $2, payload = $3, status = $4, attempts = $5, max_attempts = $6, available_at = $7, last_error = $8, updated_at = NOW()
		WHERE id = $1
	`)).WithArgs("job-404", "email", []byte(nil), StatusPending, 0, 3, testAvailableAt, "").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, ErrJobNotFound) {
					t.Fatalf("expected ErrJobNotFound, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)
			err := repo.Update(context.Background(), tt.job)
			tt.checkErr(t, err)
		})
	}
}

func TestPostgresRepositoryDelete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(sqlmock.Sqlmock)
		checkErr  func(*testing.T, error)
	}{
		{
			name: "success",
			id:   "job-1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM jobs
		WHERE id = $1
	`)).WithArgs("job-1").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("Delete() returned error: %v", err)
				}
			},
		},
		{
			name: "not found",
			id:   "job-404",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM jobs
		WHERE id = $1
	`)).WithArgs("job-404").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, ErrJobNotFound) {
					t.Fatalf("expected ErrJobNotFound, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := newMockRepository(t)
			tt.setupMock(mock)
			err := repo.Delete(context.Background(), tt.id)
			tt.checkErr(t, err)
		})
	}
}
