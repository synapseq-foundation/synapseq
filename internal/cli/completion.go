/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package cli

import (
	"fmt"
)

// completionFlags maps CLI flag names (without leading "-") to descriptions
// of 50 characters or fewer with no ":" character.
var completionFlags = map[string]string{
	"version":                    "Show version information",
	"preview":                    "Render HTML preview timeline",
	"quiet":                      "Enable quiet mode",
	"no-color":                   "Disable ANSI colors in CLI output",
	"test":                       "Validate syntax without generating output",
	"help":                       "Show help",
	"manual":                     "Show links to the canonical documentation",
	"remote-sync":                "Sync index of available sequences",
	"remote-clean":               "Clean up local cache",
	"remote-get":                 "Get remote sequence",
	"remote-list":                "List remote sequences",
	"remote-search":              "Search remote sequences",
	"remote-download":            "Download remote sequence",
	"remote-info":                "Show remote sequence information",
	"play":                       "Play audio using ffplay",
	"mp3":                        "Export to MP3 with ffmpeg",
	"install-file-association":   "Associate .spsq files (Windows only)",
	"uninstall-file-association": "Remove .spsq association (Windows only)",
	"new":                        "Generate template (meditation, focus, etc)",
	"ffmpeg-path":                "Path to ffmpeg executable",
	"ffplay-path":                "Path to ffplay executable",
	"doctor":                     "Check environment for required tools",
	"completion-bash":            "Print bash completion script",
	"completion-zsh":             "Print zsh completion script",
	"completion-args":            "Print flags with descriptions",
}

// PrintCompletionArgs prints all CLI flags with descriptions in {param}:{desc} format.
func PrintCompletionArgs() {
	for flag, desc := range completionFlags {
		fmt.Printf("%s:%s\n", flag, desc)
	}
}

// PrintBashCompletion prints a bash completion script to stdout.
func PrintBashCompletion() {
	script := `# SynapSeq bash completion
_synapseq_completion() {
    local cur opts

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"

    if [[ "$cur" == -* ]]; then
        opts=$(
            "$(basename "${COMP_WORDS[0]}")" -completion-args 2>/dev/null \
            | sed 's/:.*//' \
            | sed 's/^/-/' \
            | tr '\n' ' '
        )

        COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
        return 0
    fi

    COMPREPLY=( $(compgen -f -- "$cur") )
    return 0
}

complete -F _synapseq_completion synapseq
`
	fmt.Print(script)
}

// PrintZshCompletion prints a zsh completion script to stdout.
func PrintZshCompletion() {
	script := `# SynapSeq zsh completion
_synapseq_completion() {
    local -a opts

    if [[ "$words[CURRENT]" == -* ]]; then
        opts=(
            $(
                "$(basename "$words[1]")" -completion-args 2>/dev/null \
                | sed 's/:.*//' \
                | sed 's/^/-/'
            )
        )

        _describe 'synapseq flags' opts
        return 0
    fi

    _files
}

compdef _synapseq_completion synapseq
`
	fmt.Print(script)
}
