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
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}

func TestParseFlags(ts *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		args         []string
		expected     *CLIOptions
		expectedArgs []string
		expectError  bool
	}{
		// Version flag
		{
			args:         []string{"cmd", "-version"},
			expected:     &CLIOptions{ShowVersion: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// New template flag
		{
			args:         []string{"cmd", "-new", "meditation"},
			expected:     &CLIOptions{New: "meditation"},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-new", "focus"},
			expected:     &CLIOptions{New: "focus"},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-preview", "input.spsq"},
			expected:     &CLIOptions{Preview: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// Help flag
		{
			args:         []string{"cmd", "-help"},
			expected:     &CLIOptions{ShowHelp: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-manual"},
			expected:     &CLIOptions{ShowManual: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Quiet flag
		{
			args:         []string{"cmd", "-quiet", "input.spsq", "output.wav"},
			expected:     &CLIOptions{Quiet: true},
			expectedArgs: []string{"input.spsq", "output.wav"},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-no-color", "input.spsq"},
			expected:     &CLIOptions{NoColor: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// Test flag
		{
			args:         []string{"cmd", "-test", "input.spsq"},
			expected:     &CLIOptions{Test: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// Test and quiet combined
		{
			args:         []string{"cmd", "-test", "-quiet", "input.spsq"},
			expected:     &CLIOptions{Test: true, Quiet: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// No flags, just arguments
		{
			args:         []string{"cmd", "input.spsq", "output.wav"},
			expected:     &CLIOptions{},
			expectedArgs: []string{"input.spsq", "output.wav"},
			expectError:  false,
		},
		// No arguments at all
		{
			args:         []string{"cmd"},
			expected:     &CLIOptions{},
			expectedArgs: []string{},
			expectError:  false,
		},
		// All boolean flags enabled
		{
			args:         []string{"cmd", "-quiet", "-test"},
			expected:     &CLIOptions{Quiet: true, Test: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Hub boolean flags
		{
			args:         []string{"cmd", "-hub-update"},
			expected:     &CLIOptions{HubUpdate: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-hub-clean"},
			expected:     &CLIOptions{HubClean: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-hub-list"},
			expected:     &CLIOptions{HubList: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Hub string flags
		{
			args:         []string{"cmd", "-hub-search", "focus"},
			expected:     &CLIOptions{HubSearch: "focus"},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-hub-download", "baseline"},
			expected:     &CLIOptions{HubDownload: "baseline"},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-hub-info", "deep-sleep"},
			expected:     &CLIOptions{HubInfo: "deep-sleep"},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-hub-get", "alpha-pack"},
			expected:     &CLIOptions{HubGet: "alpha-pack"},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Combined hub options and positional args
		{
			args: []string{"cmd", "-hub-search", "relax", "-hub-download", "rain-pack", "input.spsq"},
			expected: &CLIOptions{
				HubSearch:   "relax",
				HubDownload: "rain-pack",
			},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// External tool path flags
		{
			args: []string{"cmd", "-ffmpeg-path", "/usr/local/bin/ffmpeg", "-ffplay-path", "/usr/local/bin/ffplay", "input.spsq"},
			expected: &CLIOptions{
				FFmpegPath: "/usr/local/bin/ffmpeg",
				FFplayPath: "/usr/local/bin/ffplay",
			},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// Windows association flags (parse-only validation)
		{
			args:         []string{"cmd", "-install-file-association"},
			expected:     &CLIOptions{InstallFileAssociation: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			args:         []string{"cmd", "-uninstall-file-association"},
			expected:     &CLIOptions{UninstallFileAssociation: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Invalid flag should return error
		{
			args:         []string{"cmd", "-invalid"},
			expected:     nil,
			expectedArgs: nil,
			expectError:  true,
		},
		// Unknown flag with valid flags
		{
			args:         []string{"cmd", "-quiet", "-unknown", "input.spsq"},
			expected:     nil,
			expectedArgs: nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		os.Args = test.args
		opts, args, err := ParseFlags()

		if test.expectError {
			if err == nil {
				ts.Errorf("For args %v, expected error but got none", test.args)
			}
			continue
		}

		if err != nil {
			ts.Errorf("For args %v, unexpected error: %v", test.args, err)
			continue
		}

		if opts.ShowVersion != test.expected.ShowVersion {
			ts.Errorf("For args %v, ShowVersion: expected %v but got %v", test.args, test.expected.ShowVersion, opts.ShowVersion)
		}
		if opts.New != test.expected.New {
			ts.Errorf("For args %v, New: expected %q but got %q", test.args, test.expected.New, opts.New)
		}
		if opts.ShowHelp != test.expected.ShowHelp {
			ts.Errorf("For args %v, ShowHelp: expected %v but got %v", test.args, test.expected.ShowHelp, opts.ShowHelp)
		}
		if opts.ShowManual != test.expected.ShowManual {
			ts.Errorf("For args %v, ShowManual: expected %v but got %v", test.args, test.expected.ShowManual, opts.ShowManual)
		}
		if opts.Preview != test.expected.Preview {
			ts.Errorf("For args %v, Preview: expected %v but got %v", test.args, test.expected.Preview, opts.Preview)
		}
		if opts.Quiet != test.expected.Quiet {
			ts.Errorf("For args %v, Quiet: expected %v but got %v", test.args, test.expected.Quiet, opts.Quiet)
		}
		if opts.NoColor != test.expected.NoColor {
			ts.Errorf("For args %v, NoColor: expected %v but got %v", test.args, test.expected.NoColor, opts.NoColor)
		}
		if opts.Test != test.expected.Test {
			ts.Errorf("For args %v, Test: expected %v but got %v", test.args, test.expected.Test, opts.Test)
		}
		if opts.Play != test.expected.Play {
			ts.Errorf("For args %v, Play: expected %v but got %v", test.args, test.expected.Play, opts.Play)
		}
		if opts.HubUpdate != test.expected.HubUpdate {
			ts.Errorf("For args %v, HubUpdate: expected %v but got %v", test.args, test.expected.HubUpdate, opts.HubUpdate)
		}
		if opts.HubClean != test.expected.HubClean {
			ts.Errorf("For args %v, HubClean: expected %v but got %v", test.args, test.expected.HubClean, opts.HubClean)
		}
		if opts.HubList != test.expected.HubList {
			ts.Errorf("For args %v, HubList: expected %v but got %v", test.args, test.expected.HubList, opts.HubList)
		}
		if opts.HubSearch != test.expected.HubSearch {
			ts.Errorf("For args %v, HubSearch: expected %q but got %q", test.args, test.expected.HubSearch, opts.HubSearch)
		}
		if opts.HubDownload != test.expected.HubDownload {
			ts.Errorf("For args %v, HubDownload: expected %q but got %q", test.args, test.expected.HubDownload, opts.HubDownload)
		}
		if opts.HubInfo != test.expected.HubInfo {
			ts.Errorf("For args %v, HubInfo: expected %q but got %q", test.args, test.expected.HubInfo, opts.HubInfo)
		}
		if opts.HubGet != test.expected.HubGet {
			ts.Errorf("For args %v, HubGet: expected %q but got %q", test.args, test.expected.HubGet, opts.HubGet)
		}
		if opts.FFmpegPath != test.expected.FFmpegPath {
			ts.Errorf("For args %v, FFmpegPath: expected %q but got %q", test.args, test.expected.FFmpegPath, opts.FFmpegPath)
		}
		if opts.FFplayPath != test.expected.FFplayPath {
			ts.Errorf("For args %v, FFplayPath: expected %q but got %q", test.args, test.expected.FFplayPath, opts.FFplayPath)
		}
		if opts.InstallFileAssociation != test.expected.InstallFileAssociation {
			ts.Errorf("For args %v, InstallFileAssociation: expected %v but got %v", test.args, test.expected.InstallFileAssociation, opts.InstallFileAssociation)
		}
		if opts.UninstallFileAssociation != test.expected.UninstallFileAssociation {
			ts.Errorf("For args %v, UninstallFileAssociation: expected %v but got %v", test.args, test.expected.UninstallFileAssociation, opts.UninstallFileAssociation)
		}

		if len(args) != len(test.expectedArgs) {
			ts.Errorf("For args %v, expected args %v but got %v", test.args, test.expectedArgs, args)
		} else {
			for i := range args {
				if args[i] != test.expectedArgs[i] {
					ts.Errorf("For args %v, expected args[%d] = %q but got %q", test.args, i, test.expectedArgs[i], args[i])
					break
				}
			}
		}
	}
}

func TestParseFlagsEdgeCases(ts *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with stdin input
	os.Args = []string{"cmd", "-quiet", "-", "output.wav"}
	opts, args, err := ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for stdin input: %v", err)
	}
	if !opts.Quiet {
		ts.Errorf("expected Quiet=true for stdin test")
	}
	if len(args) != 2 || args[0] != "-" || args[1] != "output.wav" {
		ts.Errorf("expected args [\"-\", \"output.wav\"], got %v", args)
	}

	// Test with stdout output
	os.Args = []string{"cmd", "-"}
	opts, args, err = ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for stdout output: %v", err)
	}
	if len(args) != 1 || args[0] != "-" {
		ts.Errorf("expected args [\"-\"], got %v", args)
	}

	// Test with URL input
	os.Args = []string{"cmd", "https://example.com/sequence.spsq", "output.wav"}
	opts, args, err = ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for URL input: %v", err)
	}
	if len(args) != 2 || args[0] != "https://example.com/sequence.spsq" || args[1] != "output.wav" {
		ts.Errorf("expected URL args, got %v", args)
	}

	// Test with ffmpeg and ffplay path options
	os.Args = []string{"cmd", "-ffmpeg-path", "ffmpeg-custom", "-ffplay-path", "ffplay-custom", "input.spsq"}
	opts, args, err = ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for ffmpeg/ffplay path options: %v", err)
	}
	if opts.FFmpegPath != "ffmpeg-custom" || opts.FFplayPath != "ffplay-custom" {
		ts.Errorf("expected custom ffmpeg/ffplay paths, got ffmpeg=%q ffplay=%q", opts.FFmpegPath, opts.FFplayPath)
	}
	if len(args) != 1 || args[0] != "input.spsq" {
		ts.Errorf("expected args [\"input.spsq\"], got %v", args)
	}
}

func TestParseFlagsUnknownFlagReturnsDiagnostic(ts *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd", "-g"}
	_, _, err := ParseFlags()
	if err == nil {
		ts.Fatalf("expected error for unknown flag")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}
	if diagnostic.Message != "unknown command-line flag" {
		ts.Fatalf("unexpected message: %q", diagnostic.Message)
	}
	if diagnostic.Found != "-g" {
		ts.Fatalf("unexpected found value: %q", diagnostic.Found)
	}
	if !strings.Contains(diagnostic.Hint, "-help") {
		ts.Fatalf("unexpected hint: %q", diagnostic.Hint)
	}
}

func TestHelpIncludesQuickStart(ts *testing.T) {
	originalStdout := os.Stdout
	originalColorOutput := color.Output
	originalNoColor := color.NoColor
	defer func() {
		color.Output = originalColorOutput
		color.NoColor = originalNoColor
		os.Stdout = originalStdout
	}()

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		ts.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = writePipe
	color.Output = writePipe
	SetColorEnabled(true)
	Help()
	writePipe.Close()
	os.Stdout = originalStdout

	output, err := io.ReadAll(readPipe)
	if err != nil {
		ts.Fatalf("failed to read help output: %v", err)
	}

	helpText := stripANSI(string(output))
	checks := []string{
		"Usage:\n  synapseq [options] <input> [output]",
		"Quick start:",
		"Generate session.wav in the current folder",
		"Print the full language and usage manual",
		"Generate session.html with a visual timeline preview",
		"defaults to <input>.wav",
		"-new TYPE         Template type: meditation, focus, sleep, relaxation",
		"-manual           Show the full manual",
		"-preview          Render an HTML preview timeline",
		"Hub examples:",
		"Run -hub-update once before using other -hub-* commands.",
		"synapseq -hub-update",
		"synapseq -hub-get calm-state calm-state.wav",
	}

	for _, expected := range checks {
		if !strings.Contains(helpText, expected) {
			ts.Errorf("help output missing %q\nfull output:\n%s", expected, helpText)
		}
	}
}

func TestHelpNoColorOmitsANSI(ts *testing.T) {
	originalStdout := os.Stdout
	originalColorOutput := color.Output
	originalNoColor := color.NoColor
	defer func() {
		color.Output = originalColorOutput
		color.NoColor = originalNoColor
		os.Stdout = originalStdout
	}()

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		ts.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = writePipe
	color.Output = writePipe
	SetColorEnabled(false)
	Help()
	writePipe.Close()

	output, err := io.ReadAll(readPipe)
	if err != nil {
		ts.Fatalf("failed to read help output: %v", err)
	}

	if ansiPattern.Match(output) {
		ts.Fatalf("expected help without ANSI escapes, got %q", string(output))
	}
	if !strings.Contains(string(output), "-no-color") {
		ts.Fatalf("expected help to mention -no-color, got:\n%s", string(output))
	}
}
