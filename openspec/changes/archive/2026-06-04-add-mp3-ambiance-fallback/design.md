## Context

SynapSeq currently resolves local ambiance options to `.wav` during sequence loading, stores the resolved path in `SequenceOptions.Ambiance`, then `internal/audio/ambiance` loads, decodes, optionally resamples, caches, and loops each WAV source. The audio runtime already keeps independent decoder state per ambiance index and restarts a source when EOF is reached.

MP3 support should preserve that rendering model. WAV remains the preferred source for loop quality because MP3 encoding can add encoder delay, padding, and frame-boundary artifacts. MP3 exists as a distribution fallback when no WAV is present.

## Goals / Non-Goals

**Goals:**

- Support MP3 ambiance sources through `github.com/gopxl/beep/v2/mp3`.
- Keep the `.spsq` syntax unchanged.
- Resolve local extensionless ambiance paths by preferring `.wav` and falling back to `.mp3` only when the WAV file is absent.
- Preserve existing behavior for remote URLs while validating that the remote resource is WAV or MP3.
- Keep looping behavior functionally equivalent to WAV, including independent decoder positions per ambiance index.
- Apply the same 20 MB size limit to WAV and MP3 ambiance assets.
- Vendor all new decoder dependencies.

**Non-Goals:**

- Do not support explicit local `.wav` or `.mp3` extensions in `@ambiance`; local paths remain extensionless.
- Do not fall back from an existing but invalid WAV to MP3.
- Do not add new track syntax, effects, or public rendering APIs.
- Do not guarantee mathematically gapless MP3 loops when the source MP3 itself contains encoder delay or padding.
- Do not change MP3 export behavior, which is handled separately through ffmpeg.

## Decisions

1. Local fallback belongs in ambiance option resolution.

   Native sequence loading should resolve each local ambiance path to a concrete existing file. It should try `<base>/<path>.wav` first and return it if present. Only when the WAV path is missing should it try `<base>/<path>.mp3`. If neither exists, the error should mention both attempted paths.

   This keeps `SequenceOptions.Ambiance` as a map from name to concrete resource path, avoiding a broader public model change.

2. Remote URLs stay explicit but gain format validation.

   Remote ambiance values should continue to be full URLs. If a URL path has `.wav` or `.mp3`, that extension determines the expected decoder. If the URL path lacks an extension, the resource loader should inspect `Content-Type` and accept WAV or MP3 MIME types. Unsupported extensions or MIME types should fail with a validation error.

   Remote URLs should not use the local `.wav` to `.mp3` fallback because probing two URL variants changes the existing explicit-URL contract.

3. Format detection should be resource metadata, not parser syntax.

   The parser should continue to treat `@ambiance <name> <path-or-url>` as a name plus raw path token. Format detection should occur after path resolution and resource loading, where local existence checks and remote headers are available.

4. Generalize the ambiance decoder while preserving runtime behavior.

   `internal/audio/ambiance.Audio` should track each source's format alongside cached bytes. Decoder creation should route to `wav.Decode` or `mp3.Decode` behind a small local helper. `ReadSamplesAt`, `restartAt`, `PrepareBuffers`, channel index management, and mixer code should remain behaviorally unchanged.

5. MP3 looping should restart from cached bytes.

   On EOF, MP3 should be looped by closing/recreating that source's decoder from cached bytes or by seeking to the beginning only when the decoder and reader are known to support seeking. Recreating from cached bytes is the most predictable path because `mp3.Decode` delegates seeking behavior to the underlying reader. This also matches the existing restart model for WAV.

6. Resampling should operate on decoded streams.

   The current WAV-only resample helper should become format-agnostic: decode the cached source with the selected decoder, resample with `beep.Resample` when needed, and cache the resampled output in a decoder-friendly representation. Encoding resampled content back to WAV is acceptable internally because downstream playback only needs decoded samples, not preservation of the original file format.

7. Rename the ambiance size limit.

   `MaxWavFileSize` should become `MaxAmbianceFileSize` and remain 20 MB. `resource.GetFile` should apply that limit to both WAV and MP3 ambiance inputs.

## Risks / Trade-offs

- MP3 encoder delay or padding can still cause audible loop gaps -> keep WAV preferred, document MP3 as fallback, and ensure implementation does not add additional decoder-position gaps.
- Falling back in sequence resolution makes missing-file behavior depend on filesystem state -> add tests for WAV-present, WAV-missing/MP3-present, both-missing, and existing-invalid-WAV cases.
- Remote MIME detection may require reading response headers before full validation -> keep URL extension authoritative when present, and inspect `Content-Type` only for extensionless URLs.
- Vendor drift can break `-mod=vendor` builds -> run dependency vendoring as part of implementation and include vendored decoder dependencies in tests.
- Re-encoding resampled MP3 to WAV increases memory use for long MP3 files -> retain the existing 20 MB input cap and keep chunked playback unchanged after cache creation.
