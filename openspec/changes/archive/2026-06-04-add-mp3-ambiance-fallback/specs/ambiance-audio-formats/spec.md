## ADDED Requirements

### Requirement: Support WAV and MP3 ambiance files
The system SHALL support WAV and MP3 files as ambiance audio sources.

#### Scenario: WAV ambiance source
- **WHEN** a sequence references an ambiance source resolved to a WAV file
- **THEN** the renderer decodes and mixes the source using the WAV decoder

#### Scenario: MP3 ambiance source
- **WHEN** a sequence references an ambiance source resolved to an MP3 file
- **THEN** the renderer decodes and mixes the source using the MP3 decoder

#### Scenario: Unsupported ambiance source format
- **WHEN** an ambiance source resolves to a format other than WAV or MP3
- **THEN** the system rejects the source with an error that identifies the unsupported format

### Requirement: Preserve ambiance playback behavior across formats
The system SHALL apply the same sample rate handling, stereo requirement, amplitude handling, effects, buffering, and channel mixing behavior to MP3 ambiance sources as it applies to WAV ambiance sources.

#### Scenario: MP3 source with mismatched sample rate
- **WHEN** an MP3 ambiance source has a sample rate different from the sequence sample rate
- **THEN** the system resamples the decoded source to the sequence sample rate before playback

#### Scenario: MP3 source on ambiance track with effects
- **WHEN** an ambiance track using an MP3 source declares supported ambiance effects
- **THEN** the renderer applies the effects during mixing the same way it applies them to WAV ambiance sources

### Requirement: Loop ambiance sources by format
The system SHALL loop WAV and MP3 ambiance sources when playback reaches the end of the source.

#### Scenario: WAV ambiance reaches end of source
- **WHEN** WAV ambiance playback reaches EOF during rendering
- **THEN** playback restarts from the beginning of the same source without changing channel assignment

#### Scenario: MP3 ambiance reaches end of source
- **WHEN** MP3 ambiance playback reaches EOF during rendering
- **THEN** playback restarts from the beginning of the same source without changing channel assignment

#### Scenario: MP3 source is shorter than render buffer
- **WHEN** an MP3 ambiance source ends before the current render buffer is filled
- **THEN** the remaining buffer samples are filled by continuing from the start of the same source

### Requirement: Validate remote ambiance format
The system SHALL determine remote ambiance format from the URL extension when present, or from the response MIME type when the URL has no file extension.

#### Scenario: Remote URL with WAV extension
- **WHEN** an ambiance URL path ends with `.wav`
- **THEN** the system treats the source as WAV

#### Scenario: Remote URL with MP3 extension
- **WHEN** an ambiance URL path ends with `.mp3`
- **THEN** the system treats the source as MP3

#### Scenario: Remote URL without extension and WAV MIME type
- **WHEN** an ambiance URL has no file extension and the response identifies a WAV MIME type
- **THEN** the system treats the source as WAV

#### Scenario: Remote URL without extension and MP3 MIME type
- **WHEN** an ambiance URL has no file extension and the response identifies an MP3 MIME type
- **THEN** the system treats the source as MP3

#### Scenario: Remote URL with unsupported extension
- **WHEN** an ambiance URL path ends with an extension other than `.wav` or `.mp3`
- **THEN** the system rejects the source with an unsupported format error

#### Scenario: Remote URL without extension and unsupported MIME type
- **WHEN** an ambiance URL has no file extension and the response MIME type is neither WAV nor MP3
- **THEN** the system rejects the source with an unsupported format error

### Requirement: Limit ambiance file size consistently
The system SHALL apply one 20 MB maximum size limit to all supported ambiance audio formats.

#### Scenario: WAV ambiance exceeds size limit
- **WHEN** a WAV ambiance source exceeds the maximum ambiance file size
- **THEN** the system reads no more than the maximum allowed ambiance file size

#### Scenario: MP3 ambiance exceeds size limit
- **WHEN** an MP3 ambiance source exceeds the maximum ambiance file size
- **THEN** the system reads no more than the maximum allowed ambiance file size
