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
	p "github.com/synapseq-foundation/synapseq/v4/internal/preset"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HasTrackOverride checks if the current line is a track override definition
func (ctx *TextParser) HasTrackOverride() bool {
	ln := ctx.Line.Raw
	if len(ln) < 3 {
		return false
	}

	if ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' ' {
		tok, ok := ctx.Line.Peek()
		if !ok || tok != t.KeywordTrack {
			return false
		}
		return true
	}

	return false
}

func (ctx *TextParser) ParseTrackOverrideDeclaration() (*p.TrackOverrideSpec, error) {
	_, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), t.KeywordTrack)
	}

	trackIdx, err := ctx.Line.NextIntStrict()
	if err != nil {
		return nil, err
	}
	trackSpan, _ := ctx.Line.LastTokenSpan()

	if trackIdx <= 0 || trackIdx >= t.NumberOfChannels {
		return nil, diag.Validation(fmt.Sprintf("track index out of range (1-%d): %d", t.NumberOfChannels-1, trackIdx)).WithSpan(trackSpan).WithFound(fmt.Sprintf("%d", trackIdx))
	}

	kind, err := ctx.Line.NextExpectOneOf(
		t.KeywordTone,
		t.KeywordBinaural,
		t.KeywordMonaural,
		t.KeywordIsochronic,
		t.KeywordWaveform,
		t.KeywordPan,
		t.KeywordModulation,
		t.KeywordDoppler,
		t.KeywordSmooth,
		t.KeywordAmplitude,
		t.KeywordIntensity)
	if err != nil {
		return nil, err
	}
	kindSpan, _ := ctx.Line.LastTokenSpan()

	decl := &p.TrackOverrideSpec{
		TrackIndex: trackIdx,
		TrackSpan:  trackSpan,
		Kind:       kind,
		KindSpan:   kindSpan,
	}

	switch kind {
	case t.KeywordTone:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordPan, t.KeywordModulation, t.KeywordDoppler:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordSmooth:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordAmplitude:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordIntensity:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordWaveform:
		waveform, err := ctx.Line.NextExpectOneOf(
			t.KeywordSine,
			t.KeywordSquare,
			t.KeywordTriangle,
			t.KeywordSawtooth)

		if err != nil {
			return nil, err
		}

		var waveformType t.WaveformType
		switch waveform {
		case t.KeywordSine:
			waveformType = t.WaveformSine
		case t.KeywordSquare:
			waveformType = t.WaveformSquare
		case t.KeywordTriangle:
			waveformType = t.WaveformTriangle
		case t.KeywordSawtooth:
			waveformType = t.WaveformSawtooth
		default:
			return nil, diag.Parse("unexpected waveform type").WithSpan(kindSpan).WithFound(waveform)
		}

		decl.Waveform = waveformType
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after track override definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	return decl, nil
}
