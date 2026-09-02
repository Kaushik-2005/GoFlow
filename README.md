# GoFlow

GoFlow is a staged learning project for building a production-style Go backend that accepts jobs, stores them, processes them asynchronously, retries transient failures, and surfaces operational behavior cleanly.

The project follows `Go_Industry_Roadmap_4_Weeks.md`. It started as a small CLI for Go fundamentals and now targets PostgreSQL for persistence.

## Current status

Module 2.4 is implemented in code. The CLI and HTTP runtime now use PostgreSQL through `database/sql`, but runtime validation still requires a live PostgreSQL instance and a configured `DATABASE_URL`.

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

The app automatically applies `migrations/001_create_jobs.sql` during startup.

## Validation

```powershell
gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/database.go ./internal/goflow/service.go ./internal/goflow/postgres_repository.go
go mod tidy
go vet ./...
go test ./...
```
