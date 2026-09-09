# GoFlow

GoFlow is a staged learning project for building a production-style Go backend that accepts jobs, stores them, processes them asynchronously, retries transient failures, and surfaces operational behavior with structured logs.

The project follows `Go_Industry_Roadmap_4_Weeks.md`. Week 3 implementation is complete: the app now has a PostgreSQL-backed worker pool with retry scheduling, failure classification, dead-letter status, deterministic worker tests, and race-detector validation.

## Current status

- Week 1 foundations completed
- Week 2 HTTP, architecture, PostgreSQL, and testing completed
- Week 3 worker pool, graceful shutdown, retry handling, and concurrency testing completed through Module 3.6
- Module 4.1 structured logging completed with JSON `slog` output
- Current next step: Module 4.2 observability

## Requirements

- Go 1.26+
- PostgreSQL
- `DATABASE_URL`

Example:

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/goflow?sslmode=disable"
```

## Run

CLI:

```powershell
go run ./cmd/goflow create email
go run ./cmd/goflow list
go run ./cmd/goflow get <job-id>
go run ./cmd/goflow process <job-id>
```

HTTP server:

```powershell
go run ./cmd/goflow serve
```

Worker command:

```powershell
go run ./cmd/goflow work
```

Test job types:

- `email` and `report` complete successfully.
- `temporary-fail` retries with `available_at` backoff until attempts are exhausted.
- `permanent-fail` moves to `dead_letter` after the failed attempt.

Current HTTP API:

- `GET /health/live`
- `GET /v1/jobs`
- `GET /v1/jobs?status=pending`
- `GET /v1/jobs/{id}`
- `POST /v1/jobs`
- `DELETE /v1/jobs/{id}`

The app automatically applies `migrations/001_create_jobs.sql` during startup.

## Validation

```powershell
gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/http_test.go ./cmd/goflow/worker.go ./cmd/goflow/worker_test.go
gofmt -w ./internal/goflow/service.go ./internal/goflow/service_test.go ./internal/goflow/postgres_repository.go ./internal/goflow/postgres_repository_test.go
go vet ./...
go test ./...
go test --% -coverprofile=coverage.out ./...
go test -race ./...
```





