/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"fmt"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestPeek(ts *testing.T) {
	// Create a sample track string line
	trLn := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()

	tests := []struct {
		line          string
		expectedToken string
		expectedOk    bool
	}{
		{trLn, t.KeywordWaveform, true},
		{fmt.Sprintf("   %s", trLn), t.KeywordWaveform, true},
		{"", "", false},
		{"   ", "", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		token, ok := ctx.Line.Peek()
		if token != test.expectedToken || ok != test.expectedOk {
			ts.Errorf("For line '%s', expected Peek() to return ('%s', %v) but got ('%s', %v)", test.line, test.expectedToken, test.expectedOk, token, ok)
		}
	}
}

func TestNextToken(ts *testing.T) {
	// Create a sample track string line
	trLn := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()
	trLnFields := strings.Fields(trLn)

	tests := []struct {
		line           string
		expectedTokens []string
	}{
		{trLn, trLnFields},
		{fmt.Sprintf("   %s", trLn), trLnFields},
		{"", []string{}},
		{"   ", []string{}},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		var tokens []string
		for {
			token, ok := ctx.Line.NextToken()
			if !ok {
				break
			}
			tokens = append(tokens, token)
		}
		if len(tokens) != len(test.expectedTokens) {
			ts.Errorf("For line '%s', expected %d tokens but got %d", test.line, len(test.expectedTokens), len(tokens))
			continue
		}
		for i, expectedToken := range test.expectedTokens {
			if tokens[i] != expectedToken {
				ts.Errorf("For line '%s', expected token %d to be '%s' but got '%s'", test.line, i, expectedToken, tokens[i])
			}
		}
	}
}

func TestTokenSpans(ts *testing.T) {
	ctx := NewTextParser("  tone 300 binaual 10 amplitude 10")

	span, ok := ctx.Line.PeekSpan()
	if !ok {
		ts.Fatal("expected span for first token")
	}
	if span.Column != 3 || span.EndColumn != 7 {
		ts.Fatalf("expected first token span 3..7, got %d..%d", span.Column, span.EndColumn)
	}

	token, ok := ctx.Line.NextToken()
	if !ok || token != t.KeywordTone {
		ts.Fatalf("expected first token %q, got %q", t.KeywordTone, token)
	}

	if _, ok := ctx.Line.NextToken(); !ok {
		ts.Fatal("expected second token")
	}
	if _, ok := ctx.Line.NextToken(); !ok {
		ts.Fatal("expected third token")
	}

	span, ok = ctx.Line.LastTokenSpan()
	if !ok {
		ts.Fatal("expected last token span")
	}
	if span.Column != 12 || span.EndColumn != 19 {
		ts.Fatalf("expected typo token span 12..19, got %d..%d", span.Column, span.EndColumn)
	}
	if span.LineText != "  tone 300 binaual 10 amplitude 10" {
		ts.Fatalf("unexpected line text in span: %q", span.LineText)
	}
}

func TestRewindToken(ts *testing.T) {
	// Create a sample track string line
	trLn := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()
	ctx := NewTextParser(trLn)

	// Read first three tokens
	for range 3 {
		_, ok := ctx.Line.NextToken()
		if !ok {
			ts.Errorf("Unexpected EOF while reading tokens")
		}
	}

	// Rewind two tokens
	ctx.Line.RewindToken(2)

	// Next token should be the second token
	token, ok := ctx.Line.NextToken()
	if !ok || token != t.KeywordSine {
		ts.Errorf("Expected '%s' after rewind, got '%s'", t.KeywordSine, token)
	}

	// Rewind more than available tokens
	ctx.Line.RewindToken(10)

	// Next token should be the first token
	token, ok = ctx.Line.NextToken()
	if !ok || token != t.KeywordWaveform {
		ts.Errorf("Expected '%s' after rewind to start, got '%s'", t.KeywordWaveform, token)
	}
}

func TestNextExpectOneOf(ts *testing.T) {
	trLnTone := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()

	trLnNoisePink := (&t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(40),
	}).String()

	trLnNoiseWhite := (&t.Track{
		Type:      t.TrackWhiteNoise,
		Amplitude: t.AmplitudePercentToRaw(5),
	}).String()

	trLnAmbiance := (&t.Track{
		Type:      t.TrackAmbiance,
		Amplitude: t.AmplitudePercentToRaw(50),
	}).String()

	tests := []struct {
		line          string
		wants         []string
		expectedToken string
		expectError   bool
	}{
		{trLnTone, []string{t.KeywordWaveform, t.KeywordNoise}, t.KeywordWaveform, false},
		{trLnNoisePink, []string{t.KeywordAmbiance, t.KeywordNoise}, t.KeywordNoise, false},
		{trLnNoiseWhite, []string{t.KeywordTriangle, t.KeywordAmbiance}, "", true},
		{trLnAmbiance, []string{t.KeywordNoise, t.KeywordWaveform}, t.KeywordWaveform, false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		token, err := ctx.Line.NextExpectOneOf(test.wants...)
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got token '%s'", test.line, token)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if token != test.expectedToken {
				ts.Errorf("For line '%s', expected token '%s' but got '%s'", test.line, test.expectedToken, token)
			}
		}
	}
}

func TestNextExpectOneOfReturnsDiagnostic(ts *testing.T) {
	ctx := NewTextParser("tone 300 binaual 10")

	for range 2 {
		if _, ok := ctx.Line.NextToken(); !ok {
			ts.Fatal("unexpected EOF preparing test")
		}
	}

	_, err := ctx.Line.NextExpectOneOf(t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic)
	if err == nil {
		ts.Fatal("expected diagnostic error")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Message != "unexpected token" {
		ts.Fatalf("expected unexpected token message, got %q", diagnostic.Message)
	}
	if diagnostic.Found != "binaual" {
		ts.Fatalf("expected found token binaual, got %q", diagnostic.Found)
	}
	if diagnostic.Span.Column != 10 || diagnostic.Span.EndColumn != 17 {
		ts.Fatalf("expected token span 10..17, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
	if diagnostic.Suggestion != "did you mean \"binaural\"?" {
		ts.Fatalf("expected typo suggestion, got %q", diagnostic.Suggestion)
	}
}

func TestFloat64Strict(ts *testing.T) {
	tests := []struct {
		line          string
		expectedValue float64
		expectError   bool
	}{
		{"440.5", 440.5, false},
		{"   123.456   ", 123.456, false},
		{"notanumber", 0, true},
		{"123abc", 0, true},
		{"", 0, true},
		{"   ", 0, true},
		{"NaN", 0, true},
		{"Inf", 0, true},
		{"-Inf", 0, true},
		{"1e10", 1e10, true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		value, err := ctx.Line.NextFloat64Strict()
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got value %f", test.line, value)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if value != test.expectedValue {
				ts.Errorf("For line '%s', expected value %f but got %f", test.line, test.expectedValue, value)
			}
		}
	}
}

func TestFloat64StrictDiagnostic(ts *testing.T) {
	ctx := NewTextParser("tone nope")
	if _, ok := ctx.Line.NextToken(); !ok {
		ts.Fatal("expected first token")
	}

	_, err := ctx.Line.NextFloat64Strict()
	if err == nil {
		ts.Fatal("expected diagnostic error")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Message != "invalid float" {
		ts.Fatalf("expected invalid float message, got %q", diagnostic.Message)
	}
	if diagnostic.Found != "nope" {
		ts.Fatalf("expected found token nope, got %q", diagnostic.Found)
	}
	if diagnostic.Span.Column != 6 || diagnostic.Span.EndColumn != 10 {
		ts.Fatalf("expected invalid float span 6..10, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}

func TestNextIntStrict(ts *testing.T) {
	tests := []struct {
		line          string
		expectedValue int
		expectError   bool
	}{
		{"48000", 48000, false},
		{"80", 80, false},
		{"   123   ", 123, false},
		{"notanumber", 0, true},
		{"123abc", 0, true},
		{"", 0, true},
		{"   ", 0, true},
		{"12.34", 0, true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		value, err := ctx.Line.NextIntStrict()
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got value %d", test.line, value)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if value != test.expectedValue {
				ts.Errorf("For line '%s', expected value %d but got %d", test.line, test.expectedValue, value)
			}
		}
	}
}

func TestNextIntStrictEOFDiagnostic(ts *testing.T) {
	ctx := NewTextParser("track")
	if _, ok := ctx.Line.NextToken(); !ok {
		ts.Fatal("expected first token")
	}

	_, err := ctx.Line.NextIntStrict()
	if err == nil {
		ts.Fatal("expected EOF diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Message != "unexpected end of line" {
		ts.Fatalf("expected EOF message, got %q", diagnostic.Message)
	}
	if diagnostic.Span.Column != 6 || diagnostic.Span.EndColumn != 7 {
		ts.Fatalf("expected EOF span 6..7, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
	if len(diagnostic.Expected) != 1 || diagnostic.Expected[0] != "integer" {
		ts.Fatalf("expected integer expectation, got %#v", diagnostic.Expected)
	}
}
