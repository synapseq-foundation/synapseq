# SynapSeq Agent Guidelines

This file provides guidance for AI agents and contributors working on SynapSeq.

## Project Overview

SynapSeq is a text-driven audio sequencer for brainwave entrainment, written in Go.

- **Module**: `github.com/synapseq-foundation/synapseq/v4`
- **Go Version**: 1.26.0
- **License**: GPL v2

## Key Commands

```bash
# Run all tests
make test

# Build CLI binary
make build

# Build WASM target for browser
make build-wasm

# Cross-platform builds
make build-macos
make build-linux-amd64
make build-linux-arm64
make build-windows-amd64
make build-windows-arm64

# Clean build artifacts
make clean
```

## Package Map

| Package | Responsibility |
|---------|----------------|
| `cmd/synapseq` | CLI entry, flag parsing, command dispatch, orchestration |
| `cmd/wasm` | WASM browser runtime - loads sequences, renders PCM, streams to JavaScript |
| `core` | Public API - `AppContext`, `LoadedContext` |
| `internal/types` | Domain model - Sequence, Period, Track, Channel, Preset (dependency leaf) |
| `internal/parser` | `.spsq` DSL parsing - lexical and syntactic interpretation |
| `internal/sequence` | Sequence loading, extends/preset resolution, building validated Sequence |
| `internal/audio` | Audio rendering - renderer, sources, effects, sync, wavetable, output |
| `internal/preview` | HTML preview generation |
| `internal/hub` | Remote sequence source - manifest, cache, download |
| `internal/cli` | CLI infrastructure - flags, help, text styling |
| `internal/diag` | Structured diagnostics and parse errors |
| `internal/timeline` | Transition math |
| `internal/preset` | Preset resolution and helpers |
| `internal/resource` | File access abstraction |
| `internal/nameref` | Name validation and reference handling |
| `external` | ffplay and ffmpeg integration |

## Architectural Invariants

These rules must be preserved when making changes:

1. **`core` is the public Go API** - External consumers should use `core` without importing internal packages
2. **`internal/types` must remain a dependency leaf** - Defines domain model, must not import other internal packages
3. **`cmd/synapseq` is the CLI shell** - Handles dispatch and output, not parser or renderer logic
4. **`internal/sequence` owns sequence loading** - Parses DSL via `internal/parser`, builds valid `types.Sequence`
5. **`internal/audio` owns synthesis and rendering** - `core` calls it, does not reimplement audio concerns
6. **Keep audio engine concrete** - Prefer focused collaborators over abstract interfaces

## Git Workflow

- **Development branch**: `development`
- **Production branch**: `main`
- **Feature branches**: `feature/*` (branched from `development`)
- **Bugfix branches**: `bugfix/*` (branched from `development`)
- **Hotfix branches**: `hotfix/*` (branched from `main`)

All PRs should target `development` except critical hotfixes.

## Commit Convention

Use Conventional Commits:

```
feat: add new waveform option
fix: correct parsing bug for noise sequences
docs: update README with usage examples
build: add Makefile for macOS
chore: clean up unused code in parser
```

## Code Conventions

- Follow Go best practices and idioms
- Keep `internal/types` pure - no dependencies on other internal packages
- Keep `core` small and stable - avoid expanding public API
- Prefer clarity over cleverness
- One way to do each task
- Less options, more focus
- Test files: `*_test.go`, use table-driven tests when appropriate
- All tests must pass before submitting PR (`make test`)

## OpenSpec Workflow

This project uses a custom workflow for proposing and implementing changes:

- **Explore**: Use `openspec-explore` skill to investigate problems and clarify requirements
- **Propose**: Use `opsx-propose` command to create a new change proposal
- **Apply**: Use `opsx-apply` command or `openspec-apply-change` skill to implement tasks
- **Archive**: Use `opsx-archive` command or `openspec-archive-change` skill to finalize completed changes

## WASM Target

The `cmd/wasm` package provides a browser-oriented WebAssembly target:

- Does not use CLI flow
- Exposes JavaScript bridge for `.spsq` content input and PCM chunk output
- Does not support Hub workflows, preview HTML, or external tools
- Entry points: `main.go`, `bridge_wasm.go`, `streamservice.go`

## Suggested Reading Order

For new contributors, the fastest way to understand the codebase:

1. `cmd/synapseq/main.go`
2. `cmd/synapseq/dispatch.go`
3. `core/context.go`, `core/sequence.go`, `core/generate.go`
4. `internal/sequence/loadtext.go` and `internal/sequence/parsecontent.go`
5. `internal/parser/*`
6. `internal/audio/renderer.go` and `internal/audio/rendercycle.go`
7. `internal/preview/preview.go`
8. `internal/hub/*`

For detailed architecture, see `docs/ARCHITECTURE.md`. For DSL syntax, see `docs/SYNTAX.md`.