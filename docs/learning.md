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

### Day 3: Module 1.3 - Arrays, Slices, Maps, Strings, and Runes

#### Learning objectives

- Explain the difference between arrays and slices.
- Understand slice length, capacity, and `append` behavior.
- Understand how slices can share backing arrays and when to use `copy`.
- Distinguish nil slices from empty slices.
- Use maps for keyed lookup and existence checks.
- Understand `make` vs `new` and basic value semantics.
- Apply these collection rules to the first in-memory GoFlow store.

#### Core concepts

- Arrays have fixed size, and the size is part of the type.
- Slices are flexible views over underlying array storage.
- `len` is current element count; `cap` is growth room within current backing storage.
- `append` returns the slice you must use afterward, because it may allocate new backing storage.
- Slice assignment shares data; `copy` can create independent storage.
- Nil slices and empty slices both have length `0`, but only nil slices compare equal to `nil`.
- Reading a missing map key returns the zero value of the value type; `value, ok := m[key]` distinguishes missing from present.
- Struct values copy by value, while slice/map values still refer to shared underlying runtime data.

#### How it works

1. Use slices for most collection work and arrays only when fixed size is semantically important.
2. Inspect `len` and `cap` to reason about current size versus growth room.
3. Always keep the result of `append`, because the resulting slice may refer to new storage.
4. Use `copy` when you need an independent slice rather than another view into shared data.
5. Use maps for direct ID-based lookup and the `ok` form when missing keys matter.
6. Use `make` to initialize ready-to-use maps and slices.
7. Remember that retrieving a struct from `map[string]Job` gives a copy that must be written back if changed.

#### Example

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[1:3]

len(s) // 2
cap(s) // 4

s = append(s, 60)
```

And for maps:

```go
counts := map[string]int{"pending": 2}
value, ok := counts["done"]
```

#### Role in our project

Day 3 created the first in-memory job model and store for GoFlow:

- `JobStatus`
- `Job`
- `Store`
- `NewStore()`
- `Create`, `Get`, `List`, `Update`, and `Delete`

The store uses `map[string]Job` because job ID lookup is the main access pattern.

#### Why it is designed this way

- Slices and maps are the normal Go tools for dynamic collections and keyed lookup.
- Day 3 teaches the underlying storage behavior before Day 4 introduces richer struct/method boundaries.
- The in-memory store is intentionally simple so the module stays focused on collection semantics.

#### Alternatives and trade-offs

- A `[]Job` store would preserve insertion order more naturally, but ID lookup/update/delete would require scanning.
- A `map[string]*Job` store would change mutation behavior, but introducing pointer semantics now would blur the Day 3 focus.
- Splitting into separate packages now would be cleaner long-term, but it would add architectural decisions before the roadmap reaches them.

#### Failure modes

- Forgetting to assign the result of `append` back to the slice variable.
- Mutating one slice and unintentionally affecting another because they share backing storage.
- Assuming `len` and `cap` mean the same thing.
- Treating nil and empty slices as identical in every context.
- Assuming map iteration order is stable.
- Mutating a retrieved struct copy and expecting the map entry to update automatically.

#### Common mistakes

- Explaining a slice only as “part of an array” instead of as a flexible view over storage.
- Saying capacity means “how many times it can grow” instead of “how many total elements fit before reallocation.”
- Forgetting that `dst := src` shares slice data while `copy(dst, src)` can create independence.
- Forgetting that `delete` on a missing map key is safe and does nothing.

#### Production implications

- Shared backing arrays create subtle bugs when data is sliced and passed around carelessly.
- Stable ordering is not guaranteed when listing map contents, which affects CLI output and future tests.
- Value semantics matter when reading/updating structs from maps.
- Correct use of `make` prevents nil-map write panics and avoids awkward initialization bugs.

#### Interview explanation

In Go, arrays are fixed-size values whose length is part of the type, while slices are flexible views over underlying array storage. `append`, `copy`, `len`, and `cap` determine how slices grow and share memory. Maps provide direct keyed lookup, but map iteration order is intentionally unstable. These rules matter immediately when building an in-memory job store.

#### Questions for revision

1. What is the difference between an array and a slice?
   Answer: An array is fixed-size storage and its length is part of the type. A slice is a flexible view over elements, usually backed by an array.

2. Why do we usually write `s = append(s, x)` instead of just `append(s, x)`?
   Answer: Because `append` returns the slice you must use afterward, and it may refer to new backing storage if reallocation happened.

3. Why can two slices affect each other's contents?
   Answer: Because both slice values can refer to the same backing array.

4. Why is `make(map[string]Job)` the right choice for the store?
   Answer: `make` returns a ready-to-use map value, so inserts are safe immediately.

5. Why does changing a `Job` returned from `Get(id)` not automatically update the stored job?
   Answer: Because `Get` returns a `Job` value, which is a copy of the struct stored in the map; the updated value must be written back.

#### Active recall review

1. Question: Why is `cap(arr[1:3])` not the same as `len(arr[1:3])`?
   Answer: The length is the number of currently visible elements, while the capacity is how much backing-array room remains from the slice start to the end of the array.

2. Question: Why is `value, ok := myMap[key]` safer than just `value := myMap[key]` when key existence matters?
   Answer: Because `value` alone cannot distinguish “missing key” from “present key with zero value.”

3. Question: Why is `List()` on `map[string]Job` not guaranteed to return jobs in a stable order?
   Answer: Because Go deliberately does not guarantee map iteration order.

#### References

- A Tour of Go: https://go.dev/tour/moretypes/1
- Go Slices: usage and internals: https://go.dev/blog/slices-intro
- Effective Go: https://go.dev/doc/effective_go

### Day 4: Module 1.4 - Structs, Methods, Pointers, and Interfaces

#### Learning objectives

- Explain why structs are better than loose grouped variables for domain modeling.
- Use named-field struct literals.
- Understand methods on structs and when behavior belongs on the type.
- Distinguish value receivers from pointer receivers.
- Understand basic pointer meaning and addressability.
- Understand implicit interface implementation and why interfaces should stay small.
- Apply these ideas to the `Job` model and store orchestration.

#### Core concepts

- A struct groups related fields into one coherent domain value.
- Methods attach behavior to a type, making code read more naturally.
- Value receivers get a copy; pointer receivers can modify the original value.
- `&value` takes an address; `*ptr` refers to the pointed-to value.
- Interfaces describe behavior, not data, and are satisfied implicitly.
- Small interfaces should be defined where they are consumed, not automatically next to every concrete type.
- Constructor functions like `NewStore()` centralize initialization and protect invariants.

#### How it works

1. Model related state as one struct value instead of several unrelated variables.
2. Use named-field literals so each value is explicit and order-independent.
3. Put read-only behavior on the type with methods like `CanRetry()`.
4. Use pointer receivers for mutating methods like `MarkRunning()`.
5. Compose small interfaces only when a consumer genuinely needs multiple behaviors.
6. Use constructor functions to guarantee valid initialization like a non-nil map.

#### Example

```go
func (j Job) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}

func (j *Job) MarkRunning() {
	j.Status = StatusRunning
}
```

And for small consumer-defined interfaces:

```go
type JobGetter interface {
	Get(id string) (Job, bool)
}

type JobUpdater interface {
	Update(job Job) bool
}

type JobReaderWriter interface {
	JobGetter
	JobUpdater
}
```

#### Role in our project

Day 4 attached the first domain behavior directly to `Job` and used that behavior through a small service-style function:

- `CanRetry()`
- `MarkRunning()`
- `startJob(store JobReaderWriter, id string) bool`

This is the first step from “data plus helpers” toward more coherent domain-oriented code.

#### Why it is designed this way

- Structs keep related job state together.
- Methods make important behavior feel like it belongs to the type.
- Pointer receivers make mutation explicit.
- Small interfaces keep consumers honest about what they actually need.
- `NewStore()` ensures `Store` starts valid instead of relying on every caller to remember map initialization.

#### Alternatives and trade-offs

- Free functions like `canRetry(job)` can work, but `job.CanRetry()` usually makes ownership of behavior clearer.
- A `map[string]*Job` store would make some mutations feel more direct, but it would also change the value-semantics lesson from Day 3 too early.
- One large `JobStore` interface is simpler to invent once, but it couples consumers to methods they may not actually need.

#### Failure modes

- Using a value receiver for a mutating method and accidentally changing only a copy.
- Duplicating type definitions across files in the same package.
- Using plain `string` for status in one place and `JobStatus` in another, weakening consistency.
- Defining wide interfaces before there is a real consumer need.
- Skipping `NewStore()` and writing into a nil map.

#### Common mistakes

- Thinking structs are only about saving declaration effort instead of modeling one coherent value.
- Forgetting that `job.MarkRunning()` on a local variable works because the variable is addressable.
- Assuming `Get(id)` returning a `Job` value means store entries update automatically.
- Making interfaces too broad or placing them with the concrete type by default.

#### Production implications

- Domain methods improve readability when state transitions become more numerous.
- Receiver choice affects both correctness and performance characteristics.
- Small interfaces make testing and refactoring easier later.
- Constructor functions help preserve invariants as types gain more setup requirements.

#### Interview explanation

Go uses structs to model domain data and methods to attach behavior directly to types without classes. Value receivers are appropriate for read-only behavior, while pointer receivers are needed for mutation. Interfaces are satisfied implicitly, and idiomatic Go keeps them small and defines them where they are consumed.

#### Questions for revision

1. Why is a struct better than several loose variables for something like a job?
   Answer: It keeps related fields together as one coherent domain value, making the data easier to pass around, store, and reason about.

2. Why is a named-field struct literal better than positional values?
   Answer: Field names make meaning explicit, reduce field-order mistakes, and keep the code readable as the struct grows.

3. Why does `MarkRunning()` need a pointer receiver while `CanRetry()` does not?
   Answer: `MarkRunning()` mutates the actual `Job`, while `CanRetry()` only reads state.

4. Why do we still need `store.Update(job)` after `job.MarkRunning()` in `startJob`?
   Answer: `Get(id)` returns a `Job` value copy; `MarkRunning()` changes that local copy, and `Update` writes it back to the map.

5. Why is a small interface like `JobGetter` better than a large `JobStore` interface for a consumer that only fetches jobs?
   Answer: The consumer depends only on the behavior it actually needs, reducing coupling and keeping contracts clearer.

#### Active recall review

1. Question: Why does `job.MarkRunning()` compile even though `MarkRunning()` has a pointer receiver and `job` was returned as a value?
   Answer: Because `job` is an addressable local variable, so Go can automatically take its address for the method call.

2. Question: Why is `NewStore()` better than forcing every caller to initialize `jobs` manually?
   Answer: It centralizes initialization and guarantees the map is usable, preventing nil-map write panics.

3. Question: When should interfaces usually be defined in Go?
   Answer: Near the consumer that needs the behavior, not automatically near the concrete type.

#### References

- Effective Go: https://go.dev/doc/effective_go
- A Tour of Go methods and interfaces: https://go.dev/tour/methods/1
- Go FAQ on interfaces: https://go.dev/doc/faq#interfaces

### Day 5: Module 1.5 - Error Handling

#### Learning objectives

- Explain what the `error` interface is.
- Understand why Go returns errors instead of using exceptions for normal failures.
- Use `errors.New` for sentinel errors.
- Understand wrapping with `%w`.
- Use `errors.Is` for known error conditions.
- Understand when custom error types plus `errors.As` are better than sentinel errors.
- Apply explicit error design to the store, service, and CLI boundary.

#### Core concepts

- `error` is an interface with an `Error() string` method.
- `nil` means no error occurred.
- In Go, normal failures are usually returned as values that callers inspect explicitly.
- Sentinel errors represent shared recognizable conditions such as `ErrJobNotFound`.
- `fmt.Errorf(... %w ...)` adds context while preserving the original error in the chain.
- `errors.Is` checks for a known error value through wrapped layers.
- `errors.As` extracts a matching custom error type from the chain.
- Error strings are for humans; sentinel values and custom types are for program logic.

#### How it works

1. Create simple shared error conditions with `errors.New` when callers need to recognize them.
2. Return `nil` when operations succeed and a non-nil error when they fail.
3. Add context with `%w` so higher layers can still inspect the underlying error.
4. Use `errors.Is(err, ErrJobNotFound)` for yes/no domain conditions.
5. Introduce a custom error type only when callers need extra structured details beyond the condition itself.
6. Let lower layers return and wrap errors; let higher layers decide user-facing behavior.

#### Example

```go
var ErrJobNotFound = errors.New("job not found")

func (s *Store) Get(id string) (Job, error) {
	job, ok := s.jobs[id]
	if !ok {
		return Job{}, ErrJobNotFound
	}
	return job, nil
}
```

And wrapping:

```go
return fmt.Errorf("start job %q: %w", id, err)
```

And a custom typed error:

```go
type InvalidJobStatusError struct {
	JobID  string
	Status JobStatus
}

func (e InvalidJobStatusError) Error() string {
	return fmt.Sprintf("cannot start job %s from status %s", e.JobID, e.Status)
}
```

#### Role in our project

Day 5 upgraded the GoFlow code from boolean failure signals to explicit error values:

- `ErrJobNotFound`
- `ErrJobAlreadyExists`
- `Get(id) (Job, error)`
- `Update(job) error`
- `Delete(id) error`
- `startJob(...) error`
- `InvalidJobStatusError`
- CLI-side `errors.Is(err, ErrJobNotFound)` handling in `main`

This is the first real error-flow design in the project and sets up later HTTP, file, and database behavior.

#### Why it is designed this way

- Sentinel errors keep common domain conditions stable and recognizable.
- `%w` preserves the underlying cause while still adding useful context.
- Higher-level code can make better user-facing decisions when lower-level code returns structured error information.
- Custom error types are used sparingly so the design stays simple until richer error data is truly necessary.

#### Alternatives and trade-offs

- Returning `bool` is shorter but loses failure meaning.
- Comparing `err.Error()` strings is easy at first but brittle and hostile to refactoring.
- Making every error a custom type is too heavy for simple yes/no conditions like “not found.”

#### Failure modes

- Returning fresh `errors.New("job not found")` values everywhere instead of one shared sentinel.
- Using `%v` instead of `%w` and losing the original error chain.
- Branching on error message text instead of `errors.Is`.
- Using a custom error type when a sentinel error would be simpler.
- Handling user-facing output too low in the stack instead of at the boundary.

#### Common mistakes

- Saying `nil` is “not an error” without stating that it means the error result is absent.
- Treating `%w` as if it were only string formatting.
- Reaching for `errors.As` when `errors.Is` is the right tool.
- Putting boundary-specific messaging inside lower-level functions like `startJob`.

#### Production implications

- Explicit error values make failure paths auditable in services and backends.
- Wrapped errors preserve debugging context without losing machine-readable meaning.
- Sentinel errors help map domain conditions cleanly to CLI, HTTP, and test behavior.
- Custom typed errors support richer policies once state machines and validation become more complex.

#### Interview explanation

Go represents normal failures with returned `error` values rather than exceptions. Sentinel errors such as `ErrJobNotFound` are good for recognizable yes/no conditions and are checked with `errors.Is`. When extra structured details are needed, a custom error type can be returned and extracted with `errors.As`. Wrapping with `%w` adds context while preserving the underlying error chain.

#### Questions for revision

1. Why is `ErrJobNotFound` better than checking for the text `"job not found"`?
   Answer: It gives the program one stable error condition to branch on with `errors.Is`, while message text is brittle and may change or be wrapped.

2. Why is `%w` better than `%v` when wrapping an error you may want to inspect later?
   Answer: `%w` preserves the original error in the chain so `errors.Is` and `errors.As` continue to work; `%v` only formats text.

3. When should you prefer a sentinel error over a custom error type?
   Answer: When callers only need to recognize a simple yes/no condition and do not need extra structured data.

4. When is a custom error type a better fit?
   Answer: When callers need structured details such as job ID or current status to decide what to do next.

5. Why should `errors.Is` checks for user-facing behavior usually happen near a boundary like `main`?
   Answer: Lower layers should return and wrap errors; higher layers should decide how to present or map them.

#### Active recall review

1. Question: Why does `errors.Is(err, ErrJobNotFound)` still work after `startJob(...)` wraps the error with context?
   Answer: Because `startJob` uses `%w`, which preserves the original error in the chain.

2. Question: Why are `ErrJobNotFound` and `ErrJobAlreadyExists` still better as sentinel errors than as custom types right now?
   Answer: They are simple recognizable conditions; callers do not need extra structured data to handle them.

3. Question: Why is `InvalidJobStatusError` a reasonable custom error type?
   Answer: Because the caller may need to know which job and which current status caused the invalid transition.

#### References

- Effective Go: https://go.dev/doc/effective_go
- Go 1.13 errors blog post: https://go.dev/blog/go1.13-errors
- Package errors: https://pkg.go.dev/errors

### Day 6: Module 1.6 - I/O, JSON, Files, and Configuration

#### Concept

This module moved GoFlow from an in-memory-only CLI to a CLI that persists jobs in `jobs.json` across separate runs. The key concepts were `io.Reader` and `io.Writer`, file helpers from `os`, `encoding/json`, JSON struct tags, `Marshal` vs `Unmarshal`, missing-file handling, and using `defer` to clean up resources when working with opened files.

#### Why it matters

A program becomes useful when it can cross the process boundary. Files and JSON are the first simple persistence layer. They teach how Go moves data between memory and the outside world, which directly prepares for HTTP request bodies, configuration loading, logs, and database drivers.

#### Mental model

- `io.Reader` and `io.Writer` represent behavior, not one concrete type.
- A file is one concrete thing that can implement those behaviors.
- `json.Marshal` converts a Go value into JSON bytes.
- `json.Unmarshal` converts JSON bytes into a Go value and needs a pointer to write into.
- A missing persistence file can be a normal startup condition rather than a failure.
- The in-memory `Store` remains the owner of job data; persistence helpers translate between `[]Job` and file bytes.

#### Syntax and APIs

```go
data, err := json.MarshalIndent(jobs, "", "  ")
if err != nil {
	return fmt.Errorf("save jobs: %w", err)
}

if err := os.WriteFile(filename, data, 0644); err != nil {
	return fmt.Errorf("save jobs: %w", err)
}
```

```go
data, err := os.ReadFile(filename)
if err != nil {
	if errors.Is(err, os.ErrNotExist) {
		return []Job{}, nil
	}
	return nil, fmt.Errorf("load jobs: %w", err)
}

var jobs []Job
if err := json.Unmarshal(data, &jobs); err != nil {
	return nil, fmt.Errorf("load jobs: %w", err)
}
```

```go
type Job struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Payload     []byte    `json:"payload"`
	Status      JobStatus `json:"status"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
}
```

#### Idiomatic Go points

- Accept the smallest behavior needed, such as `io.Reader`, instead of over-coupling to `*os.File`.
- Keep file/JSON logic in focused helpers rather than duplicating it inside every CLI branch.
- Use `%w` so the error chain stays inspectable.
- Handle user-facing branching at the boundary in `main`, not deep inside helpers.
- Return an empty slice and `nil` when a missing file is a normal condition.

#### Python comparison

- `defer file.Close()` plays a similar role to `with open(...) as f:`.
- `json.Unmarshal(data, &jobs)` is more explicit than Python's `json.loads(...)` because Go requires a typed destination to decode into.
- Go does not rely on truthiness for error or existence checks; it returns explicit values and errors.

#### Common mistakes

- Forgetting that only exported struct fields are marshaled by `encoding/json`.
- Assuming a missing file must always be treated as an error.
- Forgetting to assign the result of persistence loads back into a usable in-memory structure.
- Using `NewStore()` in read paths like `list` or `get`, which would ignore persisted data.
- Expecting stable job order from a map-backed `List()` result.

#### Project application

Day 6 added real JSON-file persistence to GoFlow:

- `cmd/goflow/persistence.go` now contains `saveJobs(...)` and `loadJobs(...)`.
- `cmd/goflow/store.go` now contains `Store.Save(...)` and `LoadStore(...)`.
- `cmd/goflow/main.go` now loads persisted state for `create`, `list`, `get`, and `process`.
- `process` now demonstrates a real `errors.As(...)` boundary by extracting `InvalidJobStatusError` and printing a specific message.

The current CLI flow now survives separate runs through `jobs.json`.

#### Production implications

- File I/O introduces real failure modes such as missing files, malformed JSON, and write failures.
- Stable JSON field names matter once other tools or services read the same data.
- Non-deterministic map order becomes important once tests or user-facing output expect stable ordering.
- The current ID generation strategy using `len(store.List())+1` is good enough for learning but not safe for production.

#### Interview questions

1. Why would a function prefer `io.Reader` over `*os.File`?
2. Why does `json.Unmarshal(...)` require a pointer?
3. Why are struct tags useful even when exported field names already work?
4. When is a missing file a normal condition instead of an error?
5. Why is Go map iteration order relevant once you persist and display data?

#### References

- Package io: https://pkg.go.dev/io
- Package os: https://pkg.go.dev/os
- Package encoding/json: https://pkg.go.dev/encoding/json
- Effective Go: https://go.dev/doc/effective_go

#### My questions and corrections

- `io.Reader` is preferred over `*os.File` when a function only needs read behavior, because it keeps the code reusable across files, buffers, and network bodies.
- `json:"id"` does not merely “make it a JSON field”; it controls the exact JSON key name used for that struct field.
- Returning `[]Job{}, nil` for a missing file is correct here because “no file yet” means “no jobs saved yet,” not “job not found.”
- The list output order is unstable because `Store.List()` ranges over a Go map.
- `errors.As(...)` is useful when the caller needs structured details like `JobID` and `Status`, not just a yes/no condition.

## Module 2: HTTP, Architecture, Databases, and Testing

### Day 7: Module 2.1 - Building HTTP Servers

#### Concept

This module moved GoFlow from a CLI-only program into its first HTTP API using the standard `net/http` package. The focus was on the server/request/handler mental model, `ServeMux` routing, method-based dispatch, JSON request and response handling, path extraction, status codes, and consistent JSON error responses.

#### Why it matters

HTTP is the normal entry point for backend systems. This layer translates external requests into application operations and translates application results back into protocol-level responses. Understanding `net/http` directly is the foundation for using any higher-level Go web framework later.

#### Mental model

- `http.Server` owns the network-facing server configuration and listens for requests.
- `ServeMux` matches request paths to handlers.
- A handler reads from `*http.Request` and writes to `http.ResponseWriter`.
- The same path can support different operations by dispatching on `r.Method`.
- Handlers should translate between HTTP and application logic, not contain every business rule themselves.
- HTTP status codes give coarse protocol-level meaning, while JSON error codes give stable application-level meaning.

#### Syntax and APIs

```go
mux := http.NewServeMux()
mux.HandleFunc("/health/live", liveHandler)
mux.HandleFunc("/v1/jobs", jobsHandler)
mux.HandleFunc("/v1/jobs/", getJobHandler)

server := &http.Server{
	Addr:    ":8080",
	Handler: mux,
}
```

```go
func liveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

```go
id := strings.TrimPrefix(r.URL.Path, "/v1/jobs/")
if id == "" {
	writeJSONError(w, http.StatusBadRequest, "JOB_ID_REQUIRED", "A job ID is required")
	return
}
```

#### Idiomatic Go points

- Start with `net/http` before adding a framework.
- Use `http.HandlerFunc` for simple handlers and early learning.
- Use `http.Server` instead of only `http.ListenAndServe(...)` so timeouts and server behavior can be configured explicitly later.
- Keep one route path and dispatch by HTTP method when the same resource supports multiple operations.
- Return structured JSON errors instead of ad hoc plain text once an API surface starts to matter.

#### Python comparison

- A Go handler plays a similar role to a Flask or FastAPI route function, but the request/response model is more explicit.
- `ServeMux` is a standard-library router rather than a framework-specific decorator system.
- JSON decoding from `r.Body` is more explicit than many Python frameworks because Go requires a typed destination struct.

#### Common mistakes

- Confusing the collection route `/v1/jobs` with the single-resource route `/v1/jobs/{id}`.
- Forgetting that the mux routes by path, not by HTTP method, so method dispatch must still happen in code.
- Writing the response body before setting headers or status.
- Returning plain-text API errors after the API has already started to stabilize around JSON.
- Assuming a JSON field rename is backward compatible with already-saved persisted data.

#### Project application

Day 7 added the first HTTP surface to GoFlow:

- `GET /health/live`
- `GET /v1/jobs`
- `GET /v1/jobs/{id}`
- `POST /v1/jobs`

The HTTP layer currently sits in `cmd/goflow/http.go` and uses the existing JSON-file store. It supports listing jobs, retrieving one job, creating a new job, and returning consistent JSON API errors.

#### Production implications

- Explicit method checks and status codes make the API predictable for clients.
- Structured JSON errors give frontends and services a stable contract.
- Plain `net/http` path handling is workable, but path extraction and route organization must stay disciplined as the API grows.
- Renaming serialized JSON fields can break existing persisted data unless compatibility or migration is handled deliberately.
- `http.Server` is the right base because timeouts and shutdown behavior become essential in later modules.

#### Interview questions

1. What is the difference between `http.Server`, `ServeMux`, and a handler?
2. Why is `http.HandlerFunc` convenient for simple handlers?
3. Why can `GET /v1/jobs` and `POST /v1/jobs` share the same path but do different things?
4. Why should an API use structured JSON errors instead of plain text?
5. Why is changing a JSON field name a compatibility concern for persisted data?

#### References

- Package net/http: https://pkg.go.dev/net/http
- Effective Go: https://go.dev/doc/effective_go
- Package encoding/json: https://pkg.go.dev/encoding/json

#### My questions and corrections

- `w` writes the response and `r` contains the incoming request; handlers read from `r` and write to `w`.
- `Content-Type: application/json` is not just “because APIs use JSON”; it explicitly tells clients how to parse the body.
- The reason to return JSON from the API and plain text from the CLI is that they serve different consumers: machines versus humans.
- `ErrJobNotFound` maps to HTTP `404`, while validation errors like a missing `type` field map to `400`.
- Keeping both HTTP status codes and JSON error codes matters because status gives the broad category and the JSON error code gives the precise application reason.
- Renaming `maxAttempts` to `max_attempts` caused old saved jobs to decode `MaxAttempts` as `0` because missing JSON fields fall back to zero values by default.


### Day 8: Module 2.2 - Middleware and API Reliability

#### Concept

This module added the first reliability layer around GoFlow's HTTP API. Instead of keeping all HTTP concerns inside handlers, the API now uses middleware for request IDs, panic recovery, request body size limits, and JSON content-type enforcement. The module also refined API error mapping so oversized request bodies return a specific error instead of being lumped into generic JSON decode failures.

#### Why it matters

A working handler is not yet a reliable HTTP service. Real APIs need consistent protections around requests: they should survive panics, reject unsupported inputs early, limit resource consumption, and make each request traceable. Middleware is the standard way to centralize those cross-cutting concerns.

#### Mental model

- Middleware is a function that wraps a handler and returns a new handler.
- Each middleware layer can run before the next handler, after it, or stop the request early.
- Middleware order matters because outer layers can observe and protect everything beneath them.
- Request IDs are useful both for clients through headers and for server-side code through request context.
- Request body limits and content-type enforcement are HTTP boundary rules, so they belong near the transport layer rather than duplicated inside every handler.

#### Syntax and APIs

```go
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeJSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An internal server error occurred")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
```

```go
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
```

```go
func requestBodyLimitMiddleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
```

```go
handler := chain(
	mux,
	requestIDMiddleware,
	recoveryMiddleware,
	requestBodyLimitMiddleware(1<<20),
	requireJSONMiddleware,
)
```

#### Idiomatic Go points

- Middleware is usually explicit handler composition, not hidden framework magic.
- Use `defer` plus `recover()` only for unexpected failures, not for normal control flow.
- Prefer typed context keys over plain strings when storing context values.
- Keep cross-cutting HTTP concerns in middleware instead of duplicating them in handlers.
- Distinguish different client errors precisely when the API contract benefits from it, such as malformed JSON versus oversized bodies.

#### Python comparison

- Middleware here plays a similar role to middleware in Django, FastAPI, or Express-style stacks, but the composition is explicit through `http.Handler` wrapping.
- `context.Context` is the Go-native place for request-scoped values, where Python frameworks often rely on request objects or framework-specific globals.

#### Common mistakes

- Treating middleware order as interchangeable.
- Returning generic invalid-body errors for oversized requests after adding body limits.
- Leaving temporary panic-testing routes in the API after verification.
- Keeping request IDs only in response headers and not making them available to server-side code.
- Using overly strict transport rules without realizing valid variants exist, such as `application/json; charset=utf-8`.

#### Project application

Day 8 added a real middleware stack to GoFlow:

- `requestIDMiddleware`
- `recoveryMiddleware`
- `requestBodyLimitMiddleware(1<<20)`
- `requireJSONMiddleware`
- `chain(...)`

The API now:
- attaches `X-Request-ID` to responses
- stores the same request ID in request context
- recovers panics into structured JSON `500` responses
- rejects oversized request bodies with `REQUEST_BODY_TOO_LARGE`
- rejects non-JSON `POST /v1/jobs` requests with `UNSUPPORTED_MEDIA_TYPE`

#### Production implications

- Recovery middleware prevents one unexpected panic from turning into a process-level outage.
- Request IDs are foundational for later structured logging and tracing.
- Centralized request limits protect server memory and resource usage.
- Content-type checks make the API contract explicit and reduce ambiguous request handling.
- Middleware becomes a key architectural boundary once the project grows and route-specific policies diverge.

#### Interview questions

1. What is middleware in `net/http`, and why is it useful?
2. Why does middleware order matter?
3. Why should request IDs exist in both response headers and request context?
4. Why is `413` more precise than `400` for oversized request bodies?
5. Why should panic recovery be middleware instead of replacing normal error handling?

#### References

- Package net/http: https://pkg.go.dev/net/http
- Package context: https://pkg.go.dev/context
- Effective Go: https://go.dev/doc/effective_go

#### My questions and corrections

- `next` in middleware is the handler being wrapped; middleware can run before it, after it, or stop the request entirely.
- Recovery middleware exists for unexpected failures, while returned errors still handle expected failure paths.
- Request body size limiting belongs in middleware because it is a shared HTTP safety rule, not per-handler business logic.
- Request IDs are more useful than only logging the path because many requests share the same path, but each request ID is unique.
- HTTP status and JSON error code both matter: status gives the broad protocol category, while the JSON code gives the exact application reason.
- `http.StatusRequestEntityTooLarge` is the more accurate status for bodies that exceed the configured limit.


### Day 9: Module 2.3 - Project Organization and Architecture

#### Concept

This module turned GoFlow from a mostly `package main` application into a program with a reusable core package. The main design lesson was that architecture is about responsibility and dependency direction, not just creating more folders. The first real boundary extracted was `internal/goflow`, which now holds core job types, store logic, persistence helpers, and service orchestration.

#### Why it matters

As soon as the project has both CLI and HTTP entry points, `package main` becomes the wrong place to keep shared logic. Future work such as PostgreSQL repositories, tests, and additional binaries becomes harder if reusable logic still lives inside the entry point package. A reusable core package gives outer layers a stable dependency target.

#### Mental model

- `main` should wire dependencies and start the app.
- Handlers should translate HTTP into application operations.
- Middleware should protect the transport boundary.
- Core job logic should live in a reusable package that outer layers can import.
- Dependency direction should point inward toward shared logic, not outward toward `main`.

#### Syntax and APIs

```go
import "goflow/internal/goflow"
```

```go
store, err := goflow.LoadStore(jobsFile)
if err != nil {
	return
}
```

```go
err = goflow.StartJob(store, jobID)
```

#### Idiomatic Go points

- Do not make other packages depend on `package main`.
- Extract packages around real responsibilities, not around arbitrary file counts.
- Keep `main` thin and focused on dependency wiring and startup.
- Export only what outer packages genuinely need, such as `StartJob(...)`.
- Use `internal/` for application-only packages that should not become public library APIs.

#### Python comparison

- This is similar to moving reusable logic out of a Python entry script and into importable modules, but Go makes the boundary stricter because `package main` is explicitly an executable package rather than a normal reusable import target.

#### Common mistakes

- Splitting folders before understanding dependency direction.
- Moving transport code out first while it still depends on logic trapped inside `main`.
- Leaving stale learning helpers in `main.go` after the project has real runtime structure.
- Exporting too much instead of only what the outer layers require.

#### Project application

Day 9 created the first reusable core package:

- `internal/goflow/store.go`
- `internal/goflow/job.go`
- `internal/goflow/service.go`
- `internal/goflow/persistence.go`

The outer app layer in `cmd/goflow` now imports `internal/goflow` for:
- store loading
- job creation
- error sentinels
- job status constants
- service orchestration via `StartJob(...)`

`main.go` was also cleaned up by removing old practice helpers that no longer belong in the entry point.

#### Production implications

- Clearer dependency direction makes storage replacement safer, which matters immediately for the PostgreSQL module.
- Shared logic becomes easier to test independently from the CLI and HTTP layers.
- Keeping `internal/` boundaries explicit discourages accidental public-library style coupling too early.
- A thinner `main` makes future binaries like `cmd/api` and `cmd/worker` easier to introduce.

#### Interview questions

1. Why should reusable application logic not stay in `package main`?
2. Why was `internal/goflow` a safer first extraction than an `internal/httpapi` package?
3. What should `main` own in a healthy Go application?
4. Why does package extraction need to respect dependency direction?
5. Why are leftover learning helpers in `main.go` an architectural smell once the app has real boundaries?

#### References

- How to Write Go Code: https://go.dev/doc/code
- Effective Go: https://go.dev/doc/effective_go
- Package layout and `internal`: https://go.dev/doc/modules/layout

#### My questions and corrections

- The key reason to extract a core package first is not that the current layout is temporary; it is that outer layers must not depend on `package main`.
- `service.go` is the clearest service-layer file because it orchestrates state transitions instead of doing direct transport or persistence work.
- `serve` still belongs in `main.go` for now because it is startup wiring; old practice helpers were the better cleanup target.
- `main.go` and `http.go` should both be able to import the reusable core package after the refactor.


### Day 10: Module 2.4 - PostgreSQL and `database/sql`

#### Concept

`database/sql` is Go's standard abstraction for relational database access. The key mental model is that your application code should talk to a long-lived `*sql.DB` handle, while repository methods perform context-aware queries through `QueryRowContext(...)`, `QueryContext(...)`, and `ExecContext(...)`.

In GoFlow, Module 2.4 replaced the old JSON-file runtime path with a PostgreSQL-backed repository and startup flow.

#### Why it matters

- File persistence is not a safe long-term storage boundary for a real API.
- A repository boundary limits how far storage changes spread through the codebase.
- `context.Context` lets database calls respect cancellation and timeouts.
- Parameterized SQL prevents SQL injection.
- `sql.DB` should be created once and reused because it manages a connection pool.

#### Mental model

Think of `*sql.DB` as a managed handle plus connection pool, not one permanently-open connection.

- `sql.Open(...)` creates the handle.
- `PingContext(...)` proves the database is reachable.
- `ExecContext(...)` is for statements like `INSERT`, `UPDATE`, and `DELETE`.
- `QueryRowContext(...)` is for one expected row.
- `QueryContext(...)` is for many rows and returns `Rows`, which must be closed.
- `Scan(...)` copies SQL columns into Go variables in positional order.

#### Syntax

```go
row := db.QueryRowContext(ctx, `
	SELECT id, job_type, status
	FROM jobs
	WHERE id = $1
`, id)

var job Job
var status string
err := row.Scan(&job.ID, &job.Type, &status)
if err != nil {
	return Job{}, err
}

job.Status = JobStatus(status)
```

```go
rows, err := db.QueryContext(ctx, `
	SELECT id, job_type, status
	FROM jobs
	ORDER BY created_at ASC, id ASC
`)
if err != nil {
	return nil, err
}
defer rows.Close()
```

```go
_, err := db.ExecContext(ctx, `
	UPDATE jobs
	SET status = $2, updated_at = NOW()
	WHERE id = $1
`, id, StatusRunning)
```

#### Idiomatic Go points

- Pass `context.Context` as the first parameter for DB-facing operations.
- Reuse one `*sql.DB` instead of opening a new handle for each query.
- `PingContext(...)` is the explicit startup connectivity check.
- Always `defer rows.Close()` after a successful `QueryContext(...)`.
- Check `rows.Err()` after iteration.
- Wrap errors with operation context.
- Use parameters like `$1` and `$2` instead of interpolating user input.

#### Python comparison

- This is similar to a repository over `psycopg`, but Go makes cancellation, cleanup, and returned errors more explicit.
- In Python, exceptions usually propagate automatically; in Go, repository methods return errors directly.
- In Python, context managers often hide cleanup; in Go, `defer rows.Close()` makes it visible.

#### Common mistakes

- Assuming `sql.Open(...)` already validated connectivity.
- Building SQL with string concatenation.
- Forgetting `rows.Close()`.
- Forgetting `rows.Err()`.
- Mismatching `Scan(...)` destinations with selected column order.
- Letting handlers know too much about SQL instead of depending on a repository boundary.

#### Project application

Module 2.4 now uses PostgreSQL as the active runtime target.

Files added or updated:

- `cmd/goflow/database.go`
- `cmd/goflow/main.go`
- `cmd/goflow/http.go`
- `internal/goflow/service.go`
- `internal/goflow/postgres_repository.go`
- `migrations/001_create_jobs.sql`

The important runtime changes are:

- `DATABASE_URL` is now required.
- `openPostgresStore(...)` opens `*sql.DB`, configures pooling, runs `PingContext(...)`, and applies the initial schema.
- CLI commands now call PostgreSQL through the repository instead of loading `jobs.json`.
- HTTP handlers now use a shared store dependency and request context instead of reopening file storage per request.
- `StartJob(...)` now accepts `context.Context` and a shared `JobStore` dependency.
- Job IDs are generated with random bytes instead of `len(list)+1`.
- PostgreSQL duplicate-key errors are mapped to `ErrJobAlreadyExists`.

The current store contract is:

```go
type JobStore interface {
	Create(ctx context.Context, job Job) error
	Get(ctx context.Context, id string) (Job, error)
	List(ctx context.Context) ([]Job, error)
	Update(ctx context.Context, job Job) error
	Delete(ctx context.Context, id string) error
}
```

#### Production implications

- The app now targets a real database boundary instead of local file persistence.
- HTTP request context can now flow into database calls.
- Pool configuration is explicit instead of hidden.
- The storage swap happened without pushing SQL details into handlers.
- Runtime validation is still pending because this environment does not currently have a live PostgreSQL instance.

#### Interview questions

1. Why should `*sql.DB` be reused instead of reopened for each operation?
2. Why does `PingContext(...)` matter after `sql.Open(...)`?
3. Why is `QueryRowContext(...)` better than `QueryContext(...)` for `Get(id)`?
4. Why must `rows.Close()` and `rows.Err()` both be handled?
5. Why should repository methods accept `context.Context`?
6. Why are parameterized queries safer than string interpolation?
7. Why is replacing file storage easier when the service depends on a repository boundary?

#### References

- Executing SQL statements: https://go.dev/doc/database/change-data
- Querying for data: https://go.dev/doc/database/querying
- Canceling in-progress operations: https://go.dev/doc/database/cancel-operations
- Managing connections: https://go.dev/doc/database/manage-connections
- `database/sql` package docs: https://pkg.go.dev/database/sql

#### My questions and corrections

- `*sql.DB` should be created once and reused; it is a pool manager, not one connection.
- `PingContext(...)` is the real connectivity check.
- `Get(id)` should use `QueryRowContext(...)`, while `List()` should use `QueryContext(...)`.
- `rows` is a live cursor and must be closed.
- SQL injection happens when untrusted input is merged into query text instead of passed as parameters.
- Replacing storage cleanly is easier when the service depends on a store interface instead of file helpers.

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









