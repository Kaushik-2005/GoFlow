# Go Learning Tracker

## Current Position

- Current week: Week 1
- Current module: Module 1.2 | 2026-08-18
- Current topic: variables, zero values, and short declaration with `:=`
- Status: Learning
- Started: 2026-08-18
- Last updated: 2026-08-18

## Overall Progress

- [ ] Week 1: Foundations and idiomatic Go
- [ ] Week 2: HTTP, architecture, databases, and testing
- [ ] Week 3: Concurrency and background processing
- [ ] Week 4: Production readiness and deployment

## Current Module: Module 1.1 | 2026-08-18

### Learning objectives

- [x] Explain packages versus modules
- [x] Explain why `package main` and `func main` are special
- [x] Use the basic Go toolchain commands
- [x] Create the initial GoFlow module structure

### Theory

- [x] Required reading completed
- [x] Core concept explained
- [x] Questions answered
- [x] Knowledge check passed

### Practice

- [x] Small exercise completed
- [x] Exercise reviewed
- [x] Mistakes documented

### Project work

- [x] Implementation task
- [ ] Tests added
- [x] Validation passed

### Completion evidence

- Exercise: explained package vs module; explained why only `package main` with `func main()` is executable; implemented CLI command dispatch and command-specific argument validation
- Files changed: `go.mod`, `cmd/goflow/main.go`, `docs/tracker.md`, `docs/learning.md`, `docs/session-log.md`, `docs/decisions.md`, `README.md`
- Tests run: `go run ./cmd/goflow`, `go build ./cmd/goflow`, `go test ./...`, `go vet ./...`, `gofmt -w ./cmd/goflow/main.go`
- Concepts demonstrated: module root vs package directory; `main` package entry point; toolchain basics; `os.Args`; safe slice bounds checks; small helper refactor
- Remaining weaknesses: Module 1.2 theory and exercises not started yet; no automated tests written yet for the CLI

## Blockers

- Repository is not initialized as a Git repository yet.

## Next Session

- Module: Module 1.2 | 2026-08-18
- Exact next topic: `var`, zero values, type inference, and `:=`
- Exact next exercise: explain the difference between `var count int` and `count := 0`, then write a small variables exercise
- Exact next project task: keep the existing CLI entry point and use Module 1.2 concepts in a small focused Go exercise before expanding project behavior
