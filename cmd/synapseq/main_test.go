//go:build !js && !wasm

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

package main

import (
	"errors"
	"strings"
	"testing"

	clistyle "github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

func TestFormatCLIErrorDiagnostic(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := diag.UnexpectedToken(
		diag.Span{
			File:      "example.spsq",
			Line:      4,
			Column:    12,
			EndColumn: 19,
			LineText:  "  tone 300 binaual 10 amplitude 10",
		},
		"binaual",
		"binaural",
		"monaural",
	)

	formatted := formatCLIError(err)

	checks := []string{
		"synapseq: example.spsq:4:12: unexpected token",
		"  tone 300 binaual 10 amplitude 10",
		"           ^^^^^^^",
		"did you mean \"binaural\"?",
	}

	for _, check := range checks {
		if !strings.Contains(formatted, check) {
			ts.Fatalf("expected formatted CLI error to contain %q, got:\n%s", check, formatted)
		}
	}
}

func TestFormatCLIErrorFallback(ts *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	err := errors.New("plain error")
	formatted := formatCLIError(err)
	if formatted != "synapseq: plain error" {
		ts.Fatalf("unexpected fallback formatting: %q", formatted)
	}
}
