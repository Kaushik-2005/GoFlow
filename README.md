# GoFlow

GoFlow is a Go backend service for asynchronous job processing. It exposes an HTTP API for submitting and inspecting jobs, stores job state in PostgreSQL, and runs a worker pool that claims ready jobs, executes them, retries temporary failures, and moves exhausted or permanent failures to a dead-letter state.

The project is intentionally built with Go's standard library for the core service surface: `net/http`, `database/sql`, `context`, `sync`, and `log/slog`.

## What GoFlow Provides

- REST API for creating, listing, reading, and deleting jobs.
- PostgreSQL-backed persistence.
- Dedicated migration command.
- Concurrent worker pool with bounded in-memory queueing.
- Atomic database claim for `pending -> running` transitions.
- Retry scheduling through `available_at`.
- Dead-letter status for permanent or exhausted failures.
- Context-aware shutdown for HTTP and worker runtimes.
- JSON structured logs with request and job correlation fields.
- Liveness, readiness, and JSON metrics endpoints.
- Dockerfile and Docker Compose stack for local production-style execution.
- Unit, HTTP handler, repository, integration, benchmark, and CI validation.

## Architecture

```mermaid
flowchart TD
    Client[HTTP Client] --> API[GoFlow API / net/http]
    API --> Middleware[Middleware]
    Middleware --> Handler[HTTP Handlers]
    Handler --> Store[JobStore interface]
    Worker[Worker Pool] --> Store
    CLI[CLI Commands] --> Store
    Migrate[Migrate Command] --> DB[(PostgreSQL)]
    Store --> Repo[PostgresRepository]
    Repo --> DB
```

GoFlow ships as one binary with multiple commands:

| Command | Purpose |
|---|---|
| `migrate` | Applies the PostgreSQL schema migration. |
| `serve` | Starts the HTTP API. |
| `work` | Starts the worker pool. |
| `create <type>` | Creates a job from the CLI. |
| `list` | Lists jobs from PostgreSQL. |
| `get <job-id>` | Reads one job from PostgreSQL. |
| `process <job-id>` | Claims one job manually from the CLI. |

## Runtime Model

1. A client submits a job with `POST /v1/jobs`.
2. The API validates the request and stores the job as `pending`.
3. The worker command polls PostgreSQL for jobs where `status = pending` and `available_at <= NOW()`.
4. A worker claims a job with a conditional SQL update, so only one worker can move it from `pending` to `running`.
5. The worker executes the job.
6. Success moves the job to `completed`.
7. Temporary failure increments `attempts`, records `last_error`, and schedules a future retry through `available_at`.
8. Permanent failure or exhausted attempts moves the job to `dead_letter`.

```mermaid
stateDiagram-v2
    [*] --> pending
    pending --> running: atomic claim
    running --> completed: execution succeeds
    running --> pending: temporary failure and retries remain
    running --> dead_letter: permanent failure
    running --> dead_letter: retry budget exhausted
```

## Requirements

- Go 1.26+
- PostgreSQL 16+ for local development and integration tests
- Docker and Docker Compose for the full local stack
- `govulncheck` for vulnerability scanning when running the full validation suite locally

## Configuration

GoFlow reads configuration from environment variables at startup.

| Variable | Required | Default | Used by | Description |
|---|---:|---|---|---|
| `DATABASE_URL` | Yes | none | all commands | PostgreSQL connection string. |
| `HTTP_ADDR` | No | `:8080` | `serve` | HTTP listen address. |
| `MIGRATION_FILE` | No | `migrations/001_create_jobs.sql` | `migrate` | SQL migration file path. |

PowerShell example:

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/goflow?sslmode=disable"
$env:HTTP_ADDR = ":8080"
$env:MIGRATION_FILE = "migrations/001_create_jobs.sql"
```

`.env.example` is provided as a local template. Do not commit real `.env` files or production credentials.

## Quick Start With Docker Compose

Start the full stack:

```powershell
docker compose up --build
```

Compose starts:

| Service | Role |
|---|---|
| `postgres` | PostgreSQL database. |
| `goflow-migrate` | Runs `goflow migrate` once after PostgreSQL is healthy. |
| `goflow-api` | Runs `goflow serve` on port `8080`. |
| `goflow-worker` | Runs `goflow work`. |

Validate the API from another terminal:

```powershell
curl.exe http://localhost:8080/health/live
curl.exe http://localhost:8080/health/ready
curl.exe http://localhost:8080/metrics
```

Create a job:

```powershell
Invoke-RestMethod -Method POST `
  -Uri "http://localhost:8080/v1/jobs" `
  -ContentType "application/json" `
  -Body '{"type":"email"}'
```

List jobs:

```powershell
curl.exe http://localhost:8080/v1/jobs
```

Stop the stack but keep PostgreSQL data:

```powershell
docker compose down
```

Reset the database volume:

```powershell
docker compose down -v
```

Inside Compose, the database hostname is `postgres`, not `localhost`. In a container, `localhost` means the current container.

## Local Development Without Compose

Start PostgreSQL locally or with Docker:

```powershell
docker run --name goflow-postgres `
  -e POSTGRES_USER=postgres `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=goflow `
  -p 5432:5432 `
  -d postgres:16
```

Set the database URL:

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/goflow?sslmode=disable"
```

Apply migrations:

```powershell
go run ./cmd/goflow migrate
```

Run the API:

```powershell
go run ./cmd/goflow serve
```

Run the worker in another terminal:

```powershell
go run ./cmd/goflow work
```

Use the CLI:

```powershell
go run ./cmd/goflow create email
go run ./cmd/goflow list
go run ./cmd/goflow get <job-id>
go run ./cmd/goflow process <job-id>
```

## HTTP API

Base URL for local Compose:

```text
http://localhost:8080
```

### Health

#### `GET /health/live`

Returns whether the GoFlow process can respond.

Response:

```json
{"status":"ok"}
```

This endpoint does not depend on PostgreSQL.

#### `GET /health/ready`

Returns whether GoFlow is ready to serve real workload.

Response when PostgreSQL is reachable:

```json
{"status":"ready"}
```

Response when PostgreSQL is unavailable:

```json
{
  "error": {
    "code": "DATABASE_NOT_READY",
    "message": "The database is not ready"
  }
}
```

### Metrics

#### `GET /metrics`

Returns an in-process JSON metrics snapshot.

Example:

```json
{
  "http_requests_total": 2,
  "http_request_duration_ms_le_10": 2,
  "http_request_duration_ms_le_100": 0,
  "http_request_duration_ms_le_1000": 0,
  "http_request_duration_ms_gt_1000": 0,
  "jobs_submitted_total": 1,
  "jobs_completed_total": 0,
  "jobs_execution_failed_total": 0,
  "jobs_retried_total": 0,
  "jobs_dead_lettered_total": 0,
  "active_workers": 0,
  "queue_depth": 0
}
```

Metrics are process-local. The API process does not expose worker-process counters when API and worker run as separate processes.

### Jobs

#### `POST /v1/jobs`

Creates a job.

Request:

```json
{"type":"email"}
```

Supported public API job types:

| Type | Behavior |
|---|---|
| `email` | Simulated successful job execution. |
| `report` | Simulated successful job execution. |

Response: `201 Created`

```json
{
  "id": "job-78cffba2dede5cef",
  "type": "email",
  "payload": null,
  "status": "pending",
  "attempts": 0,
  "max_attempts": 3,
  "available_at": "2026-09-15T17:27:24.364409Z"
}
```

#### `GET /v1/jobs`

Lists all jobs.

```powershell
curl.exe http://localhost:8080/v1/jobs
```

Optional status filter:

```powershell
curl.exe "http://localhost:8080/v1/jobs?status=pending"
```

Valid status filters:

- `pending`
- `running`
- `completed`
- `dead_letter`

#### `GET /v1/jobs/{id}`

Reads one job.

```powershell
curl.exe http://localhost:8080/v1/jobs/job-78cffba2dede5cef
```

Invalid job ID format returns `400 INVALID_JOB_ID`. A valid but missing job ID returns `404 JOB_NOT_FOUND`.

#### `DELETE /v1/jobs/{id}`

Deletes one job.

```powershell
curl.exe -X DELETE http://localhost:8080/v1/jobs/job-78cffba2dede5cef
```

Successful deletion returns `204 No Content`.

## Error Format

HTTP errors use a consistent JSON envelope:

```json
{
  "error": {
    "code": "JOB_NOT_FOUND",
    "message": "The requested job does not exist"
  }
}
```

The HTTP status code communicates the protocol-level result. The JSON `code` gives application-specific detail for clients and debugging.

Common error codes:

| Code | Meaning |
|---|---|
| `INVALID_REQUEST_BODY` | Request JSON is malformed. |
| `REQUEST_BODY_TOO_LARGE` | Request body exceeds the configured 1 MB limit. |
| `JOB_TYPE_REQUIRED` | Create request did not include a job type. |
| `INVALID_JOB_TYPE` | Job type is not supported by the public API. |
| `INVALID_JOB_ID` | Job ID format is invalid. |
| `JOB_NOT_FOUND` | Requested job does not exist. |
| `INVALID_STATUS_FILTER` | Status query parameter is not one of the supported statuses. |
| `METHOD_NOT_ALLOWED` | Endpoint does not support the HTTP method. |
| `DATABASE_NOT_READY` | Readiness check could not reach PostgreSQL. |

## Worker Behavior

The worker command starts a fixed-size worker pool.

Current worker settings:

| Setting | Value | Location |
|---|---:|---|
| Worker count | `3` | `cmd/goflow/main.go` |
| In-memory job queue size | `8` | `cmd/goflow/main.go` |
| Poll interval | `2s` | `cmd/goflow/main.go` |
| Request body limit | `1 MB` | `cmd/goflow/main.go` middleware chain |

The worker polls ready jobs, queues job IDs in a bounded channel, and workers claim jobs before execution. A local mutex-protected `queued` map reduces duplicate enqueueing inside one worker process. PostgreSQL remains the ownership boundary through `ClaimPending`.

Internal simulation job types used by the worker executor:

| Type | Behavior |
|---|---|
| `email` | Completes successfully. |
| `report` | Completes successfully. |
| `temporary-fail` | Produces a retryable failure. |
| `permanent-fail` | Produces a permanent failure. |

The public HTTP API accepts only `email` and `report`. Failure simulation types are useful through internal tests or local CLI experimentation, not as public API inputs.

## Retry And Dead-Letter Semantics

Temporary failure flow:

1. Job is `running`.
2. Worker receives a temporary execution error.
3. `attempts` is incremented.
4. `last_error` is updated.
5. If attempts remain, status returns to `pending` and `available_at` moves into the future.
6. If attempts are exhausted, status becomes `dead_letter`.

Permanent failure flow:

1. Job is `running`.
2. Worker receives a permanent execution error.
3. `attempts` is incremented.
4. `last_error` is updated.
5. Status becomes `dead_letter`.

Retry delay is currently exponential based on attempt count.

## Consistency Guarantees

GoFlow provides atomic job ownership for the claim step:

```sql
UPDATE jobs
SET status = 'running', updated_at = NOW()
WHERE id = $1
AND status = 'pending';
```

Only the worker whose conditional update affects one row owns the job.

GoFlow does not claim exactly-once external execution. A worker can perform an external side effect and crash before recording completion. Real job handlers must be idempotent or use external idempotency keys to prevent duplicate side effects.

## Database Schema

Migration file: `migrations/001_create_jobs.sql`

| Column | Type | Description |
|---|---|---|
| `id` | `TEXT PRIMARY KEY` | Stable job ID such as `job-78cffba2dede5cef`. |
| `job_type` | `TEXT NOT NULL` | Job type. |
| `payload` | `BYTEA` | Optional binary payload. |
| `status` | `TEXT NOT NULL` | `pending`, `running`, `completed`, or `dead_letter`. |
| `attempts` | `INTEGER NOT NULL DEFAULT 0` | Number of execution attempts made. |
| `max_attempts` | `INTEGER NOT NULL DEFAULT 3` | Maximum attempts before dead-lettering. |
| `available_at` | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` | Earliest time a pending job can be picked up. |
| `last_error` | `TEXT NOT NULL DEFAULT ''` | Latest execution failure message. |
| `created_at` | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` | Creation timestamp. |
| `updated_at` | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` | Last update timestamp. |

## Project Structure

```text
.
├── .github/workflows/ci.yml
├── cmd/goflow/
├── docs/
├── internal/goflow/
├── migrations/
├── .dockerignore
├── .env.example
├── .gitignore
├── CHANGELOG.md
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### Root Files

| Path | Purpose |
|---|---|
| `README.md` | Main user, operator, and developer documentation. |
| `CHANGELOG.md` | Release notes and notable changes. |
| `go.mod` / `go.sum` | Go dependency manifest and checksums. |
| `Dockerfile` | Multi-stage build for the GoFlow binary and distroless runtime image. |
| `docker-compose.yml` | Local stack with PostgreSQL, migration, API, and worker services. |
| `.dockerignore` | Keeps unnecessary files out of Docker build context. |
| `.env.example` | Safe local environment variable template. |
| `.gitignore` | Ignores binaries, local env files, and generated profiling artifacts. |
| `.github/workflows/ci.yml` | CI workflow for tests, race validation, security scan, build, Docker, and Compose checks. |

### `cmd/goflow`

This is the executable package. It owns process startup, command dispatch, dependency wiring, HTTP runtime setup, worker lifecycle, and local CLI behavior.

| File | Purpose |
|---|---|
| `main.go` | Command dispatch, process lifecycle, HTTP server setup, worker pool loop, graceful shutdown. |
| `config.go` | Environment-based configuration loading and defaults. |
| `database.go` | PostgreSQL connection setup, migration execution, random job ID generation. |
| `http.go` | HTTP handlers, API request validation, JSON responses, API error mapping. |
| `middleware.go` | Request IDs, request logging, panic recovery, request body limits, JSON content-type enforcement. |
| `metrics.go` | Mutex-protected in-process metrics collector and snapshot model. |
| `worker.go` | Single-job processing flow and simulated job executor. |
| `*_test.go` | Handler, config, metrics, worker, and benchmark tests. |

### `internal/goflow`

This is the application core package. It is importable by the executable package but hidden from packages outside this repository by Go's `internal` package rule.

| File | Purpose |
|---|---|
| `job.go` | Job methods and domain errors such as invalid status transitions. |
| `store.go` | Job model, statuses, sentinel errors, and legacy in-memory store implementation. |
| `service.go` | Business operations: claim/start, complete, fail, retry/dead-letter transitions. |
| `postgres_repository.go` | PostgreSQL implementation of job persistence using `database/sql`. |
| `*_test.go` | Service, repository, and live PostgreSQL integration tests. |

### `migrations`

| File | Purpose |
|---|---|
| `001_create_jobs.sql` | Creates and evolves the `jobs` table used by GoFlow. |

### `docs`

The `docs/` directory contains project learning and design records. These are useful for understanding design history, but they are not required to run GoFlow.

| File | Purpose |
|---|---|
| `design.md` | Architecture evolution, diagrams, and design rationale. |
| `decisions.md` | Architectural decisions and trade-offs. |
| `learning.md` | Go notes and explanations collected during development. |
| `tracker.md` | Progress and validation record. |
| `session-log.md` | Session-by-session implementation notes. |

## Build

Build a local binary:

```powershell
go build -o goflow.exe ./cmd/goflow
```

Build with version metadata:

```powershell
go build -ldflags "-X main.version=v0.1.0 -X main.commit=<commit-sha>" -o goflow.exe ./cmd/goflow
```

Show help/version:

```powershell
.\goflow.exe
```

Build Docker image:

```powershell
docker build -t goflow:dev .
```

## Testing And Validation

Normal validation:

```powershell
go fmt ./...
go vet ./...
go test ./...
```

Race validation:

```powershell
go test -race ./...
```

Coverage:

```powershell
go test -cover ./...
go test --% -coverprofile=coverage.out ./...
```

Benchmarks:

```powershell
go test -bench=. -benchmem ./...
```

Vulnerability scan:

```powershell
govulncheck ./...
```

Docker and Compose validation:

```powershell
docker build -t goflow:ci .
docker compose config
```

Generated profiling artifacts such as `cpu.out`, `mem.out`, `block.out`, `mutex.out`, and `trace.out` are ignored by Git.

## CI Pipeline

GitHub Actions runs on push and pull request to `main`.

The CI workflow validates:

- `go mod tidy` with a clean `go.mod` / `go.sum` diff check.
- `go fmt ./...` with a clean working-tree diff check.
- `go vet ./...`.
- `go test ./...` with a PostgreSQL service container.
- `go test -race ./...`.
- `govulncheck ./...`.
- `go build` with build metadata.
- `docker build`.
- `docker compose config`.

## Logging

GoFlow uses JSON structured logs through `log/slog` for `serve` and `work` commands.

HTTP request logs include fields such as:

- `request_id`
- `method`
- `path`
- `status`
- `duration_ms`

Worker logs include fields such as:

- `worker_id`
- `job_id`
- `error`

Request bodies, payloads, credentials, cookies, and `DATABASE_URL` should not be logged.

## Operational Runbook

Start stack:

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

Restart without deleting data:

```powershell
docker compose down
docker compose up --build
```

Reset database and all saved jobs:

```powershell
docker compose down -v
docker compose up --build
```

Run only migrations locally:

```powershell
go run ./cmd/goflow migrate
```

## Security Notes

- GoFlow currently has no authentication or authorization.
- Treat the HTTP API as local/internal unless it is placed behind an authenticated gateway or an auth layer is added.
- Do not expose `DELETE /v1/jobs/{id}` publicly without authorization.
- `DATABASE_URL` must come from the runtime environment or a secret manager.
- Do not commit production credentials or real `.env` files.
- SQL queries use parameters; do not build SQL by concatenating user input.
- The Compose database uses development credentials and the `postgres` superuser. Use least-privilege roles in production.

## Known Limitations

- No authentication or authorization.
- No real email/report integrations; job execution is simulated.
- Metrics are in-process and not shared between API and worker processes.
- No Prometheus/OpenTelemetry exporter yet.
- No separate dead-letter table or replay tooling.
- No distributed rate limiter.
- No least-privilege database role in the local Compose setup.
- Distroless runtime image has no shell or `curl`; health checks are expected to run from outside the container.
- Exactly-once external execution is not guaranteed; job handlers must be idempotent.

## Future Improvements

- Add authentication and authorization middleware.
- Add Prometheus metrics and OpenTelemetry traces.
- Add real job executors for email and report workloads.
- Add idempotency keys for external side effects.
- Add a separate dead-letter table and replay workflow.
- Add a least-privilege database user and separate migration role.
- Add configurable worker count, queue size, retry policy, and poll interval.
- Add jitter to retry delay.
- Add production-ready rate limiting.
- Add deployment examples for managed PostgreSQL and container platforms.
