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

func TestHasTimeline(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"00:00:00 alpha", true},
		{"23:59:59 preset-1", true},
		{" 00:00:00 alpha", false},
		{"00:00 alpha", false},
		{"alpha", false},
		{"+00:00:10", false},
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasTimeline()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasTimeline() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseTimelineDeclaration(ts *testing.T) {
	tests := []struct {
		line               string
		expectError        bool
		expectedMs         int
		expectedPreset     string
		expectedTransition t.TransitionType
		expectedSteps      int
	}{
		{"00:00:00 alpha", false, 0, "alpha", t.TransitionSteady, 0},
		{"00:00:15 alpha", false, 15_000, "alpha", t.TransitionSteady, 0},
		{"12:34:56 alpha", false, (12*3600 + 34*60 + 56) * 1000, "alpha", t.TransitionSteady, 0},
		{"00:01:00 alpha ease-out", false, 60_000, "alpha", t.TransitionEaseOut, 0},
		{"00:02:00 alpha ease-in", false, 120_000, "alpha", t.TransitionEaseIn, 0},
		{"00:03:00 alpha smooth", false, 180_000, "alpha", t.TransitionSmooth, 0},
		{"00:03:00 alpha smooth 3", false, 180_000, "alpha", t.TransitionSmooth, 3},
		{"00:03:00 alpha steady 0", false, 180_000, "alpha", t.TransitionSteady, 0},
		{"24:00:00 alpha", true, 0, "", t.TransitionSteady, 0},
		{"00:60:00 alpha", true, 0, "", t.TransitionSteady, 0},
		{"00:00:60 alpha", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha extra", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05", true, 0, "", t.TransitionSteady, 0},
		{"00:00 alpha", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha invalid-transition", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha steady -1", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha 4", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha steady extra", true, 0, "", t.TransitionSteady, 0},
		{"00:00:05 alpha steady 3 extra", true, 0, "", t.TransitionSteady, 0},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		decl, err := ctx.ParseTimelineDeclaration()
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got none", test.line)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}
		if decl.Time != test.expectedMs {
			ts.Errorf("For line '%s', expected time %d but got %d", test.line, test.expectedMs, decl.Time)
		}
		if decl.PresetName != test.expectedPreset {
			ts.Errorf("For line '%s', expected preset %q but got %q", test.line, test.expectedPreset, decl.PresetName)
		}
		if decl.Transition != test.expectedTransition {
			ts.Errorf("For line '%s', expected transition %v but got %v", test.line, test.expectedTransition, decl.Transition)
		}
		if decl.Steps != test.expectedSteps {
			ts.Errorf("For line '%s', expected steps %d but got %d", test.line, test.expectedSteps, decl.Steps)
		}
	}
}

func TestParseTime(ts *testing.T) {
	tests := []struct {
		in          string
		expectedMs  int
		expectError bool
	}{
		{"00:00:00", 0, false},
		{"00:00:01", 1_000, false},
		{"00:01:00", 60_000, false},
		{"01:00:00", 3_600_000, false},
		{"12:34:56", (12*3600 + 34*60 + 56) * 1000, false},
		{"23:59:59", (23*3600 + 59*60 + 59) * 1000, false},

		// Invalid cases
		{"0:00:00", 0, true},
		{"00:0:00", 0, true},
		{"00:00:0", 0, true},
		{"24:00:00", 0, true},
		{"00:60:00", 0, true},
		{"00:00:60", 0, true},
		{"aa:bb:cc", 0, true},
		{"00:00", 0, true},
		{"000000", 0, true},
		{"+00:01:00", 0, true},
		{"", 0, true},
		{"   ", 0, true},
	}

	for _, test := range tests {
		ms, err := parseTime(test.in)
		if test.expectError {
			if err == nil {
				ts.Errorf("For time '%s', expected error but got %d", test.in, ms)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For time '%s', unexpected error: %v", test.in, err)
			continue
		}
		if ms != test.expectedMs {
			ts.Errorf("For time '%s', expected %d but got %d", test.in, test.expectedMs, ms)
		}
	}
}

func TestParseTimelineTransitionDiagnostic(ts *testing.T) {
	ctx := NewTextParser("00:00:05 alpha smooh")
	_, err := ctx.ParseTimelineDeclaration()
	if err == nil {
		ts.Fatal("expected transition diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Found != "smooh" {
		ts.Fatalf("expected found transition smooh, got %q", diagnostic.Found)
	}
	if diagnostic.Suggestion != "did you mean \"smooth\"?" {
		ts.Fatalf("expected smooth suggestion, got %q", diagnostic.Suggestion)
	}
	if diagnostic.Span.Column != 16 || diagnostic.Span.EndColumn != 21 {
		ts.Fatalf("expected transition span 17..22, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}

func TestParseTimelineMissingPresetDiagnostic(ts *testing.T) {
	ctx := NewTextParser("00:00:05")

	_, err := ctx.ParseTimelineDeclaration()
	if err == nil {
		ts.Fatal("expected missing preset diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Message != "unexpected end of line" {
		ts.Fatalf("expected EOF message, got %q", diagnostic.Message)
	}
	if len(diagnostic.Expected) != 1 || diagnostic.Expected[0] != "preset name" {
		ts.Fatalf("expected preset name expectation, got %#v", diagnostic.Expected)
	}
	if diagnostic.Span.Column != 9 || diagnostic.Span.EndColumn != 10 {
		ts.Fatalf("expected EOF span 9..10, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}
