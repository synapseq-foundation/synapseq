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
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestHasComment(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("%s This is a comment", t.KeywordComment), true},
		{"No comment here", false},
		{fmt.Sprintf("%sComment without space", t.KeywordComment), true},
		{fmt.Sprintf("   %s Indented comment", t.KeywordComment), true},
		{fmt.Sprintf("%s Double Comment!", strings.Repeat(t.KeywordComment, 2)), true},
		{fmt.Sprintf("  %s Indented double Comment!", strings.Repeat(t.KeywordComment, 2)), true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasComment()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasComment() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseComment(ts *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{fmt.Sprintf("%s This is a comment", t.KeywordComment), ""},
		{"No comment here", ""},
		{fmt.Sprintf("%sComment without space", t.KeywordComment), ""},
		{fmt.Sprintf("   %s Indented comment", t.KeywordComment), ""},
		{fmt.Sprintf("%s Double Comment!", strings.Repeat(t.KeywordComment, 2)), "Double Comment!"},
		{fmt.Sprintf("  %s Indented double Comment!", strings.Repeat(t.KeywordComment, 2)), "Indented double Comment!"},
		{strings.Repeat(t.KeywordComment, 2), " "},
		{fmt.Sprintf("%s First part // not a comment", t.KeywordComment), ""},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.ParseComment()
		if result != test.expected {
			ts.Errorf("For line '%s', expected ParseComment() to be '%s' but got '%s'", test.line, test.expected, result)
		}
	}
}
