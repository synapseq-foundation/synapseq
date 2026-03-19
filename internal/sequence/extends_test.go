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

package sequence

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func writeExtendsFile(ts *testing.T, relPath string, content string) string {
	ts.Helper()

	dir := ts.TempDir()
	path := filepath.Join(dir, relPath)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		ts.Fatalf("mkdir temp extends dir: %v", err)
	}

	if err := os.WriteFile(path, []byte(strings.TrimLeft(content, "\n")+"\n"), 0o600); err != nil {
		ts.Fatalf("write temp extends file: %v", err)
	}

	return path
}

func TestExtendsSuccess(ts *testing.T) {
	path := writeExtendsFile(ts, "configs/base.spsc", `
@samplerate 48000
@volume 80
@ambiance forest testdata/noise

base
  noise pink amplitude 50
`)

	got, err := extends(path)
	if err != nil {
		ts.Fatalf("extends error: %v", err)
	}

	if got == nil {
		ts.Fatal("expected extends result, got nil")
	}

	if len(got.Presets) != 1 {
		ts.Fatalf("expected 1 preset, got %d", len(got.Presets))
	}

	if got.Presets[0].String() != "base" {
		ts.Fatalf("expected preset name %q, got %q", "base", got.Presets[0].String())
	}

	if got.Options == nil {
		ts.Fatal("expected parsed options, got nil")
	}

	if got.Options.Values[t.KeywordOptionSampleRate] != "48000" {
		ts.Fatalf("expected samplerate raw value 48000, got %q", got.Options.Values[t.KeywordOptionSampleRate])
	}

	if got.Options.Values[t.KeywordOptionVolume] != "80" {
		ts.Fatalf("expected volume raw value 80, got %q", got.Options.Values[t.KeywordOptionVolume])
	}

	wantAmbiancePath := filepath.Join(filepath.Dir(path), "testdata", "noise.wav")
	if got.Options.Ambiance["forest"] != wantAmbiancePath {
		ts.Fatalf("expected ambiance path %q, got %q", wantAmbiancePath, got.Options.Ambiance["forest"])
	}

	built, err := got.Options.Build()
	if err != nil {
		ts.Fatalf("build parsed options: %v", err)
	}

	if built.SampleRate != 48000 {
		ts.Fatalf("expected SampleRate=48000, got %d", built.SampleRate)
	}

	if built.Volume != 80 {
		ts.Fatalf("expected Volume=80, got %d", built.Volume)
	}

	if built.Ambiance["forest"] != wantAmbiancePath {
		ts.Fatalf("expected built ambiance path %q, got %q", wantAmbiancePath, built.Ambiance["forest"])
	}
}

func TestExtendsRejectsNestedExtendsOption(ts *testing.T) {
	path := writeExtendsFile(ts, "nested.spsc", `
@extends shared/base

alpha
  noise pink amplitude 10
`)

	_, err := extends(path)
	if err == nil {
		ts.Fatal("expected error for nested extends option, got nil")
	}

	if !strings.Contains(err.Error(), "extends option is not supported in extended files") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestExtendsRejectsOptionAfterPreset(ts *testing.T) {
	path := writeExtendsFile(ts, "late-option.spsc", `
alpha
@volume 75
`)

	_, err := extends(path)
	if err == nil {
		ts.Fatal("expected error for option after preset, got nil")
	}

	if !strings.Contains(err.Error(), "options must be defined before any presets") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestExtendsRejectsTrackBeforePreset(ts *testing.T) {
	path := writeExtendsFile(ts, "track-first.spsc", `
  noise pink amplitude 10
`)

	_, err := extends(path)
	if err == nil {
		ts.Fatal("expected error for track before preset, got nil")
	}

	if !strings.Contains(err.Error(), "track defined before any preset") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestExtendsRejectsEmptyPreset(ts *testing.T) {
	path := writeExtendsFile(ts, "empty-preset.spsc", `
alpha
`)

	_, err := extends(path)
	if err == nil {
		ts.Fatal("expected error for empty preset, got nil")
	}

	if !strings.Contains(err.Error(), "preset \"alpha\" is empty") {
		ts.Fatalf("unexpected error: %v", err)
	}
}
