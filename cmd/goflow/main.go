package main

import (
	"errors"
	"fmt"
	"os"
)

const jobsFile = "jobs.json"

func maxPriority(priorities []int) int {
	max := priorities[0]

	for _, priority := range priorities {
		if priority > max {
			max = priority
		}
	}

	return max
}

func countByStatus(statuses []string) map[string]int {
	counts := make(map[string]int)

	for _, status := range statuses {
		counts[status]++
	}

	return counts
}

func isValidJobName(name string) bool {
	if len(name) < 3 {
		return false
	}

	return true
}

func retryDelay(attempt int) int {
	return 1 << attempt
}

func filterCompleted(statuses []string) []string {
	completed := []string{}

	for _, status := range statuses {
		if status == "completed" {
			completed = append(completed, status)
		}
	}

	return completed
}

func requireArg(args []string, index int, message string) (string, bool) {
	if len(args) <= index {
		fmt.Println(message)
		return "", false
	}
	return args[index], true
}

func printHelp() {
	fmt.Println("Usage: goflow <command>")
	fmt.Println("Commands: list, create, get, process")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		store, err := LoadStore(jobsFile)
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

		store, err := LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		job := Job{
			ID:          fmt.Sprintf("job-%d", len(store.List())+1),
			Type:        jobType,
			Status:      StatusPending,
			Attempts:    0,
			MaxAttempts: 3,
		}

		if err := store.Create(job); err != nil {
			if errors.Is(err, ErrJobAlreadyExists) {
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

		store, err := LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		job, err := store.Get(jobID)
		if err != nil {
			if errors.Is(err, ErrJobNotFound) {
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

		store, err := LoadStore(jobsFile)
		if err != nil {
			fmt.Printf("failed to load jobs: %v\n", err)
			return
		}

		err = startJob(store, jobID)
		if err != nil {
			if errors.Is(err, ErrJobNotFound) {
				fmt.Printf("job not found: %s\n", jobID)
				return
			}

			var invalidStatusErr InvalidJobStatusError
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
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}
