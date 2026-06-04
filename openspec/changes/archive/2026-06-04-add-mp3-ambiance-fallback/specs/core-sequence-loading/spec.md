## ADDED Requirements

### Requirement: Resolve local ambiance files with WAV priority
The sequence loader SHALL resolve extensionless local ambiance option paths by selecting an existing WAV file before considering an MP3 file.

#### Scenario: Local WAV exists
- **WHEN** a native sequence declares `@ambiance rain audio/rain` and `audio/rain.wav` exists
- **THEN** the loaded sequence ambiance map resolves `rain` to `audio/rain.wav`

#### Scenario: Local WAV missing and MP3 exists
- **WHEN** a native sequence declares `@ambiance rain audio/rain`, `audio/rain.wav` does not exist, and `audio/rain.mp3` exists
- **THEN** the loaded sequence ambiance map resolves `rain` to `audio/rain.mp3`

#### Scenario: Local WAV and MP3 both missing
- **WHEN** a native sequence declares `@ambiance rain audio/rain` and neither `audio/rain.wav` nor `audio/rain.mp3` exists
- **THEN** sequence loading fails with an error that identifies both attempted paths

#### Scenario: Local WAV exists but is invalid
- **WHEN** a native sequence declares `@ambiance rain audio/rain`, `audio/rain.wav` exists, and `audio/rain.wav` cannot be decoded as valid WAV
- **THEN** rendering fails with the WAV decode error and does not fall back to `audio/rain.mp3`

#### Scenario: Local path includes explicit extension
- **WHEN** a native sequence declares an ambiance local path with a file extension
- **THEN** sequence loading rejects the path using the existing local path validation behavior

#### Scenario: Remote URL keeps explicit file behavior
- **WHEN** a sequence declares an ambiance URL with a full file path and extension
- **THEN** sequence loading preserves the URL value without applying local WAV-to-MP3 fallback
