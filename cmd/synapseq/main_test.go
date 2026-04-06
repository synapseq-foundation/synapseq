//go:build !js && !wasm

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

package main

import (
	"errors"
	"strings"
	"testing"

	clistyle "github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	types "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestFormatCLIErrorDiagnostic(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := diag.UnexpectedToken(
		diag.Span{
			File:      "example.spsq",
			Line:      4,
			Column:    12,
			EndColumn: 19,
			LineText:  "  tone 300 binaual 10 amplitude 10",
		},
		"binaual",
		"binaural",
		"monaural",
	)

	formatted := formatCLIError(err)

	checks := []string{
		"synapseq: example.spsq:4:12: unexpected token",
		"  tone 300 binaual 10 amplitude 10",
		"           ^^^^^^^",
		"did you mean \"binaural\"?",
	}

	for _, check := range checks {
		if !strings.Contains(formatted, check) {
			ts.Fatalf("expected formatted CLI error to contain %q, got:\n%s", check, formatted)
		}
	}
}

func TestFormatCLIErrorFallback(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := errors.New("plain error")
	formatted := formatCLIError(err)
	if formatted != "synapseq: plain error" {
		ts.Fatalf("unexpected fallback formatting: %q", formatted)
	}
}

func TestFormatCLIErrorDiagnosticCause(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := diag.Wrap(diag.KindIO, "failed while streaming audio to ffplay", errors.New("failed to load ambiance file [0] (/tmp/rain.wav): error opening file: open /tmp/rain.wav: no such file or directory"))

	formatted := formatCLIError(err)

	checks := []string{
		"synapseq: failed while streaming audio to ffplay",
		"cause: failed to load ambiance file [0] (/tmp/rain.wav): error opening file: open /tmp/rain.wav: no such file or directory",
	}

	for _, check := range checks {
		if !strings.Contains(formatted, check) {
			ts.Fatalf("expected formatted CLI error to contain %q, got:\n%s", check, formatted)
		}
	}
}

func TestFormatCLIErrorDiagnosticCauseSkipsDuplicateMessage(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := diag.Validation("ambiance not found").WithCause(errors.New("ambiance not found"))

	formatted := formatCLIError(err)

	if strings.Contains(formatted, "cause:") {
		ts.Fatalf("expected duplicate cause to be omitted, got:\n%s", formatted)
	}
}

func TestResolveOutputTargetUsesRequestedOutput(ts *testing.T) {
	opts := &clistyle.CLIOptions{}
	outputFile, outputFormat := resolveOutputTarget("session", "custom.mp3", opts)
	if outputFile != "custom.mp3" {
		ts.Fatalf("expected requested output file to be preserved, got %q", outputFile)
	}
	if outputFormat != ".mp3" {
		ts.Fatalf("expected output format .mp3, got %q", outputFormat)
	}
}

func TestResolveOutputTargetUsesOptionDefaults(ts *testing.T) {
	opts := &clistyle.CLIOptions{Preview: true}
	outputFile, outputFormat := resolveOutputTarget("session", "", opts)
	if outputFile != "session.html" {
		ts.Fatalf("expected default preview output file, got %q", outputFile)
	}
	if outputFormat != ".html" {
		ts.Fatalf("expected output format .html, got %q", outputFormat)
	}
}

func TestBuildOutputOptions(ts *testing.T) {
	opts := &clistyle.CLIOptions{Play: true, FFplayPath: "ffplay", FFmpegPath: "ffmpeg"}
	outputOpts := buildOutputOptions("out.mp3", ".mp3", opts)
	if outputOpts.OutputFile != "out.mp3" {
		ts.Fatalf("expected output file out.mp3, got %q", outputOpts.OutputFile)
	}
	if !outputOpts.Mp3 {
		ts.Fatalf("expected mp3 output to be enabled")
	}
	if !outputOpts.Play {
		ts.Fatalf("expected play to be enabled")
	}
	if outputOpts.FFplayPath != "ffplay" || outputOpts.FFmpegPath != "ffmpeg" {
		ts.Fatalf("unexpected ffmpeg/ffplay paths: %#v", outputOpts)
	}
}

func TestFindHubEntryByID(ts *testing.T) {
	manifest := &types.HubManifest{Entries: []types.HubEntry{{ID: "focus"}, {ID: "sleep"}}}
	entry := findHubEntryByID(manifest, "sleep")
	if entry == nil || entry.ID != "sleep" {
		ts.Fatalf("expected sleep entry, got %#v", entry)
	}
	if missing := findHubEntryByID(manifest, "missing"); missing != nil {
		ts.Fatalf("expected missing entry to return nil, got %#v", missing)
	}
}

func TestFormatHubTable(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	table := formatHubTable("ID\tAUTHOR\nfocus\truan\n")
	checks := []string{"ID\tAUTHOR", "focus\truan"}
	for _, check := range checks {
		if !strings.Contains(table, check) {
			ts.Fatalf("expected formatted table to contain %q, got:\n%s", check, table)
		}
	}
}

func TestFormatHubDependencies(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	formatted := formatHubDependencies([]types.HubDependency{{ID: "rain", Type: types.HubDependencyTypeAmbiance}})
	if !strings.Contains(formatted, "rain") || !strings.Contains(formatted, "ambiance") {
		ts.Fatalf("expected formatted dependencies to include dependency details, got:\n%s", formatted)
	}

	if empty := formatHubDependencies(nil); !strings.Contains(empty, "Dependencies: None") {
		ts.Fatalf("expected empty dependency text, got:\n%s", empty)
	}
}

func TestFormatHubDescription(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	formatted := formatHubDescription([]string{"first line", "second line"})
	if !strings.Contains(formatted, "Description:") || !strings.Contains(formatted, "first line") {
		ts.Fatalf("expected formatted description with comments, got:\n%s", formatted)
	}

	if empty := formatHubDescription(nil); !strings.Contains(empty, "No description available") {
		ts.Fatalf("expected empty description fallback, got:\n%s", empty)
	}
}

func TestResolveSpecialCommandPrecedence(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{
		ShowVersion: true,
		ShowManual:  true,
		HubUpdate:   true,
		HubGet:      "focus",
		New:         "sleep",
	}, []string{"out.spsq"})

	if command.Kind != clistyle.SpecialCommandShowVersion {
		ts.Fatalf("expected version to win precedence, got %q", command.Kind)
	}
}

func TestResolveSpecialCommandHubGetUsesOptionalArg(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{HubGet: "focus"}, []string{"out.wav"})
	if command.Kind != clistyle.SpecialCommandHubGet {
		ts.Fatalf("expected hub-get command, got %q", command.Kind)
	}
	if command.OptionalArg != "out.wav" {
		ts.Fatalf("expected hub-get optional arg out.wav, got %q", command.OptionalArg)
	}
}

func TestResolveSpecialCommandTemplateUsesOptionalArg(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{New: "meditation"}, []string{"custom.spsq"})
	if command.Kind != clistyle.SpecialCommandGenerateTemplate {
		ts.Fatalf("expected generate-template command, got %q", command.Kind)
	}
	if command.OptionalArg != "custom.spsq" {
		ts.Fatalf("expected template optional arg custom.spsq, got %q", command.OptionalArg)
	}
}

func TestResolveSpecialCommandNoMatch(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{}, nil)
	if command.Kind != clistyle.SpecialCommandNone {
		ts.Fatalf("expected no special command, got %q", command.Kind)
	}
}

func TestResolveSpecialCommandHubDownloadPrecedesHubInfo(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{HubDownload: "focus", HubInfo: "sleep"}, []string{"downloads"})
	if command.Kind != clistyle.SpecialCommandHubDownload {
		ts.Fatalf("expected hub-download to win precedence over hub-info, got %q", command.Kind)
	}
	if command.OptionalArg != "downloads" {
		ts.Fatalf("expected download target arg downloads, got %q", command.OptionalArg)
	}
}

func TestPrepareSequenceCommandRejectsInvalidArgCount(ts *testing.T) {
	command, err := prepareSequenceCommand([]string{"one", "two", "three"}, &clistyle.CLIOptions{})
	if err == nil {
		ts.Fatalf("expected invalid arg count error, got command %#v", command)
	}
	if !strings.Contains(err.Error(), "invalid number of flags") {
		ts.Fatalf("unexpected error text: %v", err)
	}
}

func TestPrepareSequenceCommandUsesDefaultOutput(ts *testing.T) {
	command, err := prepareSequenceCommand([]string{"sessions/focus.spsq"}, &clistyle.CLIOptions{})
	if err != nil {
		ts.Fatalf("unexpected error preparing command: %v", err)
	}
	if command.inputFile != "sessions/focus.spsq" {
		ts.Fatalf("unexpected input file: %q", command.inputFile)
	}
	if command.outputFile != "focus.wav" {
		ts.Fatalf("expected default output focus.wav, got %q", command.outputFile)
	}
	if command.outputFormat != ".wav" {
		ts.Fatalf("expected output format .wav, got %q", command.outputFormat)
	}
}

func TestPrepareSequenceCommandUsesExplicitOutput(ts *testing.T) {
	command, err := prepareSequenceCommand([]string{"sessions/focus.spsq", "custom.html"}, &clistyle.CLIOptions{Preview: true})
	if err != nil {
		ts.Fatalf("unexpected error preparing command: %v", err)
	}
	if command.outputFile != "custom.html" {
		ts.Fatalf("expected explicit output custom.html, got %q", command.outputFile)
	}
	if command.outputFormat != ".html" {
		ts.Fatalf("expected output format .html, got %q", command.outputFormat)
	}
}
