## 1. Shared Audio Package

- [x] 1.1 Create `internal/audio/audiosource` for shared WAV/MP3 external audio mechanics.
- [x] 1.2 Move decode, cache, validation, resample, and sample-reading code out of `internal/audio/ambiance/audio.go` into the neutral package.
- [x] 1.3 Replace the boolean loop flag with explicit playback modes for looped and finite playback.
- [x] 1.4 Preserve existing loader injection so ambiance still uses `resource.GetAmbianceFile` and music still uses `resource.GetMusicFile`.
- [x] 1.5 Move shared audio tests from ambiance-specific coverage to the neutral package where appropriate.

## 2. Runtime Package Ownership

- [x] 2.1 Move generic named-source indexing, period start precomputation, buffer preparation, and channel/source scope behavior into `internal/audio/audiosource`.
- [x] 2.2 Replace implicit runtime branching with explicit buffer scope policy for source-scoped and channel-scoped playback.
- [x] 2.3 Preserve source-scoped runtime behavior for ambiance.
- [x] 2.4 Preserve channel-scoped runtime behavior for music, including restart on another channel after EOF on a previous channel.

## 3. Ambiance and Music Wrappers

- [x] 3.1 Keep `internal/audio/ambiance` as a thin wrapper for looped ambiance audio and runtime construction.
- [x] 3.2 Add `internal/audio/music` as a thin wrapper for finite music audio and runtime construction.
- [x] 3.3 Remove music-specific constructor names from the ambiance package.
- [x] 3.4 Update renderer imports and construction to use `ambiance` for ambiance and `music` for music.
- [x] 3.5 Keep source mixing behavior unchanged while updating package references.

## 4. Regression Tests

- [x] 4.1 Verify ambiance still loops for WAV and MP3 sources.
- [x] 4.2 Verify music still returns silence after EOF for WAV and MP3 sources.
- [x] 4.3 Verify the same music source can restart from the beginning on a different channel.
- [x] 4.4 Verify renderer music EOF still does not stop sequence rendering.
- [x] 4.5 Verify effects and automatic crossfades involving ambiance and music continue to pass existing tests.
- [x] 4.6 Add or update a package-ownership test/check proving renderer no longer imports music helpers from the ambiance package.

## 5. Validation

- [x] 5.1 Run focused tests for `internal/audio/audiosource`, `internal/audio/ambiance`, `internal/audio/music`, `internal/audio`, and `internal/audio/sources`.
- [x] 5.2 Run parser, sequence, resource, preview, status, core, WASM, and builder tests to verify compatibility.
- [x] 5.3 Run the full project test suite with `make test`.
- [x] 5.4 Run OpenSpec validation for `refactor-external-audio-sources`.
