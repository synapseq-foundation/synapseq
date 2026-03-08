//go:build !wasm

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
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

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
	cwd, err := os.Getwd()
	if err != nil {
		ts.Fatalf("cannot get current working directory: %v", err)
	}

	basePath := filepath.Dir(cwd)

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
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{"rain": filepath.Clean(filepath.Join(basePath, "testdata", "noise.wav"))}, Extends: []string{}},
		},
		{
			fmt.Sprintf("%s%s shared/base", t.KeywordOption, t.KeywordOptionExtends),
			&t.ParseOptions{Values: map[string]string{}, Ambiance: map[string]string{}, Extends: []string{filepath.Clean(filepath.Join(basePath, "shared", "base.spsc"))}},
		},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)

		parsed, err := ctx.ParseOption(basePath)
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
	cwd, err := os.Getwd()
	if err != nil {
		ts.Fatalf("cannot get current working directory: %v", err)
	}

	basePath := filepath.Dir(cwd)

	tests := []struct {
		name        string
		line        string
		wantErrText string
	}{
		{
			name:        "ambiance path with extension rejected",
			line:        fmt.Sprintf("%s%s rain audio/river.wav", t.KeywordOption, t.KeywordOptionAmbiance),
			wantErrText: "ambiance local path must not include file extension",
		},
		{
			name:        "extends path with extension rejected",
			line:        fmt.Sprintf("%s%s shared/base.spsc", t.KeywordOption, t.KeywordOptionExtends),
			wantErrText: "extends local path must not include file extension",
		},
		{
			name:        "absolute path rejected",
			line:        fmt.Sprintf("%s%s rain /tmp/river", t.KeywordOption, t.KeywordOptionAmbiance),
			wantErrText: "absolute paths are not allowed",
		},
		{
			name:        "unexpected extra token after volume",
			line:        fmt.Sprintf("%svolume 50 extra", t.KeywordOption),
			wantErrText: "unexpected token after option definition",
		},
	}

	for _, test := range tests {
		ts.Run(test.name, func(ts *testing.T) {
			ctx := NewTextParser(test.line)

			_, err := ctx.ParseOption(basePath)
			if err == nil {
				ts.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.wantErrText) {
				ts.Fatalf("expected error containing %q, got %v", test.wantErrText, err)
			}
		})
	}
}
