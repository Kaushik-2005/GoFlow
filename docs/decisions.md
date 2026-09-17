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


## Decision: Command owners control cancellation and shutdown

- Date: 2026-09-04
- Roadmap module: Module 3.4
- Status: Accepted
- Context: GoFlow now has long-running HTTP and worker commands. Both need to react to Ctrl+C without leaking goroutines or abruptly stopping active work.
- Options considered: let each handler or worker own independent cancellation; use one command-owned signal context and pass it down; rely on process exit without graceful shutdown
- Decision: The top-level command owns the root signal context, derives short child contexts for startup or shutdown deadlines, and waits for owned goroutines before returning.
- Why: Lifecycle ownership stays clear. Startup, HTTP serving, polling, worker execution, and cleanup each receive the context appropriate to their duration.
- Trade-offs: `main.go` has more lifecycle code until the command structure is split further.
- Consequences: Future worker retry logic and shutdown behavior should continue to use command-owned cancellation and explicit waiting rather than ad hoc stop flags.

## Decision: Model dead-letter as job status before adding a separate queue

- Date: 2026-09-08
- Roadmap module: Module 3.5
- Status: Accepted
- Context: GoFlow needs to represent jobs that should no longer be retried automatically after permanent failure or retry-budget exhaustion.
- Options considered: add a separate dead-letter table immediately; model dead-letter as a `JobStatus` first
- Decision: Use `StatusDeadLetter` in the existing jobs table for the first retry implementation.
- Why: This keeps the module focused on retry semantics, failure classification, and state transitions without adding extra schema and repository complexity.
- Trade-offs: A separate dead-letter table may be cleaner later for retention, auditing, and operations.
- Consequences: Future production-readiness work can promote dead-letter jobs into a separate table or queue if operational needs justify it.

## Decision: Use available_at for persistent retry scheduling

- Date: 2026-09-08
- Roadmap module: Module 3.5
- Status: Accepted
- Context: Retryable jobs need to wait before being picked again, but workers should not sleep while holding work.
- Options considered: sleep inside the worker before retrying; store retry availability in PostgreSQL through `available_at`
- Decision: Store retry timing in `jobs.available_at` and have the poller call `ListReadyJobs(ctx)`.
- Why: The retry schedule survives restarts, remains inspectable, and lets workers process other jobs instead of blocking.
- Trade-offs: The repository needs a worker-specific query and the schema gains another operational field.
- Consequences: Future indexes should include `status` and `available_at` for efficient ready-job polling.

## Decision: Treat unauthenticated API as local/internal only

- Date: 2026-09-14
- Roadmap module: Module 4.4
- Status: Accepted
- Context: GoFlow now exposes state-changing HTTP endpoints such as `POST /v1/jobs` and `DELETE /v1/jobs/{id}`, but it does not yet have authentication or authorization.
- Options considered: expose the API publicly as-is; add a full auth system immediately; keep the API local/internal until auth is designed and implemented
- Decision: Treat the current HTTP API as local/internal only and document that it must not be exposed directly to the public internet.
- Why: The module can harden boundaries without prematurely adding a half-designed auth system.
- Trade-offs: The API remains unsuitable for public deployment until a future auth layer exists.
- Consequences: Future deployment work must include authentication, authorization, or an external gateway before public exposure.

## Decision: Own migrations with a separate command

- Date: 2026-09-15
- Roadmap module: Module 4.5
- Status: Accepted
- Context: Docker Compose starts API and worker as separate processes. When both ran migrations during startup, PostgreSQL hit a concurrent schema creation error.
- Options considered: let every process run migrations; run migrations manually only; add a dedicated `goflow migrate` command and Compose migration service
- Decision: Add `goflow migrate` and make API/worker start only after the migration service completes successfully.
- Why: Schema changes should have one clear owner, and runtime services should focus on serving requests or processing jobs.
- Trade-offs: Local and Compose startup now have one extra command/service.
- Consequences: Future deployments should run migrations as a distinct release/deployment step before starting application workloads.

## Decision: Use multi-stage Docker build with distroless runtime

- Date: 2026-09-15
- Roadmap module: Module 4.5
- Status: Accepted
- Context: GoFlow needs a container image for production-style runtime validation.
- Options considered: ship the full Go toolchain image; use Alpine with shell tools; use a multi-stage build and distroless non-root runtime
- Decision: Compile in a `golang` builder image and copy only the binary and migrations into a distroless non-root runtime image.
- Why: This keeps the runtime image smaller and reduces attack surface.
- Trade-offs: The runtime image has no shell or `curl`, so container-internal health checks need a different strategy.
- Consequences: Health validation is currently external through Compose/curl; future production health checks should be orchestrator-based or use a small dedicated healthcheck binary.



## Decision: Use GitHub Actions as the main CI gate

- Date: 2026-09-17
- Roadmap module: Module 4.6
- Status: Accepted
- Context: GoFlow now has Go code, tests, PostgreSQL integration tests, Docker packaging, and Compose runtime configuration that should be validated before merge.
- Options considered: rely on manual local commands only; add GitHub Actions CI with PostgreSQL service validation; add a heavier external CI platform
- Decision: Use GitHub Actions as the primary CI workflow and run Go, security, Docker, and Compose validation there.
- Why: GitHub Actions fits the repository workflow, supports service containers, and keeps the validation path visible with minimal extra infrastructure.
- Trade-offs: CI duration increases because it runs integration tests, race tests, vulnerability scanning, and Docker build validation.
- Consequences: Future changes should keep the workflow green and add new validation steps when new production responsibilities are added.


## Decision: Claim jobs with a conditional database update

- Date: 2026-09-18
- Roadmap module: Final review after Module 4.6
- Status: Accepted
- Context: Final review found that `StartJob` used a read-then-write sequence, which could allow two workers to observe the same pending job before either update completed.
- Options considered: keep service-level `Get` plus `Update`; protect only with the in-process queued map; add a repository-level conditional update claim
- Decision: Add `ClaimPending(ctx, id)` and implement the claim as `UPDATE jobs SET status = 'running' WHERE id = $1 AND status = 'pending'`.
- Why: PostgreSQL can enforce the ownership transition atomically, including across multiple worker processes.
- Trade-offs: The store interface gained one persistence-specific capability and test fakes need to implement it.
- Consequences: Workers now rely on a real database claim before execution, while job executors still need idempotency for crash-after-side-effect cases.
