package goflow

import (
	"context"
	"errors"
	"testing"
)

type fakeJobStore struct {
	jobs      map[string]Job
	job       Job
	getErr    error
	updateErr error
	updated   Job
	updateHit bool
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
		jobs: make(map[string]Job),
		job: Job{
			ID:     "job-1",
			Status: StatusPending,
		},
	}

	err := StartJob(context.Background(), store, "job-1")
	if err != nil {
		t.Fatalf("StartJob() returned error: %v", err)
	}

	if !store.updateHit {
		t.Fatal("expected store.Update() to be called")
	}

	if store.updated.Status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, store.updated.Status)
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

func TestStartJobUpdateFailure(t *testing.T) {
	store := &fakeJobStore{
		jobs: make(map[string]Job),
		job: Job{
			ID:     "job-3",
			Status: StatusPending,
		},
		updateErr: ErrJobNotFound,
	}

	err := StartJob(context.Background(), store, "job-3")
	if !errors.Is(err, ErrJobNotFound) {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}

	if !store.updateHit {
		t.Fatal("expected store.Update() to be called")
	}
}
