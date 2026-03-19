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
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

// TextParser holds the context for parsing
type TextParser struct {
	Line lineContext // Context for the current line
}

type tokenSpan struct {
	Value     string
	Column    int
	EndColumn int
}

// lineContext holds the context for the current line being parsed
type lineContext struct {
	Raw        string      // Raw line text
	Tokens     []string    // Tokens extracted from the line
	tokenSpans []tokenSpan // Token positions in the raw line
	tkIdx      int         // Current token index
	lastTkIdx  int         // Last consumed token index
}

// Peek retrieves the next token without advancing the index
func (ctx *lineContext) Peek() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		return ctx.Tokens[ctx.tkIdx], true
	}
	return "", false
}

// PeekSpan retrieves the span for the next token without advancing the index.
func (ctx *lineContext) PeekSpan() (diag.Span, bool) {
	if ctx.tkIdx < len(ctx.tokenSpans) {
		return ctx.spanFromToken(ctx.tkIdx), true
	}
	return diag.Span{}, false
}

// NextToken retrieves the next token from the context
func (ctx *lineContext) NextToken() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		token := ctx.Tokens[ctx.tkIdx]
		ctx.lastTkIdx = ctx.tkIdx
		ctx.tkIdx++
		return token, true
	}
	return "", false
}

// LastTokenSpan retrieves the span for the most recently consumed token.
func (ctx *lineContext) LastTokenSpan() (diag.Span, bool) {
	if ctx.lastTkIdx >= 0 && ctx.lastTkIdx < len(ctx.tokenSpans) {
		return ctx.spanFromToken(ctx.lastTkIdx), true
	}
	return diag.Span{}, false
}

// EOFSpan returns a span anchored at the next expected position.
func (ctx *lineContext) EOFSpan() diag.Span {
	if span, ok := ctx.PeekSpan(); ok {
		return span
	}
	if span, ok := ctx.LastTokenSpan(); ok {
		span.Column = span.EndColumn
		span.EndColumn = span.Column + 1
		return span
	}
	return diag.Span{
		Column:    1,
		EndColumn: 2,
		LineText:  ctx.Raw,
	}
}

// RewindToken moves the token index back by n positions
func (ctx *lineContext) RewindToken(n int) {
	if n <= 0 {
		return
	}
	if n > ctx.tkIdx {
		ctx.tkIdx = 0
		ctx.lastTkIdx = -1
		return
	}
	ctx.tkIdx -= n
	ctx.lastTkIdx = ctx.tkIdx - 1
}

// NextExpectOneOf checks if the next token is one of the expected values
func (ctx *lineContext) NextExpectOneOf(wants ...string) (string, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return "", diag.UnexpectedEOF(ctx.EOFSpan(), wants...)
	}

	if slices.Contains(wants, tok) {
		return tok, nil
	}
	span, _ := ctx.LastTokenSpan()
	return "", diag.UnexpectedToken(span, tok, wants...)
}

// NextFloat64Strict retrieves the next token as a float64, enforcing strict parsing
func (ctx *lineContext) NextFloat64Strict() (float64, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, diag.UnexpectedEOF(ctx.EOFSpan(), "float")
	}

	span, _ := ctx.LastTokenSpan()

	f, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, diag.Parse("invalid float").WithSpan(span).WithFound(tok).WithExpected("float")
	}

	// Reject NaN and Inf values
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, diag.Parse("invalid float").WithSpan(span).WithFound(tok).WithExpected("float").WithHint("NaN and Inf are not allowed")
	}

	// Reject scientific notation (e.g., 1e10)
	if strings.ContainsAny(tok, "eE") {
		return 0, diag.Parse("invalid float").WithSpan(span).WithFound(tok).WithExpected("float").WithHint("scientific notation is not allowed")
	}

	return f, nil
}

// NextIntStrict retrieves the next token as an int, enforcing strict parsing
func (ctx *lineContext) NextIntStrict() (int, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, diag.UnexpectedEOF(ctx.EOFSpan(), "integer")
	}

	span, _ := ctx.LastTokenSpan()

	i, err := strconv.Atoi(tok)
	if err != nil {
		return 0, diag.Parse("invalid integer").WithSpan(span).WithFound(tok).WithExpected("integer")
	}
	return i, nil
}

// NewTextParser creates a new TextParser for the given line
func NewTextParser(line string) *TextParser {
	tokens, spans := scanTokens(line)
	return &TextParser{
		Line: lineContext{
			Raw:        line,
			Tokens:     tokens,
			tokenSpans: spans,
			tkIdx:      0,
			lastTkIdx:  -1,
		},
	}
}

func (ctx *lineContext) spanFromToken(index int) diag.Span {
	token := ctx.tokenSpans[index]
	return diag.Span{
		Column:    token.Column,
		EndColumn: token.EndColumn,
		LineText:  ctx.Raw,
	}
}

func scanTokens(line string) ([]string, []tokenSpan) {
	tokens := make([]string, 0)
	spans := make([]tokenSpan, 0)

	for i := 0; i < len(line); {
		r, size := utf8.DecodeRuneInString(line[i:])
		if unicode.IsSpace(r) {
			i += size
			continue
		}

		start := i
		for i < len(line) {
			r, size = utf8.DecodeRuneInString(line[i:])
			if unicode.IsSpace(r) {
				break
			}
			i += size
		}

		tokens = append(tokens, line[start:i])
		spans = append(spans, tokenSpan{
			Value:     line[start:i],
			Column:    start + 1,
			EndColumn: i + 1,
		})
	}

	return tokens, spans
}
