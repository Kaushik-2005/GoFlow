# GoFlow

GoFlow is a staged learning project for building a production-style Go backend that accepts jobs, stores them, processes them asynchronously, retries transient failures, and surfaces operational behavior cleanly.

The project follows `Go_Industry_Roadmap_4_Weeks.md`. Week 2 is now complete: the app uses PostgreSQL through `database/sql` and exposes a tested HTTP API.

## Current status

- Week 1 foundations completed
- Week 2 HTTP, architecture, PostgreSQL, and testing completed
- Current next step: Week 2 checkpoint, then Week 3 concurrency modules

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
gofmt -w ./cmd/goflow/http.go ./cmd/goflow/http_test.go ./cmd/goflow/main.go
gofmt -w ./internal/goflow/service_test.go ./internal/goflow/postgres_repository_test.go ./internal/goflow/postgres_repository_integration_test.go
go vet ./...
go test ./...
go test --% -coverprofile=coverage.out ./...
```
