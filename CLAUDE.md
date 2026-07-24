# SynapSeq Agent Guidelines

This file provides guidance for AI agents and contributors working on SynapSeq.

## Project Overview

SynapSeq is a text-driven audio sequencer for brainwave entrainment, written in Go.

- **Module**: `github.com/synapseq-foundation/synapseq/v4`
- **Go Version**: 1.26.0
- **License**: GPL v3 or later

## Key Commands

```bash
# Run all tests
make test

# Build CLI binary
make build

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
| `core` | Public API - `AppContext`, `LoadedContext` |
| `internal/types` | Domain model - Sequence, Period, Track, Channel, Preset (dependency leaf) |
| `internal/parser` | `.spsq` DSL parsing - lexical and syntactic interpretation |
| `internal/sequence` | Sequence loading, extends/preset resolution, building validated Sequence |
| `internal/audio` | Audio rendering - renderer, sources, effects, sync, wavetable, output |
| `internal/remote` | Remote sequence source - index, cache, download |
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

## SynapSeq Agent Skills

The repository ships three complementary skills under `.agents/skills/`:

- **`create-spsq`** creates and validates new `.spsq` files, either from scratch or from existing read-only references. It never modifies an existing `.spsq` or `.spsc` in place.
- **`explain-spsq`** teaches SPSQ syntax and explains existing `.spsq` or `.spsc` files without creating or modifying them.
- **`review-spsq`** validates and audits existing `.spsq` files, producing a read-only technical and compositional report.

Keep their responsibilities distinct. Route creation and new versions to `create-spsq`, didactic questions to `explain-spsq`, and critical analysis to `review-spsq`.

### Mandatory synchronization for SPSQ language changes

Any change that affects accepted `.spsq` or `.spsc` syntax, parser behavior, builder semantics, validation rules, tracks, effects, options, inheritance, transitions, timeline behavior, or diagnostics must be reflected in the skills in the same change. A parser or language change is incomplete until the bundled skill references and workflows agree with the implementation.

At minimum, review and update:

1. `docs/SYNTAX.md` for grammar, placement, ranges, and validation rules;
2. `docs/HOW_IT_WORKS.md` when audible or temporal behavior changes;
3. `.agents/skills/create-spsq/references/spsq-language.md`;
4. `.agents/skills/explain-spsq/references/spsq-language.md`;
5. `.agents/skills/review-spsq/references/review-checklist.md`;
6. each affected `SKILL.md`, plus review sound-design or report references when the change alters analysis or output guidance.

Also update `README.md` when the language change affects user-facing examples or capabilities, and update `agents/openai.yaml` when a skill's routing description changes.

Before completing an SPSQ language change:

- add or update parser and sequence-builder tests;
- validate representative valid and invalid `.spsq` fixtures with `synapseq -test`;
- run `make test`;
- run the skill-creator `quick_validate.py` check for all three skill directories;
- verify that no skill documents syntax, ranges, effects, or behavior that the current implementation no longer supports.

## Git Workflow

- **Default branch**: `main`
- **Feature branches**: `feature/*` (branched from `main`)
- **Bugfix branches**: `bugfix/*` (branched from `main`)
- **Hotfix branches**: `hotfix/*` (branched from `main`)

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

## Suggested Reading Order

For new contributors, the fastest way to understand the codebase:

1. `cmd/synapseq/main.go`
2. `cmd/synapseq/dispatch.go`
3. `core/context.go`, `core/sequence.go`, `core/generate.go`
4. `internal/sequence/loadtext.go` and `internal/sequence/parsecontent.go`
5. `internal/parser/*`
6. `internal/audio/renderer.go` and `internal/audio/rendercycle.go`
7. `internal/remote/*`

For detailed architecture, see `docs/ARCHITECTURE.md`. For DSL syntax, see `docs/SYNTAX.md`.
