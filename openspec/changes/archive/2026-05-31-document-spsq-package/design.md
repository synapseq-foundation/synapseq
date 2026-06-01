## Context

The repository already exposes documented public Go packages for loading/rendering sequences (`core`) and external playback/encoding integrations (`external`). The new `spsq` package is also public, but it currently lacks a package-level documentation entry point and is not referenced from the main README or architecture guide.

The `spsq` package builds textual `.spsq` content through a fluent `Builder`. The generated text is intended to flow into the existing `core.AppContext.LoadContent` path, preserving the existing parser, sequence loading, validation, rendering, preview, and external integration boundaries.

## Goals / Non-Goals

**Goals:**
- Document `spsq` with a package comment in the same broad style as `core/doc.go` and `external/doc.go`.
- Show a concise programmatic builder example that produces `.spsq` text and loads it through `core`.
- Explain in architecture docs that `spsq` is a public construction helper for `.spsq` text, separate from parsing and rendering.
- Mention the package in README Go API documentation with a pkg.go.dev link.

**Non-Goals:**
- Changing the `spsq.Builder` public API.
- Adding validation, parsing, rendering, or file output responsibilities to `spsq`.
- Promoting the temporary `cmd/builder/main.go` helper as a supported command.
- Reworking existing `core` or `external` documentation beyond the references needed for `spsq`.

## Decisions

- Add a dedicated `spsq/doc.go` file instead of embedding the package comment in an implementation file.
  - Rationale: this matches the documented public-package pattern already used by `core` and `external`.
  - Alternative considered: add the package comment to `builder.go`; this would work for Go doc generation but mixes overview documentation with implementation.

- Position `spsq` as a textual builder that hands off to `core.LoadContent`.
  - Rationale: this keeps `core` as the loading/rendering API and preserves the parser/sequence/audio ownership boundaries.
  - Alternative considered: document `spsq` as a sequence execution API; that would blur responsibilities and imply behavior it does not own.

- Use the temporary `cmd/builder/main.go` as source material for the example, but keep the documentation example minimal.
  - Rationale: examples should demonstrate the integration path without encouraging use of the temporary command package.
  - Alternative considered: documenting every builder method; that is better left to exported method comments and pkg.go.dev indexes.

- Add README and architecture references rather than a new standalone guide.
  - Rationale: the package is small, and pkg.go.dev should be the canonical API reference for public Go symbols.
  - Alternative considered: creating a new docs page; that adds maintenance surface without a clear need.

## Risks / Trade-offs

- Documentation may overstate API stability while the builder is new -> describe current responsibilities precisely and avoid promising validation or complete DSL coverage beyond existing exported methods.
- Go doc examples can drift from the actual API -> keep examples short and based on current exported methods.
- README could become crowded with API links -> keep the Go API section concise and split loading/rendering (`core`) from programmatic sequence construction (`spsq`).
