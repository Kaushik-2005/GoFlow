# Go Learning Tracker

## Current Position

- Current week: Week 2
- Current module: Day 12 - Week 2 checkpoint
- Current topic: Week 2 checkpoint and revision
- Current task: Run the Week 2 checkpoint on 2026-09-02 now that Module 2.5 is complete
- Next milestone: Start the Week 2 checkpoint before moving into Day 13 - Module 3.1
- Active blockers: None for current testing work; live PostgreSQL validation completed on 2026-09-02

## Roadmap Progress

| Day | Module | Main Topic | Status | Deliverable | Verification | Confidence |
|---|---|---|---|---|---|---|
| 1 | Module 1.1 | Environment, packages, modules, and tooling | Completed | Runnable Go module with CLI entry point and command dispatch scaffold | `go mod init goflow`; `go run ./cmd/goflow`; `go build ./cmd/goflow`; `go test ./...`; `go vet ./...`; learner explained package vs module and why `package main` + `func main()` is executable | 4 |
| 2 | Module 1.2 | Language fundamentals | Completed | Core helper functions for max priority, status counting, validation, retry delay, and filtering | `gofmt -w ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; learner answered theory checks on zero values, `:=` vs `=`, `if`, `range`, `defer`, shadowing, and `iota` | 4 |
| 3 | Module 1.3 | Arrays, slices, maps, strings, and runes | Completed | In-memory job model and early storage helpers using slices/maps | `gofmt -w ./cmd/goflow/store.go`; `go vet ./...`; `go test ./...`; learner explained arrays vs slices, `append`, shared backing arrays, `copy`, nil vs empty slices, map lookup/delete, `make` vs `new`, and value semantics; implemented in-memory `Store` with create/get/list/update/delete | 4 |
| 4 | Module 1.4 | Structs, methods, pointers, and interfaces | Completed | `Job` methods plus a small orchestration function using small consumer-defined interfaces | `gofmt -w ./cmd/goflow/job.go`; `go vet ./...`; `go test ./...`; learner explained structs, methods, value vs pointer receivers, pointers, implicit interface implementation, small interfaces, constructor functions, and implemented `CanRetry`, `MarkRunning`, and `startJob` | 4 |
| 5 | Module 1.5 | Error handling | Completed | Sentinel errors, wrapped errors, and one custom typed error applied to store/service code | `gofmt -w ./cmd/goflow/store.go ./cmd/goflow/service.go ./cmd/goflow/job.go ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow process job-123`; learner explained `error`, `errors.New`, sentinel errors, `%w`, `errors.Is`, custom error types, and `errors.As` | 4 |
| 6 | Module 1.6 | I/O, JSON, files, and configuration | Completed | Week 1 CLI deliverable with JSON-file persistence | `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/store.go ./cmd/goflow/persistence.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow create email`; `go run ./cmd/goflow list`; `go run ./cmd/goflow get job-1`; `go run ./cmd/goflow process job-1`; learner explained `io.Reader` vs `*os.File`, JSON tags, `Marshal` vs `Unmarshal`, missing-file behavior, and map iteration order | 4 |
| 7 | Module 2.1 | HTTP servers | Completed | First `net/http` API endpoints | `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/store.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow serve`; `curl http://localhost:8080/health/live`; `curl http://localhost:8080/v1/jobs`; `curl http://localhost:8080/v1/jobs/job-1`; `Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/jobs -ContentType application/json -Body '{"type":"email"}'`; learner explained handlers, mux routing, method dispatch, status codes, JSON request/response handling, path extraction, and API error shape | 4 |
| 8 | Module 2.2 | Middleware and API reliability | Completed | Request middleware stack and reliability guards | `gofmt -w ./cmd/goflow/http.go ./cmd/goflow/middleware.go ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; `curl -i http://localhost:8080/health/live`; `Invoke-WebRequest -Method POST -Uri http://localhost:8080/v1/jobs -ContentType text/plain -Body '{"type":"email"}'`; oversized POST returned `REQUEST_BODY_TOO_LARGE`; learner explained middleware order, recovery, request IDs, body limits, and content-type enforcement | 4 |
| 9 | Module 2.3 | Project organization and architecture | Completed | Clear layered structure for API and worker code | `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/middleware.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow serve`; `curl http://localhost:8080/health/live`; learner explained why reusable core logic must leave `package main`, why handlers should not import `main`, and why `main` should stay thin | 4 |
| 10 | Module 2.4 | PostgreSQL and `database/sql` | Completed | PostgreSQL-backed CLI and HTTP runtime with automatic schema setup | `go test -run TestPostgresRepositoryIntegration ./internal/goflow -v`; `curl.exe http://localhost:8080/health/live`; `curl.exe http://localhost:8080/v1/jobs`; `Invoke-RestMethod -Method POST -Uri "http://localhost:8080/v1/jobs" -ContentType "application/json" -Body '{"type":"email"}'`; `curl.exe http://localhost:8080/v1/jobs/job-70f1178cdbe12194`; learner validated live PostgreSQL integration, list, create, and get over HTTP | 4 |
| 11 | Module 2.5 | Testing fundamentals | Completed | Tested PostgreSQL-backed REST API with create, retrieve, list, filter, delete, validation, consistent errors, and database-failure coverage | `gofmt -w ./cmd/goflow/http.go ./cmd/goflow/http_test.go ./cmd/goflow/main.go`; `gofmt -w ./internal/goflow/service_test.go ./internal/goflow/postgres_repository_test.go ./internal/goflow/postgres_repository_integration_test.go`; `go vet ./...`; `go test ./...`; `go test -cover ./...`; `go test --% -coverprofile=coverage.out ./...`; learner used handwritten fakes, `httptest`, `sqlmock`, `t.Cleanup`, repository integration testing, and added filter/delete API coverage | 4 |
| 12 | Week 2 checkpoint | Revision and spillover | Not Started | Checkpoint review and catch-up buffer | — | — |
| 13 | Module 3.1 | Goroutines and channels | Not Started | First concurrent processing pipeline | — | — |
| 14 | Module 3.2 | Synchronization | Not Started | Shared-state protection and coordination | — | — |
| 15 | Module 3.3 | Worker pool | Not Started | Configurable worker pool | — | — |
| 16 | Module 3.4 | Context and graceful shutdown | Not Started | Controlled shutdown and cancellation flow | — | — |
| 17 | Module 3.5 | Retries and failure handling | Not Started | Retry policy and dead-letter path | — | — |
| 18 | Module 3.6 | Concurrency testing | Not Started | Race-tested worker behavior | — | — |
| 19 | Module 4.1 | Structured logging | Not Started | Structured logs in API and worker paths | — | — |
| 20 | Module 4.2 | Observability | Not Started | Metrics and health signals | — | — |
| 21 | Module 4.3 | Profiling and performance | Not Started | Measured performance baseline | — | — |
| 22 | Module 4.4 | Security | Not Started | Security checklist and safer boundaries | — | — |
| 23 | Module 4.5 | Containers and configuration | Not Started | Containerized builds and startup config validation | — | — |
| 24 | Module 4.6 | CI and engineering workflow | Not Started | CI pipeline and production-ready workflow docs | — | — |

## Session Log

### 2026-08-18

- Topics studied: package vs module, `package main`, `func main`, Go toolchain basics, CLI argument handling
- Work implemented: initialized learning files, created `go.mod`, built `cmd/goflow/main.go`, added CLI command dispatch and helper functions `requireArg` and `printHelp`
- Tests executed: `go run ./cmd/goflow`; `go build ./cmd/goflow`; `go test ./...`; `go vet ./...`; `gofmt -w ./cmd/goflow/main.go`
- Results: Module 1.1 completed with working CLI scaffold and validated toolchain usage
- Problems encountered: shell sandbox helper failures required unsandboxed reads; repo was not initially a Git repository; `go` was briefly unavailable in PowerShell before becoming usable
- Decisions made: keep the first GoFlow binary small under `cmd/goflow` before introducing larger architecture
- Topics to revisit: none for Module 1.1 beyond normal review
- Next action: start Day 2, Module 1.2 with variables, zero values, and short declaration

### 2026-08-19

- Topics studied: variables, constants, zero values, `:=` vs `=`, explicit conversion, `if`, `switch`, `for`, `range`, functions, multiple returns, `defer`, shadowing, and `iota`
- Work implemented: added `maxPriority`, `countByStatus`, `isValidJobName`, `retryDelay`, and `filterCompleted` to `cmd/goflow/main.go`
- Tests executed: `gofmt -w ./cmd/goflow/main.go`; `go run ./cmd/goflow`; `go vet ./...`; `go test ./...`
- Results: Module 1.2 completed; learner handled the practice exercises and theory checks well enough to advance
- Problems encountered: patch-based edits continued failing because of the local sandbox helper; one `go run` timed out once and then succeeded on rerun
- Decisions made: keep the practice functions in `main.go` for now and defer cleaner package separation until later modules
- Topics to revisit: explain `defer` and short variable declaration with slightly sharper wording during revision
- Next action: start Day 3, Module 1.3 with arrays vs slices and `append`

### 2026-08-20

- Topics studied: arrays vs slices, slice length/capacity, `append`, shared backing arrays, `copy`, nil vs empty slices, maps, key-existence checks, map deletion, strings/UTF-8/runes, `make` vs `new`, and value semantics
- Work implemented: created `cmd/goflow/store.go` with `JobStatus`, `Job`, `Store`, `NewStore`, and create/get/list/update/delete methods; reshaped tracker and learning notes into the clearer MCP-style format
- Tests executed: `gofmt -w ./cmd/goflow/store.go`; `go vet ./...`; `go test ./...`
- Results: Module 1.3 completed; learner explained the core slice/map rules and built the first in-memory job store increment successfully
- Problems encountered: initial `store.go` used a different package name from `main.go`; fixed after review
- Decisions made: keep the store in `cmd/goflow/store.go` for now and defer package separation until later architectural modules
- Topics to revisit: stable ordering for `List()` if later CLI output or tests require deterministic ordering
- Next action: start Day 4, Module 1.4 with structs, methods, pointers, and interfaces on 2026-08-21

### 2026-08-21

- Topics studied: structs, struct literals, methods, value vs pointer receivers, pointers, interfaces, implicit interface implementation, small consumer-defined interfaces, and constructor functions
- Work implemented: added `cmd/goflow/job.go` with `CanRetry()` and `MarkRunning()` methods; added `cmd/goflow/service.go` with `startJob(store JobReaderWriter, id string) bool`
- Tests executed: `gofmt -w ./cmd/goflow/job.go`; `go vet ./...`; `go test ./...`
- Results: Module 1.4 completed; learner connected struct methods, pointer receivers, and small interfaces to the existing in-memory store
- Problems encountered: the first `job.go` duplicated `Job`/status declarations and contained a status-name typo; both were corrected during review
- Decisions made: keep orchestration in `service.go`, job behavior in `job.go`, and storage behavior in `store.go` until later package organization modules
- Topics to revisit: later compare when a value receiver is acceptable for larger structs and when pointer receivers should be used consistently
- Next action: start Day 5, Module 1.5 with explicit error handling and wrapped store failures

### 2026-08-22

- Topics studied: the `error` interface, returned errors vs exceptions, `errors.New`, sentinel errors, wrapping with `%w`, `errors.Is`, custom error types, and `errors.As`
- Work implemented: added `ErrJobNotFound` and `ErrJobAlreadyExists`; updated `Store` methods to return errors; updated `startJob` to wrap errors and reject non-pending jobs with `InvalidJobStatusError`; updated the `process` CLI path to use `errors.Is(err, ErrJobNotFound)` at the boundary
- Tests executed: `gofmt -w ./cmd/goflow/store.go ./cmd/goflow/service.go`; `gofmt -w ./cmd/goflow/job.go`; `gofmt -w ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow process job-123`
- Results: Module 1.5 completed; learner can distinguish sentinel errors from custom typed errors and knows when to use `errors.Is` vs `errors.As`
- Problems encountered: patch-based edits still required direct file rewrites outside the sandbox; `JobReaderWriter` had been referenced earlier without a definition and was added as part of the Day 5 cleanup
- Decisions made: keep `ErrJobNotFound` and `ErrJobAlreadyExists` as sentinel errors; use a custom error type only for richer state-transition failures like invalid job status
- Topics to revisit: later show `errors.As` at a boundary with a real branch on `InvalidJobStatusError` once the CLI/API path grows
- Next action: start Day 6, Module 1.6 with file I/O, JSON encoding, and persistence on 2026-08-22

### 2026-08-23

- Topics studied: `io.Reader` and `io.Writer`, file helpers from `os`, `encoding/json`, struct tags, `Marshal` vs `Unmarshal`, missing-file handling, JSON-file persistence, and a real `errors.As` boundary case
- Work implemented: added `cmd/goflow/persistence.go` with `saveJobs` and `loadJobs`; added `Store.Save(...)` and `LoadStore(...)`; added JSON tags to `Job`; wired `create`, `list`, `get`, and `process` in `cmd/goflow/main.go` to persisted `jobs.json` state
- Tests executed: `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/store.go ./cmd/goflow/persistence.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow create email`; `go run ./cmd/goflow list`; `go run ./cmd/goflow get job-1`; `go run ./cmd/goflow process job-1`
- Results: Module 1.6 completed; the Week 1 CLI deliverable now persists jobs to `jobs.json` across runs and reports invalid status transitions with `errors.As`
- Problems encountered: map-backed listing order remained non-deterministic by design; struct tags first defaulted to Go field names until explicit JSON tags were added; the Windows sandbox helper still required direct unsandboxed file rewrites for edits
- Decisions made: treat a missing `jobs.json` file as an empty job list rather than an error; keep persistence in a dedicated helper file while the project remains in `package main`
- Topics to revisit: later decide whether CLI list output should be sorted for deterministic UX and testing; later replace the temporary ID generation strategy with a safer identifier approach
- Next action: start Day 7, Module 2.1 with `net/http` handlers on 2026-08-24

### 2026-08-24

- Topics studied: `http.Server`, `http.Handler`, `http.HandlerFunc`, `ServeMux`, path-based routing, method dispatch, status codes, JSON request and response bodies, path extraction, and consistent JSON API errors
- Work implemented: added `cmd/goflow/http.go`; implemented `liveHandler`, `jobsHandler`, `listJobsHandler`, `getJobHandler`, and `createJobHandler`; added JSON API error helpers; wired `serve` mode in `cmd/goflow/main.go`; exposed `/health/live`, `GET /v1/jobs`, `GET /v1/jobs/{id}`, and `POST /v1/jobs`
- Tests executed: `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/store.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow serve`; `curl http://localhost:8080/health/live`; `curl -X POST http://localhost:8080/health/live`; `curl http://localhost:8080/v1/jobs`; `curl http://localhost:8080/v1/jobs/job-2`; `Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/jobs -ContentType application/json -Body ''{"type":"email"}''`; `Invoke-WebRequest -Method POST -Uri http://localhost:8080/v1/jobs -ContentType application/json -Body ''{}''`
- Results: Module 2.1 completed; GoFlow now serves its first HTTP API endpoints with `net/http` over the existing JSON-backed store
- Problems encountered: an HTTP route was initially missing from the mux; one dispatch bug sent `GET /v1/jobs` to the single-job handler; PowerShell quoting caused one malformed `curl` request; changing `maxAttempts` to `max_attempts` required recreating `jobs.json` because old saved data decoded the renamed field as zero
- Decisions made: keep the HTTP layer in `cmd/goflow/http.go` and `package main` for now; use one `/v1/jobs` handler that dispatches by method; prefer consistent JSON error envelopes instead of plain-text API errors
- Topics to revisit: request body size limits, content-type validation, and path parameter handling in plain `net/http` before introducing cleaner transport structure later
- Next action: start Day 8, Module 2.2 with middleware and API reliability on 2026-08-25

### 2026-08-26

- Topics studied: middleware shape and chaining, recovery middleware, request IDs, request-scoped context values, request body size limits, content-type validation, and more precise API error mapping
- Work implemented: added `cmd/goflow/middleware.go`; implemented `recoveryMiddleware`, `requestIDMiddleware`, `chain(...)`, `requestBodyLimitMiddleware(...)`, and `requireJSONMiddleware(...)`; wired middleware into the HTTP server; mapped oversized request bodies to `REQUEST_BODY_TOO_LARGE`; verified request IDs in both headers and request context
- Tests executed: `gofmt -w ./cmd/goflow/http.go ./cmd/goflow/middleware.go ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; `curl -i http://localhost:8080/health/live`; oversized `POST /v1/jobs`; `Invoke-WebRequest -Method POST -Uri http://localhost:8080/v1/jobs -ContentType text/plain -Body ''{"type":"email"}''`
- Results: Module 2.2 completed; GoFlow now has its first middleware stack and reliability guards around the current HTTP API
- Problems encountered: a temporary panic route was needed to prove recovery and then removed; exact content-type matching is currently strict and will reject variants like `application/json; charset=utf-8`; one decode path initially collapsed oversized bodies into the generic invalid-body error until `*http.MaxBytesError` handling was added
- Decisions made: keep the middleware stack simple and local to `cmd/goflow` for now; use request IDs in both response headers and request context; enforce JSON content type and body size at the HTTP boundary instead of duplicating checks inside handlers
- Topics to revisit: relax content-type parsing to allow valid JSON charset variants; later add structured logging that reads the request ID from context; revisit whether middleware should be selectively applied per route as the API grows
- Next action: start Day 9, Module 2.3 with project organization and architecture on 2026-08-27

### 2026-08-29

- Topics studied: package-boundary design, dependency direction, handler vs service vs repository responsibilities, and keeping `main` as a thin wiring layer
- Work implemented: created `internal/goflow`; moved `store.go`, `job.go`, `service.go`, and `persistence.go` into the reusable core package; exported `StartJob(...)`; updated `main.go` and `http.go` to import `goflow/internal/goflow`; removed old Day 2 practice helpers from `main.go`
- Tests executed: `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/middleware.go`; `go vet ./...`; `go test ./...`; `go run ./cmd/goflow serve`; `curl http://localhost:8080/health/live`
- Results: Module 2.3 completed; GoFlow now has a reusable core package and a clearer dependency direction from outer transport/CLI layers into shared application logic
- Problems encountered: wildcard formatting from PowerShell needed explicit file expansion; moving core logic out of `package main` required exporting `StartJob(...)`; `http.go` briefly had the wrong package declaration before being restored to `package main`
- Decisions made: extract a reusable core package before extracting a dedicated HTTP package; keep `cmd/goflow` as the current outer app layer; treat leftover practice helpers in `main.go` as cleanup targets before moving real startup code
- Topics to revisit: later extract the HTTP transport into its own internal package once the shared core boundary is stable; split CLI command logic further once additional binaries or commands justify it
- Next action: start Day 10, Module 2.4 with PostgreSQL and `database/sql` on 2026-08-30

### 2026-09-01

- Topics studied: `sql.DB`, `PingContext`, `QueryRowContext`, `QueryContext`, `ExecContext`, `Scan`, row cleanup, and parameterized queries
- Work implemented: added `internal/goflow/postgres_repository.go` with a context-aware `JobRepository` interface and PostgreSQL-backed `Create`, `Get`, `List`, `Update`, and `Delete` methods; added `migrations/001_create_jobs.sql`
- Tests executed: `gofmt -w ./internal/goflow/postgres_repository.go`; `go vet ./...`; `go test ./...`
- Results: Day 10 is in progress with the first database repository increment implemented and validated at compile time
- Problems encountered: no PostgreSQL driver or live database is wired yet, so runtime DB connectivity and queries are not validated in this step
- Decisions made: keep the current file-backed app flow temporarily while introducing the DB repository boundary behind `database/sql`
- Topics to revisit: mapping unique-constraint errors to `ErrJobAlreadyExists`; switching the app startup path from `LoadStore(...)` to a real database handle
- Next action: wire a real `*sql.DB` into startup, verify connectivity with `PingContext`, and route operations through the repository






