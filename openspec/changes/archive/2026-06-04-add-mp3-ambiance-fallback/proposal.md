## Why

Ambiance tracks currently depend on WAV files only, which makes sequence assets larger and less convenient to distribute. MP3 support gives authors a smaller fallback format while preserving WAV as the preferred format for seamless loops.

## What Changes

- Add MP3 as a supported ambiance audio format using `github.com/gopxl/beep/v2/mp3`.
- Keep local `@ambiance <name> <path>` syntax extensionless; local resolution tries `<path>.wav` first, then `<path>.mp3` only when the WAV file is missing.
- Preserve current URL behavior: remote ambiance paths must specify the full URL with extension, or provide a response MIME type that identifies WAV or MP3 when the URL has no extension.
- Reject remote ambiance URLs whose extension or MIME type is neither WAV nor MP3.
- Keep WAV as the recommended and prioritized loop format; MP3 is a fallback and must still loop correctly within the limits of MP3 encoding.
- Rename the shared 20 MB ambiance file limit to `MaxAmbianceFileSize` and apply it to both WAV and MP3 inputs.
- Vendor the MP3 decoder dependencies required by `github.com/gopxl/beep/v2/mp3`.

## Capabilities

### New Capabilities

- `ambiance-audio-formats`: Defines supported ambiance file formats, local fallback behavior, URL format detection, and looping expectations.

### Modified Capabilities

- `core-sequence-loading`: Sequence loading must resolve extensionless local ambiance paths using WAV-first, MP3-second fallback while preserving existing validation rules.

## Impact

- Affected packages: `internal/sequence`, `internal/resource`, `internal/types`, `internal/audio/ambiance`, `core`, and documentation under `docs/`.
- Dependency impact: import `github.com/gopxl/beep/v2/mp3` and vendor its indirect decoder dependencies.
- Public DSL impact: no new syntax; behavior changes only for extensionless local ambiance paths when a matching MP3 exists and WAV is absent.
- Compatibility: existing WAV ambiance files and URL-based ambiance paths continue to work.
