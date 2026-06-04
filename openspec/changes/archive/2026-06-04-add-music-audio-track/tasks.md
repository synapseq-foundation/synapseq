## 1. Domain Model and Parser

- [x] 1.1 Add music keywords and track type constants without changing existing ambiance values unexpectedly.
- [x] 1.2 Rename `Track.AmbianceName` and parser declarations to a neutral source-name field used by ambiance and music.
- [x] 1.3 Add `Music` storage to parse/sequence options while preserving existing ambiance options.
- [x] 1.4 Parse `@music <name> [path-or-url]` using the same name validation and shorthand behavior as `@ambiance`.
- [x] 1.5 Parse `music <name> amplitude ...` tracks.
- [x] 1.6 Parse `music <name> effect pan ... intensity ... amplitude ...` tracks.
- [x] 1.7 Parse `music <name> effect modulation ... intensity ... amplitude ...` tracks.
- [x] 1.8 Update string formatting, summaries, and diagnostic text for music tracks.

## 2. Path Resolution and Resource Limits

- [x] 2.1 Add `MaxMusicFileSize` set to 80 MB for music audio sources.
- [x] 2.2 Generalize resource audio loading so ambiance and music can share MP3/WAV format detection with different size limits.
- [x] 2.3 Resolve local music paths with `.mp3` priority and `.wav` fallback.
- [x] 2.4 Preserve extensionless local path validation for music options.
- [x] 2.5 Preserve explicit remote URL behavior for music and detect extensionless remote music by MIME type.
- [x] 2.6 Reject unsupported music source extensions and MIME types.
- [x] 2.7 Keep ambiance path resolution, WAV priority, and 20 MB size limit unchanged.

## 3. Audio Runtime and Rendering

- [x] 3.1 Extract or adapt the named audio decoder/cache layer so it can serve both ambiance and music sources.
- [x] 3.2 Add music source indexing and validation for named music references.
- [x] 3.3 Add renderer state for music sources alongside ambiance state.
- [x] 3.4 Add sync-engine updates for music channel source assignment.
- [x] 3.5 Add music source sampling and route music tracks through stereo effect processing.
- [x] 3.6 Implement finite EOF behavior: after a music source ends, return silence instead of restarting.
- [x] 3.7 Ensure rendering continues after music EOF until the sequence itself ends.
- [x] 3.8 Extend automatic boundary crossfade compatibility so music name changes and music/type changes behave like ambiance transitions.
- [x] 3.9 Preserve existing ambiance loop behavior after the shared audio changes.

## 4. Public API, Builder, Preview, and WASM

- [x] 4.1 Expose music metadata from loaded sequences where ambiance metadata is currently exposed.
- [x] 4.2 Add builder support for registering music and declaring music tracks.
- [x] 4.3 Update preview formatting and track analysis for music tracks.
- [x] 4.4 Update status reporting for music tracks.
- [x] 4.5 Define WASM behavior for music URLs consistently with ambiance URL-only constraints.
- [x] 4.6 Update docs and templates with `@music` examples and MP3-first local fallback guidance.

## 5. Tests and Verification

- [x] 5.1 Add parser tests for `@music` and `music` track forms.
- [x] 5.2 Add sequence loading tests for music local MP3 priority, WAV fallback, missing files, explicit extension rejection, and remote URL preservation.
- [x] 5.3 Add resource tests for the 80 MB music limit and unsupported music formats.
- [x] 5.4 Add audio runtime tests for music MP3/WAV decode, resampling, effects, and silence after EOF.
- [x] 5.5 Add renderer tests proving music EOF does not stop sequence rendering.
- [x] 5.6 Add timeline/render tests for music-to-music, music-to-other-track, and other-track-to-music automatic crossfades.
- [x] 5.7 Add regression tests proving ambiance still loops and still resolves WAV before MP3.
- [x] 5.8 Add preview/status/builder/core tests for music metadata and display.
- [x] 5.9 Run focused parser, sequence, resource, audio, preview, builder, core, and WASM tests.
- [x] 5.10 Run the full project test suite.
