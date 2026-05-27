## 1. Public Core API

- [x] 1.1 Add `AppContext.LoadFile(path string)` in `core/sequence.go` with the existing file-loading behavior moved into the core layer.
- [x] 1.2 Add `AppContext.LoadContent(content string)` in `core/sequence.go` that parses raw sequence text without file context.
- [x] 1.3 Rename repository call sites from `AppContext.Load` to `AppContext.LoadFile`, including CLI and documentation examples.

## 2. Internal Sequence Loading

- [x] 2.1 Change native `internal/sequence.LoadTextSequence` to accept raw content bytes instead of reading a file path.
- [x] 2.2 Keep source path and base directory context available to parsing for file-based loads.
- [x] 2.3 Remove the incompatible native/WASM function signature split so all build targets use byte-based sequence loading.
- [x] 2.4 Update WASM sequence loading call sites to continue passing bytes through the unified internal API.

## 3. Tests and Validation

- [x] 3.1 Update `internal/sequence` tests that currently pass file paths to pass loaded bytes plus source context, preserving diagnostics and relative-reference coverage.
- [x] 3.2 Add or update `core` tests for `LoadFile` and `LoadContent`.
- [x] 3.3 Run `go test ./...` for native validation.
- [x] 3.4 Run a WASM compile or targeted WASM package check to verify build-tag compatibility.
