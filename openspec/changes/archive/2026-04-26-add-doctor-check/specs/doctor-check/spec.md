## ADDED Requirements

### Requirement: Doctor command checks external tools
The system SHALL provide a `-doctor` flag that verifies the presence and functionality of required external tools and displays results with clear status indicators.

#### Scenario: All tools installed and working
- **WHEN** user runs `synapseq -doctor`
- **THEN** system displays ✔ for each installed tool (ffmpeg, ffplay, git, gh)
- **AND** no suggested fixes are shown

#### Scenario: Tool not installed
- **WHEN** user runs `synapseq -doctor` with a missing tool
- **THEN** system displays ✖ for the missing tool
- **AND** shows "Suggested fixes:" section with platform-specific installation commands

#### Scenario: gh installed but not authenticated
- **WHEN** user runs `synapseq -doctor` with gh installed but not authenticated
- **THEN** system displays ✖ for gh with message "gh not authenticated"
- **AND** shows "gh auth login" as a suggested fix

### Requirement: Platform-specific suggested fixes
The system SHALL display platform-specific installation commands based on the user's operating system.

#### Scenario: macOS or Linux with missing tool
- **WHEN** user runs `synapseq -doctor` on macOS or Linux with a missing tool
- **THEN** system suggests installation via `brew install <tool>`

#### Scenario: Windows with missing tool
- **WHEN** user runs `synapseq -doctor` on Windows with a missing tool
- **THEN** system suggests installation via `winget install <tool>`

### Requirement: Tool verification via exec.LookPath
The system SHALL use the standard `exec.LookPath` function to verify if each tool is available in the system PATH.

#### Scenario: Tool exists in PATH
- **WHEN** `exec.LookPath` succeeds for a tool name
- **THEN** the tool is marked as installed (✔)

#### Scenario: Tool does not exist in PATH
- **WHEN** `exec.LookPath` fails for a tool name
- **THEN** the tool is marked as not installed (✖)