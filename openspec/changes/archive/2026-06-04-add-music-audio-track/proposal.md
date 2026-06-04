## Why

SynapSeq can already layer loopable ambiance audio, but there is no equivalent track type for one-shot or finite music beds. Users need a `@music` source that can play MP3/WAV assets without automatic looping while the rest of the sequence continues after the music ends.

## What Changes

- Add a top-level `@music <name> [path-or-url]` option with the same naming and path rules as `@ambiance`.
- Add `music <name> ... amplitude <value>` preset tracks.
- Support the same effects on music tracks that ambiance tracks support.
- Resolve local `@music` paths without extensions, preferring `.mp3` first and falling back to `.wav` only when the MP3 file is missing.
- Preserve explicit URL behavior for music sources, using extension or MIME type to identify MP3/WAV.
- Apply an 80 MB maximum file size to music sources, for both MP3 and WAV.
- Render music sources without automatic looping. When a music source reaches EOF, that music channel becomes silent until the sequence changes it or rendering ends.
- Rename track/source naming internals such as `AmbianceName` to a neutral field that works for both ambiance and music.

## Capabilities

### New Capabilities

- `music-audio-tracks`: Defines `@music` options, `music` track syntax, supported formats, effects, non-looping playback behavior, and music source size limits.

### Modified Capabilities

- `core-sequence-loading`: Sequence loading must resolve extensionless local music paths using MP3-first, WAV-second fallback while preserving existing option path validation rules.

## Impact

- Affected packages: `internal/types`, `internal/parser`, `internal/sequence`, `internal/audio`, `internal/audio/ambiance` or shared audio-source package, `internal/preview`, `core`, `spsq`, docs, and templates.
- Data model impact: track source-name fields should become format-neutral so both ambiance and music can reference named audio sources cleanly.
- DSL impact: adds `@music` and `music` while preserving existing `@ambiance` and `ambiance` syntax.
- Rendering impact: introduces finite audio playback that returns silence after EOF instead of looping.
- Compatibility: existing ambiance behavior and WAV-first ambiance resolution must remain unchanged.
