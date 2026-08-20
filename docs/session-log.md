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
