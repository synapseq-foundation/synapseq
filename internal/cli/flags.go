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
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if hasNoColorArg(os.Args[1:]) {
		SetColorEnabled(false)
	}

	fs.Usage = func() {}
	bindFlags(fs, opts)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, formatFlagParseError(fs, err)
	}

	SetColorEnabled(!opts.NoColor)

	return opts, fs.Args(), err
}

func bindFlags(fs *flag.FlagSet, opts *CLIOptions) {
	fs.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	fs.StringVar(&opts.New, "new", "", "Template type: meditation, focus, sleep, relaxation, example")
	fs.BoolVar(&opts.Preview, "preview", false, "Render HTML preview timeline")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.NoColor, "no-color", false, "Disable ANSI colors in CLI output")
	fs.BoolVar(&opts.Test, "test", false, "Validate syntax without generating output")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")
	fs.BoolVar(&opts.ShowManual, "manual", false, "Show the full manual")

	fs.BoolVar(&opts.HubUpdate, "hub-update", false, "Update index of available sequences")
	fs.BoolVar(&opts.HubClean, "hub-clean", false, "Clean up local cache")
	fs.BoolVar(&opts.HubList, "hub-list", false, "List available sequences")
	fs.StringVar(&opts.HubSearch, "hub-search", "", "Search sequences")
	fs.StringVar(&opts.HubDownload, "hub-download", "", "Download sequence and dependencies")
	fs.StringVar(&opts.HubInfo, "hub-info", "", "Show information about a sequence")
	fs.StringVar(&opts.HubGet, "hub-get", "", "Get sequence")

	fs.BoolVar(&opts.Play, "play", false, "Play audio using ffplay")
	fs.BoolVar(&opts.Mp3, "mp3", false, "Export to MP3 with ffmpeg")
	fs.StringVar(&opts.FFmpegPath, "ffmpeg-path", "", "Path to ffmpeg executable")
	fs.StringVar(&opts.FFplayPath, "ffplay-path", "", "Path to ffplay executable")

	fs.BoolVar(&opts.InstallFileAssociation, "install-file-association", false, "Associate .spsq files with SynapSeq (Windows only)")
	fs.BoolVar(&opts.UninstallFileAssociation, "uninstall-file-association", false, "Remove .spsq file association (Windows only)")
}

func hasNoColorArg(args []string) bool {
	for _, arg := range args {
		if arg == "-no-color" {
			return true
		}
	}
	return false
}

func formatFlagParseError(fs *flag.FlagSet, err error) error {
	if err == nil {
		return nil
	}

	message := err.Error()
	knownFlags := flagNames(fs)

	switch {
	case strings.HasPrefix(message, "flag provided but not defined: "):
		found := strings.TrimSpace(strings.TrimPrefix(message, "flag provided but not defined: "))
		diagnostic := diag.Validation("unknown command-line flag").WithFound(found).WithHint("use -help to see valid command-line options")
		if suggestion, ok := diag.ClosestMatch(found, knownFlags, diag.DefaultSuggestionDistance(found)); ok {
			diagnostic.WithSuggestion(fmt.Sprintf("did you mean %q?", suggestion))
		}
		return diagnostic
	case strings.HasPrefix(message, "flag needs an argument: "):
		found := strings.TrimSpace(strings.TrimPrefix(message, "flag needs an argument: "))
		return diag.Validation("missing value for command-line flag").WithFound(found).WithHint("pass a value for this flag or use -help to review its syntax")
	case strings.HasPrefix(message, "invalid boolean value "):
		return diag.Validation("invalid value for command-line flag").WithHint(message + "; use -help to review accepted flag values")
	default:
		return diag.Validation("invalid command-line arguments").WithHint(message + "; use -help for usage information")
	}
}

func flagNames(fs *flag.FlagSet) []string {
	flags := make([]string, 0, 16)
	fs.VisitAll(func(f *flag.Flag) {
		flags = append(flags, "-"+f.Name)
	})
	return flags
}