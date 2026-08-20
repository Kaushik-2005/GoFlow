package main

import "errors"

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID          string
	Type        string
	Payload     []byte
	Status      JobStatus
	Attempts    int
	MaxAttempts int
}

type Store struct {
	jobs map[string]Job
}

func NewStore() *Store {
	return &Store{
		jobs: make(map[string]Job),
	}
}

func (s *Store) Create(job Job) error {
	if _, exists := s.jobs[job.ID]; exists {
		return errors.New("job already exists")
	}

	s.jobs[job.ID] = job
	return nil
}

func (s *Store) Get(id string) (Job, bool) {
	job, ok := s.jobs[id]
	return job, ok
}

func (s *Store) List() []Job {
	jobs := make([]Job, 0, len(s.jobs))

	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

func (s *Store) Update(job Job) bool {
	if _, exists := s.jobs[job.ID]; !exists {
		return false
	}

	s.jobs[job.ID] = job
	return true
}

func (s *Store) Delete(id string) bool {
	if _, exists := s.jobs[id]; !exists {
		return false
	}

	delete(s.jobs, id)
	return true
}
