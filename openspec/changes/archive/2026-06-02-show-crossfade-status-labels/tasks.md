## 1. Shared Duration Helper

- [x] 1.1 Add a shared adaptive crossfade duration helper in `internal/timeline`
- [x] 1.2 Move the 30 second maximum crossfade duration rule out of `internal/audio/renderplan.go`
- [x] 1.3 Update render planning to use the shared timeline helper without changing audio behavior
- [x] 1.4 Add timeline tests for zero, negative, short, and maximum-clamped duration inputs

## 2. Status Reporter Labels

- [x] 2.1 Add compact duration formatting for status fade labels
- [x] 2.2 Append `(fade-out <duration>)` to outgoing track lines with active `CrossfadeOut` metadata
- [x] 2.3 Append `(fade-in <duration>)` to incoming track lines with active `CrossfadeIn` metadata
- [x] 2.4 Color fade labels through the reporter color helpers using an existing palette token appropriate for transition metadata
- [x] 2.5 Ensure period-change display includes channels that have crossfade metadata even if current runtime channels would otherwise hide them

## 3. Tests

- [x] 3.1 Add reporter tests for fade-out and fade-in labels with adaptive durations
- [x] 3.2 Add reporter tests for compact duration formatting, including whole seconds, fractional seconds, and milliseconds
- [x] 3.3 Add ANSI/color test coverage confirming fade labels are colorized when colors are enabled
- [x] 3.4 Preserve render plan tests proving crossfade duration behavior is unchanged after using the shared helper

## 4. Verification

- [x] 4.1 Run focused tests for `internal/timeline`, `internal/audio/status`, and `internal/audio`
- [x] 4.2 Run full project test suite with `make test`
