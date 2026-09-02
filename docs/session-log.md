# Session Log

## 2026-08-18 — Roadmap initialization and Module 1.1 completion

### Topics covered

- Go roadmap initialization
- Module 1.1 orientation
- Packages versus modules
- `package main` and `func main`
- `go run`, `go build`, `go test`, `go vet`, and `gofmt`
- `os.Args`, `switch`, and command-specific argument validation
- small helper refactoring with `requireArg` and `printHelp`

### Work completed

- Read the canonical roadmap
- Inspected the repository state
- Confirmed the repo is a fresh start with no Go project files yet
- Created the required learning record files
- Initialized progress at the first incomplete roadmap item
- Moved learning records into `docs/` on explicit request
- Created the initial Go module with `go.mod`
- Implemented `cmd/goflow/main.go`
- Added CLI dispatch for `list`, `create`, `get`, and `process`
- Added argument validation for `create <type>`, `get <job-id>`, and `process <job-id>`
- Refactored repeated argument checks into `requireArg`
- Refactored usage output into `printHelp`
- Validated the CLI through repeated manual runs and basic toolchain commands

### Commands run

- `Get-ChildItem -Force`
- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-ChildItem -Recurse -File`
- `git status --short`
- `go mod init goflow`
- `go run ./cmd/goflow`
- `go build ./cmd/goflow`
- `go test ./...`
- `go vet ./...`
- `gofmt -w ./cmd/goflow/main.go`

### Problems encountered

- Sandbox helper failed for some read operations; recovered by re-running those reads outside the sandbox
- `git status --short` showed this directory is not yet a Git repository
- `go` was initially not found in PowerShell, then became available and the module scaffold work continued
- Help text was updated before command behavior for `get` and `process` was implemented; fixed by adding the missing `switch` cases

### What I understood well

- The correct resume point was Week 1, Module 1.1
- Package versus module
- Why only `package main` with `func main()` becomes an executable
- Why `go.mod` belongs at the module root
- Why argument length checks are required before indexing `os.Args`
- Why `go vet` complements successful compilation

### What needs revision

- Module 1.2 language fundamentals have not started yet
- Automated tests for the CLI were not added in Module 1.1

### Next session

- Start Module 1.2 with variables, zero values, and short declaration using `:=`

## 2026-08-19 — Module 1.2 completion

### Topics covered

- Variables, constants, and zero values
- `:=` versus `=`
- Explicit type conversion
- `if`, `switch`, `for`, and `range`
- Functions, multiple return values, and `defer`
- Scope, shadowing, and `iota`
- Practice exercises on loops, maps, filtering, and retry calculation

### Work completed

- Resumed from the exact Module 1.2 tracker position
- Taught and checked the main language-fundamentals concepts in sequence
- Completed the practice exercises for `maxPriority`, `countByStatus`, `isValidJobName`, `retryDelay`, and `filterCompleted`
- Added the exercise functions into `cmd/goflow/main.go`
- Ran validation commands and confirmed the package still builds and vets cleanly
- Updated learning records to mark Module 1.2 complete and set Module 1.3 as the next resume point

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `gofmt -w ./cmd/goflow/main.go`
- `go run ./cmd/goflow`
- `go vet ./...`
- `go test ./...`

### Problems encountered

- The local sandbox helper continued to fail for patch-based edits, so file updates were written directly outside the sandbox
- `go run ./cmd/goflow` timed out once unexpectedly during validation, then succeeded on rerun with the expected help output

### What I understood well

- Zero values for `int`, `string`, and `bool`
- The meaning of `:=` versus `=`
- Why Go conditions require real boolean expressions
- Why `_` is used in `range` loops to avoid unused variable errors
- The basic purpose of `defer` and why it is written near resource acquisition
- How map zero values make counting patterns concise

### What needs revision

- Use more precise wording when explaining `defer` and the exact meaning of short variable declaration
- Add automated tests for the language-fundamentals helper functions once the roadmap reaches testing in more depth

### Next session

- Start Module 1.3 with arrays versus slices, slice length/capacity, and `append`
## 2026-08-20 — Module 1.3 completion

### Topics covered

- Arrays vs slices
- Slice length, capacity, `append`, and shared backing arrays
- `copy`, nil vs empty slices, and value semantics
- Maps, key-existence checks, and `delete`
- Strings, UTF-8, bytes, and runes at the conceptual level
- `make` vs `new`
- In-memory job-store design using `map[string]Job`

### Work completed

- Resumed from the exact Day 3 tracker position
- Taught and checked the main slice/map/storage concepts in sequence
- Reviewed and corrected the first `store.go` implementation
- Completed the in-memory `JobStatus`, `Job`, and `Store` increment with create/get/list/update/delete behavior
- Updated tracker and learning notes into the clearer day-by-day format
- Marked Module 1.3 complete and set Module 1.4 as the next resume point

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/store.go'`
- `gofmt -w ./cmd/goflow/store.go`
- `go vet ./...`
- `go test ./...`

### Problems encountered

- The first `store.go` used `package store` while `main.go` used `package main`, which would have broken builds; fixed after review
- The local sandbox helper still prevented patch-based edits, so record updates continued through direct file writes outside the sandbox

### What I understood well

- Arrays have fixed size as part of the type, while slices are flexible views over storage
- `append` must usually be assigned back because it returns the slice to use afterward
- Slices can affect each other when they share backing storage
- `make(map[string]Job)` is the correct way to initialize the in-memory store map
- Mutating a `Job` returned by `Get(id)` does not update the map automatically because of value semantics

### What needs revision

- Explain slice capacity more precisely without describing it as “growing N times”
- Revisit stable ordering concerns if future CLI output or tests require deterministic `List()` results

### Next session

- Start Day 4, Module 1.4 with structs, methods, pointers, and interfaces on 2026-08-21
## 2026-08-21 — Module 1.4 completion

### Topics covered

- Structs and named-field struct literals
- Methods on structs
- Value vs pointer receivers
- Pointers and addressability
- Interfaces and implicit implementation
- Small consumer-defined interfaces
- Constructor functions
- Orchestration with `Job` methods plus store behavior

### Work completed

- Resumed from the exact Day 4 tracker position
- Taught and checked the main struct/method/pointer/interface concepts in sequence
- Reviewed and corrected the first `job.go` implementation
- Added `CanRetry()` and `MarkRunning()` to `Job`
- Added `startJob(store JobReaderWriter, id string) bool` in `service.go`
- Marked Module 1.4 complete and set Module 1.5 as the next resume point

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/store.go'`
- `Get-Content -LiteralPath 'cmd/goflow/job.go'`
- `Get-Content -LiteralPath 'cmd/goflow/service.go'`
- `gofmt -w ./cmd/goflow/job.go`
- `go vet ./...`
- `go test ./...`

### Problems encountered

- The first `job.go` duplicated `Job` and status declarations already present in `store.go`
- `MarkRunning()` initially used the misspelled `StausRunning`, which prevented compilation until corrected
- The local sandbox helper still blocked patch-based edits, so record updates continued through direct file writes outside the sandbox

### What I understood well

- Structs keep related job data together as one coherent value
- `CanRetry()` belongs on `Job` and can use a value receiver because it only reads state
- `MarkRunning()` needs a pointer receiver because it mutates the actual job
- Small interfaces like `JobGetter` and `JobUpdater` are better than large premature store interfaces
- `startJob` still needs `Update` because `Get` returns a value copy

### What needs revision

- Revisit when pointer receivers should be used consistently across a type even for some read-only methods
- Revisit how receiver choice interacts with larger structs and mutation-heavy APIs

### Next session

- Start Day 5, Module 1.5 with explicit error handling and wrapped store failures
## 2026-08-22 — Module 1.5 completion

### Topics covered

- The `error` interface
- Returned errors versus exceptions
- `errors.New` and sentinel errors
- Wrapping with `%w`
- `errors.Is` for known conditions
- Custom error types and `errors.As`
- Applying explicit error handling to store, service, and CLI boundaries

### Work completed

- Resumed from the exact Day 5 tracker position
- Taught and checked the main error-handling concepts in sequence
- Upgraded store methods from boolean failure signals to explicit error returns
- Added `ErrJobNotFound` and `ErrJobAlreadyExists`
- Added `InvalidJobStatusError`
- Updated `startJob(...)` to wrap errors and reject non-pending status transitions
- Updated the `process` CLI path to branch on `errors.Is(err, ErrJobNotFound)`
- Marked Module 1.5 complete and set Module 1.6 as the next resume point

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-ChildItem -Force`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/store.go'`
- `Get-Content -LiteralPath 'cmd/goflow/job.go'`
- `Get-Content -LiteralPath 'cmd/goflow/service.go'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `gofmt -w ./cmd/goflow/store.go ./cmd/goflow/service.go`
- `gofmt -w ./cmd/goflow/job.go`
- `gofmt -w ./cmd/goflow/main.go`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow process job-123`

### Problems encountered

- Patch-based edits still required direct file rewrites outside the sandbox because the local sandbox helper continued to fail
- `JobReaderWriter` had been referenced earlier without a concrete declaration and was added during the Day 5 cleanup
- The workspace still contains `goflow.exe`, which is ignored by `.gitignore` but remains in the working directory

### What I understood well

- `nil` as an `error` result means no error occurred
- Sentinel errors are better than checking message text for conditions like not found and already exists
- `%w` preserves the wrapped error for `errors.Is` and `errors.As`
- `errors.Is` is the right tool for sentinel errors in the current GoFlow design
- Custom error types plus `errors.As` are better when structured details are needed

### What needs revision

- Add a real boundary example for `errors.As` later, not just the type definition and return path
- Revisit whether the `process` command should construct a fresh empty store or use persisted state once Day 6 adds file-based storage

### Next session

- Start Day 6, Module 1.6 with file I/O, JSON encoding, and JSON-file persistence for the CLI on 2026-08-22

## 2026-08-23 — Module 1.6 completion

### Topics covered

- `io.Reader` and `io.Writer`
- `os.ReadFile` and `os.WriteFile`
- `encoding/json`
- JSON struct tags
- `Marshal` versus `Unmarshal`
- Missing-file handling with `errors.Is(err, os.ErrNotExist)`
- JSON-file persistence in the CLI
- Real boundary use of `errors.As` for invalid job status

### Work completed

- Added `cmd/goflow/persistence.go` with `saveJobs(...)` and `loadJobs(...)`
- Added `Store.Save(...)` and `LoadStore(...)`
- Added JSON tags to the `Job` struct
- Wired `create`, `list`, `get`, and `process` to persisted `jobs.json` state
- Added CLI handling for `InvalidJobStatusError` using `errors.As(...)`
- Completed the Week 1 CLI deliverable with persistence across runs

### Commands run

- `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/store.go ./cmd/goflow/persistence.go`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow create email`
- `go run ./cmd/goflow list`
- `go run ./cmd/goflow get job-1`
- `go run ./cmd/goflow process job-1`
- `go run ./cmd/goflow process job-999`
- `go run ./cmd/goflow process job-1`

### Problems encountered

- The Windows sandbox helper still blocked normal patch-based edits, so direct rewrites were needed for code changes
- JSON initially used Go field names until explicit struct tags were added
- `list` initially used `NewStore()` instead of `LoadStore(...)`, which would have ignored persisted state

### What I understood well

- Why `io.Reader` is better than `*os.File` when only read behavior is needed
- Why `json.Unmarshal(...)` needs a pointer destination
- Why a missing `jobs.json` file should mean an empty job list in this stage of the project
- Why map iteration order is not stable in Go
- When `errors.As(...)` is the right tool instead of `errors.Is(...)`

### What needs revision

- Revisit stable ordering if CLI output or tests later need deterministic job order
- Replace the temporary sequential ID generation strategy in a later module
- Revisit whether `Payload` should stay `null`, default to empty bytes, or evolve into a typed payload model later

### Next session

- Start Day 7, Module 2.1 with `net/http` handlers on 2026-08-24

## 2026-08-24 — Module 2.1 completion

### Topics covered

- `http.Server`
- `http.Handler` and `http.HandlerFunc`
- `ServeMux`
- Routing by path and dispatch by method
- JSON request and response bodies
- Path extraction with `strings.TrimPrefix`
- Status codes and method handling
- Consistent JSON API error envelopes

### Work completed

- Added `cmd/goflow/http.go`
- Implemented `liveHandler`, `jobsHandler`, `listJobsHandler`, `getJobHandler`, and `createJobHandler`
- Wired `serve` mode in `cmd/goflow/main.go`
- Exposed `/health/live`, `GET /v1/jobs`, `GET /v1/jobs/{id}`, and `POST /v1/jobs`
- Added `writeJSONError(...)` and a structured API error shape
- Fixed the serialized `max_attempts` field name and recreated `jobs.json` to match the new schema

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/store.go'`
- `Get-Content -LiteralPath 'cmd/goflow/service.go'`
- `Get-Content -LiteralPath 'cmd/goflow/http.go'`
- `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/store.go`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow serve`
- `curl http://localhost:8080/health/live`
- `curl -X POST http://localhost:8080/health/live`
- `curl http://localhost:8080/v1/jobs`
- `curl http://localhost:8080/v1/jobs/job-2`
- `curl http://localhost:8080/v1/jobs/job-9999`
- `Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/jobs -ContentType application/json -Body '{"type":"email"}'`
- `Invoke-WebRequest -Method POST -Uri http://localhost:8080/v1/jobs -ContentType application/json -Body '{}'`

### Problems encountered

- The `/v1/jobs` route was initially missing from the mux
- One dispatch bug pointed `GET /v1/jobs` at the single-job handler instead of the list handler
- PowerShell quoting caused one malformed `curl` POST request
- Renaming a JSON field exposed the zero-value behavior for missing fields in old persisted data

### What I understood well

- The roles of `http.Server`, `ServeMux`, and handlers
- Why the same path can support different operations by checking `r.Method`
- Why APIs should return JSON while the CLI can stay plain text
- Why `404` is the right mapping for `ErrJobNotFound`
- Why structured JSON errors need both HTTP status and application-specific error codes

### What needs revision

- Add request body size limits and content-type validation in the next reliability module
- Revisit how to keep persisted-data compatibility when serialized field names change
- Later compare manual path extraction with cleaner transport-layer organization

### Next session

- Start Day 8, Module 2.2 with middleware and API reliability on 2026-08-25

## 2026-08-26 — Module 2.2 completion

### Topics covered

- Middleware shape and chaining
- Panic recovery middleware
- Request IDs in headers and request context
- Request body size limiting
- Content-type validation
- More precise API error mapping for oversized bodies

### Work completed

- Added `cmd/goflow/middleware.go`
- Implemented `recoveryMiddleware`, `requestIDMiddleware`, `chain(...)`, `requestBodyLimitMiddleware(...)`, and `requireJSONMiddleware(...)`
- Wired the middleware stack into the HTTP server
- Added request ID context storage
- Mapped `*http.MaxBytesError` to `REQUEST_BODY_TOO_LARGE`
- Verified unsupported media type handling and removed the temporary panic-testing route after validation

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'cmd/goflow/http.go'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/middleware.go'`
- `gofmt -w ./cmd/goflow/http.go ./cmd/goflow/middleware.go ./cmd/goflow/main.go`
- `go vet ./...`
- `go test ./...`
- `curl -i http://localhost:8080/health/live`
- oversized `POST /v1/jobs`
- `Invoke-WebRequest -Method POST -Uri http://localhost:8080/v1/jobs -ContentType text/plain -Body '{"type":"email"}'`

### Problems encountered

- The first oversized-body implementation still returned `INVALID_REQUEST_BODY` until `*http.MaxBytesError` was checked explicitly
- A temporary panic endpoint was useful for proving recovery middleware, but it had to be removed afterward
- The content-type check is currently exact and will reject valid variants like `application/json; charset=utf-8`

### What I understood well

- Middleware wraps handlers and order changes behavior
- Recovery middleware should return controlled `500` responses instead of letting panics escape
- Request IDs are useful in both response headers and request context
- Oversized bodies deserve a more specific error than generic invalid JSON
- Content-type enforcement belongs at the HTTP boundary

### What needs revision

- Relax content-type validation to allow normal JSON charset variants later
- Add structured request logging that reads the request ID from context in the next stage
- Revisit selective middleware application as routes become more different from each other

### Next session

- Start Day 9, Module 2.3 with project organization and architecture on 2026-08-27

## 2026-08-29 — Module 2.3 completion

### Topics covered

- Package-boundary design
- Handler vs service vs repository responsibilities
- Dependency direction
- Why `main` should stay thin
- Using `internal/` for reusable app-only packages

### Work completed

- Created `internal/goflow`
- Moved `store.go`, `job.go`, `service.go`, and `persistence.go` into the reusable core package
- Exported `StartJob(...)`
- Updated `main.go` and `http.go` to import `goflow/internal/goflow`
- Removed old practice helpers from `main.go`
- Preserved existing CLI and HTTP behavior after the refactor

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-ChildItem -LiteralPath 'cmd/goflow'`
- `Get-ChildItem -LiteralPath 'internal/goflow'`
- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/http.go'`
- `Get-Content -LiteralPath 'cmd/goflow/middleware.go'`
- `Get-Content -LiteralPath 'internal/goflow/store.go'`
- `Get-Content -LiteralPath 'internal/goflow/job.go'`
- `Get-Content -LiteralPath 'internal/goflow/service.go'`
- `Get-Content -LiteralPath 'internal/goflow/persistence.go'`
- `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/middleware.go`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow serve`
- `curl http://localhost:8080/health/live`

### Problems encountered

- PowerShell wildcard formatting did not work directly with `gofmt`
- `StartJob(...)` had to be exported after moving service logic into another package
- `http.go` briefly had the wrong package declaration during the refactor

### What I understood well

- Why handlers and `main` should depend on a reusable core package instead of the reverse
- Why `internal/goflow` was the right first extraction target
- Why `service.go` is a clearer service-layer example than store/persistence files
- Why startup wiring should stay in `main` while leftover practice helpers should be removed

### What needs revision

- Later extract the HTTP transport into its own internal package once the core boundary is stable
- Revisit finer CLI separation if additional binaries or commands are introduced
- Continue tightening architecture as PostgreSQL replaces file persistence

### Next session

- Start Day 10, Module 2.4 with PostgreSQL and `database/sql` on 2026-08-30

## 2026-09-01 — Module 2.4 repository implementation step

### Topics covered

- `sql.DB` as a long-lived handle
- `PingContext(...)`
- `QueryRowContext(...)`, `QueryContext(...)`, and `ExecContext(...)`
- `Scan(...)`, `rows.Close()`, and `rows.Err()`
- Parameterized SQL and SQL injection prevention
- Repository boundaries for storage replacement

### Work completed

- Added `internal/goflow/postgres_repository.go`
- Added a context-aware `JobRepository` interface
- Implemented PostgreSQL-backed `Create`, `Get`, `List`, `Update`, and `Delete`
- Centralized row scanning in `scanJob(...)`
- Added `migrations/001_create_jobs.sql`
- Updated tracker, learning notes, design notes, and decisions for the Day 10 increment

### Commands run

- `Get-Content -LiteralPath 'Go_Industry_Roadmap_4_Weeks.md'`
- `Get-Content -LiteralPath 'docs/tracker.md'`
- `Get-Content -LiteralPath 'docs/learning.md'`
- `Get-Content -LiteralPath 'docs/session-log.md'`
- `Get-Content -LiteralPath 'internal/goflow/store.go'`
- `Get-Content -LiteralPath 'internal/goflow/job.go'`
- `Get-Content -LiteralPath 'internal/goflow/service.go'`
- `gofmt -w ./internal/goflow/postgres_repository.go`
- `go vet ./...`
- `go test ./...`

### Problems encountered

- The local sandbox helper kept failing, so file reads and writes had to be rerun with elevated access
- `docs/design.md` had malformed Markdown fences and had to be rebuilt cleanly
- No PostgreSQL driver or live database is wired yet, so runtime DB behavior remains unvalidated

### What I understood well

- `Get(id)` should use `QueryRowContext(...)`
- `List()` should use `QueryContext(...)`
- `Scan(...)` destinations must match query column order
- Parameterized queries prevent SQL injection better than string interpolation
- A repository boundary reduces how far persistence changes spread

### What needs revision

- Add driver setup and `PingContext(...)` in startup
- Replace the JSON-backed runtime path with the PostgreSQL repository
- Map database duplicate-key failures to `ErrJobAlreadyExists`

### Next session

- Wire a real `*sql.DB` into startup and move one app path from JSON-file persistence to PostgreSQL

## 2026-09-01 — Module 2.4 runtime switch

### Topics covered

- `sql.DB` startup and pooling
- `PingContext(...)`
- PostgreSQL driver wiring
- Context-aware repository calls
- Replacing file-backed runtime flow with PostgreSQL
- Automatic migration execution at startup

### Work completed

- Added `github.com/lib/pq`
- Added `cmd/goflow/database.go`
- Switched the CLI commands to PostgreSQL-backed persistence
- Switched HTTP handlers to a shared store dependency
- Updated `StartJob(...)` to use `context.Context`
- Added duplicate-key mapping for PostgreSQL errors
- Replaced sequential ID generation with random IDs
- Ran `go mod tidy`, `go vet ./...`, and `go test ./...`
- Updated tracker, learning notes, design notes, and README

### Commands run

- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/http.go'`
- `Get-Content -LiteralPath 'internal/goflow/postgres_repository.go'`
- `Get-Content -LiteralPath 'internal/goflow/service.go'`
- `Get-Content -LiteralPath 'go.mod'`
- `go get github.com/lib/pq`
- `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/database.go ./internal/goflow/service.go ./internal/goflow/postgres_repository.go`
- `go mod tidy`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow list`
- `docker ps -a`

### Problems encountered

- The Windows sandbox helper still blocked normal patch-based edits
- `DATABASE_URL` is not set in this environment
- `psql` is not installed here
- Docker Desktop is not running, so local PostgreSQL validation could not be started from Docker

### What I understood well

- Why `sql.DB` should be opened once and reused
- Why `PingContext(...)` belongs in startup
- Why handlers should use request context for DB work
- Why the storage swap belongs behind a repository boundary
- Why parameterized queries and duplicate-key mapping matter at the repository layer

### What needs revision

- Validate the new runtime against a real PostgreSQL instance
- Add repository and handler tests in Module 2.5
- Revisit whether `payload` should evolve from `BYTEA` to `JSONB`

### Next session

- Start a PostgreSQL instance, set `DATABASE_URL`, and validate the CLI plus HTTP paths end to end before moving to Module 2.5

## 2026-09-01 — Module 2.4 runtime switch

### Topics covered

- `sql.DB` startup and pooling
- `PingContext(...)`
- PostgreSQL driver wiring
- Context-aware repository calls
- Replacing file-backed runtime flow with PostgreSQL
- Automatic migration execution at startup

### Work completed

- Added `github.com/lib/pq`
- Added `cmd/goflow/database.go`
- Switched the CLI commands to PostgreSQL-backed persistence
- Switched HTTP handlers to a shared store dependency
- Updated `StartJob(...)` to use `context.Context`
- Added duplicate-key mapping for PostgreSQL errors
- Replaced sequential ID generation with random IDs
- Ran `go mod tidy`, `go vet ./...`, and `go test ./...`
- Updated tracker, learning notes, design notes, and README

### Commands run

- `Get-Content -LiteralPath 'cmd/goflow/main.go'`
- `Get-Content -LiteralPath 'cmd/goflow/http.go'`
- `Get-Content -LiteralPath 'internal/goflow/postgres_repository.go'`
- `Get-Content -LiteralPath 'internal/goflow/service.go'`
- `Get-Content -LiteralPath 'go.mod'`
- `go get github.com/lib/pq`
- `gofmt -w ./cmd/goflow/main.go ./cmd/goflow/http.go ./cmd/goflow/database.go ./internal/goflow/service.go ./internal/goflow/postgres_repository.go`
- `go mod tidy`
- `go vet ./...`
- `go test ./...`
- `go run ./cmd/goflow list`
- `docker ps -a`

### Problems encountered

- The Windows sandbox helper still blocked normal patch-based edits
- `DATABASE_URL` is not set in this environment
- `psql` is not installed here
- Docker Desktop is not running, so local PostgreSQL validation could not be started from Docker

### What I understood well

- Why `sql.DB` should be opened once and reused
- Why `PingContext(...)` belongs in startup
- Why handlers should use request context for DB work
- Why the storage swap belongs behind a repository boundary
- Why parameterized queries and duplicate-key mapping matter at the repository layer

### What needs revision

- Validate the new runtime against a real PostgreSQL instance
- Add repository and handler tests in Module 2.5
- Revisit whether `payload` should evolve from `BYTEA` to `JSONB`

### Next session

- Start a PostgreSQL instance, set `DATABASE_URL`, and validate the CLI plus HTTP paths end to end before moving to Module 2.5
