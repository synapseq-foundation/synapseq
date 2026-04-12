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

func TestParsePresetDeclaration(ts *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedName  string
		expectedFrom  string
		expectedTempl bool
		expectedError bool
	}{
		{name: "plain preset", line: "MyPreset", expectedName: "MyPreset"},
		{name: "template", line: "alpha as template", expectedName: "alpha", expectedTempl: true},
		{name: "from template name", line: "derived from base-template", expectedName: "derived", expectedFrom: "base-template"},
		{name: "invalid prefix", line: "%AnotherPreset", expectedError: true},
		{name: "numeric start", line: "123Preset", expectedError: true},
		{name: "empty", line: "", expectedError: true},
		{name: "whitespace", line: "   ", expectedError: true},
		{name: "underscore", line: "Preset_", expectedName: "Preset_"},
		{name: "hyphen", line: "preset-01", expectedName: "preset-01"},
		{name: "reserved name stays parser-valid", line: "silence", expectedName: "silence"},
		{name: "invalid character", line: "Pre$et", expectedError: true},
		{name: "template extra token", line: "alpha as template extra", expectedError: true},
		{name: "from missing target", line: "alpha from", expectedError: true},
		{name: "invalid as target", line: "alpha as invalid", expectedError: true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		decl, err := ctx.ParsePresetDeclaration()
		if test.expectedError {
			if err == nil {
				ts.Errorf("%s: expected an error but got none", test.name)
			}
			continue
		}
		if err != nil {
			ts.Errorf("%s: did not expect an error but got: %v", test.name, err)
			continue
		}
		if decl.Name != test.expectedName {
			ts.Errorf("%s: expected preset name '%s' but got '%s'", test.name, test.expectedName, decl.Name)
		}
		if decl.FromName != test.expectedFrom {
			ts.Errorf("%s: expected from '%s' but got '%s'", test.name, test.expectedFrom, decl.FromName)
		}
		if decl.IsTemplate != test.expectedTempl {
			ts.Errorf("%s: expected IsTemplate=%v but got %v", test.name, test.expectedTempl, decl.IsTemplate)
		}
	}
}

func TestParsePresetTemplateTypoDiagnostic(ts *testing.T) {
	ctx := NewTextParser("alpha as templat")

	_, err := ctx.ParsePresetDeclaration()
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
