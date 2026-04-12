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
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

func formatCLIError(err error) string {
	if diagnostic, ok := diag.As(err); ok {
		var lines []string
		location := formatDiagnosticLocation(diagnostic)
		if location != "" {
			lines = append(lines, cli.ErrorText("synapseq:")+" "+cli.Accent(location)+": "+cli.ErrorText(diagnostic.Message))
		} else {
			lines = append(lines, cli.ErrorText("synapseq:")+" "+cli.ErrorText(diagnostic.Message))
		}

		if diagnostic.Span.LineText != "" {
			lines = append(lines, color.New(color.FgWhite).Sprint(diagnostic.Span.LineText))
			lines = append(lines, cli.ErrorText(strings.Repeat(" ", diagnostic.Span.Column-1)+strings.Repeat("^", max(1, diagnostic.Span.EndColumn-diagnostic.Span.Column))))
		}
		if diagnostic.Found != "" {
			lines = append(lines, cli.Label("found:")+" "+cli.Accent(fmt.Sprintf("%q", diagnostic.Found)))
		}
		if len(diagnostic.Expected) > 0 {
			lines = append(lines, cli.Label("expected:")+" "+formatExpectedList(diagnostic.Expected))
		}
		if diagnostic.Suggestion != "" {
			lines = append(lines, cli.SuccessText(diagnostic.Suggestion))
		}
		if diagnostic.Hint != "" {
			lines = append(lines, cli.Muted(diagnostic.Hint))
		}
		if cause := formatDiagnosticCause(diagnostic); cause != "" {
			lines = append(lines, cli.Label("cause:")+" "+cause)
		}

		return strings.Join(lines, "\n")
	}

	return cli.ErrorText("synapseq:") + " " + err.Error()
}

func formatDiagnosticCause(diagnostic *diag.Diagnostic) string {
	if diagnostic == nil {
		return ""
	}

	seen := map[string]struct{}{}
	for cause := diagnostic.Cause; cause != nil; cause = errors.Unwrap(cause) {
		text := strings.TrimSpace(cause.Error())
		if text == "" || text == diagnostic.Message {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		return text
	}

	return ""
}

func formatDiagnosticLocation(diagnostic *diag.Diagnostic) string {
	if diagnostic == nil {
		return ""
	}

	span := diagnostic.Span
	switch {
	case span.File != "" && span.Line > 0 && span.Column > 0:
		return fmt.Sprintf("%s:%d:%d", span.File, span.Line, span.Column)
	case span.Line > 0 && span.Column > 0:
		return fmt.Sprintf("line %d:%d", span.Line, span.Column)
	case span.Line > 0:
		return fmt.Sprintf("line %d", span.Line)
	case span.File != "":
		return span.File
	default:
		return ""
	}
}

func formatExpectedList(values []string) string {
	colored := make([]string, len(values))
	for i, value := range values {
		colored[i] = cli.Accent(fmt.Sprintf("%q", value))
	}
	if len(colored) == 1 {
		return colored[0]
	}
	if len(colored) == 2 {
		return colored[0] + " or " + colored[1]
	}
	return strings.Join(colored[:len(colored)-1], ", ") + ", or " + colored[len(colored)-1]
}