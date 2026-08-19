# Go Learning Notes

## Week 1: Foundations and Idiomatic Go

### Module 1.1: Environment, Packages, and Modules

#### Concept

Go organizes code at two levels:

- A `package` is a directory of `.go` files compiled together.
- A `module` is a versioned collection of packages defined by `go.mod`.

To build an executable, a package must be named `main` and must define `func main()`.

#### Why it matters

This is the first structural idea in Go. If you do not understand package boundaries, module roots, and the toolchain, the rest of the language feels arbitrary. Real Go services are assembled from many packages inside one module, then built into one or more binaries.

#### Mental model

Think in layers:

- File: one source file
- Package: one compile unit made from files in the same directory
- Module: one dependency and version boundary made from related packages
- Binary: the final executable produced from a `main` package

Python comparison:

- A rough Python analogy is that a package is like an importable module/package directory, while a Go module is closer to the whole distributable project boundary managed by `pyproject.toml` plus dependency lock behavior.
- The important difference is that Go's module system is built into the language toolchain, not bolted on through separate packaging conventions.

#### Syntax

Minimal executable:

```go
package main

import "fmt"

func main() {
	fmt.Println("GoFlow")
}
```

Key points:

- `package main` marks this package as executable.
- `func main()` is the program entry point.
- Imports must be used, or compilation fails.

#### Idiomatic Go

- Keep `main` packages thin. Their job is wiring, startup, and shutdown.
- Put reusable logic in non-`main` packages.
- Name packages by responsibility, not by type of file.
- Let `go.mod` define the module root early.

#### Python comparison

Python often lets structure remain loose for a while. Go pushes structure earlier:

- unused imports are errors
- circular imports are not allowed
- project layout affects buildability more directly
- dependency management is standardized through `go.mod` and `go.sum`

This strictness buys faster builds, clearer ownership, and fewer packaging surprises in production.

#### Common mistakes

- Confusing a package with a module
- Assuming every package can be run directly
- Putting unrelated responsibilities into `main`
- Treating `go.mod` as optional
- Expecting Go to allow unused variables or imports during early experimentation

#### Project application

GoFlow will eventually have multiple binaries, likely:

- `cmd/api`
- `cmd/worker`

Each binary will be its own `main` package, while shared logic will live under `internal/...`.

For the first increment, we started smaller with one CLI entry point under `cmd/goflow`.

#### Interview questions

- What is the difference between a Go package and a Go module?
- Why are `package main` and `func main()` special?
- Why does Go enforce unused import and variable checks?
- Why do many Go repos place binaries under `cmd/`?

#### Official references

- Required: [How to Write Go Code](https://go.dev/doc/code) — understand packages, modules, workspace layout, and the build model
- Required: [Go Modules Reference](https://go.dev/ref/mod) — understand what `go.mod` and `go.sum` are for
- Optional: [Effective Go](https://go.dev/doc/effective_go) — package naming and layout conventions

#### My questions and corrections

- Correction: a module is not only the project boundary; it is also the dependency and versioning boundary defined by `go.mod`.
- Confirmed understanding: only a package named `main` with `func main()` becomes an executable entry point.
- Confirmed understanding: `go.mod` belongs at the module root, not inside `cmd/goflow`.
- Correction from exercise review: changing help text is not the same as implementing command behavior; the `switch` cases must be updated too.
- Confirmed understanding: `len(os.Args) < 2` protects access to `os.Args[1]`, while `len(os.Args) < 3` protects access to command-specific arguments like `os.Args[2]`.

#### Small example

Completed CLI shape:

```go
package main

import (
	"fmt"
	"os"
)

func requireArg(args []string, index int, message string) (string, bool) {
	if len(args) <= index {
		fmt.Println(message)
		return "", false
	}
	return args[index], true
}

func printHelp() {
	fmt.Println("Usage: goflow <command>")
	fmt.Println("Commands: list, create, get, process")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "list":
		fmt.Println("listing jobs")
	case "create":
		jobType, ok := requireArg(os.Args, 2, "missing job type")
		if !ok {
			return
		}
		fmt.Printf("creating job of type: %s\n", jobType)
	case "get":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}
		fmt.Printf("getting job: %s\n", jobID)
	case "process":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}
		fmt.Printf("processing job: %s\n", jobID)
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
	}
}
```

#### Production implications

- `main()` should coordinate startup and dispatch, not accumulate all business logic.
- Exact CLI output matters once tests are added, so string precision is part of correctness.
- Small helper functions are useful when they remove real repetition without introducing unnecessary abstraction.

#### Project application

- Module 1.1 produced the first runnable GoFlow CLI scaffold under `cmd/goflow`.
- The project now has a module root with `go.mod` and one executable package.
- This layout is the basis for later expansion into richer CLI behavior and, later in the roadmap, separate binaries.

### Module 1.2: Language Fundamentals

#### Concept

Module 1.2 covers the core control-flow and state-building tools in Go: variables, constants, zero values, type conversion, branching, loops, functions, multiple return values, `defer`, scope, shadowing, and `iota`.

#### Why it matters

These are the mechanics behind almost every line of ordinary Go code. Before building data structures and storage layers, you need to be comfortable reading and writing small functions with predictable state and control flow.

#### Mental model

- `var` declares a variable with an explicit type or inferred zero value.
- `:=` declares and initializes a new variable.
- `=` updates an existing variable.
- Go conditions must be real `bool` expressions, not truthy/falsy values.
- `for` is the only loop keyword; `range` is a common iteration form built on top of it.
- Multiple returns let a function return data plus status information.
- `defer` schedules cleanup for function exit.
- Inner scopes can shadow outer variables if you accidentally use `:=`.

#### Syntax

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

#### Idiomatic Go

- Use zero values intentionally; they are part of the design, not an accident.
- Use `:=` for short local declarations when it improves clarity.
- Use `=` when you mean update, especially to avoid shadowing.
- Keep branching explicit; write boolean expressions directly.
- Use `_` to ignore a `range` value you do not need.
- Prefer small functions that do one clear thing.

#### Python comparison

- Python uses truthy/falsy rules; Go requires a real `bool` in conditions.
- Python allows unused variables freely; Go rejects them.
- Python has separate `for` and `while`; Go uses `for` for both roles.
- Python has no direct equivalent to Go's `iota` constant counter.
- Python cleanup often depends on context managers; Go often uses `defer`.

#### Common mistakes

- Saying `:=` is just assignment instead of new-variable declaration.
- Forgetting that declared variables already have zero values.
- Writing `if count {}` when `count` is an `int`.
- Forgetting `_` and leaving an unused index variable in a `range` loop.
- Using `:=` in an inner scope when you meant to update an outer variable.
- Explaining `defer` only as “runs last” instead of “runs at function return for cleanup.”

#### Project application

Module 1.2 concepts were applied by adding small helper functions to the GoFlow CLI package:

- `maxPriority(priorities []int) int`
- `countByStatus(statuses []string) map[string]int`
- `isValidJobName(name string) bool`
- `retryDelay(attempt int) int`
- `filterCompleted(statuses []string) []string`

These are still small exercises, but they map directly to later GoFlow responsibilities like retries, status counting, filtering, and validation.

#### Interview questions

- What is the difference between `var count int` and `count := 0`?
- Why is `if count {}` invalid when `count` is an `int`?
- What does `value, ok := map[key]` represent?
- What problem does `defer` solve?
- What is variable shadowing in Go?
- How does `iota` work inside a `const` block?

#### Official references

- Required: [A Tour of Go](https://go.dev/tour/) — review variables, flow control, methods of iteration, and functions
- Required: [Effective Go](https://go.dev/doc/effective_go) — review declarations, control structures, and idiomatic style
- Optional: [Go Language Specification](https://go.dev/ref/spec) — reference for declarations, statements, and constants

#### My questions and corrections

- Correction: `:=` declares and initializes a new variable; `=` updates an existing one.
- Confirmed understanding: zero value of `string` is `""`, of `bool` is `false`, and of `int` is `0`.
- Confirmed understanding: `counts[status]++` works on a missing map key because the zero value of `int` is `0`.
- Correction: the second return value in `divide(a, b)` represented whether the operation was valid, not whether the division was even.
- Correction: `defer` is primarily about cleanup at function return, not just “running last.”
- Confirmed understanding: `_` is used in `range` to intentionally ignore an unused value.
- Confirmed understanding: `iota` starts at `0`, increments per line in a `const` block, and resets in a new block.

#### Small examples

```go
func maxPriority(priorities []int) int {
	max := priorities[0]
	for _, priority := range priorities {
		if priority > max {
			max = priority
		}
	}
	return max
}

func countByStatus(statuses []string) map[string]int {
	counts := make(map[string]int)
	for _, status := range statuses {
		counts[status]++
	}
	return counts
}

func retryDelay(attempt int) int {
	return 1 << attempt
}
```

#### Production implications

- Zero values simplify initialization but must still be interpreted carefully in domain logic.
- Shadowing bugs are easy to write and often subtle in larger functions.
- `defer` keeps cleanup next to acquisition and makes early-return code safer.
- Explicit boolean conditions and explicit type conversions make code noisier than Python, but much less ambiguous.

#### Project application

- The CLI package now contains several small language-fundamentals functions alongside the command dispatcher.
- These are stepping stones toward the Week 1 CLI deliverable and later job-store logic.
- The next module should separate storage/data behavior more intentionally rather than growing everything inside `main.go` forever.
