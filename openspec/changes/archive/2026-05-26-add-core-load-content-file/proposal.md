## Why

The public `core` API currently exposes `Load(path string)` while the internal sequence loader has different contracts for native and WASM builds. This makes the loading boundary unclear: file I/O is mixed into `internal/sequence` on native builds, while WASM already correctly passes raw content bytes.

## What Changes

- Add `(*core.AppContext).LoadContent(content string)` to load a sequence from raw `.spsq` text.
- Rename the existing public file-loading API from `Load(path string)` to `LoadFile(path string)`.
- **BREAKING**: callers of `core.AppContext.Load(path)` must migrate to `LoadFile(path)`.
- Unify `internal/sequence.LoadTextSequence` so it accepts raw content bytes for all build targets.
- Move file reading and path/base directory resolution to the `core` layer for file-based loading.
- Preserve WASM behavior where sequence content is supplied as bytes instead of read from disk.

## Capabilities

### New Capabilities
- `core-sequence-loading`: Public sequence loading APIs for raw content and files, with a unified internal content parser boundary.

### Modified Capabilities

## Impact

- Affects `core/sequence.go` public API and any code/tests calling `AppContext.Load`.
- Affects `internal/sequence/loadtext.go` and `internal/sequence/loadtext_wasm.go` by consolidating their API shape.
- May affect CLI, WASM, or tests that call `internal/sequence.LoadTextSequence` directly.
- No new external dependencies are expected.
