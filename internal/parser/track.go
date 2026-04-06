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

type ParsedTrackDeclaration struct {
	Type                   t.TrackType
	Carrier                float64
	Resonance              float64
	AmplitudePercent       float64
	NoiseSmooth            float64
	Waveform               t.WaveformType
	AmbianceName           string
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

	if tok, ok := ctx.Line.Peek(); ok && tok == t.KeywordWaveform {
		ctx.Line.NextToken() // skip "waveform"

		wfTok, err := ctx.Line.NextExpectOneOf(t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth)
		if err != nil {
			return nil, err
		}

		switch wfTok {
		case t.KeywordSine:
			waveform = t.WaveformSine
		case t.KeywordSquare:
			waveform = t.WaveformSquare
		case t.KeywordTriangle:
			waveform = t.WaveformTriangle
		case t.KeywordSawtooth:
			waveform = t.WaveformSawtooth
		}

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordAmbiance); err != nil {
			return nil, err
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance)
	}

	decl := &ParsedTrackDeclaration{
		Waveform:   waveform,
		EffectType: t.EffectOff,
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
	case t.KeywordAmbiance:
		decl.Type = t.TrackAmbiance

		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "ambiance name")
		}

		if err := nr.IsValid(name); err != nil {
			span, _ := ctx.Line.LastTokenSpan()
			return nil, diag.Validation(err.Error()).WithSpan(span).WithFound(name).WithCause(err)
		}

		if name == "" {
			span, _ := ctx.Line.LastTokenSpan()
			return nil, diag.Validation("ambiance name cannot be empty").WithSpan(span)
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

		decl.AmbianceName = name
	default:
		span, _ := ctx.Line.LastTokenSpan()
		return nil, diag.UnexpectedToken(span, first, t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance, t.KeywordTrack)
	}

	var err error
	if decl.AmplitudePercent, err = ctx.Line.NextFloat64Strict(); err != nil {
		return nil, err
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		span, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after track definition").WithSpan(span).WithFound(unknown)
	}

	return decl, nil
}
