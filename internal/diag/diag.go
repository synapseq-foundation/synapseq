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
	"fmt"
	"strings"
)

// Kind identifies the high-level category of a diagnostic.
type Kind string

const (
	KindParse      Kind = "parse"
	KindValidation Kind = "validation"
	KindIO         Kind = "io"
	KindInternal   Kind = "internal"
)

// Span identifies a location in source text.
// Columns are 1-based and EndColumn is exclusive.
type Span struct {
	File      string
	Line      int
	Column    int
	EndColumn int
	LineText  string
}

// HasLocation reports whether the span contains a usable source position.
func (s Span) HasLocation() bool {
	return s.Line > 0 || s.Column > 0 || s.File != ""
}

func (s Span) normalized() Span {
	if s.Column < 1 {
		s.Column = 1
	}
	if s.EndColumn < s.Column+1 {
		s.EndColumn = s.Column + 1
	}
	return s
}

// Width returns the visual width of the span.
func (s Span) Width() int {
	s = s.normalized()
	return s.EndColumn - s.Column
}

// Diagnostic is a structured user-facing error with source location.
type Diagnostic struct {
	Kind       Kind
	Message    string
	Span       Span
	Found      string
	Expected   []string
	Suggestion string
	Hint       string
	Cause      error
}

// New creates a new diagnostic of the given kind.
func New(kind Kind, message string) *Diagnostic {
	return &Diagnostic{
		Kind:    kind,
		Message: message,
	}
}

// Parse creates a new parser diagnostic.
func Parse(message string) *Diagnostic {
	return New(KindParse, message)
}

// Validation creates a new validation diagnostic.
func Validation(message string) *Diagnostic {
	return New(KindValidation, message)
}

// Wrap creates a diagnostic that keeps the original error attached.
func Wrap(kind Kind, message string, err error) *Diagnostic {
	return New(kind, message).WithCause(err)
}

// As extracts a diagnostic from err if present.
func As(err error) (*Diagnostic, bool) {
	var diagnostic *Diagnostic
	if errors.As(err, &diagnostic) {
		return diagnostic, true
	}
	return nil, false
}

// UnexpectedToken creates a parser diagnostic for an invalid token.
func UnexpectedToken(span Span, found string, expected ...string) *Diagnostic {
	diagnostic := Parse("unexpected token").WithSpan(span).WithFound(found).WithExpected(expected...)
	if suggestion, ok := ClosestMatch(found, expected, DefaultSuggestionDistance(found)); ok {
		diagnostic.WithSuggestion(fmt.Sprintf("did you mean %q?", suggestion))
	}
	return diagnostic
}

// UnexpectedEOF creates a parser diagnostic for a missing token.
func UnexpectedEOF(span Span, expected ...string) *Diagnostic {
	return Parse("unexpected end of line").WithSpan(span).WithExpected(expected...)
}

// WithSpan attaches source location to the diagnostic.
func (d *Diagnostic) WithSpan(span Span) *Diagnostic {
	d.Span = span
	return d
}

// WithFound sets the token or value that was encountered.
func (d *Diagnostic) WithFound(found string) *Diagnostic {
	d.Found = found
	return d
}

// WithExpected sets the acceptable values for this diagnostic.
func (d *Diagnostic) WithExpected(expected ...string) *Diagnostic {
	d.Expected = append([]string(nil), expected...)
	return d
}

// WithSuggestion adds a fix suggestion.
func (d *Diagnostic) WithSuggestion(suggestion string) *Diagnostic {
	d.Suggestion = suggestion
	return d
}

// WithHint adds an extra hint.
func (d *Diagnostic) WithHint(hint string) *Diagnostic {
	d.Hint = hint
	return d
}

// WithCause attaches the wrapped cause.
func (d *Diagnostic) WithCause(err error) *Diagnostic {
	d.Cause = err
	return d
}

// Unwrap returns the wrapped cause.
func (d *Diagnostic) Unwrap() error {
	return d.Cause
}

// Error returns a concise single-line representation.
func (d *Diagnostic) Error() string {
	if d == nil {
		return "<nil>"
	}

	parts := make([]string, 0, 4)
	if location := formatLocation(d.Span); location != "" {
		parts = append(parts, location)
	}
	if d.Message != "" {
		parts = append(parts, d.Message)
	}
	if d.Found != "" {
		parts = append(parts, fmt.Sprintf("found %q", d.Found))
	}
	if len(d.Expected) > 0 {
		parts = append(parts, fmt.Sprintf("expected %s", quoteList(d.Expected)))
	}
	if len(parts) == 0 && d.Cause != nil {
		return d.Cause.Error()
	}
	return strings.Join(parts, ": ")
}

// FormatHuman returns a multi-line user-facing rendering.
func (d *Diagnostic) FormatHuman() string {
	if d == nil {
		return "<nil>"
	}

	var lines []string
	if location := formatLocation(d.Span); location != "" {
		lines = append(lines, fmt.Sprintf("%s: %s", location, d.Message))
	} else if d.Message != "" {
		lines = append(lines, d.Message)
	}

	if d.Span.LineText != "" {
		lines = append(lines, d.Span.LineText)
		lines = append(lines, caretLine(d.Span))
	}

	if d.Found != "" {
		lines = append(lines, fmt.Sprintf("found: %q", d.Found))
	}
	if len(d.Expected) > 0 {
		lines = append(lines, fmt.Sprintf("expected: %s", quoteList(d.Expected)))
	}
	if d.Suggestion != "" {
		lines = append(lines, d.Suggestion)
	}
	if d.Hint != "" {
		lines = append(lines, d.Hint)
	}
	if len(lines) == 0 && d.Cause != nil {
		return d.Cause.Error()
	}
	return strings.Join(lines, "\n")
}

func formatLocation(span Span) string {
	span = span.normalized()
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

func caretLine(span Span) string {
	span = span.normalized()
	return strings.Repeat(" ", span.Column-1) + strings.Repeat("^", span.Width())
}

func quoteList(values []string) string {
	if len(values) == 0 {
		return ""
	}

	quoted := make([]string, len(values))
	for i, value := range values {
		quoted[i] = fmt.Sprintf("%q", value)
	}

	if len(quoted) == 1 {
		return quoted[0]
	}
	if len(quoted) == 2 {
		return quoted[0] + " or " + quoted[1]
	}

	return strings.Join(quoted[:len(quoted)-1], ", ") + ", or " + quoted[len(quoted)-1]
}
