## 1. Timeline Behavior

- [x] 1.1 Add helper logic to detect bidirectional `EffectOff <-> active effect` transitions on otherwise compatible active tracks
- [x] 1.2 Update `AdjustPeriods()` to prepare the off endpoint with the concrete effect type and value at zero intensity
- [x] 1.3 Ensure these cases do not create `CrossfadeIn` or `CrossfadeOut` metadata
- [x] 1.4 Preserve automatic amplitude crossfade for transitions between different concrete effect types

## 2. Tests

- [x] 2.1 Add timeline test for `EffectOff -> EffectPan` intensity fade-in endpoint preparation
- [x] 2.2 Add timeline test for active effect to `EffectOff` intensity fade-out endpoint preparation
- [x] 2.3 Add regression test confirming active effect type changes still create crossfade metadata
- [x] 2.4 Add render plan or audio tests confirming effect intensity interpolates from zero to target and from target to zero

## 3. Verification

- [x] 3.1 Run focused tests for `internal/timeline` and `internal/audio`
- [x] 3.2 Run full project test suite with `make test`
