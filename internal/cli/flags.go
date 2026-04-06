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

type boolFlagBinding struct {
	Name  string
	Usage string
	Bind  func(*CLIOptions) *bool
}

type stringFlagBinding struct {
	Name  string
	Usage string
	Bind  func(*CLIOptions) *string
}

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
	for _, binding := range boolFlagBindings() {
		fs.BoolVar(binding.Bind(opts), binding.Name, false, binding.Usage)
	}

	for _, binding := range stringFlagBindings() {
		fs.StringVar(binding.Bind(opts), binding.Name, "", binding.Usage)
	}
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

func boolFlagBindings() []boolFlagBinding {
	return []boolFlagBinding{
		{Name: "version", Usage: "Show version information", Bind: func(opts *CLIOptions) *bool { return &opts.ShowVersion }},
		{Name: "preview", Usage: "Render HTML preview timeline", Bind: func(opts *CLIOptions) *bool { return &opts.Preview }},
		{Name: "quiet", Usage: "Enable quiet mode", Bind: func(opts *CLIOptions) *bool { return &opts.Quiet }},
		{Name: "no-color", Usage: "Disable ANSI colors in CLI output", Bind: func(opts *CLIOptions) *bool { return &opts.NoColor }},
		{Name: "test", Usage: "Validate syntax without generating output", Bind: func(opts *CLIOptions) *bool { return &opts.Test }},
		{Name: "help", Usage: "Show help", Bind: func(opts *CLIOptions) *bool { return &opts.ShowHelp }},
		{Name: "manual", Usage: "Show the full manual", Bind: func(opts *CLIOptions) *bool { return &opts.ShowManual }},
		{Name: "hub-update", Usage: "Update index of available sequences", Bind: func(opts *CLIOptions) *bool { return &opts.HubUpdate }},
		{Name: "hub-clean", Usage: "Clean up local cache", Bind: func(opts *CLIOptions) *bool { return &opts.HubClean }},
		{Name: "hub-list", Usage: "List available sequences", Bind: func(opts *CLIOptions) *bool { return &opts.HubList }},
		{Name: "play", Usage: "Play audio using ffplay", Bind: func(opts *CLIOptions) *bool { return &opts.Play }},
		{Name: "mp3", Usage: "Export to MP3 with ffmpeg", Bind: func(opts *CLIOptions) *bool { return &opts.Mp3 }},
		{Name: "install-file-association", Usage: "Associate .spsq files with SynapSeq (Windows only)", Bind: func(opts *CLIOptions) *bool { return &opts.InstallFileAssociation }},
		{Name: "uninstall-file-association", Usage: "Remove .spsq file association (Windows only)", Bind: func(opts *CLIOptions) *bool { return &opts.UninstallFileAssociation }},
	}
}

func stringFlagBindings() []stringFlagBinding {
	return []stringFlagBinding{
		{Name: "new", Usage: "Template type: meditation, focus, sleep, relaxation, example", Bind: func(opts *CLIOptions) *string { return &opts.New }},
		{Name: "hub-search", Usage: "Search sequences", Bind: func(opts *CLIOptions) *string { return &opts.HubSearch }},
		{Name: "hub-download", Usage: "Download sequence and dependencies", Bind: func(opts *CLIOptions) *string { return &opts.HubDownload }},
		{Name: "hub-info", Usage: "Show information about a sequence", Bind: func(opts *CLIOptions) *string { return &opts.HubInfo }},
		{Name: "hub-get", Usage: "Get sequence", Bind: func(opts *CLIOptions) *string { return &opts.HubGet }},
		{Name: "ffmpeg-path", Usage: "Path to ffmpeg executable", Bind: func(opts *CLIOptions) *string { return &opts.FFmpegPath }},
		{Name: "ffplay-path", Usage: "Path to ffplay executable", Bind: func(opts *CLIOptions) *string { return &opts.FFplayPath }},
	}
}
