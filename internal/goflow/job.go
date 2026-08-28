package goflow

import "fmt"

type InvalidJobStatusError struct {
	JobID  string
	Status JobStatus
}

func (e InvalidJobStatusError) Error() string {
	return fmt.Sprintf("cannot start job %s from status %s", e.JobID, e.Status)
}

func (j Job) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}

func (j *Job) MarkRunning() {
	j.Status = StatusRunning
}
