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
	"version":                     "Show version information",
	"preview":                     "Render HTML preview timeline",
	"quiet":                       "Enable quiet mode",
	"no-color":                    "Disable ANSI colors in CLI output",
	"test":                        "Validate syntax without generating output",
	"help":                        "Show help",
	"manual":                      "Show links to the canonical documentation",
	"hub-update":                  "Update index of available sequences",
	"hub-clean":                   "Clean up local cache",
	"hub-get":                     "Get sequence",
	"hub-list":                    "List available sequences",
	"hub-search":                  "Search sequences",
	"hub-download":                "Download sequence and dependencies",
	"hub-info":                    "Show information about a sequence",
	"play":                        "Play audio using ffplay",
	"mp3":                         "Export to MP3 with ffmpeg",
	"install-file-association":    "Associate .spsq files (Windows only)",
	"uninstall-file-association":  "Remove .spsq association (Windows only)",
	"new":                         "Generate template (meditation, focus, etc)",
	"ffmpeg-path":                 "Path to ffmpeg executable",
	"ffplay-path":                 "Path to ffplay executable",
	"doctor":                      "Check environment for required tools",
	"completion-bash":             "Print bash completion script",
	"completion-zsh":              "Print zsh completion script",
	"completion-args":             "Print flags with descriptions",
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
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    if [[ ${cur} == -* ]]; then
        opts=$($(basename ${COMP_WORDS[0]}) -completion-args 2>/dev/null | sed 's/:.*//' | sed 's/^/-/' | tr '\n' ' ')
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    fi
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
    opts=($($(basename $words[1]) -completion-args 2>/dev/null | sed 's/:.*//' | sed 's/^/-/' ))
    _describe 'synapseq flags' opts
}

compdef _synapseq_completion synapseq
`
	fmt.Print(script)
}
