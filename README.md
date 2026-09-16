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
- Module 4.6 CI and engineering workflow completed with GitHub Actions validation and production-ready documentation
- Current next step: final roadmap review and portfolio polish

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

Schema setup is owned by `goflow migrate`. API and worker processes open the database but do not run migrations.


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

## Architecture

```mermaid
flowchart LR
    Client[HTTP client] --> API[goflow-api / net/http]
    API --> Repo[PostgresRepository]
    Worker[goflow-worker pool] --> Repo
    Migrate[goflow-migrate] --> DB[(PostgreSQL)]
    Repo --> DB
```

GoFlow uses one Go binary with multiple commands:

- `migrate`: applies database schema changes.
- `serve`: runs the HTTP API.
- `work`: runs the worker pool.
- `create`, `list`, `get`, `process`: local CLI operations.

## Configuration Reference

| Variable | Required | Default | Purpose |
|---|---:|---|---|
| `DATABASE_URL` | Yes | none | PostgreSQL connection string. |
| `HTTP_ADDR` | No | `:8080` | HTTP listen address for `serve`. |
| `MIGRATION_FILE` | No | `migrations/001_create_jobs.sql` | SQL migration file path. |

`DATABASE_URL` should come from the runtime environment or a secret manager. Do not commit real credentials.

## Database Schema

GoFlow stores jobs in PostgreSQL:

| Column | Purpose |
|---|---|
| `id` | Stable job identifier, primary key. |
| `job_type` | Job domain type, such as `email` or `report`. |
| `payload` | Optional binary payload. |
| `status` | Current state: `pending`, `running`, `completed`, or `dead_letter`. |
| `attempts` | Number of execution attempts made. |
| `max_attempts` | Retry budget. |
| `available_at` | When a pending job becomes eligible for workers. |
| `last_error` | Most recent execution failure message. |
| `created_at` / `updated_at` | Audit timestamps. |

## Retry Semantics

Workers claim ready jobs and move them from `pending` to `running`. Successful execution moves a job to `completed`. Temporary failures increment `attempts`; if retries remain, the job returns to `pending` with a future `available_at`. Permanent failures, or exhausted retry budget, move the job to `dead_letter`.

## Consistency Guarantees

GoFlow reduces duplicate work by claiming jobs before execution and by polling only ready `pending` jobs. The current implementation is suitable for learning and local operation, but true distributed exactly-once processing is not guaranteed. Production systems still need idempotent job handlers because crashes can happen after external side effects.

## Testing Guide

Run normal validation:

```powershell
go vet ./...
go test ./...
```

Run race validation when the local Windows toolchain supports it:

```powershell
go test -race ./...
```

Run benchmark validation:

```powershell
go test -bench=. -benchmem ./...
```

Run vulnerability scanning:

```powershell
govulncheck ./...
```

Run PostgreSQL integration tests by setting `DATABASE_URL` before `go test ./...`.

## CI

GitHub Actions runs:

- `go mod tidy` plus `git diff --exit-code`
- `go fmt ./...` plus `git diff --exit-code`
- `go vet ./...`
- `go test ./...` with PostgreSQL service container
- `go test -race ./...`
- `govulncheck ./...`
- `go build` with build metadata
- `docker build`
- `docker compose config`

## Engineering Workflow

Use small branches or pull requests for changes. Before opening or merging a PR, run the same local validation that CI runs where practical:

```powershell
go mod tidy
go fmt ./...
go vet ./...
go test ./...
govulncheck ./...
go build -ldflags "-X main.version=dev -X main.commit=local" -o goflow.exe ./cmd/goflow
docker build -t goflow:ci .
docker compose config
```

Use semantic version tags such as `v0.1.0`, `v0.2.0`, and `v1.0.0` for releases. Update `CHANGELOG.md` before tagging a release.

## Operational Runbook

Start locally with Compose:

```powershell
docker compose up --build
```

Check health:

```powershell
curl.exe http://localhost:8080/health/live
curl.exe http://localhost:8080/health/ready
```

Inspect logs:

```powershell
docker compose logs -f goflow-api
docker compose logs -f goflow-worker
```

Restart without deleting jobs:

```powershell
docker compose down
docker compose up --build
```

Reset the database:

```powershell
docker compose down -v
docker compose up --build
```

## Known Limitations

- No authentication or authorization yet; keep the API local/internal.
- In-process metrics are not shared between separate API and worker processes.
- Worker execution examples are simulated, not real email/report integrations.
- No distributed rate limiter.
- No least-privilege production database user in Compose yet.
- Distroless runtime has no shell or `curl`; health validation is external.

## Future Improvements

- Add authentication and authorization.
- Add Prometheus/OpenTelemetry metrics and traces.
- Add a separate dead-letter table or operational replay tooling.
- Add a least-privilege PostgreSQL user and separate migration role.
- Add release tagging and changelog automation.
- Add idempotency keys for external side effects.
