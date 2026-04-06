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
	"io"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
)

const (
	syntaxDocURL       = info.REPOSITORY + "/blob/main/SYNTAX.md"
	architectureDocURL = info.REPOSITORY + "/blob/main/ARCHITECTURE.md"
	contributingDocURL = info.REPOSITORY + "/blob/main/CONTRIBUTING.md"
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

type helpLink struct {
	Target      string
	Description string
}

// Help prints the help message
func Help() {
	writer := color.Output
	writeHelpHeader(writer)
	writeUsageSection(writer)
	writeExamplesSection(writer, "Quick start:", quickStartExamples())
	writeExamplesSection(writer, "Next steps:", nextStepExamples())
	writeInputSection(writer)
	writeOutputSection(writer)
	writeOptionsSection(writer, "Most common options:", commonHelpOptions())
	writeMutedLeadSection(writer, "Hub:", "Run -hub-update first to initialize the local Hub index.")
	writeOptionsList(writer, hubHelpOptions())
	fmt.Fprintln(writer)
	writeCommandListSection(writer, "Hub quick start:", hubQuickStartCommands())
	writeOptionsSection(writer, "Advanced:", advancedHelpOptions())

	if runtime.GOOS == "windows" {
		writeOptionsSection(writer, "Windows-specific options:", windowsHelpOptions())
	}

	writeLinkSection(writer)
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

// ShowManual prints documentation links for users who discover the project from the CLI.
func ShowManual() {
	writer := color.Output

	fmt.Fprintf(writer, "%s\n\n", Title("SynapSeq Documentation"))
	fmt.Fprintf(writer, "  %s\n", Muted("The built-in manual was retired to avoid duplicated documentation."))
	fmt.Fprintf(writer, "  %s\n\n", Muted("Use the documents below as the single source of truth."))

	lines := []struct {
		label string
		url   string
		desc  string
	}{
		{label: "Syntax reference", url: syntaxDocURL, desc: "Full .spsq and .spsc language reference, examples, and semantic rules"},
		{label: "Architecture guide", url: architectureDocURL, desc: "Package boundaries, runtime flow, and design trade-offs"},
		{label: "Contribution guide", url: contributingDocURL, desc: "How to contribute changes and where to read first"},
	}

	for _, line := range lines {
		fmt.Fprintf(writer, "  %s\n", Label(line.label))
		fmt.Fprintf(writer, "    %s\n", Command(line.url))
		fmt.Fprintf(writer, "      %s\n", Muted(line.desc))
	}

	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "  %s\n", Label("Tip"))
	fmt.Fprintf(writer, "    %s\n", Muted(strings.Join([]string{
		"Keep -help for command usage,",
		"use this command for documentation links,",
		"and treat the repository documents as canonical.",
	}, " ")))
}

func writeHelpHeader(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n\n", Title(fmt.Sprintf("SynapSeq %s - Text-Driven Audio Sequencer for Brainwave Entrainment", info.VERSION)))
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
	fmt.Fprintf(writer, "  standard output   %s\n\n", Muted("-   raw PCM (16-bit stereo)"))
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

func writeCommandListSection(writer io.Writer, title string, commands []string) {
	fmt.Fprintf(writer, "%s\n", Section(title))
	for _, commandText := range commands {
		fmt.Fprintf(writer, "  %s\n", Command(commandText))
	}
	fmt.Fprintln(writer)
}

func writeLinkSection(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", Section("For more information:"))
	for _, link := range moreInfoLinks() {
		fmt.Fprintf(writer, "  %s\n", Command(link.Target))
		fmt.Fprintf(writer, "    %s\n", Muted(link.Description))
	}
}

func quickStartExamples() []helpExample {
	return []helpExample{
		{Label: "1. Create a starter file", CommandText: "synapseq -new meditation starter.spsq", Description: "Create starter.spsq from the meditation template"},
		{Label: "2. Render audio", CommandText: "synapseq starter.spsq", Description: "Generate starter.wav in the current folder"},
		{Label: "Available templates", Description: "meditation, focus, sleep, relaxation, example"},
	}
}

func nextStepExamples() []helpExample {
	return []helpExample{
		{CommandText: "synapseq -test starter.spsq", Description: "Validate syntax and semantics without generating audio"},
		{CommandText: "synapseq -preview starter.spsq", Description: "Generate starter.html with a visual timeline preview"},
		{CommandText: "synapseq -play starter.spsq", Description: "Play the sequence directly with ffplay"},
		{CommandText: "synapseq starter.spsq starter.mp3", Description: "Export to MP3 with ffmpeg"},
		{CommandText: "synapseq -manual", Description: "Show links to syntax, architecture, and contribution docs"},
	}
}

func commonHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-new TYPE", ColumnWidth: 18, Description: "Template type: meditation, focus, sleep, relaxation, example"},
		{FlagText: "-test", ColumnWidth: 18, Description: "Check syntax only"},
		{FlagText: "-preview", ColumnWidth: 18, Description: "Render an HTML preview timeline"},
		{FlagText: "-play", ColumnWidth: 18, Description: "Play audio using ffplay"},
		{FlagText: "-mp3", ColumnWidth: 18, Description: "Export to MP3 with ffmpeg"},
		{FlagText: "-quiet", ColumnWidth: 18, Description: "Suppress non-error output"},
		{FlagText: "-no-color", ColumnWidth: 18, Description: "Disable ANSI colors in CLI output"},
		{FlagText: "-manual", ColumnWidth: 18, Description: "Show links to the canonical docs"},
		{FlagText: "-version", ColumnWidth: 18, Description: "Show version information"},
		{FlagText: "-help", ColumnWidth: 18, Description: "Show this help message"},
	}
}

func hubHelpOptions() []helpOption {
	return []helpOption{
		{FlagText: "-hub-update", ColumnWidth: 24, Description: "Update the local Hub index"},
		{FlagText: "-hub-list", ColumnWidth: 24, Description: "List available sequences"},
		{FlagText: "-hub-search WORD", ColumnWidth: 24, Description: "Search the Hub"},
		{FlagText: "-hub-info NAME", ColumnWidth: 24, Description: "Show information about a sequence"},
		{FlagText: "-hub-download NAME [DIR]", ColumnWidth: 24, Description: "Download a sequence and dependencies"},
		{FlagText: "-hub-get NAME [OUTPUT]", ColumnWidth: 24, Description: "Download and generate in one step"},
		{FlagText: "-hub-clean", ColumnWidth: 24, Description: "Clean up local Hub cache"},
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

func hubQuickStartCommands() []string {
	return []string{
		"synapseq -hub-update",
		"synapseq -hub-list",
		"synapseq -hub-search calm-state",
		"synapseq -hub-get calm-state calm-state.wav",
		"synapseq -hub-get calm-state calm-state.mp3",
	}
}

func moreInfoLinks() []helpLink {
	return []helpLink{
		{Target: "synapseq -manual", Description: "Show links to syntax, architecture, and contribution docs"},
		{Target: syntaxDocURL, Description: "Canonical syntax reference for .spsq and .spsc"},
		{Target: architectureDocURL, Description: "Architecture guide, package boundaries, and design trade-offs"},
		{Target: "https://synapseq.org", Description: "Visit the website for documentation, examples, and the latest updates"},
	}
}
