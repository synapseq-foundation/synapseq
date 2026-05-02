## ADDED Requirements

### Requirement: CLI flag metadata map
The system SHALL provide a `map[string]string` in `internal/cli/completion.go` mapping each CLI flag name (without leading `-`) to a description of 50 characters or fewer, containing no `:` character.

#### Scenario: Map contains all completable flags
- **WHEN** the completion metadata map is initialized
- **THEN** it SHALL contain every flag exposed by `CLIOptions` with a valid description (≤50 chars, no `:`)

#### Scenario: Description validation
- **WHEN** a flag description is added to the map
- **THEN** the description MUST be 50 characters or fewer and MUST NOT contain the `:` character

### Requirement: -completion-args outputs flag list
The system SHALL print all CLI flags with their descriptions in `{param}:{desc}` format when `-completion-args` is passed.

#### Scenario: Output format correctness
- **WHEN** the user runs `synapseq -completion-args`
- **THEN** the system SHALL print each flag as `{flag_name}:{description}` with the flag name excluding the leading `-`, one per line, to stdout

#### Scenario: Flag names exclude hyphen in output
- **WHEN** `-completion-args` is invoked
- **THEN** the output SHALL NOT include the leading `-` in the parameter name (e.g., `help:print help message`, not `-help:...`)

### Requirement: -completion-bash outputs bash completion script
The system SHALL print a bash-compatible completion function to stdout when `-completion-bash` is passed.

#### Scenario: Valid bash script output
- **WHEN** the user runs `synapseq -completion-bash`
- **THEN** the system SHALL output a bash function named `_synapseq_completion` that uses `compgen` to complete flags, sourcing flag data from `synapseq -completion-args`

#### Scenario: Redirect to ~/.bashrc
- **WHEN** the user runs `synapseq -completion-bash >> ~/.bashrc`
- **THEN** the appended script SHALL enable tab-completion for `synapseq` in new bash sessions

### Requirement: -completion-zsh outputs zsh completion script
The system SHALL print a zsh-compatible completion function to stdout when `-completion-zsh` is passed.

#### Scenario: Valid zsh script output
- **WHEN** the user runs `synapseq -completion-zsh`
- **THEN** the system SHALL output a zsh completion definition using `compdef` that completes flags with descriptions from `synapseq -completion-args`

#### Scenario: Redirect to ~/.zshrc
- **WHEN** the user runs `synapseq -completion-zsh >> ~/.zshrc`
- **THEN** the appended script SHALL enable tab-completion for `synapseq` in new zsh sessions

### Requirement: Completion scripts prepend hyphen to flags
The bash and zsh completion scripts SHALL prepend `-` to flag names when presenting completions to the user.

#### Scenario: Hyphen added in bash completion
- **WHEN** the user types `synapseq ` and presses TAB in bash
- **THEN** the completion SHALL display flags with a leading `-` (e.g., `-help`, `-version`)

#### Scenario: Hyphen added in zsh completion
- **WHEN** the user types `synapseq ` and presses TAB in zsh
- **THEN** the completion SHALL display flags with a leading `-` (e.g., `-help`, `-version`)

### Requirement: CLIOptions extended with completion flags
The `CLIOptions` struct SHALL include boolean fields `CompletionBash`, `CompletionZsh`, and `CompletionArgs` to support the new flags.

#### Scenario: Struct fields present
- **WHEN** the `CLIOptions` struct is inspected
- **THEN** it SHALL contain `CompletionBash bool`, `CompletionZsh bool`, and `CompletionArgs bool` fields

#### Scenario: Flags parsed correctly
- **WHEN** the user passes `-completion-bash`, `-completion-zsh`, or `-completion-args`
- **THEN** the corresponding field in `CLIOptions` SHALL be set to `true`
