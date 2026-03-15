/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
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
	color.NoColor = !enabled
}

// Title formats top-level headings.
func Title(text string) string {
	return color.New(color.FgCyan, color.Bold).Sprint(text)
}

// Section formats section headings.
func Section(text string) string {
	return color.New(color.FgYellow, color.Bold).Sprint(text)
}

// Command formats shell commands and paths.
func Command(text string) string {
	return color.New(color.FgGreen).Sprint(text)
}

// Flag formats command-line flags.
func Flag(text string) string {
	return color.New(color.FgMagenta, color.Bold).Sprint(text)
}

// Muted formats secondary explanatory text.
func Muted(text string) string {
	return color.New(color.FgHiBlack).Sprint(text)
}

// ErrorText formats error text.
func ErrorText(text string) string {
	return color.New(color.FgRed, color.Bold).Sprint(text)
}

// SuccessText formats success text.
func SuccessText(text string) string {
	return color.New(color.FgGreen, color.Bold).Sprint(text)
}

// Label formats field labels.
func Label(text string) string {
	return color.New(color.FgBlue, color.Bold).Sprint(text)
}

// Accent formats highlighted values.
func Accent(text string) string {
	return color.New(color.FgCyan).Sprint(text)
}

// Help prints the help message
func Help() {
	fmt.Fprintf(color.Output, "%s\n\n", Title(fmt.Sprintf("SynapSeq %s - Synapse-Sequenced Brainwave Generator", info.VERSION)))

	fmt.Fprintf(color.Output, "%s\n", Section("Usage:"))
	fmt.Fprintf(color.Output, "  %s\n\n", Command("synapseq [options] <input> [output]"))

	fmt.Fprintf(color.Output, "%s\n", Section("Quick start:"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq session.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Generate session.wav in the current folder"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq session.spsq relax.wav"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Generate a WAV file with a custom name"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -test session.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Validate the sequence without generating audio"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -play session.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Play the sequence directly with ffplay"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -preview session.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Generate session.html with a visual timeline preview"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq session.spsq relax.mp3"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Export to MP3 with ffmpeg"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -new meditation"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Create a new sequence file from the meditation template"))

	fmt.Fprintf(color.Output, "%s\n", Section("Input:"))
	fmt.Fprintf(color.Output, "  local file        %s\n", Command("path/to/sequence.spsq"))
	fmt.Fprintf(color.Output, "  URL               %s\n", Command("https://example.com/sequence.spsq"))
	fmt.Fprintf(color.Output, "  standard input    %s\n\n", Command("-"))

	fmt.Fprintf(color.Output, "%s\n", Section("Output:"))
	fmt.Fprintf(color.Output, "  omitted           %s\n", Muted("defaults to <input>.wav"))
	fmt.Fprintf(color.Output, "  WAV file          %s\n", Command("path/to/output.wav"))
	fmt.Fprintf(color.Output, "  MP3 file          %s\n", Command("path/to/output.mp3"))
	fmt.Fprintf(color.Output, "  standard output   %s\n\n", Muted("-   raw PCM (16-bit stereo)"))

	fmt.Fprintf(color.Output, "%s\n", Section("Most common options:"))
	fmt.Fprintf(color.Output, "  %s TYPE         Template type: meditation, focus, sleep, relaxation, example\n", Flag("-new"))
	fmt.Fprintf(color.Output, "  %s             Check syntax only\n", Flag("-test"))
	fmt.Fprintf(color.Output, "  %s          Render an HTML preview timeline\n", Flag("-preview"))
	fmt.Fprintf(color.Output, "  %s             Play audio using ffplay\n", Flag("-play"))
	fmt.Fprintf(color.Output, "  %s            Suppress non-error output\n", Flag("-quiet"))
	fmt.Fprintf(color.Output, "  %s         Disable ANSI colors in CLI output\n", Flag("-no-color"))
	fmt.Fprintf(color.Output, "  %s          Show version information\n", Flag("-version"))
	fmt.Fprintf(color.Output, "  %s             Show this help message\n\n", Flag("-help"))

	fmt.Fprintf(color.Output, "%s\n", Section("Hub:"))
	fmt.Fprintf(color.Output, "  %s\n\n", Muted("Run -hub-update once before using other -hub-* commands."))
	fmt.Fprintf(color.Output, "  %s                     List available sequences\n", Flag("-hub-list"))
	fmt.Fprintf(color.Output, "  %s WORD              Search the Hub\n", Flag("-hub-search"))
	fmt.Fprintf(color.Output, "  %s NAME                Show information about a sequence\n", Flag("-hub-info"))
	fmt.Fprintf(color.Output, "  %s NAME [DIR]      Download a sequence and dependencies\n", Flag("-hub-download"))
	fmt.Fprintf(color.Output, "  %s NAME [OUTPUT]        Download and generate in one step\n", Flag("-hub-get"))
	fmt.Fprintf(color.Output, "  %s                   Update the local Hub index\n", Flag("-hub-update"))
	fmt.Fprintf(color.Output, "  %s                    Clean up local Hub cache\n\n", Flag("-hub-clean"))

	fmt.Fprintf(color.Output, "%s\n", Section("Hub examples:"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-update"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-search calm-state"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-download calm-state"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-get calm-state calm-state.wav"))
	fmt.Fprintf(color.Output, "  %s\n\n", Command("synapseq -hub-get calm-state calm-state.mp3"))

	fmt.Fprintf(color.Output, "%s\n", Section("Advanced:"))
	fmt.Fprintf(color.Output, "  %s PATH   Path to ffmpeg executable\n", Flag("-ffmpeg-path"))
	fmt.Fprintf(color.Output, "  %s PATH   Path to ffplay executable\n\n", Flag("-ffplay-path"))

	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, "%s\n", Section("Windows-specific options:"))
		fmt.Fprintf(color.Output, "  %s    Associate .spsq files with SynapSeq\n", Flag("-install-file-association"))
		fmt.Fprintf(color.Output, "  %s  Remove .spsq file association\n\n", Flag("-uninstall-file-association"))
	}

	fmt.Fprintf(color.Output, "%s\n", Section("Docs:"))
	fmt.Fprintf(color.Output, "  %s\n", Command(info.DOC_URL))
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

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use -help flag for usage information.\n")
	}

	// General options
	fs.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	fs.StringVar(&opts.New, "new", "", "Template type: meditation, focus, sleep, relaxation, example")
	fs.BoolVar(&opts.Preview, "preview", false, "Render HTML preview timeline")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.NoColor, "no-color", false, "Disable ANSI colors in CLI output")
	fs.BoolVar(&opts.NoColor, "no-colors", false, "Disable ANSI colors in CLI output")
	fs.BoolVar(&opts.Test, "test", false, "Validate syntax without generating output")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")

	// Hub options
	fs.BoolVar(&opts.HubUpdate, "hub-update", false, "Update index of available sequences")
	fs.BoolVar(&opts.HubClean, "hub-clean", false, "Clean up local cache")
	fs.BoolVar(&opts.HubList, "hub-list", false, "List available sequences")
	fs.StringVar(&opts.HubSearch, "hub-search", "", "Search sequences")
	fs.StringVar(&opts.HubDownload, "hub-download", "", "Download sequence and dependencies")
	fs.StringVar(&opts.HubInfo, "hub-info", "", "Show information about a sequence")
	fs.StringVar(&opts.HubGet, "hub-get", "", "Get sequence")

	// External tool options
	fs.BoolVar(&opts.Play, "play", false, "Play audio using ffplay")
	fs.StringVar(&opts.FFmpegPath, "ffmpeg-path", "", "Path to ffmpeg executable")
	fs.StringVar(&opts.FFplayPath, "ffplay-path", "", "Path to ffplay executable")

	// Windows-specific options
	fs.BoolVar(&opts.InstallFileAssociation, "install-file-association", false, "Associate .spsq files with SynapSeq (Windows only)")
	fs.BoolVar(&opts.UninstallFileAssociation, "uninstall-file-association", false, "Remove .spsq file association (Windows only)")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, err
	}

	SetColorEnabled(!opts.NoColor)

	return opts, fs.Args(), err
}
