package main

import (
	"context"
	"errors"
	"fmt"
	"goflow/internal/goflow"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

const jobQueueSize = 8
const workerCount = 3

var (
	version = "dev"
	commit  = "none"
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
	fmt.Println("Commands: list, create, get, process, migrate, serve, work")
	fmt.Println("Environment: DATABASE_URL must point to PostgreSQL")
	fmt.Printf("Version: %s (%s)\n", version, commit)
}

func commandNeedsConfig(command string) bool {
	switch command {
	case "list", "create", "get", "process", "migrate", "serve", "work":
		return true
	default:
		return false
	}
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	var cfg config
	if commandNeedsConfig(command) {
		loadedConfig, err := loadConfig()
		if err != nil {
			fmt.Printf("invalid config: %v\n", err)
			return
		}
		cfg = loadedConfig
	}

	switch command {
	case "list":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(ctx, cfg)
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

		store, db, err := openPostgresStore(ctx, cfg)
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
			AvailableAt: time.Now(),
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

		store, db, err := openPostgresStore(ctx, cfg)
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

		store, db, err := openPostgresStore(ctx, cfg)
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
	case "migrate":
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := migratePostgres(ctx, cfg); err != nil {
			fmt.Printf("migration failed: %v\n", err)
			return
		}

		fmt.Println("migration completed")
	case "serve":
		serverCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		startupCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)
		store, db, err := openPostgresStore(startupCtx, cfg)
		cancel()

		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}

		defer db.Close()

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		metrics := newMetrics()

		api := newAPIHandler(store)
		api.ready = store
		api.logger = logger
		api.metrics = metrics

		mux := http.NewServeMux()
		mux.HandleFunc("/health/live", api.liveHandler)
		mux.HandleFunc("/health/ready", api.readyHandler)
		mux.HandleFunc("/metrics", api.metricsHandler)
		mux.HandleFunc("/v1/jobs", api.jobsHandler)
		mux.HandleFunc("/v1/jobs/", api.jobByIDHandler)

		handler := chain(
			mux,
			requestIDMiddleware,
			requestLoggingMiddleware(logger, metrics),
			recoveryMiddleware,
			requestBodyLimitMiddleware(1<<20), // 1 MB limit
			requireJSONMiddleware,
		)

		server := &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		}

		go func() {
			fmt.Printf("starting server on %s\n", cfg.HTTPAddr)

			if err := server.ListenAndServe(); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("server error: %v\n", err)
			}
		}()

		<-serverCtx.Done()

		fmt.Println("server shutting down")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("server shutdown error: %v\n", err)
			return
		}

		fmt.Println("server stopped")

	case "work":
		workCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		startupCtx, cancel := context.WithTimeout(workCtx, 5*time.Second)
		defer cancel()

		store, db, err := openPostgresStore(startupCtx, cfg)
		if err != nil {
			fmt.Printf("failed to open store: %v\n", err)
			return
		}
		defer db.Close()

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		workerMetrics := newMetrics()

		queued := make(map[string]struct{})
		var queueMu sync.Mutex

		jobs := make(chan string, jobQueueSize)

		var wg sync.WaitGroup

		for workerID := 1; workerID <= workerCount; workerID++ {
			wg.Add(1)

			go func(id int) {
				defer wg.Done()

				for {
					select {
					case jobID, ok := <-jobs:
						if !ok {
							return
						}

						queueMu.Lock()
						delete(queued, jobID)
						queueMu.Unlock()

						workerMetrics.decrementQueueDepth()
						workerMetrics.incrementActiveWorkers()
						logger.InfoContext(workCtx, "job processing", "worker_id", id, "job_id", jobID)

						err := processQueuedJob(workCtx, store, jobID, executeJob, time.Now)
						workerMetrics.decrementActiveWorkers()
						if err != nil {
							if errors.Is(err, context.Canceled) {
								logger.InfoContext(workCtx, "worker stopping", "worker_id", id)
								return
							}

							workerMetrics.incrementJobsExecutionFailed()
							updatedJob, loadErr := store.Get(workCtx, jobID)
							if loadErr == nil {
								switch updatedJob.Status {
								case goflow.StatusPending:
									workerMetrics.incrementJobsRetried()
								case goflow.StatusDeadLetter:
									workerMetrics.incrementJobsDeadLettered()
								}
							}

							logger.WarnContext(workCtx, "job failed", "worker_id", id, "job_id", jobID, "error", err)
							continue
						}

						workerMetrics.incrementJobsCompleted()
						logger.InfoContext(workCtx, "job completed", "worker_id", id, "job_id", jobID)

					case <-workCtx.Done():
						logger.InfoContext(workCtx, "worker stopping", "worker_id", id)
						return
					}
				}
			}(workerID)
		}

		stopWorkers := func() {
			close(jobs)
			wg.Wait()
		}

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			jobsList, err := store.ListReadyJobs(workCtx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logger.InfoContext(workCtx, "work command stopping")
					stopWorkers()
					return
				}
				logger.ErrorContext(workCtx, "failed to list ready jobs", "error", err)
			} else {
				for _, job := range jobsList {

					queueMu.Lock()

					_, alreadyQueued := queued[job.ID]
					if alreadyQueued {
						queueMu.Unlock()
						continue
					}

					queued[job.ID] = struct{}{}
					queueMu.Unlock()

					select {
					case jobs <- job.ID:
						workerMetrics.incrementQueueDepth()
					case <-workCtx.Done():
						logger.InfoContext(workCtx, "work command stopping")
						stopWorkers()
						return
					}
				}
			}

			select {
			case <-ticker.C:
			case <-workCtx.Done():
				logger.InfoContext(workCtx, "work command stopping")
				stopWorkers()
				return
			}
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}
