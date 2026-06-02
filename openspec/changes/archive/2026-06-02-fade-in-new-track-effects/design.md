## Context

The adaptive crossfade behavior currently treats any active effect type mismatch as an incompatible boundary. That is correct for transitions between two concrete effects, such as `pan -> modulation`, because the runtime effect identity changes.

The `EffectOff -> active effect` and `active effect -> EffectOff` cases are different. In both directions, the concrete effect identity and value are known, and the perceptual transition should be the effect intensity moving between zero and the target value. This can be represented by normal track interpolation if the endpoint with `EffectOff` is prepared with the concrete effect type/value and zero intensity.

## Goals / Non-Goals

**Goals:**
- Make `EffectOff -> EffectPan`, `EffectOff -> EffectModulation`, and `EffectOff -> EffectDoppler` use normal interpolation when track type and ambiance source are otherwise compatible.
- Make `EffectPan`, `EffectModulation`, and `EffectDoppler` to `EffectOff` use normal interpolation when track type and ambiance source are otherwise compatible.
- For fade-in, copy the incoming effect type and effect value onto the previous endpoint with zero intensity.
- For fade-out, preserve the previous effect type and effect value on the next endpoint with zero intensity.
- Keep automatic amplitude crossfade for transitions between different concrete effect types.

**Non-Goals:**
- Do not change crossfade behavior for track type changes or ambiance source changes.
- Do not add effect-specific syntax or user-configurable effect fade duration.
- Do not change `TrackSilence` or `TrackOff` boundary behavior.

## Decisions

### Treat Effect On/Off As Compatible Endpoint Preparation

`AdjustPeriods()` should check for effect on/off transitions before deciding that an effect type mismatch requires boundary crossfade. When the track type and ambiance source are compatible, it should prepare the endpoint that would otherwise be `EffectOff` with the concrete effect type and value but zero intensity.

Rationale: the existing render path already interpolates effect intensity. Preparing the off endpoint keeps the transition inside the normal interpolation model and avoids unnecessary amplitude crossfade.

Alternative considered: keep amplitude crossfade and add a secondary effect-intensity fade. That would mix two behaviors for the same perceptual change and could reduce channel loudness unnecessarily.

### Preserve Crossfade For Concrete Effect Type Changes

Transitions such as `pan -> modulation`, `modulation -> doppler`, or `doppler -> pan` should continue to create automatic boundary crossfade metadata.

Rationale: changing between two concrete effect identities is an incompatible runtime state change, not just an effect being introduced.

### Keep Existing Numeric Interpolation Ownership

The renderer should continue to interpolate effect intensity through `interpolateTrack()` and compile effect state as it does today. The timeline adjustment should only prepare endpoints so the renderer receives compatible effect identity.

Rationale: this keeps the change local and preserves existing audio/render ownership.

## Risks / Trade-offs

- Effect on/off transitions may reset effect runtime state when the effect type changes from off to active or active to off -> this is acceptable because the audible intensity is zero at the prepared off endpoint.
- If future effects have non-intensity parameters beyond `Value`, this endpoint-copy rule may need to copy those parameters too.
