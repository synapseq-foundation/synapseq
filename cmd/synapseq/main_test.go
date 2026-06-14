// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
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

func TestResolveOutputTargetUsesDumpDefault(ts *testing.T) {
	opts := &clistyle.CLIOptions{Dump: true}
	outputFile, outputFormat := resolveOutputTarget("session", "", opts)
	if outputFile != "session.json" {
		ts.Fatalf("expected default dump output file, got %q", outputFile)
	}
	if outputFormat != ".json" {
		ts.Fatalf("expected output format .json, got %q", outputFormat)
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

func TestProcessSequenceOutputDumpWritesJSON(ts *testing.T) {
	loaded, err := synapseq.NewAppContext().LoadContent(`
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
`)
	if err != nil {
		ts.Fatalf("LoadContent error: %v", err)
	}

	outputFile := filepath.Join(ts.TempDir(), "dump.json")
	if err := processSequenceOutput(loaded, &outputOptions{OutputFile: outputFile, Dump: true, Quiet: true}); err != nil {
		ts.Fatalf("processSequenceOutput error: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		ts.Fatalf("read dump: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(content, &got); err != nil {
		ts.Fatalf("invalid JSON dump: %v\n%s", err, content)
	}
	if _, ok := got["timeline"].([]any); !ok {
		ts.Fatalf("expected timeline array in dump: %#v", got["timeline"])
	}
}

func TestFindRemoteEntryByID(ts *testing.T) {
	index := &types.RemoteIndex{Entries: []types.RemoteEntry{{ID: "focus"}, {ID: "sleep"}}}
	entry := findRemoteEntryByID(index, "sleep")
	if entry == nil || entry.ID != "sleep" {
		ts.Fatalf("expected sleep entry, got %#v", entry)
	}
	if missing := findRemoteEntryByID(index, "missing"); missing != nil {
		ts.Fatalf("expected missing entry to return nil, got %#v", missing)
	}
}

func TestFormatRemoteTable(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	table := formatRemoteTable("ID\tDURATION\nfocus\t15\n")
	checks := []string{"ID\tDURATION", "focus\t15"}
	for _, check := range checks {
		if !strings.Contains(table, check) {
			ts.Fatalf("expected formatted table to contain %q, got:\n%s", check, table)
		}
	}
}

func TestFormatRemoteDescription(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	formatted := formatRemoteDescription("first line")
	if !strings.Contains(formatted, "Description:") || !strings.Contains(formatted, "first line") {
		ts.Fatalf("expected formatted description, got:\n%s", formatted)
	}

	if empty := formatRemoteDescription(""); !strings.Contains(empty, "No description available") {
		ts.Fatalf("expected empty description fallback, got:\n%s", empty)
	}
}

func TestSortedRemoteEntries(ts *testing.T) {
	entries := sortedRemoteEntries([]types.RemoteEntry{
		{ID: "old", CreatedAt: "2026-01-01T00:00:00Z"},
		{ID: "new", CreatedAt: "2026-02-01T00:00:00Z"},
	})
	if entries[0].ID != "new" || entries[1].ID != "old" {
		ts.Fatalf("expected newest entry first, got %#v", entries)
	}
}

func TestShortDate(ts *testing.T) {
	if got := shortDate("2026-01-02T03:04:05Z"); got != "2026-01-02" {
		ts.Fatalf("expected date prefix, got %q", got)
	}
	if got := shortDate("short"); got != "short" {
		ts.Fatalf("expected short value unchanged, got %q", got)
	}
}

func TestResolveSpecialCommandPrecedence(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{
		ShowVersion: true,
		RemoteSync:  true,
		RemoteGet:   "focus",
		New:         "sleep",
	}, []string{"out.spsq"})

	if command.Kind != clistyle.SpecialCommandShowVersion {
		ts.Fatalf("expected version to win precedence, got %q", command.Kind)
	}
}

func TestResolveSpecialCommandRemoteGetUsesOptionalArg(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{RemoteGet: "focus"}, []string{"out.wav"})
	if command.Kind != clistyle.SpecialCommandRemoteGet {
		ts.Fatalf("expected remote-get command, got %q", command.Kind)
	}
	if command.OptionalArg != "out.wav" {
		ts.Fatalf("expected remote-get optional arg out.wav, got %q", command.OptionalArg)
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

func TestResolveSpecialCommandRemoteDownloadPrecedesRemoteInfo(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{RemoteDownload: "focus", RemoteInfo: "sleep"}, []string{"downloads"})
	if command.Kind != clistyle.SpecialCommandRemoteDownload {
		ts.Fatalf("expected remote-download to win precedence over remote-info, got %q", command.Kind)
	}
	if command.OptionalArg != "downloads" {
		ts.Fatalf("expected download target arg downloads, got %q", command.OptionalArg)
	}
}

func TestResolveSpecialCommandIgnoresQuietBeforeRemoteList(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{Quiet: true, RemoteList: true}, nil)
	if command.Kind != clistyle.SpecialCommandRemoteList {
		ts.Fatalf("expected remote-list command, got %q", command.Kind)
	}
}

func TestResolveSpecialCommandIgnoresNoColorBeforeRemoteGet(ts *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{NoColor: true, RemoteGet: "calm-state"}, []string{"out.wav"})
	if command.Kind != clistyle.SpecialCommandRemoteGet {
		ts.Fatalf("expected remote-get command, got %q", command.Kind)
	}
	if command.OptionalArg != "out.wav" {
		ts.Fatalf("expected remote-get optional arg out.wav, got %q", command.OptionalArg)
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
