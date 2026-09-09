package main

import (
	"context"
	"errors"
	"goflow/internal/goflow"
	"testing"
	"time"
)

type fakeWorkerStore struct {
	jobs map[string]goflow.Job
}

func newFakeWorkerStore(jobs ...goflow.Job) *fakeWorkerStore {
	store := &fakeWorkerStore{jobs: make(map[string]goflow.Job)}
	for _, job := range jobs {
		store.jobs[job.ID] = job
	}
	return store
}

func (f *fakeWorkerStore) Create(ctx context.Context, job goflow.Job) error {
	if _, exists := f.jobs[job.ID]; exists {
		return goflow.ErrJobAlreadyExists
	}
	f.jobs[job.ID] = job
	return nil
}

func (f *fakeWorkerStore) Get(ctx context.Context, id string) (goflow.Job, error) {
	job, exists := f.jobs[id]
	if !exists {
		return goflow.Job{}, goflow.ErrJobNotFound
	}
	return job, nil
}

func (f *fakeWorkerStore) List(ctx context.Context) ([]goflow.Job, error) {
	jobs := make([]goflow.Job, 0, len(f.jobs))
	for _, job := range f.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (f *fakeWorkerStore) Update(ctx context.Context, job goflow.Job) error {
	if _, exists := f.jobs[job.ID]; !exists {
		return goflow.ErrJobNotFound
	}
	f.jobs[job.ID] = job
	return nil
}

func (f *fakeWorkerStore) Delete(ctx context.Context, id string) error {
	if _, exists := f.jobs[id]; !exists {
		return goflow.ErrJobNotFound
	}
	delete(f.jobs, id)
	return nil
}

func TestProcessQueuedJobCompletesJob(t *testing.T) {
	store := newFakeWorkerStore(goflow.Job{
		ID:          "job-1",
		Type:        "email",
		Status:      goflow.StatusPending,
		MaxAttempts: 3,
	})

	executorCalled := false
	executor := func(ctx context.Context, job goflow.Job) error {
		executorCalled = true
		if job.Status != goflow.StatusRunning {
			t.Fatalf("expected executor to receive running job, got %q", job.Status)
		}
		return nil
	}

	err := processQueuedJob(context.Background(), store, "job-1", executor, time.Now)
	if err != nil {
		t.Fatalf("processQueuedJob() returned error: %v", err)
	}
	if !executorCalled {
		t.Fatal("expected executor to be called")
	}

	job, err := store.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if job.Status != goflow.StatusCompleted {
		t.Fatalf("expected status %q, got %q", goflow.StatusCompleted, job.Status)
	}
}

func TestProcessQueuedJobRetriesTemporaryFailure(t *testing.T) {
	now := time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC)
	store := newFakeWorkerStore(goflow.Job{
		ID:          "job-1",
		Type:        "temporary-fail",
		Status:      goflow.StatusPending,
		Attempts:    1,
		MaxAttempts: 3,
	})

	executor := func(ctx context.Context, job goflow.Job) error {
		return goflow.JobExecutionError{Temporary: true, Message: "smtp unavailable"}
	}

	err := processQueuedJob(context.Background(), store, "job-1", executor, func() time.Time { return now })
	if err == nil {
		t.Fatal("expected processQueuedJob() to return execution error")
	}

	job, err := store.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if job.Status != goflow.StatusPending {
		t.Fatalf("expected status %q, got %q", goflow.StatusPending, job.Status)
	}
	if job.Attempts != 2 {
		t.Fatalf("expected attempts 2, got %d", job.Attempts)
	}
	if job.LastError != "smtp unavailable" {
		t.Fatalf("expected last error to be recorded, got %q", job.LastError)
	}
	if !job.AvailableAt.Equal(now.Add(4 * time.Second)) {
		t.Fatalf("expected available_at to be delayed, got %v", job.AvailableAt)
	}
}

func TestProcessQueuedJobDeadLettersPermanentFailure(t *testing.T) {
	store := newFakeWorkerStore(goflow.Job{
		ID:          "job-1",
		Type:        "permanent-fail",
		Status:      goflow.StatusPending,
		MaxAttempts: 3,
	})

	executor := func(ctx context.Context, job goflow.Job) error {
		return goflow.JobExecutionError{Temporary: false, Message: "invalid payload"}
	}

	err := processQueuedJob(context.Background(), store, "job-1", executor, time.Now)
	if err == nil {
		t.Fatal("expected processQueuedJob() to return execution error")
	}

	job, err := store.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if job.Status != goflow.StatusDeadLetter {
		t.Fatalf("expected status %q, got %q", goflow.StatusDeadLetter, job.Status)
	}
	if job.Attempts != 1 {
		t.Fatalf("expected attempts 1, got %d", job.Attempts)
	}
	if job.LastError != "invalid payload" {
		t.Fatalf("expected last error to be recorded, got %q", job.LastError)
	}
}

func TestProcessQueuedJobReturnsContextCancellation(t *testing.T) {
	store := newFakeWorkerStore(goflow.Job{
		ID:          "job-1",
		Type:        "email",
		Status:      goflow.StatusPending,
		MaxAttempts: 3,
	})

	executor := func(ctx context.Context, job goflow.Job) error {
		return context.Canceled
	}

	err := processQueuedJob(context.Background(), store, "job-1", executor, time.Now)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	job, err := store.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if job.Status != goflow.StatusRunning {
		t.Fatalf("expected canceled job to remain %q, got %q", goflow.StatusRunning, job.Status)
	}
}

func TestExecuteJobSuccess(t *testing.T) {
	err := executeJob(context.Background(), goflow.Job{Type: "email"})
	if err != nil {
		t.Fatalf("executeJob() returned error: %v", err)
	}
}

func TestExecuteJobTemporaryFailure(t *testing.T) {
	err := executeJob(context.Background(), goflow.Job{Type: "temporary-fail"})

	var jobErr goflow.JobExecutionError
	if !errors.As(err, &jobErr) {
		t.Fatalf("expected JobExecutionError, got %v", err)
	}
	if !jobErr.Temporary {
		t.Fatal("expected temporary execution error")
	}
}

func TestExecuteJobPermanentFailure(t *testing.T) {
	err := executeJob(context.Background(), goflow.Job{Type: "permanent-fail"})

	var jobErr goflow.JobExecutionError
	if !errors.As(err, &jobErr) {
		t.Fatalf("expected JobExecutionError, got %v", err)
	}
	if jobErr.Temporary {
		t.Fatal("expected permanent execution error")
	}
}
