## Context

Timeline loading currently builds each period from a preset with `TrackStart` and `TrackEnd` initialized to the same 16-channel track array. When a later timeline entry is parsed, `internal/timeline.AdjustPeriods()` mutates the previous period's `TrackEnd` to match the next period's `TrackStart`, allowing the renderer to interpolate values across the interval.

That model only supports direct interpolation when the same channel contains compatible track identity: same active track type, same effect type, and same ambiance source. Incompatible changes currently produce validation errors, which forces users to manually align preset channel order or insert explicit `silence` entries.

The requested behavior is automatic smoothing: incompatible channel boundaries should render as a fade-out/fade-in around the boundary using up to 30 seconds from each side, without adding visible timeline entries.

## Goals / Non-Goals

**Goals:**
- Allow consecutive timeline entries to use different track types, effect types, or ambiance names on the same channel without validation failure.
- Treat `TrackOff` as fade-compatible with active tracks at period boundaries, ending inactive after fade-out or starting from silence before fade-in.
- Preserve existing `TrackSilence` fade-in/fade-out semantics.
- Apply crossfade behavior independently per channel.
- Use an adaptive duration that uses 30 seconds as the maximum on each side of the boundary and clamps to available adjacent period duration.
- Keep user-visible timeline periods unchanged for duration, diagnostics, and preview structure.
- Make preview output communicate automatic crossfade boundaries.

**Non-Goals:**
- Do not introduce new `.spsq` syntax for configuring crossfade duration in this change.
- Do not add new public `core` APIs.
- Do not remap or reorder preset tracks during parsing.
- Do not create actual hidden periods in the `Sequence.Periods` array.

## Decisions

### Model Crossfades As Boundary Metadata

Add an internal representation for per-channel boundary crossfades instead of inserting synthetic periods. The boundary metadata should describe the previous track, next track, channel index, boundary time, and the resolved fade-out/fade-in durations.

Rationale: synthetic periods would change period counts and would complicate step validation, preview grouping, and duration semantics. Boundary metadata keeps the user's timeline intact while allowing the renderer to mix two independent track states around the boundary.

Alternative considered: mutate `TrackStart` and `TrackEnd` only. This cannot represent a true crossfade between incompatible track identities because the current interpolation path holds one track identity at a time.

### Keep Normal Interpolation For Compatible Active Tracks

When consecutive channel states are compatible, keep the current interpolation path: `AdjustPeriods()` carries the next track into the previous period's `TrackEnd`, and `renderPlan` interpolates numeric values over the full period using the configured transition and steps.

Rationale: this preserves behavior for existing sequences and avoids changing the sound of valid compatible transitions.

Alternative considered: route every transition through the new crossfade path. That would simplify branching but would change existing long-form parameter interpolation behavior.

### Treat `off` As Fade-Compatible Boundary State

Replace the current active-to-off and off-to-active errors with fade-compatible behavior. Active-to-off fades the active track out near the boundary and leaves the next side inactive. Off-to-active fades the active track in after the boundary.

Rationale: users expect an inactive channel to behave like silence at a boundary, while still preserving `TrackOff` as the steady-state result outside the fade window.

Alternative considered: convert `off` to `silence` in parsed periods. That would lose the difference between an intentionally inactive channel and a silent-but-shaped channel.

### Resolve Adaptive Duration Per Side

For each boundary, use up to 30 seconds of fade-out before the boundary and up to 30 seconds of fade-in after the boundary. Clamp each side independently to the available adjacent interval:

```text
fadeOutMs = min(30000, boundaryTime - previousPeriod.Time)
fadeInMs  = min(30000, nextPeriodEndTime - boundaryTime)
```

If a side has no positive duration, that side contributes no fade. The final timeline period has no following period, so only boundaries between declared periods can create crossfades.

Rationale: independent clamping avoids consuming more time than exists and handles short periods predictably.

Alternative considered: `min(30000, intervalBetweenPeriods / 2)` for both sides. That works for a single interval but does not account for the duration after the boundary, which matters for fade-in.

### Execute Crossfades In The Audio Render Plan

The render plan should detect active crossfade windows at cue time and compile both the base channel signal and any overlay crossfade signal needed for that channel. The renderer then mixes the fade-out and fade-in contributions with per-window amplitude scaling.

Rationale: true crossfade requires two track identities to sound during the boundary window. This belongs near cue planning/rendering rather than sequence parsing.

Alternative considered: render two full sequences and mix them externally. That would be expensive and would not fit the current concrete audio engine.

### Preview Uses The Same Boundary Semantics

Preview data should indicate automatic crossfades for affected segment items and graph transitions. It does not need to create user-visible periods, but it should prevent the preview from implying a hard invalid change.

Rationale: users need to understand why a transition is smooth even when tracks differ.

## Risks / Trade-offs

- Crossfade mixing can increase peak amplitude near boundaries -> clamp or preserve existing output limiting behavior and add render tests around boundary windows.
- Short periods may produce very short fades -> document adaptive clamping and test intervals shorter than the 30-second maximum.
- Ambiance crossfades may require two ambiance sources at once -> ensure source lookup/playback can run both sources independently during the overlap window.
- Existing preview assumptions may rely on one track state per channel at a time -> keep timeline nodes unchanged and add explicit crossfade annotations instead of reshaping preview periods.
- Steps and transition curves do not naturally apply to incompatible crossfades -> keep steps/transition for compatible interpolation and use a simple fade curve for incompatible boundary windows unless a later change introduces configurable crossfade curves.
