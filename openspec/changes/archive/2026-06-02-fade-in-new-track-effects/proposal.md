## Why

Automatic crossfade currently treats any active effect type mismatch as an incompatible boundary. When a track moves between no effect and a concrete effect, the smoother behavior is to fade effect intensity in or out while preserving normal track interpolation instead of crossfading the whole channel amplitude.

## What Changes

- Treat both `EffectOff -> active effect` and `active effect -> EffectOff` as compatible for normal track interpolation.
- For `EffectOff -> active effect`, copy the incoming effect type and effect value onto the outgoing period endpoint with zero intensity so the existing interpolation fades effect intensity in.
- For `active effect -> EffectOff`, keep the previous effect type and effect value on the outgoing period endpoint with zero intensity so the existing interpolation fades effect intensity out completely.
- Do not create automatic amplitude crossfade metadata for these effect-on/effect-off cases when the track type and ambiance source are otherwise compatible.
- Preserve existing automatic crossfade behavior for transitions between different active effect types.

## Capabilities

### New Capabilities

### Modified Capabilities
- `adaptive-track-crossfade`: Refine effect compatibility so effects fade in or out by intensity instead of triggering amplitude crossfade.

## Impact

- Affected packages: `internal/timeline`, `internal/audio`, and related tests.
- Existing incompatible active effect changes still use automatic crossfade.
- No new `.spsq` syntax or public API changes are expected.
