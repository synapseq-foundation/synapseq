# SynapSeq Agent Guidelines

This file provides guidance for AI agents and contributors working on SynapSeq.

## Project Overview

SynapSeq is a text-driven audio sequencer for brainwave entrainment, written in Go.

- **Module**: `github.com/synapseq-foundation/synapseq/v4`
- **Go Version**: 1.27.0
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
| `core` | Public runtime API - `AppContext`, `LoadedContext` |
| `spsq` | Public fluent API for constructing and validating `.spsq` content |
| `sbg` | Public SBaGen conversion API backed by the `spsq` builder |
| `internal/types` | Domain model - Sequence, Period, Track, Channel, Preset (dependency leaf) |
| `internal/parser` | `.spsq` DSL parsing - lexical and syntactic interpretation |
| `internal/sequence` | Sequence loading, extends/preset resolution, building validated Sequence |
| `internal/audio` | Rendering root - render plans, signal state, mixing, and orchestration |
| `internal/audio/audiosource` | External WAV/MP3 loading, decoding, caching, resampling, and playback cursors |
| `internal/audio/{ambiance,music,effects,sources,sync,wavetable,output,pcm,status}` | Focused audio rendering collaborators |
| `internal/remote` | Remote sequence source - index, cache, download |
| `internal/cli` | CLI infrastructure - flags, help, text styling |
| `internal/textstyle` | Terminal text styling primitives |
| `internal/diag` | Structured diagnostics and parse errors |
| `internal/timeline` | Transition math |
| `internal/preset` | Preset resolution and helpers |
| `internal/resource` | File access abstraction |
| `internal/nameref` | Name validation and reference handling |
| `external` | ffplay and ffmpeg integration |

## Architectural Invariants

These rules must be preserved when making changes:

1. **`core` is the public Go API** - External consumers should use `core` without importing internal packages.
2. **`internal/types` must remain a dependency leaf** - It defines the domain model and must not import other internal packages.
3. **`cmd/synapseq` is the CLI shell** - It handles dispatch and output, not parser or renderer logic.
4. **`internal/sequence` owns sequence loading** - It parses DSL via `internal/parser` and builds valid `types.Sequence` values.
5. **`internal/audio` owns synthesis and rendering** - `core` calls it and does not reimplement audio concerns.
6. **Keep the audio engine concrete** - Prefer focused collaborators over abstract interfaces.
7. **`spsq` remains a builder** - It generates text and returns `core.LoadedContext`; parsing, validation, and rendering stay in their owning packages.
8. **`sbg` converts through `spsq`** - It must not create a parallel loading or rendering pipeline.

## External Audio

- `@ambiance` and `@music` are lazy-loaded only when their tracks first become active. Do not reintroduce eager renderer-wide loading.
- Loaded external assets are cached for the render. Ambiance loops; music has finite playback.
- Remote source failures, including non-2xx HTTP responses, invalid MIME types, decode failures, and later read failures, must propagate as render errors. Never silently replace a failed source with silence.
- Keep source loading in `internal/resource` and external-audio decoding/playback in `internal/audio/audiosource`.

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

Also update `README.md` when the language change affects user-facing examples or capabilities. When a skill's routing description changes, update its front matter and every repository manifest that references it.

Before completing an SPSQ language change:

- add or update parser and sequence-builder tests;
- validate representative valid and invalid `.spsq` fixtures with `bin/synapseq -test` or `go run ./cmd/synapseq -test`;
- run `make test`;
- validate each affected skill using the checks shipped with that skill, when available;
- verify that no skill documents syntax, ranges, effects, or behavior that the current implementation no longer supports.

## Git Workflow

- **Default branch**: `main`
- **Feature branches**: `feature/*` (branched from `main`)
- **Bugfix branches**: `bugfix/*` (branched from `main`)
- **Hotfix branches**: `hotfix/*` (branched from `main`)
- **Release branches**: `release/*` (branched from `main`)
- **Maintenance branches**: `chore/*` (branched from `main`)

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

- Follow Go best practices and idioms.
- Keep `internal/types` pure - no dependencies on other internal packages.
- Keep `core` small and stable - avoid expanding public API.
- Prefer clarity over cleverness.
- One way to do each task.
- Less options, more focus.
- Test files: `*_test.go`; use table-driven tests when appropriate.
- All tests must pass before submitting a PR (`make test`).

## Suggested Reading Order

For new contributors, the fastest way to understand the codebase:

1. `cmd/synapseq/main.go`
2. `cmd/synapseq/dispatch.go`
3. `core/context.go`, `core/sequence.go`, `core/generate.go`
4. `spsq/builder.go` and `sbg/sbg.go`
5. `internal/sequence/loadtext.go` and `internal/sequence/parsecontent.go`
6. `internal/parser/*`
7. `internal/audio/renderer.go`, `internal/audio/rendercycle.go`, and `internal/audio/audiosource/*`
8. `internal/remote/*`

For detailed architecture, see `docs/ARCHITECTURE.md`. For DSL syntax, see `docs/SYNTAX.md`.
