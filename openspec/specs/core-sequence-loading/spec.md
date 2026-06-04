## Requirements

### Requirement: Load sequence from file
The `core` package SHALL expose `AppContext.LoadFile(path string)` to load a text sequence from a filesystem path.

#### Scenario: Loading a valid file
- **WHEN** a caller invokes `LoadFile` with the path to a valid `.spsq` file
- **THEN** the system returns a loaded context for the parsed sequence

#### Scenario: Preserving file context
- **WHEN** `LoadFile` loads a sequence from a relative or absolute path
- **THEN** sequence parsing receives source path and base directory context equivalent to the existing native file-loading behavior

#### Scenario: Loading a missing file
- **WHEN** a caller invokes `LoadFile` with a path that cannot be read
- **THEN** the system returns an error instead of attempting to parse empty content

### Requirement: Load sequence from raw content
The `core` package SHALL expose `AppContext.LoadContent(content string)` to load a text sequence from raw `.spsq` content.

#### Scenario: Loading valid raw content
- **WHEN** a caller invokes `LoadContent` with valid `.spsq` text
- **THEN** the system returns a loaded context for the parsed sequence

#### Scenario: Loading raw content without file context
- **WHEN** `LoadContent` parses sequence text
- **THEN** parsing receives empty source path and base directory context

### Requirement: Use explicit public loading methods
The existing public `AppContext.Load(path string)` file-loading method SHALL be renamed to `AppContext.LoadFile(path string)`.

#### Scenario: File-loading call sites
- **WHEN** repository code loads a sequence through the public `core` API
- **THEN** file-based call sites use `LoadFile` and raw-content call sites use `LoadContent`

### Requirement: Parse text sequences from bytes internally
The `internal/sequence` package SHALL expose a unified text sequence loading function that accepts raw content bytes on all build targets.

#### Scenario: Native sequence parsing
- **WHEN** native code asks `internal/sequence` to load a text sequence
- **THEN** it passes already-read content bytes instead of a file path

#### Scenario: WASM sequence parsing
- **WHEN** WASM code asks `internal/sequence` to load a text sequence
- **THEN** it continues to pass content bytes without requiring filesystem access

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
