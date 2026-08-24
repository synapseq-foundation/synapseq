// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
			&t.ParseOptions{Values: map[string]string{t.KeywordOptionVolume: "50"}, Ambiance: map[string]string{}, Music: map[string]string{}, Extends: []string{}, Waveforms: []t.WaveformDefinition{}},
		},
		{
			fmt.Sprintf("%ssamplerate 48000", t.KeywordOption),
			&t.ParseOptions{Values: map[string]string{t.KeywordOptionSampleRate: "48000"}, Ambiance: map[string]string{}, Music: map[string]string{}, Extends: []string{}, Waveforms: []t.WaveformDefinition{}},
		},
		{
			fmt.Sprintf("%s%s rain testdata/noise", t.KeywordOption, t.KeywordOptionAmbiance),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{"rain": "testdata/noise"}, Music: map[string]string{}, Extends: []string{}, Waveforms: []t.WaveformDefinition{}},
		},
		{
			fmt.Sprintf("%s%s meditation audio/meditation", t.KeywordOption, t.KeywordOptionMusic),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{}, Music: map[string]string{"meditation": "audio/meditation"}, Extends: []string{}, Waveforms: []t.WaveformDefinition{}},
		},
		{
			fmt.Sprintf("%s%s shared/base", t.KeywordOption, t.KeywordOptionExtends),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{}, Music: map[string]string{}, Extends: []string{"shared/base"}, Waveforms: []t.WaveformDefinition{}},
		},
		{
			"@waveform softpulse 0 20 50 100",
			&t.ParseOptions{
				Values:   map[string]string{},
				Ambiance: map[string]string{},
				Music:    map[string]string{},
				Extends:  []string{},
				Waveforms: []t.WaveformDefinition{{
					Name: "softpulse",
					Points: []float64{
						t.NormalizeWaveformPoint(0),
						t.NormalizeWaveformPoint(20),
						t.NormalizeWaveformPoint(50),
						t.NormalizeWaveformPoint(100),
					},
				}},
			},
		},
		{
			"@transition soft-land 0 20 50 100",
			&t.ParseOptions{
				Values:    map[string]string{},
				Ambiance:  map[string]string{},
				Music:     map[string]string{},
				Extends:   []string{},
				Waveforms: []t.WaveformDefinition{},
				Transitions: []t.TransitionDefinition{{
					Name:   "soft-land",
					Points: []float64{0, 0.2, 0.5, 1},
				}},
			},
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
		{name: "missing waveform name", line: "@waveform", wantErrText: "waveform name"},
		{name: "transition requires endpoints", line: "@transition soft 0 80", wantErrText: "transition must start at 0 and end at 100"},
		{name: "transition points are monotonic", line: "@transition soft 0 80 70 100", wantErrText: "transition points must be monotonic"},
		{name: "transition name is reserved", line: "@transition smooth 0 100", wantErrText: "reserved for a built-in transition"},
		{name: "empty waveform", line: "@waveform pulse", wantErrText: "at least 2 points"},
		{name: "one waveform point", line: "@waveform pulse 50", wantErrText: "at least 2 points"},
		{name: "non-numeric waveform point", line: "@waveform pulse 0 nope", wantErrText: "invalid float"},
		{name: "negative waveform point", line: "@waveform pulse -1 100", wantErrText: "between 0 and 100"},
		{name: "high waveform point", line: "@waveform pulse 0 101", wantErrText: "between 0 and 100"},
		{name: "built-in waveform name", line: "@waveform sine 0 100", wantErrText: "reserved for a built-in waveform"},
		{name: "case-insensitive built-in name", line: "@waveform SQUARE 0 100", wantErrText: "reserved for a built-in waveform"},
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

func TestParseOptionRejectsTooManyWaveformPoints(ts *testing.T) {
	line := "@waveform pulse " + strings.TrimSpace(strings.Repeat("0 ", t.MaxWaveformPoints+1))
	_, err := NewTextParser(line).ParseOption("")
	if err == nil || !strings.Contains(err.Error(), "cannot contain more than") {
		ts.Fatalf("expected maximum waveform point error, got %v", err)
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
