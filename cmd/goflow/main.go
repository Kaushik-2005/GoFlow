package main

import (
	"errors"
	"fmt"
	"goflow/internal/goflow"
	"net/http"
	"os"
)

const jobsFile = "jobs.json"

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
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		store, err := goflow.LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		jobs := store.List()
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

		store, err := goflow.LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		job := goflow.Job{
			ID:          fmt.Sprintf("job-%d", len(store.List())+1),
			Type:        jobType,
			Status:      goflow.StatusPending,
			Attempts:    0,
			MaxAttempts: 3,
		}

		if err := store.Create(job); err != nil {
			if errors.Is(err, goflow.ErrJobAlreadyExists) {
				fmt.Printf("job already exists: %s\n", job.ID)
				return
			}
			fmt.Printf("failed to create job: %v\n", err)
			return
		}

		if err := store.Save(jobsFile); err != nil {
			fmt.Printf("failed to save jobs: %v\n", err)
			return
		}

		fmt.Printf("created job: %s\n", job.ID)
	case "get":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}

		store, err := goflow.LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		job, err := store.Get(jobID)
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

		store, err := goflow.LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		err = goflow.StartJob(store, jobID)
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

		if err := store.Save(jobsFile); err != nil {
			fmt.Printf("failed to save jobs: %v\n", err)
			return
		}

		fmt.Printf("processing job: %s\n", jobID)
	case "serve":
		mux := http.NewServeMux()
		mux.HandleFunc("/health/live", liveHandler)
		mux.HandleFunc("/v1/jobs", jobsHandler)
		mux.HandleFunc("/v1/jobs/", getJobHandler)

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
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("server error: %v\n", err)
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}
