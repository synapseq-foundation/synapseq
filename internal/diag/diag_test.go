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

package diag

import (
	"errors"
	"strings"
	"testing"
)

func TestUnexpectedTokenHumanFormatting(ts *testing.T) {
	diagnostic := UnexpectedToken(
		Span{
			File:      "example.spsq",
			Line:      7,
			Column:    10,
			EndColumn: 17,
			LineText:  "tone 300 binaual 10 amplitude 10",
		},
		"binaual",
		"binaural",
		"monaural",
		"isochronic",
		"effect",
		"amplitude",
	)

	text := diagnostic.FormatHuman()

	checks := []string{
		"example.spsq:7:10: unexpected token",
		"tone 300 binaual 10 amplitude 10",
		"         ^^^^^^^",
		"found: \"binaual\"",
		"expected: \"binaural\", \"monaural\", \"isochronic\", \"effect\", or \"amplitude\"",
		"did you mean \"binaural\"?",
	}

	for _, check := range checks {
		if !strings.Contains(text, check) {
			ts.Fatalf("expected formatted diagnostic to contain %q, got:\n%s", check, text)
		}
	}
}

func TestDiagnosticError(ts *testing.T) {
	diagnostic := Parse("unexpected token").
		WithSpan(Span{Line: 3, Column: 5, EndColumn: 9}).
		WithFound("binaual").
		WithExpected("binaural")

	got := diagnostic.Error()
	want := "line 3:5: unexpected token: found \"binaual\": expected \"binaural\""
	if got != want {
		ts.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDiagnosticUnwrap(ts *testing.T) {
	root := errors.New("root cause")
	diagnostic := Wrap(KindParse, "invalid track", root)

	if !errors.Is(diagnostic, root) {
		ts.Fatal("expected diagnostic to unwrap to root cause")
	}

	extracted, ok := As(diagnostic)
	if !ok || extracted != diagnostic {
		ts.Fatal("expected As to extract the same diagnostic")
	}
}

func TestClosestMatch(ts *testing.T) {
	match, ok := ClosestMatch("binaual", []string{"binaural", "monaural", "isochronic"}, 2)
	if !ok {
		ts.Fatal("expected a typo suggestion")
	}
	if match != "binaural" {
		ts.Fatalf("expected binaural suggestion, got %q", match)
	}

	if _, ok := ClosestMatch("tone", []string{"pan", "noise"}, 1); ok {
		ts.Fatal("did not expect an unrelated suggestion")
	}
}
