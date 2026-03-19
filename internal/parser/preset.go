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
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

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

// ParsePreset extracts and returns a Preset from the current line context
func (ctx *TextParser) ParsePreset(presets *[]t.Preset) (*t.Preset, error) {
	presetName, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "preset name")
	}
	presetSpan, _ := ctx.Line.LastTokenSpan()

	var fromPreset *t.Preset
	isTemplate := false

	tok, ok := ctx.Line.NextToken()
	if ok {
		switch tok {
		case t.KeywordFrom:
			fromPresetName, ok := ctx.Line.NextToken()
			if !ok {
				return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "preset name")
			}
			fromSpan, _ := ctx.Line.LastTokenSpan()

			fromPreset = s.FindPreset(fromPresetName, *presets)
			if fromPreset == nil {
				return nil, diag.Validation(fmt.Sprintf("unknown preset to inherit from: %q", fromPresetName)).WithSpan(fromSpan).WithFound(fromPresetName)
			}

			if !fromPreset.IsTemplate {
				return nil, diag.Validation(fmt.Sprintf("can only inherit from a template preset, but %q is not a template", fromPresetName)).WithSpan(fromSpan).WithFound(fromPresetName)
			}
		case t.KeywordAs:
			// "as template" clause
			_, err := ctx.Line.NextExpectOneOf(t.KeywordTemplate)
			if err != nil {
				return nil, err
			}
			isTemplate = true
		default:
			ctx.Line.RewindToken(1) // Un-consume the token
		}
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after preset definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	if err := s.IsValidNamedRef(presetName); err != nil {
		return nil, diag.Validation(err.Error()).WithSpan(presetSpan).WithFound(presetName).WithCause(err)
	}

	preset, err := t.NewPreset(presetName, isTemplate, fromPreset)
	if err != nil {
		return nil, diag.Validation(err.Error()).WithSpan(presetSpan).WithFound(presetName).WithCause(err)
	}
	return preset, nil
}
