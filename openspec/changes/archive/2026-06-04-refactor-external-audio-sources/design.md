## Context

The current music implementation reuses `internal/audio/ambiance/audio.go` and `runtime.go` for both looped ambiance and finite music. That reuse is technically useful, but the package name now hides two different concerns:

- external audio mechanics that are shared by WAV/MP3 ambiance and music;
- playback policy that differs between ambiance and music.

Ambiance is looped and currently uses source-scoped runtime buffers. Music is finite and now requires channel-scoped playback state so the same source can restart on a different channel while an ended channel remains silent.

## Goals / Non-Goals

**Goals:**

- Move format-agnostic WAV/MP3 decode, cache, resample, and sample reading into a neutral package.
- Keep ambiance-specific behavior in an ambiance package.
- Keep music-specific behavior in a music package.
- Make playback mode and runtime buffer scope explicit in code.
- Preserve all current behavior, including ambiance looping, music finite EOF, music restart on another channel, path resolution, status output, preview output, and renderer results.
- Keep public APIs, DSL syntax, and resource resolution unchanged.

**Non-Goals:**

- No changes to `.spsq` syntax.
- No changes to `core` or `spsq` public API shape.
- No changes to WAV/MP3 decoding libraries or vendored dependencies.
- No changes to path resolution priority or file size limits.
- No audio quality, crossfade, or effect algorithm changes.

## Decisions

1. Create a neutral shared package under `internal/audio`.

   Proposed name: `internal/audio/audiosource`.

   This package owns:
   - WAV/MP3 decode helpers;
   - resampling;
   - byte cache;
   - `ReadSamplesAt`;
   - playback mode handling for looped vs finite sources;
   - generic named runtime indexing and buffer preparation.

   Rationale: the code is audio-source infrastructure, not ambiance-specific. A neutral package makes future external source types easier to add without expanding `ambiance`.

2. Represent playback policy explicitly.

   Add a small enum-like type:

   ```go
   type PlaybackMode int

   const (
       PlaybackLoop PlaybackMode = iota
       PlaybackFinite
   )
   ```

   Ambiance uses `PlaybackLoop`; music uses `PlaybackFinite`.

   Rationale: a boolean `loop` works but hides intent at call sites. Named modes make regressions easier to spot.

3. Represent runtime buffer/read scope explicitly.

   Add a scope policy:

   ```go
   type BufferScope int

   const (
       BufferScopeSource BufferScope = iota
       BufferScopeChannel
   )
   ```

   Ambiance uses source-scoped buffers/read state. Music uses channel-scoped buffers/read state.

   Rationale: this captures the important behavioral difference found after music implementation: finite EOF is per channel, not per source.

4. Keep ambiance and music packages as thin policy wrappers.

   `internal/audio/ambiance` should expose ambiance-specific constructors and avoid music-named functions.

   `internal/audio/music` should expose music-specific constructors and avoid ambiance-named functions.

   Rationale: renderer code should read as `ambiance.NewRuntime(...)` and `music.NewRuntime(...)`, making ownership obvious.

5. Preserve tests first, then move code.

   Existing behavior tests should continue to pass after each extraction step. Add focused package-ownership tests or static checks where useful, but avoid brittle tests that over-constrain private implementation details.

   Rationale: this refactor is valuable only if it does not change audio semantics.

## Risks / Trade-offs

- Package churn may create noisy diffs → Keep moves mechanical and avoid unrelated renames.
- Runtime policy extraction could accidentally change read position behavior → Preserve and extend tests for ambiance loop, music EOF, and music restart on another channel.
- Introducing a generic package can become over-abstracted → Keep the API concrete and local to WAV/MP3 external sources; do not introduce broad interfaces beyond the existing loader callback and sample reader.
- Import cycles are possible if wrappers are not kept thin → Put shared code below both `ambiance` and `music`, and keep `audiosource` independent from those packages.
