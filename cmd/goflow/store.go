package main

import "errors"

var ErrJobNotFound = errors.New("job not found")
var ErrJobAlreadyExists = errors.New("job already exists")

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Payload     []byte    `json:"payload"`
	Status      JobStatus `json:"status"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
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
		return ErrJobAlreadyExists
	}

	s.jobs[job.ID] = job
	return nil
}

func (s *Store) Get(id string) (Job, error) {
	job, ok := s.jobs[id]
	if !ok {
		return Job{}, ErrJobNotFound
	}

	return job, nil
}

func (s *Store) List() []Job {
	jobs := make([]Job, 0, len(s.jobs))

	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

func (s *Store) Update(job Job) error {
	if _, exists := s.jobs[job.ID]; !exists {
		return ErrJobNotFound
	}

	s.jobs[job.ID] = job
	return nil
}

func (s *Store) Delete(id string) error {
	if _, exists := s.jobs[id]; !exists {
		return ErrJobNotFound
	}

	delete(s.jobs, id)
	return nil
}

func (s *Store) Save(filename string) error {
	jobs := s.List()
	return saveJobs(filename, jobs)
}

func LoadStore(filename string) (*Store, error) {
	jobs, err := loadJobs(filename)
	if err != nil {
		return nil, err
	}

	store := NewStore()

	for _, job := range jobs {
		store.jobs[job.ID] = job
	}

	return store, nil
}
