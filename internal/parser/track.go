// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	nr "github.com/synapseq-foundation/synapseq/v4/internal/nameref"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type ParsedTrackDeclaration struct {
	Type                   t.TrackType
	Carrier                float64
	Resonance              float64
	AmplitudePercent       float64
	RightAmplitudePercent  float64
	HasRightAmplitude      bool
	NoiseSmooth            float64
	Waveform               t.WaveformName
	WaveformSpan           diag.Span
	SourceName             string
	EffectType             t.EffectType
	EffectValue            float64
	EffectIntensityPercent float64
}

// HasTrack checks if the current line is a track definition
func (ctx *TextParser) HasTrack() bool {
	ln := ctx.Line.Raw

	if len(ln) < 3 {
		return false
	}

	if ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' ' {
		tok, ok := ctx.Line.Peek()
		if !ok || tok == t.KeywordTrack {
			return false
		}
		return true
	}

	return false
}

func (ctx *TextParser) ParseTrackDeclaration() (*ParsedTrackDeclaration, error) {
	waveform := t.WaveformSine
	var waveformSpan diag.Span

	if tok, ok := ctx.Line.Peek(); ok && tok == t.KeywordWaveform {
		ctx.Line.NextToken() // skip "waveform"

		wfTok, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "waveform name")
		}
		wfSpan, _ := ctx.Line.LastTokenSpan()
		if err := nr.IsValid(wfTok); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(wfSpan).WithFound(wfTok).WithCause(err)
		}
		waveform = t.WaveformName(wfTok)
		waveformSpan = wfSpan

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordAmbiance, t.KeywordMusic); err != nil {
			return nil, err
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance, t.KeywordMusic)
	}

	decl := &ParsedTrackDeclaration{
		Waveform:   waveform,
		EffectType: t.EffectOff,
	}
	if waveformSpan.HasLocation() {
		decl.WaveformSpan = waveformSpan
	}
	switch first {
	case t.KeywordTone:
		var err error
		if decl.Carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, err
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, err
		}

		switch kind {
		case t.KeywordBinaural:
			decl.Type = t.TrackBinauralBeat
		case t.KeywordMonaural:
			decl.Type = t.TrackMonauralBeat
		case t.KeywordIsochronic:
			decl.Type = t.TrackIsochronicBeat
		default:
			decl.Type = t.TrackPureTone
		}

		if decl.Type == t.TrackBinauralBeat ||
			decl.Type == t.TrackMonauralBeat ||
			decl.Type == t.TrackIsochronicBeat {
			if decl.Resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, err
			}

			kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
			if err != nil {
				return nil, err
			}
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation, t.KeywordDoppler)
			if err != nil {
				return nil, err
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectValue = effectValue

			switch effectKind {
			case t.KeywordPan:
				decl.EffectType = t.EffectPan
			case t.KeywordModulation:
				decl.EffectType = t.EffectModulation
			case t.KeywordDoppler:
				decl.EffectType = t.EffectDoppler
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectIntensityPercent = intensity

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, err
			}
		}
	case t.KeywordNoise:
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordWhite, t.KeywordPink, t.KeywordBrown)
		if err != nil {
			return nil, err
		}

		switch kind {
		case t.KeywordWhite:
			decl.Type = t.TrackWhiteNoise
		case t.KeywordPink:
			decl.Type = t.TrackPinkNoise
		case t.KeywordBrown:
			decl.Type = t.TrackBrownNoise
		}

		kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordSmooth, t.KeywordAmplitude)
		if err != nil {
			return nil, err
		}

		if kind == t.KeywordSmooth {
			decl.NoiseSmooth, err = ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
			if err != nil {
				return nil, err
			}
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation)
			if err != nil {
				return nil, err
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectValue = effectValue

			switch effectKind {
			case t.KeywordPan:
				decl.EffectType = t.EffectPan
			case t.KeywordModulation:
				decl.EffectType = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectIntensityPercent = intensity

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, err
			}
		}
	case t.KeywordAmbiance, t.KeywordMusic:
		trackKeyword := first
		if first == t.KeywordAmbiance {
			decl.Type = t.TrackAmbiance
		} else {
			decl.Type = t.TrackMusic
		}

		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), trackKeyword+" name")
		}

		if err := nr.IsValid(name); err != nil {
			span, _ := ctx.Line.LastTokenSpan()
			return nil, diag.Validation(err.Error()).WithSpan(span).WithFound(name).WithCause(err)
		}

		if name == "" {
			span, _ := ctx.Line.LastTokenSpan()
			return nil, diag.Validation(trackKeyword + " name cannot be empty").WithSpan(span)
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, err
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation)
			if err != nil {
				return nil, err
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectValue = effectValue

			switch effectKind {
			case t.KeywordPan:
				decl.EffectType = t.EffectPan
			case t.KeywordModulation:
				decl.EffectType = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			decl.EffectIntensityPercent = intensity

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, err
			}
		}

		decl.SourceName = name
	default:
		span, _ := ctx.Line.LastTokenSpan()
		return nil, diag.UnexpectedToken(span, first, t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance, t.KeywordMusic, t.KeywordTrack)
	}

	if token, ok := ctx.Line.Peek(); ok && token == t.KeywordLeft {
		ctx.Line.NextToken()

		left, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.AmplitudePercent = left

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordRight); err != nil {
			return nil, err
		}

		right, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		if left != right {
			decl.RightAmplitudePercent = right
			decl.HasRightAmplitude = true
		}
	} else {
		left, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return nil, err
		}
		decl.AmplitudePercent = left
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		span, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after track definition").WithSpan(span).WithFound(unknown)
	}

	return decl, nil
}
