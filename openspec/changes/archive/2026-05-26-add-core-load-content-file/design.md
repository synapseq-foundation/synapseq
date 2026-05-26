## Context

`core.AppContext.Load(path string)` currently loads a sequence from a file by delegating to `internal/sequence.LoadTextSequence`. On native builds, `internal/sequence.LoadTextSequence` reads the file and resolves the source/base directory. On WASM builds, the same function name accepts raw bytes and parses content directly.

This creates two incompatible contracts behind build tags. The desired boundary is for `core` to decide whether input is a file path or raw text, and for `internal/sequence` to parse already-loaded bytes consistently.

## Goals / Non-Goals

**Goals:**

- Expose explicit public APIs: `LoadFile(path string)` for file paths and `LoadContent(content string)` for raw `.spsq` text.
- Move file I/O out of `internal/sequence.LoadTextSequence`.
- Make `internal/sequence.LoadTextSequence` accept `[]byte` on all build targets.
- Preserve source path and base directory behavior for file-based loading so includes, extends, diagnostics, or relative resolution continue to work.
- Preserve WASM behavior where callers supply content bytes directly.

**Non-Goals:**

- Redesign sequence parsing or validation.
- Add new DSL syntax.
- Expand remote loading, preview generation, or audio rendering APIs.
- Add compatibility shims beyond the planned `Load` to `LoadFile` rename unless implementation discovers a required internal migration path.

## Decisions

1. Public API names will distinguish source type.

   `LoadFile(path string)` will own file reads and path normalization. `LoadContent(content string)` will convert the string to bytes and load without file path or base directory context. This makes call sites explicit and avoids interpreting arbitrary content as a filesystem path.

   Alternative considered: keep `Load(path string)` and add `LoadContent`. That would avoid a breaking change but leaves the ambiguous API name in place.

2. `internal/sequence.LoadTextSequence` will parse bytes, not read files.

   The internal sequence package should receive `rawContent []byte` plus source context needed by parsing. The native and WASM implementations should no longer have incompatible function signatures. This matches the existing WASM behavior and keeps `internal/sequence` focused on sequence construction rather than file access.

   Alternative considered: add a second internal function for byte loading while keeping file loading in native `loadtext.go`. That would preserve existing internals but keep two loading boundaries.

3. File source metadata remains available for file loads.

   `LoadFile` will resolve the absolute input path and base directory before calling the internal loader, preserving current behavior from native `internal/sequence/loadtext.go`.

   Alternative considered: parse file content without source metadata. That would simplify the function signature but risks breaking relative reference behavior and diagnostic source locations.

## Risks / Trade-offs

- **Breaking public API rename** -> Update all repository call sites and tests from `Load` to `LoadFile`; document the break in the proposal and tasks.
- **Internal function signature churn** -> Search for all `LoadTextSequence` callers and migrate them in one pass.
- **Raw content lacks filesystem context** -> `LoadContent` will use empty source path and base directory; features requiring relative file context must use `LoadFile`.
- **Build tag divergence can hide compile failures** -> Run normal tests and at least compile/check the WASM package after implementation.
