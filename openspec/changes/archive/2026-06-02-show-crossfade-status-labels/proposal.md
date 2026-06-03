## Why

Automatic channel crossfades are currently applied silently during rendering, so the status report does not tell users when a track line is fading in or out. The status output should make these hidden transitions visible without adding noisy extra lines or requiring users to inspect preview data.

## What Changes

- Add compact crossfade labels to status period-change track lines:
  - outgoing crossfade tracks append `(fade-out <duration>)`
  - incoming crossfade tracks append `(fade-in <duration>)`
- Resolve the displayed duration with the same adaptive crossfade rule used by rendering: 30 seconds maximum, clamped to the available period duration.
- Move the adaptive crossfade duration helper into `internal/timeline` so render planning and status reporting share one implementation.
- Color the fade label with an existing palette token that distinguishes transition metadata from track text while staying consistent with current status colors.

## Capabilities

### New Capabilities
- `status-reporting`: Covers user-facing render status output, including period-change track lines and status metadata labels.

### Modified Capabilities
- `adaptive-track-crossfade`: Defines that the adaptive crossfade duration is shared by render behavior and status reporting.

## Impact

- Affected code:
  - `internal/timeline`: add shared adaptive crossfade duration helper and tests.
  - `internal/audio/renderplan.go`: reuse the shared helper instead of local duration logic.
  - `internal/audio/status/reporter.go`: append compact fade labels to affected track strings and color them through the existing palette.
  - `internal/audio/status/reporter_test.go`, `internal/audio/renderplan_test.go`, and timeline tests.
- No public API changes.
- No changes to sequence loading behavior, audio synthesis, preview structure, or `.spsq` syntax.
