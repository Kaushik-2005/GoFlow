# GoFlow Design

## Purpose

This document records how GoFlow's design evolved across the roadmap. It is organized day by day so changes in architecture, boundaries, and trade-offs stay easy to review.

## Design Principles

- Keep external behavior simple while the internal architecture is still being learned.
- Separate transport, business rules, persistence, and cross-cutting reliability concerns.
- Prefer the standard library first so the Go model stays visible.
- Keep interfaces small and define them where they are consumed.
- Add reliability protections at the HTTP boundary instead of duplicating them in handlers.
- Accept temporary learning-stage shortcuts only when they are clearly documented and easy to replace.

## Current Architecture Snapshot

### Layer Summary

- CLI layer: command dispatch in `cmd/goflow/main.go`.
- HTTP transport layer: handlers in `cmd/goflow/http.go`.
- Middleware layer: request IDs, panic recovery, request body limits, and content-type enforcement in `cmd/goflow/middleware.go`.
- Core domain layer: `Job`, `JobStatus`, store logic, persistence helpers, service orchestration, and repository contracts in `internal/goflow`.
- Persistence layer: PostgreSQL-backed repository access through `database/sql` plus the initial SQL migration; `jobs.json` is now legacy learning-stage data rather than the active runtime store.

### Current High-Level Diagram

```mermaid
flowchart TD
    CLI[CLI in cmd/goflow] --> Core[internal/goflow]
    HTTP[HTTP handlers in cmd/goflow] --> Core
    Middleware[Middleware] --> HTTP
    Core --> Repo[PostgresRepository]
    Repo --> DB[(PostgreSQL)]
```

## Day 1 - Module 1.1: Executable Entry Point

### Design impact

- Started with one executable under `cmd/goflow`.
- Kept the app intentionally flat to focus on packages, modules, and the Go toolchain.

### Diagram

```mermaid
flowchart TD
    User[Terminal user] --> Main[cmd/goflow/main.go]
    Main --> Output[CLI output]
```

## Day 2 - Module 1.2: Language Practice in `main.go`

### Design impact

- Kept small practice helpers local to `main.go`.
- Deferred package separation until responsibilities became real.

### Diagram

```mermaid
flowchart TD
    Main[main.go] --> Helpers[Practice helpers]
    Helpers --> Output[Printed results]
```

## Day 3 - Module 1.3: In-Memory Store

### Design impact

- Introduced `Job`, `JobStatus`, and `Store`.
- Established the first reusable domain and storage boundary.

### Diagram

```mermaid
flowchart TD
    CLI[CLI command] --> Store[Store]
    Store --> JobsMap[(map[string]Job)]
    Store --> ListSlice[[[]Job]]
```

## Day 4 - Module 1.4: Job Behavior and Small Interfaces

### Design impact

- Moved behavior onto `Job` with methods like `CanRetry()` and `MarkRunning()`.
- Added orchestration with a small consumer-defined interface.

### Diagram

```mermaid
flowchart TD
    CLI[CLI process command] --> Service[Start job logic]
    Service --> JobMethods[Job methods]
    Service --> StoreInterface[Small interface]
    StoreInterface --> StoreImpl[Store]
```

## Day 5 - Module 1.5: Explicit Error Model

### Design impact

- Replaced boolean-style failures with explicit errors.
- Added sentinel errors and a custom typed error for invalid state transitions.

### Diagram

```mermaid
flowchart TD
    Store[Store methods] --> Sentinel[Sentinel errors]
    Service[Service layer] --> Wrapped[Wrapped errors]
    CLI[CLI boundary] --> Inspect[errors.Is / errors.As]
```

## Day 6 - Module 1.6: JSON File Persistence

### Design impact

- Added `saveJobs(...)`, `loadJobs(...)`, `Store.Save(...)`, and `LoadStore(...)`.
- Moved from ephemeral CLI behavior to persisted local state in `jobs.json`.

### Diagram

```mermaid
flowchart TD
    CLI[CLI command] --> Core[Store]
    Core --> Persistence[persistence helpers]
    Persistence --> JSON[(jobs.json)]
```

## Day 7 - Module 2.1: First HTTP API

### Design impact

- Added the first `net/http` transport layer.
- Reused the same job model and store behavior across CLI and HTTP.

### Diagram

```mermaid
flowchart TD
    Client[HTTP client] --> Server[http.Server]
    Server --> Mux[ServeMux]
    Mux --> Handlers[HTTP handlers]
    Handlers --> Core[Store / service]
    Core --> JSON[(jobs.json)]
```

## Day 8 - Module 2.2: Middleware and Reliability Guards

### Design impact

- Added a distinct middleware layer.
- Centralized request IDs, recovery, body limits, and JSON content-type enforcement.

### Diagram

```mermaid
flowchart TD
    Client[HTTP client] --> ReqID[requestIDMiddleware]
    ReqID --> Recovery[recoveryMiddleware]
    Recovery --> Limit[requestBodyLimitMiddleware]
    Limit --> JSON[requireJSONMiddleware]
    JSON --> Mux[ServeMux]
    Mux --> Handler[HTTP handler]
```

## Day 9 - Module 2.3: Core Package Extraction

### Design impact

- Extracted reusable application logic into `internal/goflow`.
- Corrected dependency direction so CLI and HTTP depend on the core package rather than the reverse.

### Diagram

```mermaid
flowchart TD
    CLI[cmd/goflow main.go CLI] --> Core[internal/goflow]
    HTTP[cmd/goflow http.go handlers] --> Core
    Middleware[cmd/goflow middleware.go] --> HTTP
    Core --> JSON[(jobs.json)]
```

### Flow

```mermaid
sequenceDiagram
    participant Main as cmd/goflow main.go
    participant HTTP as http.go
    participant Core as internal/goflow
    participant File as jobs.json

    Main->>Core: LoadStore / StartJob / Job creation
    HTTP->>Core: LoadStore / Store methods / errors
    Core->>File: read and write persisted jobs
    Core-->>Main: shared logic results
    Core-->>HTTP: shared logic results
```

## Day 10 - Module 2.4: PostgreSQL Runtime Switch

### Goal

Replace the old file-based runtime path with PostgreSQL while keeping the higher-level service boundary clean.

### Design change

- Added `cmd/goflow/database.go` for database startup helpers.
- Added the PostgreSQL driver dependency with `github.com/lib/pq`.
- Added `openPostgresStore(...)` to read `DATABASE_URL`, open `*sql.DB`, configure pooling, run `PingContext(...)`, and apply the first migration.
- Switched CLI commands from `LoadStore(jobs.json)` to the PostgreSQL-backed repository.
- Switched HTTP handlers from file reloads to one shared store dependency.
- Updated `StartJob(...)` to use `context.Context` and the shared `JobStore` abstraction.
- Replaced sequential ID generation with random ID generation.

### Diagram

```mermaid
flowchart TD
    CLI[CLI in cmd/goflow] --> Main[main.go]
    HTTP[HTTP requests] --> Handlers[http.go handlers]
    Main --> DBSetup[database.go startup helpers]
    Handlers --> Store[JobStore]
    Main --> Store
    Store --> PG[PostgresRepository]
    DBSetup --> PG
    PG --> DB[(PostgreSQL)]
    DBSetup --> Migration[001_create_jobs.sql]
```

### Flow

```mermaid
sequenceDiagram
    participant Caller as CLI or HTTP layer
    participant Main as startup/wiring
    participant Repo as PostgresRepository
    participant DB as PostgreSQL

    Caller->>Main: start command or request path
    Main->>DB: sql.Open + PingContext + migration
    Main->>Repo: create shared repository
    Caller->>Repo: Create/Get/List/Update/Delete with context
    Repo->>DB: ExecContext / QueryRowContext / QueryContext
    DB-->>Repo: result or rows
    Repo-->>Caller: job data or error
```

### Why this matters

- The app now targets a real relational persistence boundary.
- The service layer still depends on a store abstraction instead of SQL details.
- HTTP request context can flow into DB work.
- Startup owns connectivity checks and schema setup instead of handlers.

## Known Current Limitations

- Live PostgreSQL validation completed on 2026-09-02 through the repository integration test and real HTTP requests.
- The current schema still uses `BYTEA` for `payload`, which is simpler than the roadmap's suggested `JSONB` shape.
- The runtime was validated against a Docker-based PostgreSQL 16 container named `goflow-postgres`.
- Repository behavior now has `sqlmock` tests plus a real integration-test entry point.
- Content-type validation is still strict and does not yet allow variants like `application/json; charset=utf-8`.

## Related Documents

- `docs/tracker.md`
- `docs/learning.md`
- `docs/session-log.md`
- `docs/decisions.md`

## Day 11 - Module 2.5: Testing Layer Introduction

### Goal

Complete the first real testing layer for the PostgreSQL-backed API and ensure the Week 2 HTTP surface matches the roadmap deliverable.

### Design change

- Added `internal/goflow/service_test.go` for service-layer behavior.
- Added `cmd/goflow/http_test.go` for handler-level HTTP behavior.
- Added `internal/goflow/postgres_repository_test.go` for repository unit tests with `sqlmock`.
- Added `internal/goflow/postgres_repository_integration_test.go` for the real PostgreSQL path.
- Extended the HTTP API to support `DELETE /v1/jobs/{id}` and `GET /v1/jobs?status=...`.
- Introduced `jobByIDHandler` so one path can dispatch by method for single-job operations.

### Testing boundaries

- Service tests validate state-transition rules.
- Handler tests validate status codes, JSON bodies, validation failures, filtering, and deletion behavior.
- Repository tests validate SQL interactions and error mapping.
- Integration tests validate the live PostgreSQL path.

### Diagram

```mermaid
flowchart TD
    ServiceTests[service_test.go] --> Service[StartJob]
    Service --> FakeStore[fakeJobStore]
    HandlerTests[http_test.go] --> Handlers[HTTP handlers]
    Handlers --> FakeAPIStore[fakeAPIStore]
    HandlerTests --> HTTPTest[httptest request/recorder]
    RepoTests[postgres_repository_test.go] --> SQLMock[sqlmock]
    Integration[TestPostgresRepositoryIntegration] --> Postgres[(PostgreSQL)]
```

### Flow

```mermaid
sequenceDiagram
    participant Test as test case
    participant Boundary as service, handler, or repository
    participant Double as fake, sqlmock, or PostgreSQL

    Test->>Boundary: call function or handler
    Boundary->>Double: use dependency
    Double-->>Boundary: controlled or real result
    Boundary-->>Test: response or error
    Test-->>Test: assert status, body, state change, or DB effect
```

### Why this matters

- Testing now matches the architecture boundaries introduced earlier.
- The HTTP API now matches the Week 2 roadmap surface: create, retrieve, list, filter, and delete.
- Repository tests reduce the risk of silent SQL regressions.
- Live integration validation proves the PostgreSQL runtime, not just the test doubles.

## Day 11 Validation Update

- Live PostgreSQL integration was validated on 2026-09-02 with Docker PostgreSQL, `TestPostgresRepositoryIntegration`, and real HTTP requests to `/health/live`, `/v1/jobs`, and `/v1/jobs/{id}`.
- The test strategy now has four useful levels: service fakes, handler tests with `httptest`, repository tests with `sqlmock`, and a real database integration path.
- The Week 2 API surface now includes both `status` filtering on `GET /v1/jobs` and deletion through `DELETE /v1/jobs/{id}`.

## Day 13 - Module 3.1: First Worker Pipeline

### Goal

Introduce the first in-process concurrency flow for GoFlow without yet adding a full worker pool or atomic claim logic.

### Design change

- Added a new CLI command: `work`.
- `work` opens the PostgreSQL-backed store once.
- `work` creates `jobs chan string` carrying job IDs.
- One worker goroutine receives IDs and calls `goflow.StartJob(...)`.
- The poll loop lists pending jobs and enqueues their IDs.
- The worker now listens to both the jobs channel and `workCtx.Done()`.
- `signal.NotifyContext(...)` is used so the command can stop on interrupt.

### Diagram

```mermaid
flowchart TD
    Work[go run ./cmd/goflow work] --> Poller[poll loop]
    Poller --> Jobs[jobs chan string]
    Jobs --> Worker[single worker goroutine]
    Worker --> Service[StartJob]
    Service --> Repo[PostgresRepository]
    Repo --> DB[(PostgreSQL)]
```

### Flow

```mermaid
sequenceDiagram
    participant Work as work command
    participant Poller as poll loop
    participant Ch as jobs channel
    participant Worker as worker goroutine
    participant DB as PostgreSQL

    Work->>Poller: start loop
    Poller->>DB: list jobs
    Poller->>Ch: send pending job ID
    Worker->>Ch: receive job ID
    Worker->>DB: StartJob -> Get + Update
    DB-->>Worker: job claimed as running
```

### Why this matters

- This is the first real producer/consumer pipeline in GoFlow.
- The command now demonstrates goroutine ownership, channel handoff, and context-based shutdown in real project code.
- The design still has duplicate-enqueue risk, which is a useful lead-in to later synchronization and claiming work.

### Day 13 completion update

- The `jobs` channel is now bounded with a small buffer instead of being unbuffered.
- This makes the queue behavior closer to a real work buffer while still preserving backpressure once the buffer fills.
- Shutdown now propagates through one owned context to both the poller and the worker.

## Day 14 - Module 3.2: Shared Bookkeeping Synchronization

### Goal

Protect the first shared in-memory bookkeeping used by the worker pipeline as GoFlow moves from basic concurrency into coordinated concurrency.

### Design change

- Added a mutex-protected `queued` set local to the `work` command.
- The poller records a job ID in the set before enqueueing it.
- The worker removes the job ID after the claim attempt finishes.
- The channel remains the communication path; the mutex protects only the separate shared map.

### Diagram

```mermaid
flowchart TD
    Poller[poll loop] --> QueueSet[queued set + mutex]
    Poller --> Jobs[jobs channel]
    Jobs --> Worker[worker goroutine]
    Worker --> QueueSet
    Worker --> Service[StartJob]
```

### Why this matters

- This is the first explicit shared-memory synchronization step in GoFlow.
- It demonstrates the boundary between channel-based communication and mutex-protected bookkeeping.
- It reduces duplicate enqueueing within the current process and prepares the design for a later worker pool.

## Day 15 - Module 3.3: Worker Pool Expansion

### Goal

Expand the single-worker pipeline into a small bounded worker pool while keeping ownership, queue behavior, and shutdown coordination explicit.

### Design change

- Added `workerCount = 3`.
- Replaced one worker goroutine with a loop that starts three workers.
- Added `sync.WaitGroup` so the `work` command waits for all workers to stop.
- Added worker IDs to processing, failure, and shutdown logs.
- Kept the same shared `jobs` channel and mutex-protected `queued` set.

### Diagram

```mermaid
flowchart TD
    Poller[poll loop] --> Jobs[jobs channel]
    Jobs --> W1[worker 1]
    Jobs --> W2[worker 2]
    Jobs --> W3[worker 3]
    W1 --> Service[StartJob]
    W2 --> Service
    W3 --> Service
    Service --> Repo[PostgresRepository]
```

### Why this matters

- GoFlow now has bounded concurrent consumers instead of a single worker.
- The shared queue plus fixed worker count demonstrates the standard worker-pool pattern.
- Coordinated shutdown with `WaitGroup` becomes visible and necessary once the pool has several workers.

## Day 16 - Module 3.4: Context and Graceful Shutdown

### Goal

Give GoFlow's long-running `serve` and `work` commands explicit cancellation ownership and graceful shutdown behavior.

### Design change

- `serve` now owns a signal context for server lifetime.
- database startup uses a short child timeout context.
- `ListenAndServe()` runs in a goroutine so the main goroutine can wait for Ctrl+C.
- `server.Shutdown(...)` uses a fresh 5-second timeout context.
- `work` now treats the poller as the owner of the jobs channel and worker shutdown sequence.
- worker shutdown uses `close(jobs)` plus `WaitGroup` coordination.

### HTTP shutdown flow

```mermaid
flowchart TD
    Start[go run ./cmd/goflow serve] --> SignalCtx[signal.NotifyContext]
    SignalCtx --> StartupCtx[5s startup context]
    StartupCtx --> DB[open PostgreSQL store]
    DB --> Server[http.Server]
    Server --> Listen[ListenAndServe goroutine]
    SignalCtx --> Interrupt[Ctrl+C cancels serverCtx]
    Interrupt --> ShutdownCtx[fresh 5s shutdown context]
    ShutdownCtx --> Shutdown[server.Shutdown]
    Shutdown --> Stop[server stopped]
```

### Worker shutdown flow

```mermaid
flowchart TD
    Work[go run ./cmd/goflow work] --> WorkCtx[signal.NotifyContext]
    WorkCtx --> Poller[poll pending jobs]
    Poller --> Queue[jobs channel]
    Queue --> Workers[worker pool]
    WorkCtx --> Interrupt[Ctrl+C cancels workCtx]
    Interrupt --> StopPoller[poller stops]
    StopPoller --> CloseQueue[close jobs channel]
    CloseQueue --> WorkerExit[workers exit]
    WorkerExit --> Wait[WaitGroup waits]
    Wait --> Done[work command exits]
```

### Why this matters

- The application now has explicit owners for cancellation and cleanup.
- The HTTP path can stop without abruptly killing active requests.
- The worker path can stop polling, signal workers, wait for them, and then release resources.
- Concurrent log output is intentionally unordered, so correctness is based on shutdown coordination rather than printed order.

## Day 17 - Module 3.5: Retries and Failure Handling

### Goal

Add a real worker failure state machine so GoFlow can complete successful jobs, retry temporary failures, and preserve exhausted or permanent failures as dead-letter jobs.

### Design change

- Added `available_at` as persistent retry scheduling metadata.
- Added `last_error` as operational failure evidence.
- Added `completed` and `dead_letter` terminal states.
- Added `ListReadyJobs(ctx)` so the database selects only pending jobs whose retry delay has elapsed.
- Kept retry decisions in the service layer through `CompleteJob(...)` and `FailJob(...)`.
- Added a small simulated executor in the command layer so worker orchestration can exercise success, temporary failure, and permanent failure paths.

### Retry state flow

```mermaid
flowchart TD
    Pending[pending and available_at <= now] --> Running[running]
    Running --> Success{execution result}
    Success -->|success| Completed[completed]
    Success -->|temporary failure and attempts remain| Retry[attempts incremented, last_error set, available_at moved forward]
    Retry --> Pending
    Success -->|temporary failure and attempts exhausted| Dead[dead_letter]
    Success -->|permanent failure| Dead
```

### Worker flow

```mermaid
sequenceDiagram
    participant Poller as Poller
    participant Repo as PostgresRepository
    participant Worker as Worker
    participant Service as Service
    participant Exec as executeJob

    Poller->>Repo: ListReadyJobs(ctx)
    Repo-->>Poller: ready pending jobs
    Poller->>Worker: send job ID on channel
    Worker->>Service: StartJob(ctx, id)
    Worker->>Repo: Get(ctx, id)
    Worker->>Exec: executeJob(ctx, job)
    alt success
        Worker->>Service: CompleteJob(ctx, id)
    else temporary failure
        Worker->>Service: FailJob(ctx, id, temporary error, now)
    else permanent failure
        Worker->>Service: FailJob(ctx, id, permanent error, now)
    end
```

### Why this matters

- Retry state now survives process restarts because it is stored in PostgreSQL.
- Workers are not blocked sleeping for future retries.
- Dead-letter jobs remain queryable for diagnosis.
- Service-layer transition functions prevent worker code from spreading business rules across the command layer.

## Day 18 - Module 3.6: Concurrency Testing

### Goal

Make GoFlow's worker behavior testable under deterministic conditions and validate the code with the Go race detector.

### Design change

- Extracted single-job worker behavior into `processQueuedJob(...)`.
- Injected the job executor so tests can force success, temporary failure, permanent failure, or cancellation.
- Injected the clock so retry scheduling can be asserted without relying on real time.
- Kept the long-running polling loop in `main.go`, but moved the state-transition path into a focused test seam.

### Test seam flow

```mermaid
flowchart TD
    Test[Test case] --> FakeStore[Fake JobStore]
    Test --> FakeExecutor[Injected executor]
    Test --> FakeClock[Injected now function]
    FakeStore --> Process[processQueuedJob]
    FakeExecutor --> Process
    FakeClock --> Process
    Process --> Start[StartJob pending to running]
    Start --> Execute[executor result]
    Execute -->|success| Complete[CompleteJob running to completed]
    Execute -->|temporary failure| Retry[FailJob running to pending with available_at]
    Execute -->|permanent failure| Dead[FailJob running to dead_letter]
    Execute -->|context canceled| Canceled[return context.Canceled]
```

### Race validation flow

```mermaid
flowchart LR
    Tests[go test ./...] --> Normal[Behavior validation]
    Race[go test -race ./...] --> Instrumented[Runtime race instrumentation]
    Instrumented --> Access[Observe executed shared-memory accesses]
    Access --> Result[Report unsafe read/write races]
```

### Why this matters

- The race detector catches data races only on code paths that tests execute.
- Deterministic worker tests prove important job transitions without relying on scheduler timing.
- Data-race freedom and logical correctness are separate concerns.
- The worker loop is still simple, while its important behavior is now testable in isolation.

## Week 3 Checkpoint: Worker Correctness Review

### Goal

Confirm that the worker design is understandable before moving into production-readiness modules.

### Concurrency correctness map

```mermaid
flowchart TD
    Poller[Poller goroutine] -->|send job IDs| JobsChannel[jobs channel]
    JobsChannel --> Worker1[Worker 1]
    JobsChannel --> Worker2[Worker 2]
    JobsChannel --> Worker3[Worker 3]
    Poller -->|read/write| Queued[queued map]
    Worker1 -->|delete after claim attempt| Queued
    Worker2 -->|delete after claim attempt| Queued
    Worker3 -->|delete after claim attempt| Queued
    Mutex[sync.Mutex] --> Queued
    Repo[PostgreSQL] -->|persistent status| Claim[StartJob pending to running]
    Claim --> Idempotency[Idempotent external execution required]
```

### Why this matters

- The channel handles communication between poller and workers.
- The mutex protects shared in-memory queue bookkeeping.
- PostgreSQL status transitions provide persistent claim state.
- Idempotency is still required because external effects and database updates are not one atomic operation.

## Day 19 - Module 4.1: Structured Logging

### Goal

Make GoFlow's important runtime events machine-readable with standard-library JSON structured logs.

### Design change

- Added a JSON `slog.Logger` in the `serve` and `work` command boundaries.
- Added `requestLoggingMiddleware` to log HTTP request completion.
- Added a response-writer wrapper to capture response status codes.
- Added handler-level `job created` logs after successful `POST /v1/jobs`.
- Kept service and repository layers log-free so errors are logged once at the boundary.

### HTTP logging flow

```mermaid
sequenceDiagram
    participant Client as Client
    participant ReqID as requestIDMiddleware
    participant LogMW as requestLoggingMiddleware
    participant Handler as HTTP Handler
    participant Logger as slog JSON logger

    Client->>ReqID: HTTP request
    ReqID->>ReqID: assign request_id
    ReqID->>LogMW: request with context
    LogMW->>Handler: wrapped ResponseWriter
    Handler-->>LogMW: response status/body
    LogMW->>Logger: http request completed with request_id, method, path, status, duration_ms
```

### Worker logging flow

```mermaid
flowchart TD
    Worker[Worker goroutine] --> Processing[log job processing]
    Processing --> Execute[processQueuedJob]
    Execute -->|success| Completed[log job completed]
    Execute -->|failure| Failed[log job failed with error]
    Execute -->|context canceled| Stopping[log worker stopping]
```

### Why this matters

- `request_id` connects response headers to server-side logs.
- `job_id` connects worker events for a single job.
- JSON logs prepare the project for log aggregation and later observability work.
- Logging at boundaries avoids duplicated service/repository logs.

## Day 20 - Module 4.2: Observability

### Goal

Expose basic runtime signals for health and metrics without introducing a production metrics dependency yet.

### Design change

- Added `/health/ready` for dependency readiness.
- Added `PostgresRepository.Ping(ctx)` for PostgreSQL readiness checks.
- Added `/metrics` returning a JSON snapshot from an in-process metrics collector.
- Added mutex protection around metrics state.
- Added HTTP request counters and duration buckets.
- Added worker-process counters and gauges while documenting that they are not visible from the separate server process.

### Health flow

```mermaid
flowchart TD
    Live[/GET /health/live/] --> LiveResult[200 if process can respond]
    Ready[/GET /health/ready/] --> Ping[PostgresRepository.Ping]
    Ping -->|success| ReadyOK[200 ready]
    Ping -->|failure| ReadyFail[503 DATABASE_NOT_READY]
```

### Metrics flow

```mermaid
flowchart TD
    Request[HTTP request] --> Middleware[requestLoggingMiddleware]
    Middleware --> Count[http_requests_total counter]
    Middleware --> Duration[duration bucket observation]
    Create[POST /v1/jobs success] --> Submitted[jobs_submitted_total counter]
    Metrics[/GET /metrics/] --> Snapshot[mutex-protected metrics snapshot]
    Snapshot --> JSON[JSON metrics response]
```

### Worker metrics limitation

```mermaid
flowchart LR
    Serve[serve process] --> MetricsEndpoint[/metrics endpoint/]
    MetricsEndpoint --> HTTPMetrics[HTTP in-process metrics]
    Work[work process] --> WorkerMetrics[worker in-process metrics]
    WorkerMetrics -. not visible to .-> MetricsEndpoint
```

### Why this matters

- Health endpoints support safe orchestration decisions.
- Metrics show operational trends that logs alone cannot summarize.
- Mutex-protected snapshots avoid data races when handlers and goroutines read/write metrics.
- The separate-process limitation is explicit, preventing misleading assumptions about worker visibility.
