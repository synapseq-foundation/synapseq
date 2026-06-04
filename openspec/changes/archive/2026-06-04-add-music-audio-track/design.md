## Context

Ambiance already provides named external audio sources with WAV/MP3 decoding, resampling, effects, buffering, and channel mixing. Its key assumption is automatic looping: EOF restarts the source from the beginning. Music needs nearly the same pipeline, but with finite playback semantics and a different local format preference.

The current domain model also carries `AmbianceName` on `Track`. Adding music with another named source would make that name too narrow, so the implementation should move toward a neutral source-name field while preserving existing ambiance behavior.

## Goals / Non-Goals

**Goals:**

- Add `@music <name> [path-or-url]` as a top-level option.
- Add `music <name>` track declarations with the same amplitude and supported effects as ambiance.
- Resolve local music paths extensionlessly with `.mp3` priority and `.wav` fallback.
- Support MP3 and WAV music sources, including remote URL extension/MIME detection.
- Apply an 80 MB maximum music source size.
- Render music without automatic looping; after EOF, the channel contributes silence.
- Refactor source naming so ambiance and music use a neutral track source-name field.

**Non-Goals:**

- Do not change ambiance looping, WAV-first priority, or 20 MB size limit.
- Do not add seeking, start offsets, trimming, ducking, or music-specific effects.
- Do not make music restart automatically when it reaches EOF inside a single period.
- Do not add explicit local file extensions to `.spsq` option paths.
- Do not guarantee MP3 gapless looping because music is intentionally non-looping.

## Decisions

1. Represent music as a separate track type.

   Add a `TrackMusic` value instead of overloading `TrackAmbiance`. This keeps render and preview behavior explicit: ambiance is loopable background audio; music is finite external audio.

2. Rename source-name model fields.

   Replace `Track.AmbianceName` and related parser declarations with a neutral field such as `SourceName` or `AudioName`. Existing ambiance parsing should populate the same field. Any user-facing strings should continue to say "ambiance" or "music" based on the track type.

3. Keep options separate by source class.

   Sequence options should contain distinct maps for ambiance and music. This prevents accidental name collisions and allows different resolution policy and size limits:

   - ambiance: `.wav` first, `.mp3` fallback, 20 MB;
   - music: `.mp3` first, `.wav` fallback, 80 MB.

4. Generalize the audio source loader.

   The current ambiance loader can become a shared named-audio loader or have a sibling music loader built from the same lower-level decoder pieces. The lower layer should handle path, format, cached bytes, decoder state, sample rate, channel count, and resampling. Runtime policy should control EOF:

   - loop mode for ambiance: EOF restarts from cached bytes;
   - finite mode for music: EOF marks the source ended and fills remaining/current samples with silence.

5. Music playback position is per source index, not per channel instance.

   Preserve the current one decoder position per named source index model unless implementation shows a need for independent channel instances. This matches current ambiance behavior and keeps initial scope small.

6. Track changes should reset music assignment intentionally.

   When the sync engine assigns a channel to a music track for a new period/source, the runtime should select the source for that channel. If a channel remains on the same music source across periods, playback should continue from the current decoder position rather than restart just because the timeline advanced. If a channel switches away from music and later switches back, implementation should reset to the beginning unless the existing sync model makes continuation more coherent.

7. Reuse effects and mixing path.

   Music sources should produce stereo samples and pass through the same effect handling as ambiance. Silence after EOF should happen before effect processing so a finished music channel remains silent.

8. Music participates in automatic boundary crossfades like ambiance.

   Track compatibility should treat music source identity the same way ambiance source identity is treated. A change from one music source name to another, or from music to another active track type on the same channel, should create the existing automatic boundary crossfade metadata. The difference from ambiance remains EOF policy only: during a crossfade, music contributes decoded samples while available and silence after natural EOF.

## Risks / Trade-offs

- Field rename touches parser, types, timeline, preview, renderer, and tests -> keep the rename mechanical and covered by existing ambiance tests before adding music behavior.
- Shared loader refactor could destabilize ambiance -> add regression tests for existing ambiance WAV/MP3 loop behavior.
- Music EOF inside a render buffer can create partial-buffer behavior -> explicitly test that samples after EOF are zero and rendering continues.
- Music crossfade can overlap two finite sources -> use the same overlapping-source approach as ambiance crossfade and treat ended music sources as silence.
- Per-source decoder state means two channels using the same music name may share playback position -> document or test current behavior; defer per-channel instances unless product behavior requires it.
- Large MP3/WAV music files increase memory use because sources are cached -> enforce the 80 MB limit and keep the existing cache model for now.
