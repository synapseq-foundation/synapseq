# Contributing to fmgo

Thanks for contributing to fmgo. Keep changes small, focused, and aligned with
the library's purpose: a Go interface to Apple's native `fm` command.

## Scope

fmgo is a focused Go library. Do not add a CLI, provider abstractions, Swift or
cgo bridges, private framework bindings, or unrelated AI features.

## Branches

Use one branch for one responsibility. Separate features from bug fixes and
avoid mixing unrelated cleanup into either.

Use descriptive branch names:

```text
feature/streaming-timeout
fix/server-startup-error
docs/usage-guide
test/schema-recursion
refactor/process-cleanup
```

## Code Principles

- Give each type, function, and package a single responsibility.
- Prefer the simplest implementation that satisfies the current requirement.
- Write idiomatic, self-explanatory Go; use clear names and straightforward
  control flow instead of indirection.
- Add documentation when a public contract, invariant, or non-obvious decision
  needs explanation.
- Avoid speculative abstractions, unnecessary dependencies, and unrelated
  refactors.

## Validation

Run these commands before opening a pull request:

```bash
make check
```

The test suite uses a fake `fm` executable and does not require a Foundation
Model download.

Native integration tests require macOS 27+ on Apple Silicon, the Foundation
Models CLI license accepted, and the system model available. Run them only on
a prepared machine:

```bash
make integration
```

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/). Keep each
commit focused on one logical change.

```text
feat: add structured response helper
fix: report server startup failures
docs: clarify streaming lifecycle
test: cover unavailable model errors
refactor: simplify process cleanup
chore: update development metadata
```

## Pull Requests

Describe the problem and the chosen solution, include relevant tests, and keep
each pull request limited to one feature, bug fix, or documentation change.
