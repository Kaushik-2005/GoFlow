package main

import (
	"context"
	"errors"
	"fmt"
	"goflow/internal/goflow"
	"time"
)

type jobExecutor func(context.Context, goflow.Job) error
type clockFunc func() time.Time

func processQueuedJob(
	ctx context.Context,
	store goflow.JobStore,
	jobID string,
	executor jobExecutor,
	now clockFunc,
) error {
	if err := goflow.StartJob(ctx, store, jobID); err != nil {
		return fmt.Errorf("start job %q: %w", jobID, err)
	}

	job, err := store.Get(ctx, jobID)
	if err != nil {
		return fmt.Errorf("load job %q: %w", jobID, err)
	}

	if err := executor(ctx, job); err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}

		if failErr := goflow.FailJob(ctx, store, jobID, err, now()); failErr != nil {
			return fmt.Errorf("record job failure %q: %w", jobID, failErr)
		}

		return err
	}

	if err := goflow.CompleteJob(ctx, store, jobID); err != nil {
		return fmt.Errorf("complete job %q: %w", jobID, err)
	}

	return nil
}

func executeJob(ctx context.Context, job goflow.Job) error {
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}

	switch job.Type {
	case "email", "report":
		return nil
	case "temporary-fail":
		return goflow.JobExecutionError{Temporary: true, Message: "temporary job failure"}
	case "permanent-fail":
		return goflow.JobExecutionError{Temporary: false, Message: "permanent job failure"}
	default:
		return goflow.JobExecutionError{Temporary: false, Message: fmt.Sprintf("unsupported job type %q", job.Type)}
	}
}
