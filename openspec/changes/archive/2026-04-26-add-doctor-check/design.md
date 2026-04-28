## Context

SynapSeq depends on external tools (`ffmpeg`, `ffplay`) for audio encoding and playback. Users often encounter issues when these tools are not installed or not in the system PATH. A diagnostic command helps users quickly identify and fix these environment issues.

The `-doctor` flag provides a self-service diagnostic that:
- Checks each required tool independently
- Reports clear pass/fail status with visual indicators
- Provides actionable, platform-specific remediation steps

## Goals / Non-Goals

**Goals:**
- Provide quick environment diagnostics without loading a sequence
- Display clear visual indicators (✔/✖) for each tool check
- Show platform-appropriate installation commands (brew for macOS/Linux, winget for Windows)
- Check both installed status and basic functionality for gh

**Non-Goals:**
- Fix missing tools automatically
- Support authentication for gh (placeholder check only, not implementing auth flow)
- Check for specific versions or capabilities of tools
- Integration with hub or preview functionality

## Decisions

### Location: `cmd/synapseq` package

The doctor functionality lives in the main CLI package because:
- It's a user-facing diagnostic command that outputs directly to terminal
- No other packages need to import this functionality
- Uses existing `cli` package styling for consistent output

**No changes to `core` package** - the implementation stays entirely within `cmd/synapseq`.

### Output Styling: `internal/cli` package

Using the existing `cli` package for output formatting:
- Uses `cli.SuccessText()` for installed tools (green)
- Uses `cli.ErrorText()` for missing tools (red)
- Consistent with existing CLI output style

### Platform Detection: `runtime.GOOS`

Using `runtime.GOOS` to determine the platform:
- `darwin` → suggest `brew install`
- `linux` → suggest `brew install`
- `windows` → suggest `winget install`

## Risks / Trade-offs

- **Limited gh auth check**: Only checking if gh is installed, not authenticated. Future enhancement possible.
- **Tool version detection**: Not checking for minimum versions, which could lead to false positives if an incompatible version is installed.
- **Output format lock-in**: Terminal output with ANSI colors may not suit all environments (could add `--json` flag later).

## Open Questions

1. Should the doctor command exit with non-zero status if any tool is missing?
2. Add `--json` flag for machine-readable output?