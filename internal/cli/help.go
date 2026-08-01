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
	writeAISection(writer)
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
	writeIndentedOptionsList(writer, "  ", options)
}

func writeIndentedOptionsList(writer io.Writer, indent string, options []helpOption) {
	for _, option := range options {
		fmt.Fprintf(writer, "%s%s%s\n", indent, FlagColumn(option.FlagText, option.ColumnWidth), option.Description)
	}
}

func writeAISection(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", Section("AI:"))
	fmt.Fprintf(writer, "  %s\n\n", Muted("Requires SYNAPSEQ_AI_API_KEY."))
	writeAISubsection(writer, "Command:", aiCommandHelpOptions())
	writeAISubsection(writer, "Options:", aiConfigurationHelpOptions())
	writeAISubsection(writer, "Environment:", aiEnvironmentHelpOptions())
	fmt.Fprintf(writer, "  %s %s\n\n", Label("Priority:"), Muted("-ai-* option > SYNAPSEQ_AI_* environment variable > CLI default"))
}

func writeAISubsection(writer io.Writer, label string, options []helpOption) {
	fmt.Fprintf(writer, "  %s\n", Label(label))
	writeIndentedOptionsList(writer, "    ", options)
	fmt.Fprintln(writer)
}

func quickStartExamples() []helpExample {
	return []helpExample{
		{Label: "1. Render audio", CommandText: "synapseq session.spsq", Description: "Generate session.wav in the current folder"},
		{Label: "2. Play audio", CommandText: "synapseq -play session.spsq", Description: "Play the sequence directly with ffplay"},
		{Label: "3. Export to MP3", CommandText: "synapseq session.spsq session.mp3", Description: "Export to MP3 with ffmpeg"},
		{Label: "4. Convert SBaGen", CommandText: "synapseq -sbg session.sbg [session.spsq]", Description: "Convert an SBaGen sequence to SPSQ"},
		{Label: "5. Generate SPSQ", CommandText: "synapseq -ai \"10 minutes of relaxation\" [session.spsq]", Description: "Generate an SPSQ sequence with an OpenAI-compatible model"},
	}
}

func commonHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-test", ColumnWidth: 19, Description: "Check syntax only"},
		{FlagText: "-sbg FILE [OUTPUT]", ColumnWidth: 19, Description: "Convert an SBaGen file to SPSQ"},
		{FlagText: "-dump", ColumnWidth: 19, Description: "Render JSON sequence data"},
		{FlagText: "-play", ColumnWidth: 19, Description: "Play audio using ffplay"},
		{FlagText: "-mp3", ColumnWidth: 19, Description: "Export to MP3 with ffmpeg"},
		{FlagText: "-quiet", ColumnWidth: 19, Description: "Suppress non-error output"},
		{FlagText: "-no-color", ColumnWidth: 19, Description: "Disable ANSI colors in CLI output"},
		{FlagText: "-version", ColumnWidth: 19, Description: "Show version information"},
		{FlagText: "-doctor", ColumnWidth: 19, Description: "Run the doctor check for tool dependencies"},
		{FlagText: "-completion-bash", ColumnWidth: 19, Description: "Generate bash completion script"},
		{FlagText: "-completion-zsh", ColumnWidth: 19, Description: "Generate zsh completion script"},
		{FlagText: "-help", ColumnWidth: 19, Description: "Show this help message"},
	}
}

func aiCommandHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-ai PROMPT [OUTPUT]", ColumnWidth: 28, Description: "Generate an SPSQ sequence; use - for standard output"},
	}
}

func aiConfigurationHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-ai-model MODEL", ColumnWidth: 28, Description: "Model name"},
		{FlagText: "-ai-base-url URL", ColumnWidth: 28, Description: "OpenAI-compatible API host"},
		{FlagText: "-ai-temperature VALUE", ColumnWidth: 28, Description: "Sampling temperature from 0 to 2"},
		{FlagText: "-ai-timeout DURATION", ColumnWidth: 28, Description: "Request timeout"},
	}
}

func aiEnvironmentHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "SYNAPSEQ_AI_API_KEY", ColumnWidth: 30, Description: "Required API key"},
		{FlagText: "SYNAPSEQ_AI_MODEL", ColumnWidth: 30, Description: "Model name; default gpt-4.1-mini"},
		{FlagText: "SYNAPSEQ_AI_BASE_URL", ColumnWidth: 30, Description: "OpenAI-compatible API host"},
		{FlagText: "SYNAPSEQ_AI_TEMPERATURE", ColumnWidth: 30, Description: "Sampling temperature; default 1"},
		{FlagText: "SYNAPSEQ_AI_TIMEOUT", ColumnWidth: 30, Description: "Request timeout; default 5m"},
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
		{FlagText: "-install-file-association", ColumnWidth: 30, Description: "Associate .spsq files and add .sbg conversion"},
		{FlagText: "-uninstall-file-association", ColumnWidth: 30, Description: "Remove .spsq association and .sbg conversion"},
	}
}
