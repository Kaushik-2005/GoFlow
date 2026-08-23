package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func saveJobs(filename string, jobs []Job) error {
	data, err := json.MarshalIndent(jobs, "", "  ")

	if err != nil {
		return fmt.Errorf("save jobs: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("save jobs: %w", err)
	}

	return nil
}

func loadJobs(filename string) ([]Job, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Job{}, nil
		}
		return nil, fmt.Errorf("load jobs: %w", err)
	}

	var jobs []Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, fmt.Errorf("load jobs: %w", err)
	}

	return jobs, nil
}
