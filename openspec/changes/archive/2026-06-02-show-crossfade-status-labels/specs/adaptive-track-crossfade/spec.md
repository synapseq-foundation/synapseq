## MODIFIED Requirements

### Requirement: Use adaptive crossfade duration
The system SHALL use 30 seconds as the maximum before and after an incompatible boundary and clamp each side to the available adjacent period duration. The system SHALL expose this duration resolution as shared timeline behavior so render planning and status reporting use the same adaptive duration.

#### Scenario: Sufficient adjacent duration
- **WHEN** both adjacent sides of an incompatible boundary have at least 30 seconds available
- **THEN** the crossfade uses 30 seconds of fade-out before the boundary and 30 seconds of fade-in after the boundary

#### Scenario: Short previous side
- **WHEN** the previous period has less than 30 seconds available before an incompatible boundary
- **THEN** the fade-out duration equals the available previous period duration

#### Scenario: Short next side
- **WHEN** the next period has less than 30 seconds available after an incompatible boundary
- **THEN** the fade-in duration equals the available next period duration

#### Scenario: Status and render use the same resolved duration
- **WHEN** an automatic crossfade is present on a period boundary
- **THEN** both render planning and status reporting resolve the crossfade duration through the shared timeline duration rule
