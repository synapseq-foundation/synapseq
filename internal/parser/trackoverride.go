// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	nr "github.com/synapseq-foundation/synapseq/v4/internal/nameref"
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
		t.KeywordLeft,
		t.KeywordRight,
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
		decl.RawAmplitude[0], _ = ctx.Line.Peek()
		decl.Amplitude[0], err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.AmplitudeSpan[0], _ = ctx.Line.LastTokenSpan()
		decl.RelativeAmplitude[0] = isRelative(decl.RawAmplitude[0])
		decl.Amplitude[1] = decl.Amplitude[0]
		decl.RawAmplitude[1] = decl.RawAmplitude[0]
		decl.AmplitudeSpan[1] = decl.AmplitudeSpan[0]
		decl.RelativeAmplitude[1] = decl.RelativeAmplitude[0]
		decl.Value = decl.Amplitude[0]
		decl.RawValue = decl.RawAmplitude[0]
		decl.Relative = decl.RelativeAmplitude[0]
		decl.ValueSpan = decl.AmplitudeSpan[0]
	case t.KeywordLeft, t.KeywordRight:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = isRelative(decl.RawValue)
	case t.KeywordIntensity:
		decl.RawValue, _ = ctx.Line.Peek()
		decl.Value, err = ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		decl.Relative = decl.RawValue != "" && (decl.RawValue[0] == '+' || decl.RawValue[0] == '-')
	case t.KeywordWaveform:
		waveform, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "waveform name")
		}
		decl.ValueSpan, _ = ctx.Line.LastTokenSpan()
		if err := nr.IsValid(waveform); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(decl.ValueSpan).WithFound(waveform).WithCause(err)
		}
		decl.Waveform = t.WaveformName(waveform)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after track override definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	return decl, nil
}

func isRelative(value string) bool {
	return value != "" && (value[0] == '+' || value[0] == '-')
}
