package goflow

import "fmt"

type JobReaderWriter interface {
	Get(id string) (Job, error)
	Update(job Job) error
}

func StartJob(store JobReaderWriter, id string) error {
	job, err := store.Get(id)
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

	if err := store.Update(job); err != nil {
		return fmt.Errorf("start job %q: %w", id, err)
	}

	return nil
}
