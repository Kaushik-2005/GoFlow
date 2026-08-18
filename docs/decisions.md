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
