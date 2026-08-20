# Go Learning Notes

## Table of Contents

- Module 1: Foundations and Idiomatic Go
- Module 2: HTTP, Architecture, Databases, and Testing
- Module 3: Concurrency and Reliable Background Processing
- Module 4: Production Readiness, Performance, Security, and Deployment

## Module 1: Foundations and Idiomatic Go

### Day 1: Module 1.1 - Environment, Packages, Modules, and Tooling

#### Learning objectives

- Distinguish a Go package from a Go module.
- Explain why only `package main` with `func main()` builds an executable.
- Use the basic Go toolchain commands: `go run`, `go build`, `go test`, `go vet`, and `gofmt`.
- Build the first GoFlow CLI scaffold under `cmd/goflow`.

#### Core concepts

- A `package` is a directory of `.go` files compiled together.
- A `module` is a versioned project boundary defined by `go.mod`.
- A runnable Go program comes from a package named `main` that defines `func main()`.
- The Go toolchain owns dependency resolution, builds, tests, formatting, and basic validation.

#### How it works

1. Create the module root with `go mod init`.
2. Put the executable entry point in a `main` package, typically under `cmd/...`.
3. Use `go run` for quick execution, `go build` for a binary, `go test` for tests, `go vet` for suspicious code, and `gofmt` for standard formatting.
4. Keep reusable logic out of `main` so the entry package stays small and focused.

#### Example

```go
package main

import "fmt"

func main() {
	fmt.Println("GoFlow")
}
```

#### Role in our project

- Day 1 created the first GoFlow scaffold.
- The project now has `go.mod` at the root and one executable package under `cmd/goflow`.
- This is the starting point for all later CLI, API, and worker binaries.

#### Why it is designed this way

- Go makes structure explicit early so builds, imports, and dependencies stay predictable.
- `cmd/...` is a common layout because it separates runnable binaries from reusable internal code.
- A thin `main` package is easier to maintain and test around.

#### Alternatives and trade-offs

- Putting everything in one root `main.go` is faster at first, but it becomes messy sooner.
- Using a framework or larger architecture immediately adds more moving parts before the fundamentals are clear.

#### Failure modes

- Confusing package boundaries with module boundaries.
- Putting `go.mod` inside a nested directory like `cmd/goflow`.
- Forgetting that only `package main` plus `func main()` is runnable.
- Treating `gofmt` as optional style polish instead of normal Go practice.

#### Common mistakes

- Saying a module is only a folder boundary instead of a dependency and versioning boundary too.
- Updating help text without updating actual CLI behavior.
- Accessing `os.Args[1]` or `os.Args[2]` without length checks.

#### Production implications

- Small CLI and service entry points keep startup logic understandable.
- Standardized toolchain commands reduce environment-specific behavior.
- Strict compiler and tooling rules catch a lot of issues early.

#### Interview explanation

A Go package is the compile unit; a Go module is the dependency and version boundary defined by `go.mod`. To build an executable, the package must be named `main` and define `func main()`. The normal development loop uses `go run`, `go build`, `go test`, `go vet`, and `gofmt`.

#### Questions for revision

1. What is the difference between a Go package and a Go module?
   Answer: A package is a directory of Go files compiled together. A module is the versioned project boundary defined by `go.mod` and can contain multiple packages.

2. Why can every package not become an executable?
   Answer: Only a package named `main` with `func main()` gives Go a program entry point to run.

3. Why does `go.mod` belong at the project root here?
   Answer: It defines the module boundary for the whole project, including `cmd/goflow` and later packages.

4. Why is `len(os.Args)` checked before indexing command arguments?
   Answer: Because reading a missing element would panic with an index-out-of-range error.

#### Active recall review

1. Question: Why is `go test ./...` useful even before tests exist?
   Answer: It still loads and checks all packages and reports `[no test files]` rather than failing if there are simply no `_test.go` files.

2. Question: Why is `go vet` useful when the code already compiles?
   Answer: It catches suspicious patterns and likely mistakes that still pass the compiler.

#### References

- How to Write Go Code: https://go.dev/doc/code
- Go Modules Reference: https://go.dev/ref/mod
- Effective Go: https://go.dev/doc/effective_go

### Day 2: Module 1.2 - Language Fundamentals

#### Learning objectives

- Explain variables, constants, and zero values.
- Use `:=` versus `=` correctly.
- Understand explicit type conversion.
- Use `if`, `switch`, `for`, and `range` correctly.
- Write small functions, including multiple return values.
- Explain `defer`, scope, shadowing, and `iota`.

#### Core concepts

- `var` declares a variable with an explicit type or with a zero value.
- `:=` declares and initializes a new variable.
- `=` updates an existing variable.
- Go conditions require real boolean expressions.
- `for` is the only loop keyword; `range` is a common iteration form.
- Multiple returns commonly express data plus status information.
- `defer` schedules cleanup at function return.
- `iota` generates incrementing constants inside a `const` block.

#### How it works

1. Use zero values deliberately rather than assuming uninitialized state is unsafe.
2. Use `:=` when creating a new local variable and `=` when updating an existing one.
3. Write direct boolean conditions like `attempts < maxRetries` instead of relying on truthy/falsy rules.
4. Use `range` with `_` when you intentionally ignore a returned value.
5. Use multiple returns for patterns like `value, ok` or `result, err`.
6. Place `defer` close to the resource acquisition it should clean up.

#### Example

```go
var count int
name := "goflow"
const maxRetries = 3

if count == 0 {
	fmt.Println("zero")
}

for _, status := range []string{"pending", "completed"} {
	fmt.Println(status)
}

func divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}
```

#### Role in our project

Module 1.2 produced the first small helper functions in the GoFlow CLI package:

- `maxPriority(priorities []int) int`
- `countByStatus(statuses []string) map[string]int`
- `isValidJobName(name string) bool`
- `retryDelay(attempt int) int`
- `filterCompleted(statuses []string) []string`

These are still focused exercises, but they map directly to later job validation, filtering, counting, and retry logic.

#### Why it is designed this way

- Go favors explicit control flow and explicit state changes over implicit magic.
- Zero values reduce boilerplate but still require careful interpretation.
- Multiple returns and `defer` support simple, safe code patterns without exceptions.

#### Alternatives and trade-offs

- Python-style truthiness is shorter, but Go prefers clarity over convenience in conditions.
- A larger abstraction around these functions would be premature at this stage.

#### Failure modes

- Shadowing a variable by using `:=` when `=` was intended.
- Writing `if count {}` where `count` is an `int`.
- Forgetting `_` in `range` and creating an unused variable compile error.
- Explaining `defer` only as “runs last” without understanding that it is about cleanup at function return.

#### Common mistakes

- Saying `:=` is just assignment.
- Forgetting that missing map keys return the zero value of the map’s value type.
- Thinking the second return value in a function like `divide(a, b)` means “evenly divisible” rather than “operation valid.”

#### Production implications

- Explicit conversions and boolean conditions make code noisier than Python, but much less ambiguous.
- Shadowing bugs are subtle and common in larger Go functions.
- `defer` keeps cleanup next to acquisition and makes early returns safer.

#### Interview explanation

Go language fundamentals are about explicit state and explicit control flow. `var` plus zero values, `:=` versus `=`, direct boolean conditions, `for`/`range`, multiple returns, and `defer` are the patterns behind most ordinary Go code.

#### Questions for revision

1. What is the difference between `var count int` and `count := 0`?
   Answer: `var count int` declares `count` with explicit type `int` and zero value `0`. `count := 0` declares a new variable and infers its type from the right-hand side.

2. Why is `if count {}` invalid when `count` is an `int`?
   Answer: Go requires an actual `bool` expression in conditions and does not use integer truthiness.

3. Why does `counts[status]++` work for a missing key in `map[string]int`?
   Answer: Missing keys return the zero value of `int`, which is `0`, so incrementing starts from `0`.

4. What problem does `defer` solve?
   Answer: It schedules cleanup at function return, which keeps cleanup near resource acquisition and protects early-return paths.

5. How does `iota` work?
   Answer: It starts at `0`, increments per line in a `const` block, and resets in a new `const` block.

#### Active recall review

1. Question: Why do we use `_` in `for _, job := range jobs`?
   Answer: To intentionally ignore the index and avoid an unused variable compile error.

2. Question: Why is `max := priorities[0]` safe only under a specific assumption?
   Answer: It assumes the slice is non-empty; otherwise indexing `priorities[0]` would panic.

3. Question: What does `1 << attempt` mean in the retry delay function?
   Answer: It left-shifts `1` by `attempt` bits, producing powers of two such as `1, 2, 4, 8`.

#### References

- A Tour of Go: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Go Language Specification: https://go.dev/ref/spec

## Module 2: HTTP, Architecture, Databases, and Testing

### Day 7: Module 2.1 - Building HTTP Servers

### Day 8: Module 2.2 - Middleware and API Reliability

### Day 9: Module 2.3 - Project Organization and Architecture

### Day 10: Module 2.4 - PostgreSQL and `database/sql`

### Day 11: Module 2.5 - Testing Fundamentals

## Module 3: Concurrency and Reliable Background Processing

### Day 13: Module 3.1 - Goroutines and Channels

### Day 14: Module 3.2 - Synchronization

### Day 15: Module 3.3 - Worker Pool

### Day 16: Module 3.4 - Context, Cancellation, and Graceful Shutdown

### Day 17: Module 3.5 - Retries and Failure Handling

### Day 18: Module 3.6 - Concurrency Testing

## Module 4: Production Readiness, Performance, Security, and Deployment

### Day 19: Module 4.1 - Structured Logging

### Day 20: Module 4.2 - Observability

### Day 21: Module 4.3 - Profiling and Performance

### Day 22: Module 4.4 - Security

### Day 23: Module 4.5 - Containers and Configuration

### Day 24: Module 4.6 - CI and Engineering Workflow
