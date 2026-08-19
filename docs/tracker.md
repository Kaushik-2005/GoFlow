# Go Learning Tracker

## Current Position

- Current week: Week 1
- Current day and module: Day 3 | Module 1.3 | 2026-08-19
- Current topic: arrays vs slices, slice length/capacity, and `append`
- Status: Learning
- Started: 2026-08-18
- Last updated: 2026-08-19

## Overall Progress

- [ ] Week 1: Foundations and idiomatic Go
- [ ] Week 2: HTTP, architecture, databases, and testing
- [ ] Week 3: Concurrency and background processing
- [ ] Week 4: Production readiness and deployment

## Current Module: Day 2 | Module 1.2 | 2026-08-19

### Learning objectives

- [x] Explain variables and constants
- [x] Use short declaration with `:=`
- [x] Explain zero values for common types
- [x] Use `if`, `switch`, `for`, and `range`
- [x] Write basic functions and explain multiple returns
- [x] Explain `defer`, scope, shadowing, and `iota`

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

- Exercise: completed focused exercises for `maxPriority`, `countByStatus`, `isValidJobName`, `retryDelay`, and `filterCompleted`; answered concept checks on zero values, `:=` vs `=`, `if`, `range`, multiple returns, `defer`, shadowing, and `iota`
- Files changed: `cmd/goflow/main.go`, `docs/tracker.md`, `docs/learning.md`, `docs/session-log.md`
- Tests run: `gofmt -w ./cmd/goflow/main.go`, `go vet ./...`, `go test ./...`
- Concepts demonstrated: zero values; constants vs variables; explicit boolean conditions; safe `range` usage with `_`; multiple returns; deferred cleanup; shadowing; `iota`; simple loop-based algorithms over slices and maps
- Remaining weaknesses: needs more repetition for precise wording around `:=` vs `=` and the exact purpose of `defer`; no automated tests written yet for the practice functions

## Blockers

- None

## Next Session

- Day and module: Day 3 | Module 1.3 | 2026-08-19
- Exact next topic: arrays vs slices, slice length/capacity, and how `append` grows slices
- Exact next exercise: explain the difference between an array and a slice, then inspect how filtering with `append` built a new slice in `filterCompleted`
- Exact next project task: start an in-memory job model and storage helpers using slices or maps before building the full Module 1.3 job store
