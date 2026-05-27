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
