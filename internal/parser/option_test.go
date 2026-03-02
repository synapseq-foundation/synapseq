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
	backgroundFile := "noise.wav"

	cwd, err := os.Getwd()
	if err != nil {
		ts.Fatalf("cannot get current working directory: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		ts.Fatalf("cannot get user home directory: %v", err)
	}

	basePath := filepath.Dir(cwd)

	tests := []struct {
		line     string
		expected t.SequenceOptions
	}{
		{
			fmt.Sprintf("%svolume 50", t.KeywordOption),
			t.SequenceOptions{Volume: 50, Ambiance: map[string]string{}},
		},
		{
			fmt.Sprintf("%ssamplerate 48000", t.KeywordOption),
			t.SequenceOptions{SampleRate: 48000, Ambiance: map[string]string{}},
		},
		{
			fmt.Sprintf("%s%s rain testdata/%s", t.KeywordOption, t.KeywordOptionAmbiance, backgroundFile),
			t.SequenceOptions{Ambiance: map[string]string{"rain": filepath.Clean(filepath.Join(basePath, "testdata", backgroundFile))}},
		},
		{
			fmt.Sprintf("%s%s river ~/Downloads/%s", t.KeywordOption, t.KeywordOptionAmbiance, backgroundFile),
			t.SequenceOptions{Ambiance: map[string]string{"river": filepath.Clean(filepath.Join(homeDir, "Downloads", backgroundFile))}},
		},
	}

	for _, test := range tests {
		option := t.SequenceOptions{Ambiance: map[string]string{}}
		ctx := NewTextParser(test.line)

		if err := ctx.ParseOption(&option, basePath); err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}

		if !reflect.DeepEqual(option, test.expected) {
			ts.Errorf("For line '%s', expected option %+v but got %+v",
				test.line, test.expected, option)
		}
	}
}
