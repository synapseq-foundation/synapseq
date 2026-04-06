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
	style "github.com/synapseq-foundation/synapseq/v4/internal/textstyle"
)

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
	// Show full manual and exit
	ShowManual bool
	// Disable ANSI colors in CLI output
	NoColor bool
	// Windows file association installation
	InstallFileAssociation bool
	// Clean Windows file association removal
	UninstallFileAssociation bool
	// Play (with ffplay)
	Play bool
	// Hub update index of available sequences
	HubUpdate bool
	// Hub clean up local cache
	HubClean bool
	// Hub list available sequences
	HubList bool
	// Hub search sequences
	HubSearch string
	// Hub download sequences
	HubDownload string
	// Hub info of sequence
	HubInfo string
	// Hub get sequence
	HubGet string
	// Mp3 export with ffmpeg
	Mp3 bool
	// Path to ffplay executable
	FFplayPath string
	// Path to ffmpeg executable
	FFmpegPath string
	// Path to ffprobe executable
	FFprobePath string
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
