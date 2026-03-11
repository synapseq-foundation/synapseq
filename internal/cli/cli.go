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

// Help prints the help message
func Help() {
	fmt.Printf("SynapSeq %s - Synapse-Sequenced Brainwave Generator\n\n", info.VERSION)

	fmt.Printf("Usage:\n")
	fmt.Printf("  synapseq [options] <input> [output]\n\n")

	fmt.Printf("Quick start:\n")
	fmt.Printf("  synapseq session.spsq\n")
	fmt.Printf("    Generate session.wav in the current folder\n\n")
	fmt.Printf("  synapseq session.spsq relax.wav\n")
	fmt.Printf("    Generate a WAV file with a custom name\n\n")
	fmt.Printf("  synapseq -test session.spsq\n")
	fmt.Printf("    Validate the sequence without generating audio\n\n")
	fmt.Printf("  synapseq -play session.spsq\n")
	fmt.Printf("    Play the sequence directly with ffplay\n\n")
	fmt.Printf("  synapseq -preview session.spsq\n")
	fmt.Printf("    Generate session.html with a visual timeline preview\n\n")
	fmt.Printf("  synapseq session.spsq relax.mp3\n")
	fmt.Printf("    Export to MP3 with ffmpeg\n\n")
	fmt.Printf("  synapseq -new meditation\n")
	fmt.Printf("    Create a new sequence file from the meditation template\n\n")

	fmt.Printf("Input:\n")
	fmt.Printf("  local file        path/to/sequence.spsq\n")
	fmt.Printf("  URL               https://example.com/sequence.spsq\n")
	fmt.Printf("  standard input    -\n\n")

	fmt.Printf("Output:\n")
	fmt.Printf("  omitted           defaults to <input>.wav\n")
	fmt.Printf("  WAV file          path/to/output.wav\n")
	fmt.Printf("  MP3 file          path/to/output.mp3\n")
	fmt.Printf("  standard output   -   raw PCM (16-bit stereo)\n\n")

	fmt.Printf("Most common options:\n")
	fmt.Printf("  -new TYPE         Template type: meditation, focus, sleep, relaxation, example\n")
	fmt.Printf("  -test             Check syntax only\n")
	fmt.Printf("  -preview          Render an HTML preview timeline\n")
	fmt.Printf("  -play             Play audio using ffplay\n")
	fmt.Printf("  -quiet            Suppress non-error output\n")
	fmt.Printf("  -version          Show version information\n")
	fmt.Printf("  -help             Show this help message\n\n")

	fmt.Printf("Hub:\n")
	fmt.Printf("  Run -hub-update once before using other -hub-* commands.\n\n")
	fmt.Printf("  -hub-list                     List available sequences\n")
	fmt.Printf("  -hub-search WORD              Search the Hub\n")
	fmt.Printf("  -hub-info NAME                Show information about a sequence\n")
	fmt.Printf("  -hub-download NAME [DIR]      Download a sequence and dependencies\n")
	fmt.Printf("  -hub-get NAME [OUTPUT]        Download and generate in one step\n")
	fmt.Printf("  -hub-update                   Update the local Hub index\n")
	fmt.Printf("  -hub-clean                    Clean up local Hub cache\n\n")

	fmt.Printf("Hub examples:\n")
	fmt.Printf("  synapseq -hub-update\n")
	fmt.Printf("  synapseq -hub-search calm-state\n")
	fmt.Printf("  synapseq -hub-download calm-state\n")
	fmt.Printf("  synapseq -hub-get calm-state calm-state.wav\n")
	fmt.Printf("  synapseq -hub-get calm-state calm-state.mp3\n\n")

	fmt.Printf("Advanced:\n")
	fmt.Printf("  -ffmpeg-path PATH   Path to ffmpeg executable\n")
	fmt.Printf("  -ffplay-path PATH   Path to ffplay executable\n\n")

	if runtime.GOOS == "windows" {
		fmt.Printf("Windows-specific options:\n")
		fmt.Printf("  -install-file-association    Associate .spsq files with SynapSeq\n")
		fmt.Printf("  -uninstall-file-association  Remove .spsq file association\n\n")
	}

	fmt.Printf("Docs:\n")
	fmt.Printf("  %s\n", info.DOC_URL)
}

// ShowVersion prints the version information
func ShowVersion() {
	fmt.Printf("SynapSeq %s (%s) built %s for %s/%s\n",
		info.VERSION,
		info.GIT_COMMIT,
		info.BUILD_DATE,
		runtime.GOOS,
		runtime.GOARCH,
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

	return opts, fs.Args(), err
}
