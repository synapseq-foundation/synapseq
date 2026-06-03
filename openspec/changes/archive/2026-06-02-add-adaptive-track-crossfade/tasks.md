## 1. Timeline Semantics

- [x] 1.1 Add internal per-channel boundary metadata for automatic crossfades without changing declared period count
- [x] 1.2 Update `AdjustPeriods()` compatibility logic so different active track types, effect types, and ambiance names create crossfade metadata instead of validation errors
- [x] 1.3 Update boundary handling so `TrackOff` fades like silence when adjacent to an active track while remaining inactive outside the fade window
- [x] 1.4 Preserve existing `TrackSilence` fade-in and fade-out endpoint mutation behavior
- [x] 1.5 Add timeline tests for incompatible type, effect, ambiance, off-to-active, active-to-off, silence, and compatible interpolation cases

## 2. Adaptive Duration

- [x] 2.1 Implement adaptive crossfade duration resolution using 30000 ms as the maximum per side
- [x] 2.2 Clamp fade-out duration to available time before the boundary and fade-in duration to available time after the boundary
- [x] 2.3 Add tests for full-length crossfades and short adjacent periods

## 3. Audio Rendering

- [x] 3.1 Extend render plan cues so a channel can carry base signal state plus active crossfade fade-out/fade-in signal state
- [x] 3.2 Render incompatible crossfades by mixing previous and next track identities during their boundary windows
- [x] 3.3 Ensure compatible tracks continue to use existing full-period numeric interpolation with transition and steps
- [x] 3.4 Add render plan and renderer tests for crossfade amplitude behavior, incompatible track identity, and unchanged compatible interpolation
- [x] 3.5 Add ambiance crossfade coverage to verify two different ambiance sources can overlap during the boundary window

## 4. Preview And Documentation

- [x] 4.1 Update preview data/view models to mark automatic crossfade transitions without inserting extra periods
- [x] 4.2 Update preview graph or segment rendering tests for automatic crossfade boundaries
- [x] 4.3 Document automatic crossfade behavior, adaptive duration, and `off` boundary behavior in syntax or how-it-works docs

## 5. Verification

- [x] 5.1 Run focused tests for `internal/timeline`, `internal/audio`, `internal/preview`, and `internal/sequence`
- [x] 5.2 Run full project test suite with `make test`
