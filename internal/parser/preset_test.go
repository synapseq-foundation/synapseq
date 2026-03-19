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

func TestHasPreset(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"MyPreset", true},
		{" AnotherPreset", false},
		{"preset1", true},
		{"  preset2", false},
		{"123Preset", false},
		{"", false},
		{"   ", false},
		{"%Preset", false},
		{"Preset_", true},
		{"preset-01", true},
		{"preset-", true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasPreset()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasPreset() to return %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParsePreset(ts *testing.T) {
	tests := []struct {
		line          string
		expectedName  string
		expectedError bool
	}{
		{"MyPreset", "mypreset", false},
		{"%AnotherPreset", "", true},
		{"preset1", "preset1", false},
		{"123Preset", "", true},
		{"", "", true},
		{"   ", "", true},
		{"Preset_", "preset_", false},
		{"preset-01", "preset-01", false},
		{"preset-", "preset-", false},
		{"silence", "", true}, // reserved name
		{"Pre$et", "", true},  // invalid character
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		preset, err := ctx.ParsePreset(nil)
		if test.expectedError {
			if err == nil {
				ts.Errorf("For line '%s', expected an error but got none", test.line)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For line '%s', did not expect an error but got: %v", test.line, err)
			continue
		}
		if preset.String() != test.expectedName {
			ts.Errorf("For line '%s', expected preset name '%s' but got '%s'", test.line, test.expectedName, preset.String())
		}
	}
}

func TestParsePreset_WithTemplate(ts *testing.T) {
	// Create base presets for inheritance tests
	var presets []t.Preset

	// Create a template preset
	templatePreset, err := t.NewPreset("base-template", true, nil)
	if err != nil {
		ts.Fatalf("failed to create template preset: %v", err)
	}
	templatePreset.Track[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	templatePreset.Track[1] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	presets = append(presets, *templatePreset)

	// Create a regular preset for testing inheritance errors
	regularPreset, err := t.NewPreset("regular", false, nil)
	if err != nil {
		ts.Fatalf("failed to create regular preset: %v", err)
	}
	presets = append(presets, *regularPreset)

	tests := []struct {
		name          string
		line          string
		expectedName  string
		expectedError bool
		checkTemplate bool
		isTemplate    bool
		checkFrom     bool
		hasFrom       bool
	}{
		{
			name:          "template declaration",
			line:          "alpha as template",
			expectedName:  "alpha",
			expectedError: false,
			checkTemplate: true,
			isTemplate:    true,
		},
		{
			name:          "inherit from template",
			line:          "derived from base-template",
			expectedName:  "derived",
			expectedError: false,
			checkFrom:     true,
			hasFrom:       true,
		},
		{
			name:          "inherit from non-template should fail",
			line:          "bad from regular",
			expectedError: true,
		},
		{
			name:          "inherit from unknown preset should fail",
			line:          "bad from unknown-preset",
			expectedError: true,
		},
		{
			name:          "template with extra tokens should fail",
			line:          "alpha as template extra",
			expectedError: true,
		},
		{
			name:          "from with missing preset name should fail",
			line:          "alpha from",
			expectedError: true,
		},
		{
			name:          "invalid 'as' keyword usage should fail",
			line:          "alpha as invalid",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		preset, err := ctx.ParsePreset(&presets)

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

		if preset.String() != tt.expectedName {
			ts.Errorf("%s: expected name %q, got %q", tt.name, tt.expectedName, preset.String())
		}

		if tt.checkTemplate && preset.IsTemplate != tt.isTemplate {
			ts.Errorf("%s: expected IsTemplate=%v, got %v", tt.name, tt.isTemplate, preset.IsTemplate)
		}

		if tt.checkFrom {
			if tt.hasFrom && preset.From == nil {
				ts.Errorf("%s: expected From to be set", tt.name)
			} else if !tt.hasFrom && preset.From != nil {
				ts.Errorf("%s: expected From to be nil", tt.name)
			}
		}
	}
}

func TestParsePresetTemplateTypoDiagnostic(ts *testing.T) {
	ctx := NewTextParser("alpha as templat")

	_, err := ctx.ParsePreset(nil)
	if err == nil {
		ts.Fatal("expected preset diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Found != "templat" {
		ts.Fatalf("expected found token templat, got %q", diagnostic.Found)
	}
	if diagnostic.Suggestion != "did you mean \"template\"?" {
		ts.Fatalf("expected template suggestion, got %q", diagnostic.Suggestion)
	}
	if diagnostic.Span.Column != 10 || diagnostic.Span.EndColumn != 17 {
		ts.Fatalf("expected template span 10..17, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}
