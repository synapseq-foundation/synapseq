## ADDED Requirements

### Requirement: Display automatic crossfade labels in period status
The status reporter SHALL append compact fade labels to period-change track lines for automatic channel crossfades.

#### Scenario: Outgoing track shows fade-out duration
- **WHEN** a displayed period has active `CrossfadeOut` metadata for a channel
- **THEN** the period-change line for that channel's outgoing track includes `(fade-out <duration>)`

#### Scenario: Incoming track shows fade-in duration
- **WHEN** a displayed period has active `CrossfadeIn` metadata for a channel
- **THEN** the period-change line for that channel's incoming track includes `(fade-in <duration>)`

#### Scenario: No label without automatic crossfade metadata
- **WHEN** a displayed track line has no active `CrossfadeOut` or `CrossfadeIn` metadata for its channel
- **THEN** the period-change output does not append a fade label to that track line

#### Scenario: Labels are colored as status metadata
- **WHEN** status colors are enabled and a fade label is displayed
- **THEN** the label uses the status reporter's palette-based coloring for transition metadata

### Requirement: Format crossfade label durations compactly
The status reporter SHALL format automatic crossfade durations in compact human-readable units.

#### Scenario: Whole-second duration
- **WHEN** the adaptive crossfade duration is an exact whole number of seconds
- **THEN** the fade label displays the duration as `<seconds>s`

#### Scenario: Fractional-second duration
- **WHEN** the adaptive crossfade duration is at least one second but not an exact whole number of seconds
- **THEN** the fade label displays the duration as a fractional second value without unnecessary trailing zeroes

#### Scenario: Subsecond duration
- **WHEN** the adaptive crossfade duration is less than one second
- **THEN** the fade label displays the duration in milliseconds
