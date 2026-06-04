## Why


## What Changes

- Add `(*core.AppContext).LoadContent(content string)` to load a sequence from raw `.spsq` text.
- Rename the existing public file-loading API from `Load(path string)` to `LoadFile(path string)`.
- **BREAKING**: callers of `core.AppContext.Load(path)` must migrate to `LoadFile(path)`.
- Unify `internal/sequence.LoadTextSequence` so it accepts raw content bytes for all build targets.
- Move file reading and path/base directory resolution to the `core` layer for file-based loading.

## Capabilities

### New Capabilities
- `core-sequence-loading`: Public sequence loading APIs for raw content and files, with a unified internal content parser boundary.

### Modified Capabilities

## Impact

- Affects `core/sequence.go` public API and any code/tests calling `AppContext.Load`.
- No new external dependencies are expected.
