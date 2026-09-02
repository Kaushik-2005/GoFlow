package goflow

import (
	"context"
	"fmt"
)

type JobStore interface {
	Create(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
	List(ctx context.Context) ([]Job, error)
	Update(ctx context.Context, job Job) error
	Delete(ctx context.Context, id string) error
}

func StartJob(ctx context.Context, store JobStore, id string) error {
	job, err := store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("start job %q: %w", id, err)
	}

	if job.Status != StatusPending {
		return InvalidJobStatusError{
			JobID:  job.ID,
			Status: job.Status,
		}
	}

	job.MarkRunning()

	if err := store.Update(ctx, job); err != nil {
		return fmt.Errorf("start job %q: %w", id, err)
	}

	return nil
}
