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
