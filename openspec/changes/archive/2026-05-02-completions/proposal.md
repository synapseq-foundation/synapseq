## Why

Add shell completion support for SynapSeq CLI to improve user experience when interacting with bash/zsh shells, eliminating the need to memorize all available flags.

## What Changes

- Add `-completion-bash` flag to output a bash completion script for users to append to `~/.bashrc`
- Add `-completion-zsh` flag to output a zsh completion script for users to append to `~/.zshrc`
- Add `-completion-args` flag to print all CLI parameters with concise descriptions in `{param}:{desc}` format
- Create `internal/cli/completion.go` containing a `map[string]string` of CLI flags to their truncated descriptions (max 50 characters, no colon allowed in description)
- Completion scripts will prepend `-` to all commands for correct autocomplete behavior
- Extend `CLIOptions` struct with new boolean fields for the three completion flags

## Capabilities

### New Capabilities
- `shell-completions`: Encompasses shell completion script generation and CLI argument metadata for bash/zsh autocomplete support

### Modified Capabilities

## Impact

- New file: `internal/cli/completion.go` (completion logic, flag metadata map, script generation)
- Modified file: `internal/cli` package (CLIOptions struct extended with new completion flags)
- Modified file: `cmd/synapseq` (flag parsing to handle new completion flags, dispatch logic to trigger completion output)
- No new external dependencies required
