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
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	nr "github.com/synapseq-foundation/synapseq/v4/internal/nameref"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type ParsedPresetDeclaration struct {
	Name       string
	IsTemplate bool
	FromName   string
}

// HasPreset checks if the current line is a preset definition
func (ctx *TextParser) HasPreset() bool {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	ch := tok[0]
	if ln[0] != ' ' && ((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
		return true
	}

	return false
}

func (ctx *TextParser) ParsePresetDeclaration() (*ParsedPresetDeclaration, error) {
	presetName, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "preset name")
	}
	presetSpan, _ := ctx.Line.LastTokenSpan()

	decl := &ParsedPresetDeclaration{Name: presetName}

	tok, ok := ctx.Line.NextToken()
	if ok {
		switch tok {
		case t.KeywordFrom:
			fromPresetName, ok := ctx.Line.NextToken()
			if !ok {
				return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "preset name")
			}
			decl.FromName = fromPresetName
		case t.KeywordAs:
			_, err := ctx.Line.NextExpectOneOf(t.KeywordTemplate)
			if err != nil {
				return nil, err
			}
			decl.IsTemplate = true
		default:
			ctx.Line.RewindToken(1)
		}
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after preset definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	if err := nr.IsValid(presetName); err != nil {
		return nil, diag.Validation(err.Error()).WithSpan(presetSpan).WithFound(presetName).WithCause(err)
	}

	return decl, nil
}
