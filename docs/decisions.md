# Decisions

## Decision: Start GoFlow as a single CLI binary

- Date: 2026-08-18
- Roadmap module: Module 1.1
- Status: Accepted
- Context: The roadmap begins with packages, modules, and tooling, and the repository currently contains no Go project files.
- Options considered: start with the full eventual API and worker layout; start with one CLI binary under `cmd/goflow`
- Decision: Start with one CLI entry point and add the larger architecture later when the roadmap reaches it.
- Why: This keeps the first module focused on Go structure and tooling rather than premature architecture.
- Trade-offs: The first scaffold will be intentionally simple and will later be reorganized.
- Consequences: Early project code should optimize for learning module/package structure, not for final production layout.

## Decision: Add a PostgreSQL repository behind a context-aware interface

- Date: 2026-09-01
- Roadmap module: Module 2.4
- Status: Accepted
- Context: GoFlow currently persists jobs to `jobs.json`, but the roadmap now requires a move toward PostgreSQL using `database/sql` without letting database details leak through the whole application.
- Options considered: replace the file store everywhere in one step; add a PostgreSQL repository boundary first and wire runtime startup later
- Decision: Add a context-aware `JobRepository` plus a PostgreSQL implementation in `internal/goflow` before switching the app startup path.
- Why: This keeps the storage transition controlled, teaches the correct `database/sql` patterns in isolation, and avoids mixing database concerns into handlers and startup code all at once.
- Trade-offs: The app is temporarily in a mixed state where the database repository exists but the runtime still uses the JSON-backed store.
- Consequences: The next increment must add driver setup, `PingContext(...)`, and real wiring from `main.go` into the database-backed repository.

