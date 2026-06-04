## ADDED Requirements

### Requirement: Resolve local music files with MP3 priority
The sequence loader SHALL resolve extensionless local music option paths by selecting an existing MP3 file before considering a WAV file.

#### Scenario: Local MP3 exists
- **WHEN** a native sequence declares `@music meditation audio/meditation` and `audio/meditation.mp3` exists
- **THEN** the loaded sequence music map resolves `meditation` to `audio/meditation.mp3`

#### Scenario: Local MP3 missing and WAV exists
- **WHEN** a native sequence declares `@music meditation audio/meditation`, `audio/meditation.mp3` does not exist, and `audio/meditation.wav` exists
- **THEN** the loaded sequence music map resolves `meditation` to `audio/meditation.wav`

#### Scenario: Local MP3 and WAV both missing
- **WHEN** a native sequence declares `@music meditation audio/meditation` and neither `audio/meditation.mp3` nor `audio/meditation.wav` exists
- **THEN** sequence loading fails with an error that identifies both attempted paths

#### Scenario: Local music path includes explicit extension
- **WHEN** a native sequence declares a music local path with a file extension
- **THEN** sequence loading rejects the path using the existing local option path validation behavior

#### Scenario: Remote music URL keeps explicit file behavior
- **WHEN** a sequence declares a music URL with a full file path and extension
- **THEN** sequence loading preserves the URL value without applying local MP3-to-WAV fallback

#### Scenario: Remote music URL without extension uses MIME type
- **WHEN** a sequence declares a music URL without a file extension
- **THEN** resource loading determines whether the music source is MP3 or WAV from the response MIME type
