package main

import (
	"context"
	"errors"
	"goflow/internal/goflow"
	"testing"
)

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
