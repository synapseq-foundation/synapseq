## Why

Music support introduced finite external audio playback alongside looped ambiance playback, but both currently live inside `internal/audio/ambiance`. This makes package ownership misleading and risks future changes coupling music-specific behavior to ambiance-specific behavior.

This change separates format-agnostic external audio loading/runtime code from the policies that are specific to ambiance and music, while preserving the current rendered output and public syntax.

## What Changes

- Extract WAV/MP3 loading, decoding, caching, resampling, and sample reading into a neutral internal audio package.
- Keep ambiance-specific policy in an ambiance package: looped playback, WAV-first option resolution already performed by sequence loading, and source-scoped runtime behavior.
- Move music-specific policy into a music package: finite playback and channel-scoped runtime behavior.
- Update the renderer to depend on explicit ambiance and music packages instead of using music helpers from the ambiance package.
- Preserve existing DSL, public API, file resolution behavior, size limits, crossfade behavior, status output, preview output, and rendered audio semantics.
- Add or retain regression tests proving ambiance and music behavior remains unchanged after the package split.

## Capabilities

### New Capabilities

- `external-audio-source-architecture`: internal package ownership and regression constraints for shared WAV/MP3 external audio source code used by ambiance and music.

### Modified Capabilities

- None. This is an implementation-only refactor with no requirement-level behavior change.

## Impact

- Affected packages:
  - `internal/audio/ambiance`
  - new neutral package under `internal/audio`
  - new `internal/audio/music` package
  - `internal/audio` renderer integration
  - audio runtime/source tests
- No public DSL, `core`, or `spsq` API changes are intended.
- No dependency changes are intended.
