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

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

func TestHasOption(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("%svolume 50", t.KeywordOption), true},
		{fmt.Sprintf("%ssamplerate 48000", t.KeywordOption), true},
		{fmt.Sprintf("   %sgainlevel medium", t.KeywordOption), false},
		{fmt.Sprintf("background file.wav %s", t.KeywordComment), false},
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
			t.SequenceOptions{Volume: 50},
		},
		{
			fmt.Sprintf("%ssamplerate 48000", t.KeywordOption),
			t.SequenceOptions{SampleRate: 48000},
		},
		{
			fmt.Sprintf("%sgainlevel low", t.KeywordOption),
			t.SequenceOptions{GainLevel: t.GainLevelLow},
		},
		{
			fmt.Sprintf("%sbackground testdata/%s", t.KeywordOption, backgroundFile),
			t.SequenceOptions{BackgroundList: []string{filepath.Clean(filepath.Join(basePath, "testdata", backgroundFile))}},
		},
		{
			fmt.Sprintf("%sbackground ~/Downloads/%s", t.KeywordOption, backgroundFile),
			t.SequenceOptions{BackgroundList: []string{filepath.Clean(filepath.Join(homeDir, "Downloads", backgroundFile))}},
		},
	}

	for _, test := range tests {
		option := t.SequenceOptions{}
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
