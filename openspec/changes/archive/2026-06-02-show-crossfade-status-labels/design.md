## Context

Status reporting currently prints period-change summaries from `internal/audio/status/reporter.go`. Each active channel emits the period start track, and if `TrackStart` differs from `TrackEnd`, the reporter emits an arrow line for the end track. Automatic crossfades are stored on `types.Period` as `CrossfadeOut` and `CrossfadeIn`, but the reporter does not surface that metadata.

Render planning resolves crossfade duration in `internal/audio/renderplan.go` using a 30 second maximum clamped to the available period duration. The status reporter needs to display the same duration, so the duration rule should not remain private to render planning.

## Goals / Non-Goals

**Goals:**

- Append compact fade metadata to the affected period-change track line:
  - `(fade-out 30s)` on the outgoing track line
  - `(fade-in 30s)` on the incoming track line
- Use the adaptive duration that the renderer actually applies.
- Share the adaptive crossfade duration helper from `internal/timeline` so both render planning and status reporting use one rule.
- Color the fade label through the existing status color plumbing and palette tokens.
- Keep output compact and readable for multiple channels.

**Non-Goals:**

- Do not add real-time "crossfade started" events during rendering.
- Do not change audio behavior, sequence loading, preview behavior, or `.spsq` syntax.
- Do not add new palette colors unless existing tokens are insufficient.
- Do not display labels for effect intensity fade-in/fade-out cases that do not create `CrossfadeIn` or `CrossfadeOut` metadata.

## Decisions

### Put the duration helper in `internal/timeline`

The adaptive duration rule is tied to timeline boundaries rather than signal generation. `timeline.CrossfadeDuration(availableMs int) int` should return `0` for non-positive durations and otherwise clamp to the existing 30 second maximum. `renderplan.go` will call this helper instead of a local `crossfadeDuration`.

Alternative considered: keep a duplicate helper in `status`. That is simpler locally but risks status reporting `30s` while rendering uses a different value after future changes.

Alternative considered: export the helper from `internal/audio`. That creates an awkward dependency direction for `internal/audio/status`, because `audio` already owns the runtime that imports the status package.

### Append labels to existing track strings

The reporter should keep the current visual shape:

```text
       <track-start> (fade-out 30s)
   ->  <track-end> (fade-in 30s)
```

This is less noisy than adding dedicated crossfade lines and aligns the label with the track that is actually being faded.

### Format durations for status labels

Status labels should be compact:

- whole seconds: `30s`
- fractional seconds when needed: `7.5s`
- milliseconds for durations below 1000 ms: `750ms`

The first implementation can use deterministic formatting with tests. The display should avoid unnecessary trailing `.0`.

### Color the fade label as transition metadata

Track text currently uses `palette.Green`, arrows and steps use `palette.Ochre`, and transition text uses `palette.MutedWarm`. The fade label is transition metadata attached to a track, so `palette.Ochre` is the best existing token: it is already used for timing/step annotations and visually separates the label from the green track body.

If this reads too strong in terminal output, `palette.MutedWarm` is the fallback; it keeps the label secondary but may be too subtle.

## Risks / Trade-offs

- Duplicating duration formatting between status tests and implementation could make tests brittle. Mitigation: test a small matrix of important values only.
- `CountActiveChannels(view.Channels)` can hide a crossfade label for a channel that is off in the current runtime channel state but has period metadata. Mitigation: compute the period-change display channel count from both runtime channels and period metadata, or ensure active crossfade metadata expands the display range.
- Labels append to already long track strings. Mitigation: keep labels short and avoid source/destination descriptions in this change.
- Moving the helper from `audio` to `timeline` touches render behavior. Mitigation: preserve existing render plan tests and add focused helper tests.
