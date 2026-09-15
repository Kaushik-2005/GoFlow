# GoFlow

GoFlow is a staged learning project for building a production-style Go backend that accepts jobs, stores them, processes them asynchronously, retries transient failures, and surfaces operational behavior with structured logs.

The project follows `Go_Industry_Roadmap_4_Weeks.md`. Week 3 implementation is complete: the app now has a PostgreSQL-backed worker pool with retry scheduling, failure classification, dead-letter status, deterministic worker tests, and race-detector validation.

## Current status

- Week 1 foundations completed
- Week 2 HTTP, architecture, PostgreSQL, and testing completed
- Week 3 worker pool, graceful shutdown, retry handling, and concurrency testing completed through Module 3.6
- Module 4.1 structured logging completed with JSON `slog` output
- Module 4.2 observability completed with readiness and JSON metrics
- Module 4.3 profiling and performance completed with benchmark/profile baseline
- Module 4.4 security completed with stronger input validation, HTTP timeouts, dependency scanning, and security notes
- Module 4.5 containers and configuration completed with explicit config, Dockerfile, Compose stack, migration command, and runtime validation
- Current next step: Module 4.6 CI and engineering workflow

## Requirements

- Go 1.26+
- PostgreSQL
- `DATABASE_URL`

Example:

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/goflow?sslmode=disable"
$env:HTTP_ADDR = ":8080"
$env:MIGRATION_FILE = "migrations/001_create_jobs.sql"
```

Copy `.env.example` when you need a local template, but do not commit real `.env` files.


## Security Notes

- GoFlow currently has no authentication or authorization.
- Treat the HTTP API as local/internal only.
- Do not expose it directly to the public internet.
- Store `DATABASE_URL` in the runtime environment or a secret manager.
- Do not commit production credentials.
- The local `DATABASE_URL` example is for development only.

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
go run ./cmd/goflow migrate
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
- `GET /health/ready`
- `GET /metrics`
- `GET /v1/jobs`
- `GET /v1/jobs?status=pending`
- `GET /v1/jobs/{id}`
- `POST /v1/jobs`
- `DELETE /v1/jobs/{id}`

The app automatically applies `migrations/001_create_jobs.sql` during startup.


## Docker

Build the image:

```powershell
docker build -t goflow:dev .
```

## Docker Compose

Start the full stack:

```powershell
docker compose config
docker compose up --build
```

Validate from another terminal:

```powershell
curl.exe http://localhost:8080/health/live
curl.exe http://localhost:8080/health/ready
curl.exe http://localhost:8080/metrics
```

Stop containers but keep PostgreSQL data:

```powershell
docker compose down
```

Reset containers and delete PostgreSQL data:

```powershell
docker compose down -v
```

`goflow-migrate` runs schema setup once before the API and worker start. Inside Compose, GoFlow uses `postgres` as the database hostname because `localhost` would refer to the GoFlow container itself. The API address is configured with `HTTP_ADDR`.

## Validation

```powershell
gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/http_test.go ./cmd/goflow/worker.go ./cmd/goflow/worker_test.go
gofmt -w ./internal/goflow/service.go ./internal/goflow/service_test.go ./internal/goflow/postgres_repository.go ./internal/goflow/postgres_repository_test.go
go vet ./...
go test ./...
go test --% -coverprofile=coverage.out ./...
go test -race ./...
go test -bench=. -benchmem ./...
```

Generated profiling artifacts such as `cpu.out`, `mem.out`, `block.out`, `mutex.out`, and `trace.out` are ignored by Git.
