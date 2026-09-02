package main

import (
	"context"
	"errors"
	"fmt"
	"goflow/internal/goflow"
	"net/http"
	"os"
	"time"
)

func requireArg(args []string, index int, message string) (string, bool) {
	if len(args) <= index {
		fmt.Println(message)
		return "", false
	}
	return args[index], true
}

func printHelp() {
	fmt.Println("Usage: goflow <command>")
	fmt.Println("Commands: list, create, get, process, serve")
	fmt.Println("Environment: DATABASE_URL must point to PostgreSQL")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(ctx)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		jobs, err := store.List(ctx)
		if err != nil {
			fmt.Printf("failed to list jobs: %v\n", err)
			return
		}
		if len(jobs) == 0 {
			fmt.Println("no jobs found")
			return
		}
		for _, job := range jobs {
			fmt.Printf("%s %s %s\n", job.ID, job.Type, job.Status)
		}
	case "create":
		jobType, ok := requireArg(os.Args, 2, "missing job type")
		if !ok {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(ctx)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		jobID, err := newJobID()
		if err != nil {
			fmt.Printf("failed to create job id: %v\n", err)
			return
		}

		job := goflow.Job{
			ID:          jobID,
			Type:        jobType,
			Status:      goflow.StatusPending,
			Attempts:    0,
			MaxAttempts: 3,
		}

		if err := store.Create(ctx, job); err != nil {
			if errors.Is(err, goflow.ErrJobAlreadyExists) {
				fmt.Printf("job already exists: %s\n", job.ID)
				return
			}
			fmt.Printf("failed to create job: %v\n", err)
			return
		}

		fmt.Printf("created job: %s\n", job.ID)
	case "get":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(ctx)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		job, err := store.Get(ctx, jobID)
		if err != nil {
			if errors.Is(err, goflow.ErrJobNotFound) {
				fmt.Printf("job not found: %s\n", jobID)
				return
			}
			fmt.Printf("failed to get job %s: %v\n", jobID, err)
			return
		}

		fmt.Printf("%s %s %s\n", job.ID, job.Type, job.Status)
	case "process":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(ctx)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		err = goflow.StartJob(ctx, store, jobID)
		if err != nil {
			if errors.Is(err, goflow.ErrJobNotFound) {
				fmt.Printf("job not found: %s\n", jobID)
				return
			}

			var invalidStatusErr goflow.InvalidJobStatusError
			if errors.As(err, &invalidStatusErr) {
				fmt.Printf("cannot process job %s: current status is %s\n", invalidStatusErr.JobID, invalidStatusErr.Status)
				return
			}

			fmt.Printf("failed to start job %s: %v\n", jobID, err)
			return
		}

		fmt.Printf("processing job: %s\n", jobID)
	case "serve":
		startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(startupCtx)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		api := newAPIHandler(store)

		mux := http.NewServeMux()
		mux.HandleFunc("/health/live", api.liveHandler)
		mux.HandleFunc("/v1/jobs", api.jobsHandler)
		mux.HandleFunc("/v1/jobs/", api.jobByIDHandler)

		handler := chain(
			mux,
			requestIDMiddleware,
			recoveryMiddleware,
			requestBodyLimitMiddleware(1<<20),
			requireJSONMiddleware,
		)

		server := &http.Server{
			Addr:    ":8080",
			Handler: handler,
		}

		fmt.Println("starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server error: %v\n", err)
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}
