# Agent Instructions

## Scope

fmgo is a focused Go library for Apple's native `fm` command. Do not add a CLI,
Swift or cgo bridge, private framework binding, alternate model provider, or
generic AI abstraction.

## Repository Structure

- The repository root is the public `fmgo` package.
- `internal/` contains only reusable implementation helpers.
- `cmd/` contains integration examples, not a distributed CLI.
- `tests/` contains public behavior tests.

## Implementation Rules

- Keep code in clear, idiomatic English.
- Give each type, function, and package one responsibility.
- Prefer the simplest direct implementation that satisfies the current task.
- Keep public APIs small and preserve existing contracts unless the task requires
  a change.
- Use the standard library unless a dependency provides substantial current
  value.
- Accept `context.Context` as the first argument of blocking operations.
- Return errors; do not introduce panics into library flows.
- Pass values to `fm` as process arguments. Never use a shell, invoke `sudo`,
  accept Apple terms, log prompts, or weaken native guardrails.
- Add the Apache-2.0 SPDX header to every new Go file and GoDoc comments to
  exported identifiers.

## fmgo Invariants

- `New` returns `ErrUnsupportedPlatform` or `ErrUnsupportedVersion` for an
  incompatible machine.
- `fm` availability is operational: return `ErrFMNotFound` instead of panicking
  when the executable cannot be found.
- Process lifecycle code must respect cancellation, capture diagnostics, and
  reap child processes.
- Keep platform checks in `internal/platform` and process mechanics in
  `internal/process`.

## Tests and Documentation

- Use `WithExecutable` and the existing fake executable pattern for process
  tests. Native tests run only with `FMGO_INTEGRATION=1` on a prepared machine.
- Run `go fmt ./...`, `go vet ./...`, and `go test ./...` after Go changes.
- Use `make integration` only on macOS 27+ Apple Silicon with the fm license
  accepted and the system model available.
- Update README, GoDoc, examples, and tests whenever a public contract changes.

## Workflow

- Inspect the worktree before editing and preserve unrelated user changes.
- Keep a branch and pull request limited to one feature, bug fix, refactor,
  test change, or documentation change.
- When asked to commit, use Conventional Commits such as `feat:`, `fix:`,
  `docs:`, `test:`, `refactor:`, or `chore:`.
- Do not create commits or push changes unless explicitly requested.
