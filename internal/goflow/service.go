package goflow

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const baseRetryDelay = time.Second

type JobStore interface {
	Create(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
	List(ctx context.Context) ([]Job, error)
	Update(ctx context.Context, job Job) error
	Delete(ctx context.Context, id string) error
	ClaimPending(ctx context.Context, id string) (bool, error)
}

type JobExecutionError struct {
	Temporary bool
	Message   string
}

func (e JobExecutionError) Error() string {
	return e.Message
}

func StartJob(ctx context.Context, store JobStore, id string) error {
	claimed, err := store.ClaimPending(ctx, id)
	if err != nil {
		return fmt.Errorf("start job %q: %w", id, err)
	}
	if claimed {
		return nil
	}

	job, err := store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("start job %q: %w", id, err)
	}

	return InvalidJobStatusError{
		JobID:  job.ID,
		Status: job.Status,
	}
}

func CompleteJob(ctx context.Context, store JobStore, id string) error {
	job, err := store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("complete job %q: %w", id, err)
	}

	if job.Status != StatusRunning {
		return InvalidJobStatusError{
			JobID:  job.ID,
			Status: job.Status,
		}
	}

	job.Status = StatusCompleted
	job.LastError = ""

	if err := store.Update(ctx, job); err != nil {
		return fmt.Errorf("complete job %q: %w", id, err)
	}

	return nil
}

func FailJob(ctx context.Context, store JobStore, id string, executionErr error, now time.Time) error {
	if executionErr == nil {
		return fmt.Errorf("fail job %q: missing execution error", id)
	}

	job, err := store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("fail job %q: %w", id, err)
	}

	if job.Status != StatusRunning {
		return InvalidJobStatusError{
			JobID:  job.ID,
			Status: job.Status,
		}
	}

	job.Attempts++
	job.LastError = executionErr.Error()

	var jobErr JobExecutionError
	isTemporary := errors.As(executionErr, &jobErr) && jobErr.Temporary
	if isTemporary && job.Attempts < job.MaxAttempts {
		job.Status = StatusPending
		job.AvailableAt = now.Add(retryDelay(job.Attempts))
	} else {
		job.Status = StatusDeadLetter
	}

	if err := store.Update(ctx, job); err != nil {
		return fmt.Errorf("fail job %q: %w", id, err)
	}

	return nil
}

func retryDelay(attempt int) time.Duration {
	return baseRetryDelay * time.Duration(1<<attempt)
}
