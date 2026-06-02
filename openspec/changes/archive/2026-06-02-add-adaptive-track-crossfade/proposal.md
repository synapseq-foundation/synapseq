## Why

Users currently need to keep compatible track types, effects, and ambiance sources on the same channel across consecutive presets. When presets evolve independently, timeline loading fails or requires manual `.spsq` channel editing even though the intended listening behavior is a smooth transition.

## What Changes

- Add automatic adaptive crossfade handling for incompatible consecutive channel states.
- Stop rejecting direct transitions between different active track types, different effect types, or different ambiance names.
- Treat `track off` as fade-compatible at timeline boundaries, using the same fade-out/fade-in behavior currently used for `silence`, with the channel ending inactive after fade-out.
- Preserve the current `silence` fade-in and fade-out behavior.
- Use an adaptive crossfade duration that uses up to 30 seconds from each side of the boundary and clamps to available timeline duration for shorter periods.
- Keep crossfade behavior local to affected channels instead of inserting user-visible timeline periods.
- Update preview and documentation so automatic crossfades are visible and understandable.

## Capabilities

### New Capabilities
- `adaptive-track-crossfade`: Covers automatic per-channel crossfades for incompatible consecutive timeline track states.

### Modified Capabilities

## Impact

- Affected packages: `internal/timeline`, `internal/audio`, `internal/preview`, `internal/sequence`, and related tests.
- The public `.spsq` language becomes more tolerant of channel ordering differences between presets.
- No new external dependencies are expected.
- Existing valid sequences should continue to render with the same behavior unless they previously depended on direct `off` transition errors.
