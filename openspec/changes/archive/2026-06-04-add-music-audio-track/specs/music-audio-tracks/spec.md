## ADDED Requirements

### Requirement: Register named music sources
The system SHALL support top-level `@music <name> [path-or-url]` options for named music audio sources.

#### Scenario: Music option with explicit path
- **WHEN** a sequence declares `@music meditation audio/meditation`
- **THEN** sequence loading registers `meditation` as a named music source using the provided path

#### Scenario: Music option shorthand
- **WHEN** a sequence declares `@music tibetan-bells`
- **THEN** sequence loading registers `tibetan-bells` as a named music source using the name as the path

#### Scenario: Invalid music name
- **WHEN** a sequence declares `@music` with an invalid name
- **THEN** sequence loading rejects the option using the same name validation rules as ambiance

### Requirement: Parse music tracks
The parser SHALL support `music <name>` track declarations under presets.

#### Scenario: Music track with amplitude
- **WHEN** a preset contains `music meditation amplitude 50`
- **THEN** the parsed track is a music track referencing the named source `meditation` with amplitude 50

#### Scenario: Music track with pan effect
- **WHEN** a preset contains `music meditation effect pan 0.5 intensity 60 amplitude 50`
- **THEN** the parsed track is a music track with the same pan effect semantics supported by ambiance

#### Scenario: Music track with modulation effect
- **WHEN** a preset contains `music meditation effect modulation 2.5 intensity 40 amplitude 50`
- **THEN** the parsed track is a music track with the same modulation effect semantics supported by ambiance

#### Scenario: Unknown music source name
- **WHEN** a sequence uses `music meditation amplitude 50` without registering `@music meditation`
- **THEN** sequence rendering or loading fails with an unknown music source error

### Requirement: Support WAV and MP3 music files
The system SHALL support MP3 and WAV files as music audio sources.

#### Scenario: MP3 music source
- **WHEN** a music source resolves to an MP3 file
- **THEN** the renderer decodes and mixes the source using the MP3 decoder

#### Scenario: WAV music source
- **WHEN** a music source resolves to a WAV file
- **THEN** the renderer decodes and mixes the source using the WAV decoder

#### Scenario: Unsupported music source format
- **WHEN** a music source resolves to a format other than MP3 or WAV
- **THEN** the system rejects the source with an error that identifies the unsupported format

### Requirement: Render music without looping
The renderer SHALL play music sources without automatic looping and SHALL continue rendering the sequence after music reaches EOF.

#### Scenario: Music ends before sequence render ends
- **WHEN** a music source reaches EOF while the sequence still has remaining render time
- **THEN** the renderer continues rendering and the music channel contributes silence

#### Scenario: Music ends within a render buffer
- **WHEN** a music source reaches EOF before the current render buffer is filled
- **THEN** samples after EOF for that music channel are silent for the rest of the buffer

#### Scenario: Ambiance still loops independently
- **WHEN** ambiance and music are both active and the music source reaches EOF
- **THEN** ambiance continues using its existing loop behavior while music remains silent

### Requirement: Apply ambiance-compatible effects to music
The renderer SHALL apply the same supported effects to music tracks that it applies to ambiance tracks.

#### Scenario: Music pan effect
- **WHEN** a music track has a pan effect
- **THEN** the renderer applies the pan effect to the decoded stereo music source

#### Scenario: Music modulation effect
- **WHEN** a music track has a modulation effect
- **THEN** the renderer applies the modulation effect to the decoded stereo music source

### Requirement: Crossfade music track transitions
The system SHALL apply automatic boundary crossfades to music track transitions using the same compatibility rules used for ambiance source changes and track type changes.

#### Scenario: Music source name changes on same channel
- **WHEN** consecutive timeline periods use active music tracks with different music source names on the same channel
- **THEN** sequence loading succeeds and the renderer crossfades the outgoing and incoming music sources around the boundary

#### Scenario: Music changes to another active track type
- **WHEN** a timeline boundary moves a channel from an active music track to another active track type
- **THEN** sequence loading succeeds and the renderer crossfades the outgoing music track with the incoming track

#### Scenario: Another active track type changes to music
- **WHEN** a timeline boundary moves a channel from another active track type to an active music track
- **THEN** sequence loading succeeds and the renderer crossfades the outgoing track with the incoming music track

#### Scenario: Music source reaches EOF during crossfade
- **WHEN** a music source reaches EOF while it participates in an automatic crossfade
- **THEN** the music side of the crossfade contributes silence after EOF while the remaining crossfade continues

### Requirement: Limit music file size
The system SHALL apply an 80 MB maximum size limit to supported music audio formats.

#### Scenario: MP3 music exceeds size limit
- **WHEN** an MP3 music source exceeds the maximum music file size
- **THEN** the system reads no more than the maximum allowed music file size

#### Scenario: WAV music exceeds size limit
- **WHEN** a WAV music source exceeds the maximum music file size
- **THEN** the system reads no more than the maximum allowed music file size
