package main

import (
	"context"
	"fmt"
	"goflow/internal/goflow"
	"time"
)

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
