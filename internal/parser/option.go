//go:build !wasm

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

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	nr "github.com/synapseq-foundation/synapseq/v4/internal/nameref"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HasOption checks if the first element is an option.
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and returns raw parsed option values.
func (ctx *TextParser) ParseOption(_ string) (*t.ParseOptions, error) {
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "option")
	}
	span, _ := ctx.Line.LastTokenSpan()

	if string(tok[0]) != t.KeywordOption {
		return nil, diag.Parse("expected option").WithSpan(span).WithFound(tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return nil, diag.Parse("expected option name").WithSpan(span).WithFound(tok)
	}

	parsed := t.NewParseOptions()
	validOptions := []string{
		t.KeywordOptionSampleRate,
		t.KeywordOptionVolume,
		t.KeywordOptionAmbiance,
		t.KeywordOptionExtends,
	}

	switch option {
	case t.KeywordOptionSampleRate:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "samplerate value")
		}

		parsed.Values[t.KeywordOptionSampleRate] = value
	case t.KeywordOptionVolume:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "volume value")
		}

		parsed.Values[t.KeywordOptionVolume] = value
	case t.KeywordOptionAmbiance:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "ambiance name")
		}
		nameSpan, _ := ctx.Line.LastTokenSpan()

		if err := nr.IsValid(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			content = name // allow shorthand ambiance name as path
		}

		parsed.Ambiance[name] = content
	case t.KeywordOptionExtends:
		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "extends path")
		}

		parsed.Extends = append(parsed.Extends, content)
	default:
		diagnostic := diag.Parse("invalid option").WithSpan(span).WithFound(option).WithExpected(validOptions...)
		if suggestion, ok := diag.ClosestMatch(option, validOptions, diag.DefaultSuggestionDistance(option)); ok {
			diagnostic.WithSuggestion(fmt.Sprintf("did you mean %q?", suggestion))
		}
		return nil, diagnostic
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after option definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	return parsed, nil
}
