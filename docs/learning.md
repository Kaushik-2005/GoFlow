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
