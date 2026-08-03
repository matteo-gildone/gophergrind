### Interaction Rules

* The reader is learning Go. Be honest and push back when something is wrong,
  un-idiomatic, or a bad idea — do not soften it with filler praise.
* Explain the reasoning behind corrections so they teach, not just fix.
* Ask clarifying questions if input is unclear.
* Explain why and suggest alternatives if a task is not feasible.
* When disagreeing or proposing a different approach, back it with a concrete
  standard-library code example, not just an assertion.
* Use structured, readable formatting (headings, lists, code blocks).
* Follow instructions closely and explain clearly what you have done.
* Don't modify code unrelated to the current task.
* Match the style of the code you are touching.

### Idiomatic Go (primary authority)

Follow Effective Go, the Go Proverbs, and the Google Go Style Guide. When any
rule below conflicts with these, these win.

* Code must be `gofmt`-clean. Never hand-format.
* Code must pass `go vet` and GoLand's integrated inspections with no new
  warnings.
* Handle every error explicitly. Wrap with context: `fmt.Errorf("doing X: %w", err)`.
  Do not discard errors with `_` unless there is a stated reason.
* Return errors; do not `panic` in library code. Reserve `panic` for truly
  unrecoverable programmer errors.
* Accept interfaces, return concrete types. Keep interfaces small and define
  them at the consumer, not the producer.
* Pass `context.Context` as the first parameter for anything that does I/O,
  blocks, or may need cancellation.
* Make zero values useful where practical.
* Prefer clear, explicit code over cleverness. Avoid reflection and
  `any`/`interface{}` unless there is no reasonable alternative.

### Dependencies

* Standard library only. Do not add external dependencies.
* If a task seems to require a third-party package, stop and ask first.

### Coding Standards

* Write meaningful tests with assertions for all code.
* Prefer table-driven tests.
* Avoid duplicated test assertions.
* Keep functions short and single-purpose, but do not fragment code to hit a
  line count.
* Simple design, in priority order:
  1. Code works (tests pass).
  2. Reveals intent (clear names, obvious control flow).
  3. No needless duplication.
  4. Minimal elements.
* State and mutation:
  * Prefer pure functions where it reads naturally.
  * Avoid unnecessary shared mutable state; guard shared state that must exist.
  * Idiomatic pointer mutation is fine — do not force immutability against
    Go's grain.

### Concurrency

* Code must pass `go test -race`.
* Do not leak goroutines; ensure every goroutine has a clear exit path.
* Prefer channels or `sync` primitives deliberately; do not add concurrency
  without a reason.

### Architecture

* Organise packages by capability, not by technical layer. Package names should
  say what they provide (`auth`, `billing`), not `models`, `controllers`, or a
  `utils` grab-bag.
* One clear responsibility per package.
* Low coupling between packages. No import cycles.
* Avoid premature abstraction: no interface until there is a second
  implementation or a real test seam.

### Testing

* Use the standard `testing` package only. No external test/assertion libraries.
* Run tests with `-race`.
* Use black-box `_test` packages for public API tests where practical.
* Use `t.Parallel()` when tests are independent.
* Test behavior, not coverage numbers. Add tests for each bug fixed.

### Workflow

* For non-trivial changes, read `spec.md` before coding and log changes to it
  after. Small, obvious changes do not require a spec update.
* Write and pass tests before finalising.
* Keep `README.md` with setup/run info current.
* Store docs and specs in Markdown.

### Safe Practices

* Do not change test assertions during refactoring.
* Do not skip or delete failing tests to make a build pass.
* Do not invent unknown APIs; ask if you are unsure.

### Goal

Produce consistent, idiomatic, testable, and maintainable Go. Stick to the
rules — no shortcuts.