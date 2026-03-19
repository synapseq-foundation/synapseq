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

package sequence

import (
	"testing"
)

func TestNewSequenceFile(ts *testing.T) {
	content := []byte("line 1\nline 2\nline 3\n")
	sf := NewSequenceFile(content)

	if sf == nil {
		ts.Fatal("NewSequenceFile returned nil")
	}
	if sf.scanner == nil {
		ts.Fatal("scanner not initialized")
	}
	if sf.CurrentLineNumber() != 0 {
		ts.Errorf("expected CurrentLineNumber=0, got %d", sf.CurrentLineNumber())
	}
}

func TestSequenceFile_NextLine(ts *testing.T) {
	tests := []struct {
		name      string
		content   []byte
		wantLines []string
	}{
		{
			name:      "simple lines",
			content:   []byte("a\nb\nc\n"),
			wantLines: []string{"a", "b", "c"},
		},
		{
			name:      "empty lines",
			content:   []byte("a\n\nb\n"),
			wantLines: []string{"a", "", "b"},
		},
		{
			name:      "no trailing newline",
			content:   []byte("a\nb"),
			wantLines: []string{"a", "b"},
		},
		{
			name:      "empty content",
			content:   []byte(""),
			wantLines: []string{},
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func(ts *testing.T) {
			sf := NewSequenceFile(tt.content)

			var got []string
			for sf.NextLine() {
				got = append(got, sf.CurrentLine())
			}

			if len(got) != len(tt.wantLines) {
				ts.Fatalf("expected %d lines, got %d", len(tt.wantLines), len(got))
			}

			for i, want := range tt.wantLines {
				if got[i] != want {
					ts.Errorf("line %d: expected %q, got %q", i+1, want, got[i])
				}
			}

			if sf.CurrentLineNumber() != len(tt.wantLines) {
				ts.Errorf("expected CurrentLineNumber=%d, got %d", len(tt.wantLines), sf.CurrentLineNumber())
			}
		})
	}
}

func TestSequenceFile_NextLine_NilScanner(ts *testing.T) {
	sf := &SequenceFile{scanner: nil}

	if sf.NextLine() {
		ts.Error("expected NextLine()=false with nil scanner")
	}
}

func TestSequenceFile_LineNumberIncrement(ts *testing.T) {
	content := []byte("1\n2\n3\n")
	sf := NewSequenceFile(content)

	expectedLineNums := []int{1, 2, 3}
	for i := 0; sf.NextLine(); i++ {
		if sf.CurrentLineNumber() != expectedLineNums[i] {
			ts.Errorf("iteration %d: expected line number %d, got %d",
				i, expectedLineNums[i], sf.CurrentLineNumber())
		}
	}
}
