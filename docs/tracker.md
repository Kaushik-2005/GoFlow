# Go Learning Tracker

## Current Position

- Current week: Week 1
- Current module: Day 4 - Module 1.4
- Current topic: Structs, methods, pointers, and interfaces
- Current task: Start the Job struct and method-oriented design for GoFlow after completing the Day 3 collection/storage foundations
- Next milestone: Explain structs vs plain grouped values and begin method-based behavior on the job model
- Active blockers: None

## Roadmap Progress

| Day | Module | Main Topic | Status | Deliverable | Verification | Confidence |
|---|---|---|---|---|---|---|
| 1 | Module 1.1 | Environment, packages, modules, and tooling | Completed | Runnable Go module with CLI entry point and command dispatch scaffold | `go mod init goflow`; `go run ./cmd/goflow`; `go build ./cmd/goflow`; `go test ./...`; `go vet ./...`; learner explained package vs module and why `package main` + `func main()` is executable | 4 |
| 2 | Module 1.2 | Language fundamentals | Completed | Core helper functions for max priority, status counting, validation, retry delay, and filtering | `gofmt -w ./cmd/goflow/main.go`; `go vet ./...`; `go test ./...`; learner answered theory checks on zero values, `:=` vs `=`, `if`, `range`, `defer`, shadowing, and `iota` | 4 |
| 3 | Module 1.3 | Arrays, slices, maps, strings, and runes | Completed | In-memory job model and early storage helpers using slices/maps | `gofmt -w ./cmd/goflow/store.go`; `go vet ./...`; `go test ./...`; learner explained arrays vs slices, `append`, shared backing arrays, `copy`, nil vs empty slices, map lookup/delete, `make` vs `new`, and value semantics; implemented in-memory `Store` with create/get/list/update/delete | 4 |
| 4 | Module 1.4 | Structs, methods, pointers, and interfaces | Not Started | Job struct behavior and first clean boundaries | Resume point set for 2026-08-21 | — |
| 5 | Module 1.5 | Error handling | Not Started | Explicit job-store errors and wrapped failures | — | — |
| 6 | Module 1.6 | I/O, JSON, files, and configuration | Not Started | Week 1 CLI deliverable with JSON-file persistence | — | — |
| 7 | Module 2.1 | HTTP servers | Not Started | First `net/http` API endpoints | — | — |
| 8 | Module 2.2 | Middleware and API reliability | Not Started | Request middleware stack and reliability guards | — | — |
| 9 | Module 2.3 | Project organization and architecture | Not Started | Clear layered structure for API and worker code | — | — |
| 10 | Module 2.4 | PostgreSQL and `database/sql` | Not Started | PostgreSQL-backed repository implementation | — | — |
| 11 | Module 2.5 | Testing fundamentals | Not Started | Tested REST API with repository coverage | — | — |
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
