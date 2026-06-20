// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cli

import (
	"fmt"
	"io"
	"runtime"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
)

type helpExample struct {
	Label       string
	CommandText string
	Description string
}

type helpOption struct {
	FlagText    string
	ColumnWidth int
	Description string
}

// Help prints the help message
func Help() {
	writer := color.Output
	writeHelpHeader(writer)
	writeUsageSection(writer)
	writeExamplesSection(writer, "Quick start:", quickStartExamples())
	writeInputSection(writer)
	writeOutputSection(writer)
	writeOptionsSection(writer, "Most common options:", commonHelpOptions())
	writeMutedLeadSection(writer, "Remote:", "Run -sync first to initialize the local Remote index.")
	writeOptionsList(writer, remoteHelpOptions())
	fmt.Fprintln(writer)
	writeOptionsSection(writer, "Advanced:", advancedHelpOptions())

	if runtime.GOOS == "windows" {
		writeOptionsSection(writer, "Windows-specific options:", windowsHelpOptions())
	}
}

// ShowVersion prints the version information
func ShowVersion() {
	fmt.Fprintf(
		color.Output,
		"%s %s %s %s %s\n",
		Title("SynapSeq"),
		Accent(info.VERSION),
		Muted(fmt.Sprintf("(%s)", info.GIT_COMMIT)),
		Label("built"),
		Command(fmt.Sprintf("%s for %s/%s", info.BUILD_DATE, runtime.GOOS, runtime.GOARCH)),
	)
}

func writeHelpHeader(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n\n", Title("SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment"))
}

func writeUsageSection(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", Section("Usage:"))
	fmt.Fprintf(writer, "  %s\n\n", Command("synapseq [options] <input> [output]"))
}

func writeExamplesSection(writer io.Writer, title string, examples []helpExample) {
	fmt.Fprintf(writer, "%s\n", Section(title))
	for _, example := range examples {
		if example.Label != "" {
			fmt.Fprintf(writer, "  %s\n", Label(example.Label))
			fmt.Fprintf(writer, "     %s\n", Command(example.CommandText))
			fmt.Fprintf(writer, "       %s\n", Muted(example.Description))
			continue
		}

		fmt.Fprintf(writer, "  %s\n", Command(example.CommandText))
		if example.Description != "" {
			fmt.Fprintf(writer, "    %s\n", Muted(example.Description))
		}
	}
	fmt.Fprintln(writer)
}

func writeInputSection(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", Section("Input:"))
	fmt.Fprintf(writer, "  local file        %s\n", Command("path/to/sequence.spsq"))
	fmt.Fprintf(writer, "  URL               %s\n", Command("https://example.com/sequence.spsq"))
	fmt.Fprintf(writer, "  standard input    %s\n\n", Command("-"))
}

func writeOutputSection(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", Section("Output:"))
	fmt.Fprintf(writer, "  omitted           %s\n", Muted("defaults to <input>.wav"))
	fmt.Fprintf(writer, "  WAV file          %s\n", Command("path/to/output.wav"))
	fmt.Fprintf(writer, "  MP3 file          %s\n", Command("path/to/output.mp3"))
	fmt.Fprintf(writer, "  JSON file         %s\n", Command("path/to/output.json"))
	fmt.Fprintf(writer, "  standard output   %s\n\n", Muted("-   raw PCM or JSON depending on mode"))
}

func writeOptionsSection(writer io.Writer, title string, options []helpOption) {
	fmt.Fprintf(writer, "%s\n", Section(title))
	writeOptionsList(writer, options)
	fmt.Fprintln(writer)
}

func writeMutedLeadSection(writer io.Writer, title, lead string) {
	fmt.Fprintf(writer, "%s\n", Section(title))
	fmt.Fprintf(writer, "  %s\n\n", Muted(lead))
}

func writeOptionsList(writer io.Writer, options []helpOption) {
	for _, option := range options {
		fmt.Fprintf(writer, "  %s%s\n", FlagColumn(option.FlagText, option.ColumnWidth), option.Description)
	}
}

func quickStartExamples() []helpExample {
	return []helpExample{
		{Label: "1. Render audio", CommandText: "synapseq session.spsq", Description: "Generate session.wav in the current folder"},
		{Label: "2. Play audio", CommandText: "synapseq -play session.spsq", Description: "Play the sequence directly with ffplay"},
		{Label: "3. Export to MP3", CommandText: "synapseq session.spsq session.mp3", Description: "Export to MP3 with ffmpeg"},
	}
}

func commonHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-test", ColumnWidth: 18, Description: "Check syntax only"},
		{FlagText: "-dump", ColumnWidth: 18, Description: "Render JSON sequence data"},
		{FlagText: "-play", ColumnWidth: 18, Description: "Play audio using ffplay"},
		{FlagText: "-mp3", ColumnWidth: 18, Description: "Export to MP3 with ffmpeg"},
		{FlagText: "-quiet", ColumnWidth: 18, Description: "Suppress non-error output"},
		{FlagText: "-no-color", ColumnWidth: 18, Description: "Disable ANSI colors in CLI output"},
		{FlagText: "-version", ColumnWidth: 18, Description: "Show version information"},
		{FlagText: "-doctor", ColumnWidth: 18, Description: "Run the doctor check for tool dependencies"},
		{FlagText: "-completion-bash", ColumnWidth: 18, Description: "Generate bash completion script"},
		{FlagText: "-completion-zsh", ColumnWidth: 18, Description: "Generate zsh completion script"},
		{FlagText: "-help", ColumnWidth: 18, Description: "Show this help message"},
	}
}

func remoteHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-sync", ColumnWidth: 28, Description: "Sync the local Remote index"},
		{FlagText: "-list", ColumnWidth: 28, Description: "List available remote sequences"},
		{FlagText: "-search WORD", ColumnWidth: 28, Description: "Search remote sequences"},
		{FlagText: "-info NAME", ColumnWidth: 28, Description: "Show information about a remote sequence"},
		{FlagText: "-download NAME [DIR]", ColumnWidth: 28, Description: "Download a remote sequence"},
		{FlagText: "-get NAME [OUTPUT]", ColumnWidth: 28, Description: "Download and generate in one step"},
		{FlagText: "-clean", ColumnWidth: 28, Description: "Clean up local Remote cache"},
	}
}

func advancedHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-ffmpeg-path PATH", ColumnWidth: 22, Description: "Path to ffmpeg executable"},
		{FlagText: "-ffplay-path PATH", ColumnWidth: 22, Description: "Path to ffplay executable"},
	}
}

func windowsHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-install-file-association", ColumnWidth: 30, Description: "Associate .spsq files with SynapSeq"},
		{FlagText: "-uninstall-file-association", ColumnWidth: 30, Description: "Remove .spsq file association"},
	}
}
