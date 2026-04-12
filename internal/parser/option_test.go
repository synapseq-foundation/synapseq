//go:build !wasm

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
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestHasOption(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("%svolume 50", t.KeywordOption), true},
		{fmt.Sprintf("%ssamplerate 48000", t.KeywordOption), true},
		{fmt.Sprintf("   %sambiance rain file.wav", t.KeywordOption), false},
		{fmt.Sprintf("ambiance rain file.wav %s", t.KeywordComment), false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasOption()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasOption() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseOption(ts *testing.T) {
	tests := []struct {
		line     string
		expected *t.ParseOptions
	}{
		{
			fmt.Sprintf("%svolume 50", t.KeywordOption),
			&t.ParseOptions{Values: map[string]string{t.KeywordOptionVolume: "50"}, Ambiance: map[string]string{}, Extends: []string{}},
		},
		{
			fmt.Sprintf("%ssamplerate 48000", t.KeywordOption),
			&t.ParseOptions{Values: map[string]string{t.KeywordOptionSampleRate: "48000"}, Ambiance: map[string]string{}, Extends: []string{}},
		},
		{
			fmt.Sprintf("%s%s rain testdata/noise", t.KeywordOption, t.KeywordOptionAmbiance),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{"rain": "testdata/noise"}, Extends: []string{}},
		},
		{
			fmt.Sprintf("%s%s shared/base", t.KeywordOption, t.KeywordOptionExtends),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{}, Extends: []string{"shared/base"}},
		},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)

		parsed, err := ctx.ParseOption("")
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}

		if !reflect.DeepEqual(parsed, test.expected) {
			ts.Errorf("For line '%s', expected option %+v but got %+v",
				test.line, test.expected, parsed)
		}
	}
}

func TestParseOptionErrors(ts *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantErrText string
	}{
		{
			name:        "unexpected extra token after volume",
			line:        fmt.Sprintf("%svolume 50 extra", t.KeywordOption),
			wantErrText: "unexpected token after option definition",
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func(ts *testing.T) {
			ctx := NewTextParser(test.line)

				_, err := ctx.ParseOption("")
			if err == nil {
				ts.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.wantErrText) {
				ts.Fatalf("expected error containing %q, got %v", test.wantErrText, err)
			}
		})
	}
}

func TestParseOptionTypoDiagnostic(ts *testing.T) {
	ctx := NewTextParser("@volum 50")

	_, err := ctx.ParseOption("")
	if err == nil {
		ts.Fatal("expected option diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Message != "invalid option" {
		ts.Fatalf("expected invalid option message, got %q", diagnostic.Message)
	}
	if diagnostic.Found != "volum" {
		ts.Fatalf("expected found option volum, got %q", diagnostic.Found)
	}
	if diagnostic.Suggestion != "did you mean \"volume\"?" {
		ts.Fatalf("expected volume suggestion, got %q", diagnostic.Suggestion)
	}
	if diagnostic.Span.Column != 1 || diagnostic.Span.EndColumn != 7 {
		ts.Fatalf("expected option span 1..7, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}
