package goflow

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeJobStore struct {
	jobs        map[string]Job
	job         Job
	getErr      error
	updateErr   error
	claimErr    error
	claimResult bool
	claimHit    bool
	updated     Job
	updateHit   bool
}

func (f *fakeJobStore) Get(ctx context.Context, id string) (Job, error) {
	if f.getErr != nil {
		return Job{}, f.getErr
	}
	if f.job.ID == id {
		return f.job, nil
	}
	return Job{}, ErrJobNotFound
}

func (f *fakeJobStore) Update(ctx context.Context, job Job) error {
	f.updated = job
	f.updateHit = true
	if f.updateErr != nil {
		return f.updateErr
	}
	return nil
}

func (f *fakeJobStore) ClaimPending(ctx context.Context, id string) (bool, error) {
	f.claimHit = true
	if f.claimErr != nil {
		return false, f.claimErr
	}
	return f.claimResult, nil
}

func (f *fakeJobStore) Create(ctx context.Context, job Job) error {
	if _, exists := f.jobs[job.ID]; exists {
		return ErrJobAlreadyExists
	}

	f.jobs[job.ID] = job
	return nil
}

func (f *fakeJobStore) List(ctx context.Context) ([]Job, error) {
	jobs := make([]Job, 0, len(f.jobs))
	for _, job := range f.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (f *fakeJobStore) Delete(ctx context.Context, id string) error {
	if _, exists := f.jobs[id]; !exists {
		return ErrJobNotFound
	}

	delete(f.jobs, id)
	return nil
}

func TestStartJobSuccess(t *testing.T) {
	store := &fakeJobStore{
		jobs:        make(map[string]Job),
		claimResult: true,
	}

	err := StartJob(context.Background(), store, "job-1")
	if err != nil {
		t.Fatalf("StartJob() returned error: %v", err)
	}

	if !store.claimHit {
		t.Fatal("expected store.ClaimPending() to be called")
	}
	if store.updateHit {
		t.Fatal("expected store.Update() not to be called")
	}
}

func TestStartJobInvalidStatus(t *testing.T) {
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:     "job-2",
			Status: StatusRunning,
		},
	}

	err := StartJob(context.Background(), store, "job-2")
	if err == nil {
		t.Fatal("expected StartJob() to return an error")
	}

	if !store.claimHit {
		t.Fatal("expected store.ClaimPending() to be called")
	}
	if store.updateHit {
		t.Fatal("expected store.Update() not to be called")
	}

	var invalidStatusErr InvalidJobStatusError
	if !errors.As(err, &invalidStatusErr) {
		t.Fatalf("expected InvalidJobStatusError, got %v", err)
	}

	if invalidStatusErr.JobID != "job-2" {
		t.Errorf("expected job id %q, got %q", "job-2", invalidStatusErr.JobID)
	}

	if invalidStatusErr.Status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, invalidStatusErr.Status)
	}
}

func TestStartJobNotFound(t *testing.T) {
	store := &fakeJobStore{
		jobs:   make(map[string]Job),
		getErr: ErrJobNotFound,
	}

	err := StartJob(context.Background(), store, "missing")
	if !errors.Is(err, ErrJobNotFound) {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}

	if store.updateHit {
		t.Fatal("expected store.Update() not to be called")
	}
}

func TestStartJobClaimFailure(t *testing.T) {
	store := &fakeJobStore{
		jobs:     make(map[string]Job),
		claimErr: ErrJobNotFound,
	}

	err := StartJob(context.Background(), store, "job-3")
	if !errors.Is(err, ErrJobNotFound) {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}

	if !store.claimHit {
		t.Fatal("expected store.ClaimPending() to be called")
	}
	if store.updateHit {
		t.Fatal("expected store.Update() not to be called")
	}
}

func TestCompleteJobSuccess(t *testing.T) {
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:        "job-1",
			Status:    StatusRunning,
			LastError: "previous error",
		},
	}

	err := CompleteJob(context.Background(), store, "job-1")
	if err != nil {
		t.Fatalf("CompleteJob() returned error: %v", err)
	}

	if !store.updateHit {
		t.Fatal("expected store.Update() to be called")
	}
	if store.updated.Status != StatusCompleted {
		t.Fatalf("expected status %q, got %q", StatusCompleted, store.updated.Status)
	}
	if store.updated.LastError != "" {
		t.Fatalf("expected last error to be cleared, got %q", store.updated.LastError)
	}
}

func TestCompleteJobInvalidStatus(t *testing.T) {
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:     "job-1",
			Status: StatusPending,
		},
	}

	err := CompleteJob(context.Background(), store, "job-1")
	if err == nil {
		t.Fatal("expected CompleteJob() to return an error")
	}
	if store.updateHit {
		t.Fatal("expected store.Update() not to be called")
	}
}

func TestFailJobTemporaryFailureWithRetriesLeft(t *testing.T) {
	now := time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:          "job-1",
			Status:      StatusRunning,
			Attempts:    1,
			MaxAttempts: 3,
		},
	}

	err := FailJob(context.Background(), store, "job-1", JobExecutionError{Temporary: true, Message: "smtp unavailable"}, now)
	if err != nil {
		t.Fatalf("FailJob() returned error: %v", err)
	}

	if store.updated.Status != StatusPending {
		t.Fatalf("expected status %q, got %q", StatusPending, store.updated.Status)
	}
	if store.updated.Attempts != 2 {
		t.Fatalf("expected attempts 2, got %d", store.updated.Attempts)
	}
	if store.updated.LastError != "smtp unavailable" {
		t.Fatalf("expected last error to be recorded, got %q", store.updated.LastError)
	}
	if !store.updated.AvailableAt.Equal(now.Add(4 * time.Second)) {
		t.Fatalf("expected available_at to use retry delay, got %v", store.updated.AvailableAt)
	}
}

func TestFailJobTemporaryFailureMaxAttemptsReached(t *testing.T) {
	now := time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:          "job-1",
			Status:      StatusRunning,
			Attempts:    2,
			MaxAttempts: 3,
		},
	}

	err := FailJob(context.Background(), store, "job-1", JobExecutionError{Temporary: true, Message: "smtp unavailable"}, now)
	if err != nil {
		t.Fatalf("FailJob() returned error: %v", err)
	}

	if store.updated.Status != StatusDeadLetter {
		t.Fatalf("expected status %q, got %q", StatusDeadLetter, store.updated.Status)
	}
	if store.updated.Attempts != 3 {
		t.Fatalf("expected attempts 3, got %d", store.updated.Attempts)
	}
}

func TestFailJobPermanentFailure(t *testing.T) {
	now := time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:          "job-1",
			Status:      StatusRunning,
			Attempts:    0,
			MaxAttempts: 3,
		},
	}

	err := FailJob(context.Background(), store, "job-1", JobExecutionError{Temporary: false, Message: "invalid payload"}, now)
	if err != nil {
		t.Fatalf("FailJob() returned error: %v", err)
	}

	if store.updated.Status != StatusDeadLetter {
		t.Fatalf("expected status %q, got %q", StatusDeadLetter, store.updated.Status)
	}
	if store.updated.Attempts != 1 {
		t.Fatalf("expected attempts 1, got %d", store.updated.Attempts)
	}
	if store.updated.LastError != "invalid payload" {
		t.Fatalf("expected last error to be recorded, got %q", store.updated.LastError)
	}
}
