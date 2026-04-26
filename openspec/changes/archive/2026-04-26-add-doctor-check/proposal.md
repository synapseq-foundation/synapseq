## Why

Users need a way to quickly diagnose whether their environment has all the required external tools for SynapSeq to function properly. A doctor command provides immediate feedback with actionable fix suggestions, reducing debugging time and improving user experience.

## What Changes

- Add `-doctor` flag to CLI that checks for required external tools
- Create doctor functionality in `cmd/synapseq` package
- Display check results with visual indicators (✔/✖)
- Show platform-specific suggested fixes when tools are missing
- Use `internal/cli` package for output styling
- Check: `ffmpeg` and `ffplay` (currently supported), `git` and `gh` (future support)

## Capabilities

### New Capabilities
- `doctor-check`: Health diagnostic tool that verifies external tool installation status and provides platform-specific remediation suggestions

### Modified Capabilities
- None

## Impact

- New file: `cmd/synapseq/doctor.go` - Doctor check implementation
- CLI changes: New `-doctor` flag handling in `cmd/synapseq`
- Uses existing `internal/cli` package for output styling
- No changes to `core` or other internal packages