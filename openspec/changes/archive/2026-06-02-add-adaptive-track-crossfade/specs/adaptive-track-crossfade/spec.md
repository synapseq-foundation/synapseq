## ADDED Requirements

### Requirement: Crossfade incompatible active channel boundaries
The system SHALL allow consecutive timeline entries to transition between incompatible active track states on the same channel by applying an automatic per-channel crossfade.

#### Scenario: Different track types on the same channel
- **WHEN** a timeline boundary moves a channel from one active track type to a different active track type
- **THEN** sequence loading succeeds and the renderer fades the previous track out before the boundary while fading the next track in after the boundary

#### Scenario: Different effect types on the same channel
- **WHEN** a timeline boundary moves a channel between active tracks with different effect types
- **THEN** sequence loading succeeds and the renderer uses automatic crossfade behavior instead of direct parameter interpolation

#### Scenario: Different ambiance sources on the same channel
- **WHEN** a timeline boundary moves a channel between active ambiance tracks with different ambiance names
- **THEN** sequence loading succeeds and the renderer crossfades the two ambiance sources around the boundary

### Requirement: Preserve compatible track interpolation
The system SHALL continue to use existing full-period interpolation for compatible active track states on the same channel.

#### Scenario: Compatible active tracks
- **WHEN** a timeline boundary moves a channel between active tracks with the same track type, effect type, and ambiance name
- **THEN** the previous period interpolates numeric track values toward the next period according to the period transition and steps

### Requirement: Treat off as a fade boundary state
The system SHALL treat `TrackOff` as fade-compatible at timeline boundaries while preserving inactive steady-state behavior outside the fade window.

#### Scenario: Active track fades out to off
- **WHEN** a timeline boundary moves a channel from an active track to `TrackOff`
- **THEN** sequence loading succeeds, the active track fades out before the boundary, and the channel is inactive after the boundary

#### Scenario: Off fades in to active track
- **WHEN** a timeline boundary moves a channel from `TrackOff` to an active track
- **THEN** sequence loading succeeds, the channel is inactive before the boundary, and the active track fades in after the boundary

### Requirement: Preserve silence fade behavior
The system SHALL preserve the existing fade-in and fade-out behavior for `TrackSilence` at timeline boundaries.

#### Scenario: Silence fades into active track
- **WHEN** a timeline boundary moves a channel from `TrackSilence` to an active track
- **THEN** the previous side uses the next track shape with zero amplitude so the active track fades in

#### Scenario: Active track fades into silence
- **WHEN** a timeline boundary moves a channel from an active track to `TrackSilence`
- **THEN** the next side preserves the previous track shape with zero amplitude so the active track fades out

### Requirement: Use adaptive crossfade duration
The system SHALL use 30 seconds as the maximum before and after an incompatible boundary and clamp each side to the available adjacent period duration.

#### Scenario: Sufficient adjacent duration
- **WHEN** both adjacent sides of an incompatible boundary have at least 30 seconds available
- **THEN** the crossfade uses 30 seconds of fade-out before the boundary and 30 seconds of fade-in after the boundary

#### Scenario: Short previous side
- **WHEN** the previous period has less than 30 seconds available before an incompatible boundary
- **THEN** the fade-out duration equals the available previous period duration

#### Scenario: Short next side
- **WHEN** the next period has less than 30 seconds available after an incompatible boundary
- **THEN** the fade-in duration equals the available next period duration

### Requirement: Keep automatic crossfades hidden from timeline structure
The system SHALL apply automatic crossfades without inserting additional user-visible periods into the loaded sequence.

#### Scenario: Loaded period count remains unchanged
- **WHEN** sequence loading creates automatic crossfade behavior for an incompatible boundary
- **THEN** the loaded sequence contains only the periods declared by the `.spsq` timeline

### Requirement: Preview automatic crossfades
The preview system SHALL represent automatic crossfade boundaries without converting them into user-visible timeline periods.

#### Scenario: Preview shows crossfade transition
- **WHEN** preview data is generated for a sequence with an automatic crossfade boundary
- **THEN** the affected channel transition is identified as an automatic crossfade rather than a direct compatible interpolation
