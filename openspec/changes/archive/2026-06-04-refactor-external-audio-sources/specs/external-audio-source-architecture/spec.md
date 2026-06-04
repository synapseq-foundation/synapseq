## ADDED Requirements

### Requirement: Separate shared external audio mechanics from source-specific policies
The implementation SHALL keep WAV/MP3 external audio mechanics in a neutral internal package that is not named after ambiance or music.

#### Scenario: Shared audio code is neutral
- **WHEN** the renderer constructs ambiance and music external audio sources
- **THEN** shared decoding, caching, resampling, and sample-reading code is provided by a neutral internal audio package

#### Scenario: Ambiance wrapper remains ambiance-specific
- **WHEN** ambiance runtime or audio constructors are used
- **THEN** they expose ambiance behavior without music-specific constructor names in the ambiance package

#### Scenario: Music wrapper remains music-specific
- **WHEN** music runtime or audio constructors are used
- **THEN** they expose music behavior from a music package instead of from the ambiance package

### Requirement: Preserve external audio playback behavior
The refactor SHALL NOT change existing ambiance or music playback behavior.

#### Scenario: Ambiance continues to loop
- **WHEN** ambiance playback reads beyond the end of a source
- **THEN** playback restarts using the existing loop behavior

#### Scenario: Music remains finite on the same channel
- **WHEN** music playback reaches EOF on a channel
- **THEN** subsequent samples for that same channel are silent until the channel changes source or becomes inactive

#### Scenario: Music restarts independently on another channel
- **WHEN** a music source reaches EOF on one channel and the same source is later used on another channel
- **THEN** the later channel plays the source from the beginning

#### Scenario: Existing renderer behavior remains stable
- **WHEN** existing sequences using ambiance, music, effects, crossfades, and EOF behavior are rendered
- **THEN** the tests for those behaviors continue to pass without changing expected output semantics

### Requirement: Preserve public and DSL compatibility
The refactor SHALL NOT change public APIs, `.spsq` syntax, resource resolution behavior, file size limits, or supported audio formats.

#### Scenario: DSL remains unchanged
- **WHEN** existing valid `@ambiance`, `ambiance`, `@music`, and `music` declarations are parsed
- **THEN** they continue to parse and load with the same semantics

#### Scenario: Resource resolution remains unchanged
- **WHEN** ambiance or music sources are resolved from local paths or URLs
- **THEN** existing WAV/MP3 priority, MIME handling, and size-limit behavior is preserved
