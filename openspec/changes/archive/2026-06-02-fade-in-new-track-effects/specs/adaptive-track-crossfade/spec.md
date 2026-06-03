## MODIFIED Requirements

### Requirement: Crossfade incompatible active channel boundaries
The system SHALL allow consecutive timeline entries to transition between incompatible active track states on the same channel by applying an automatic per-channel crossfade. The system SHALL NOT use amplitude crossfade when the only effect incompatibility is a transition between `EffectOff` and an active effect; that case SHALL use normal interpolation with effect intensity fading between zero and the active effect intensity.

#### Scenario: Different track types on the same channel
- **WHEN** a timeline boundary moves a channel from one active track type to a different active track type
- **THEN** sequence loading succeeds and the renderer fades the previous track out before the boundary while fading the next track in after the boundary

#### Scenario: Different concrete effect types on the same channel
- **WHEN** a timeline boundary moves a channel between active tracks with different non-off effect types
- **THEN** sequence loading succeeds and the renderer uses automatic crossfade behavior instead of direct parameter interpolation

#### Scenario: Effect fades in from off
- **WHEN** a timeline boundary moves a compatible active channel from `EffectOff` to an active effect type
- **THEN** sequence loading succeeds, no automatic amplitude crossfade is created for that channel, and the previous endpoint uses the incoming effect type and value with zero intensity

#### Scenario: Effect fades out to off
- **WHEN** a timeline boundary moves a compatible active channel from an active effect type to `EffectOff`
- **THEN** sequence loading succeeds, no automatic amplitude crossfade is created for that channel, and the next endpoint uses the previous effect type and value with zero intensity

#### Scenario: Different ambiance sources on the same channel
- **WHEN** a timeline boundary moves a channel between active ambiance tracks with different ambiance names
- **THEN** sequence loading succeeds and the renderer crossfades the two ambiance sources around the boundary

### Requirement: Preserve compatible track interpolation
The system SHALL continue to use existing full-period interpolation for compatible active track states on the same channel, including compatible tracks where an effect is introduced from `EffectOff` or removed to `EffectOff`.

#### Scenario: Compatible active tracks
- **WHEN** a timeline boundary moves a channel between active tracks with the same track type, effect type, and ambiance name
- **THEN** the previous period interpolates numeric track values toward the next period according to the period transition and steps

#### Scenario: Compatible active track gains an effect
- **WHEN** a timeline boundary moves a channel between compatible active tracks and the next track introduces an active effect from `EffectOff`
- **THEN** the previous period interpolates track values and effect intensity toward the next period according to the period transition and steps

#### Scenario: Compatible active track loses an effect
- **WHEN** a timeline boundary moves a channel between compatible active tracks and the next track removes an active effect to `EffectOff`
- **THEN** the previous period interpolates track values and effect intensity down to zero according to the period transition and steps
