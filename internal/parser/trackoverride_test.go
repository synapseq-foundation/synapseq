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

package parser

import (
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestHasTrackOverride(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"  track 1 amplitude 10", true},
		{"  track 2 tone 440", true},
		{"  track 3 binaural 8", true},
		{" track 1 amplitude 10", false},               // only 1 space
		{"   track 1 amplitude 10", false},             // 3 spaces
		{"track 1 amplitude 10", false},                // no indentation
		{"  tone 440 binaural 10 amplitude 20", false}, // regular track, not override
		{"  noise pink amplitude 30", false},           // regular track
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasTrackOverride()
		if result != test.expected {
			ts.Errorf("For line %q, expected HasTrackOverride()=%v, got %v",
				test.line, test.expected, result)
		}
	}
}

func TestParseTrackOverrideDeclaration(ts *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedKind  string
		expectedIndex int
		expectedValue float64
		expectedRaw   string
		expectedRel   bool
		expectedWave  t.WaveformType
		expectedError bool
	}{
		{name: "tone absolute", line: "  track 1 tone 350", expectedKind: t.KeywordTone, expectedIndex: 1, expectedValue: 350, expectedRaw: "350"},
		{name: "tone relative", line: "  track 1 tone +50", expectedKind: t.KeywordTone, expectedIndex: 1, expectedValue: 50, expectedRaw: "+50", expectedRel: true},
		{name: "amplitude", line: "  track 1 amplitude -5", expectedKind: t.KeywordAmplitude, expectedIndex: 1, expectedValue: -5, expectedRaw: "-5", expectedRel: true},
		{name: "waveform", line: "  track 2 waveform sawtooth", expectedKind: t.KeywordWaveform, expectedIndex: 2, expectedWave: t.WaveformSawtooth},
		{name: "smooth", line: "  track 4 smooth 45", expectedKind: t.KeywordSmooth, expectedIndex: 4, expectedValue: 45, expectedRaw: "45"},
		{name: "missing track index", line: "  track amplitude 10", expectedError: true},
		{name: "invalid track index", line: "  track abc amplitude 10", expectedError: true},
		{name: "track index out of range", line: "  track 20 amplitude 10", expectedError: true},
		{name: "missing value", line: "  track 1 amplitude", expectedError: true},
		{name: "invalid value", line: "  track 1 amplitude abc", expectedError: true},
		{name: "invalid waveform value", line: "  track 1 waveform pulse", expectedError: true},
		{name: "extra tokens", line: "  track 1 amplitude 10 extra", expectedError: true},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		decl, err := ctx.ParseTrackOverrideDeclaration()
		if tt.expectedError {
			if err == nil {
				ts.Errorf("%s: expected error but got none", tt.name)
			}
			continue
		}
		if err != nil {
			ts.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		if decl.TrackIndex != tt.expectedIndex {
			ts.Errorf("%s: expected track index %d, got %d", tt.name, tt.expectedIndex, decl.TrackIndex)
		}
		if decl.Kind != tt.expectedKind {
			ts.Errorf("%s: expected kind %q, got %q", tt.name, tt.expectedKind, decl.Kind)
		}
		if decl.Value != tt.expectedValue {
			ts.Errorf("%s: expected value %v, got %v", tt.name, tt.expectedValue, decl.Value)
		}
		if decl.RawValue != tt.expectedRaw {
			ts.Errorf("%s: expected raw value %q, got %q", tt.name, tt.expectedRaw, decl.RawValue)
		}
		if decl.Relative != tt.expectedRel {
			ts.Errorf("%s: expected relative=%v, got %v", tt.name, tt.expectedRel, decl.Relative)
		}
		if decl.Waveform != tt.expectedWave {
			ts.Errorf("%s: expected waveform %v, got %v", tt.name, tt.expectedWave, decl.Waveform)
		}
	}
}

func TestParseTrackOverrideKeywordTypoDiagnostic(ts *testing.T) {
	ctx := NewTextParser("  track 1 amplitud 10")
	_, err := ctx.ParseTrackOverrideDeclaration()
	if err == nil {
		ts.Fatal("expected override diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Found != "amplitud" {
		ts.Fatalf("expected found token amplitud, got %q", diagnostic.Found)
	}
	if diagnostic.Suggestion != "did you mean \"amplitude\"?" {
		ts.Fatalf("expected amplitude suggestion, got %q", diagnostic.Suggestion)
	}
	if diagnostic.Span.Column != 11 || diagnostic.Span.EndColumn != 19 {
		ts.Fatalf("expected override span 11..19, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}
