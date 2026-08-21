package main

func (j Job) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}

func (j *Job) MarkRunning() {
	j.Status = StatusRunning
}
