// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cli

import (
	style "github.com/synapseq-foundation/synapseq/v4/internal/textstyle"
)

type SpecialCommandKind string

const (
	SpecialCommandNone                     SpecialCommandKind = ""
	SpecialCommandShowVersion              SpecialCommandKind = "show-version"
	SpecialCommandRemoteSync               SpecialCommandKind = "remote-sync"
	SpecialCommandRemoteClean              SpecialCommandKind = "remote-clean"
	SpecialCommandRemoteGet                SpecialCommandKind = "remote-get"
	SpecialCommandRemoteList               SpecialCommandKind = "remote-list"
	SpecialCommandRemoteSearch             SpecialCommandKind = "remote-search"
	SpecialCommandRemoteDownload           SpecialCommandKind = "remote-download"
	SpecialCommandRemoteInfo               SpecialCommandKind = "remote-info"
	SpecialCommandInstallFileAssociation   SpecialCommandKind = "install-file-association"
	SpecialCommandUninstallFileAssociation SpecialCommandKind = "uninstall-file-association"
	SpecialCommandGenerateTemplate         SpecialCommandKind = "generate-template"
	SpecialCommandDoctor                   SpecialCommandKind = "doctor"
)

type SpecialCommand struct {
	Kind        SpecialCommandKind
	OptionalArg string
}

// CLIOptions holds command-line options
type CLIOptions struct {
	// Show version information and exit
	ShowVersion bool
	// New starter sequence template type
	New string
	// Preview mode, renders HTML timeline instead of audio
	Preview bool
	// Quiet mode, suppress non-error output
	Quiet bool
	// Test mode, validate syntax without generating output
	Test bool
	// Show help message and exit
	ShowHelp bool
	// Disable ANSI colors in CLI output
	NoColor bool
	// Windows file association installation
	InstallFileAssociation bool
	// Clean Windows file association removal
	UninstallFileAssociation bool
	// Play (with ffplay)
	Play bool
	// Remote sync index of available sequences
	RemoteSync bool
	// Remote clean up local cache
	RemoteClean bool
	// Remote list available sequences
	RemoteList bool
	// Remote search sequences
	RemoteSearch string
	// Remote download sequence
	RemoteDownload string
	// Remote info of sequence
	RemoteInfo string
	// Remote get sequence
	RemoteGet string
	// Mp3 export with ffmpeg
	Mp3 bool
	// Path to ffplay executable
	FFplayPath string
	// Path to ffmpeg executable
	FFmpegPath string
	// Path to ffprobe executable
	FFprobePath string
	// Show doctor diagnostic information
	ShowDoctor bool
	// Print bash completion script
	CompletionBash bool
	// Print zsh completion script
	CompletionZsh bool
	// Print completion args (param:desc format)
	CompletionArgs bool
}

func init() {
	SetColorEnabled(true)
}

// SetColorEnabled enables or disables ANSI colors for CLI output.
func SetColorEnabled(enabled bool) {
	style.SetColorEnabled(enabled)
}

// Title formats top-level headings.
func Title(text string) string {
	return style.Title(text)
}

// Section formats section headings.
func Section(text string) string {
	return style.Section(text)
}

// Command formats shell commands and paths.
func Command(text string) string {
	return style.Command(text)
}

// Flag formats command-line flags.
func Flag(text string) string {
	return style.Flag(text)
}

func FlagColumn(text string, width int) string {
	return style.FlagColumn(text, width)
}

// Muted formats secondary explanatory text.
func Muted(text string) string {
	return style.Muted(text)
}

// ErrorText formats error text.
func ErrorText(text string) string {
	return style.ErrorText(text)
}

// SuccessText formats success text.
func SuccessText(text string) string {
	return style.SuccessText(text)
}

// Label formats field labels.
func Label(text string) string {
	return style.Label(text)
}

// Accent formats highlighted values.
func Accent(text string) string {
	return style.Accent(text)
}
