## 1. Completion Metadata and Logic

- [x] 1.1 Create `internal/cli/completion.go` with `completionFlags` map[string]string containing all CLI flags and their descriptions (≤50 chars, no `:`)
- [x] 1.2 Implement `PrintCompletionArgs()` function that outputs `{param}:{desc}` lines to stdout, excluding leading `-` in param names
- [x] 1.3 Implement `PrintBashCompletion()` function that outputs a bash completion script using `_synapseq_completion` function and `synapseq -completion-args` as data source
- [x] 1.4 Implement `PrintZshCompletion()` function that outputs a zsh completion script using `compdef` and `synapseq -completion-args` as data source
- [x] 1.5 Ensure both completion scripts prepend `-` to flag names when displaying completions

## 2. Extend CLIOptions Struct

- [x] 2.1 Add `CompletionBash bool`, `CompletionZsh bool`, and `CompletionArgs bool` fields to `CLIOptions` struct in `internal/cli` package

## 3. CLI Flag Parsing and Dispatch

- [x] 3.1 Register `-completion-bash`, `-completion-zsh`, and `-completion-args` flags in the CLI flag parsing logic in `cmd/synapseq`
- [x] 3.2 Add dispatch logic to call the appropriate completion function when any completion flag is set, exiting afterward without running other commands

## 4. Verification

- [x] 4.1 Run `synapseq -completion-args` and verify output format `{param}:{desc}` with no leading `-` on params
- [x] 4.2 Run `synapseq -completion-bash` and verify valid bash script is printed
- [x] 4.3 Run `synapseq -completion-zsh` and verify valid zsh script is printed
- [x] 4.4 Run `make test` to ensure all existing tests pass
