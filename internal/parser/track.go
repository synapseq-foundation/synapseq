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
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

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

// ParseTrack extracts and returns a Track from the current line context
func (ctx *TextParser) ParseTrack() (*t.Track, error) {
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

	var (
		carrier, resonance, amplitude, smooth float64
		trackType                             t.TrackType
		ambianceName                          string
	)

	effect := t.Effect{
		Type:      t.EffectOff,
		Value:     0.0,
		Intensity: 0.0,
	}

	switch first {
	case t.KeywordTone:
		var err error
		if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, err
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, err
		}

		switch kind {
		case t.KeywordBinaural:
			trackType = t.TrackBinauralBeat
		case t.KeywordMonaural:
			trackType = t.TrackMonauralBeat
		case t.KeywordIsochronic:
			trackType = t.TrackIsochronicBeat
		default:
			trackType = t.TrackPureTone
		}

		if trackType == t.TrackBinauralBeat ||
			trackType == t.TrackMonauralBeat ||
			trackType == t.TrackIsochronicBeat {
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
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
			effect.Value = effectValue

			switch effectKind {
			case t.KeywordPan:
				effect.Type = t.EffectPan
			case t.KeywordModulation:
				effect.Type = t.EffectModulation
			case t.KeywordDoppler:
				effect.Type = t.EffectDoppler
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

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
			trackType = t.TrackWhiteNoise
		case t.KeywordPink:
			trackType = t.TrackPinkNoise
		case t.KeywordBrown:
			trackType = t.TrackBrownNoise
		}

		kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordSmooth, t.KeywordAmplitude)
		if err != nil {
			return nil, err
		}

		if kind == t.KeywordSmooth {
			smooth, err = ctx.Line.NextFloat64Strict()
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
			effect.Value = effectValue

			switch effectKind {
			case t.KeywordPan:
				effect.Type = t.EffectPan
			case t.KeywordModulation:
				effect.Type = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, err
			}
		}
	case t.KeywordAmbiance:
		trackType = t.TrackAmbiance

		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "ambiance name")
		}

		if err := s.IsValidNamedRef(name); err != nil {
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
			effect.Value = effectValue

			switch effectKind {
			case t.KeywordPan:
				effect.Type = t.EffectPan
			case t.KeywordModulation:
				effect.Type = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, err
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, err
			}
		}

		// convert to 0-based index
		ambianceName = name
	default:
		span, _ := ctx.Line.LastTokenSpan()
		return nil, diag.UnexpectedToken(span, first, t.KeywordTone, t.KeywordNoise, t.KeywordAmbiance, t.KeywordTrack)
	}

	var err error
	if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
		return nil, err
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		span, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after track definition").WithSpan(span).WithFound(unknown)
	}

	track := t.Track{
		Type:         trackType,
		Carrier:      carrier,
		Resonance:    resonance,
		Amplitude:    t.AmplitudePercentToRaw(amplitude),
		AmbianceName: ambianceName,
		NoiseSmooth:  smooth,
		Waveform:     waveform,
		Effect:       effect,
	}
	if err := track.Validate(); err != nil {
		if span, ok := ctx.Line.LastTokenSpan(); ok {
			return nil, diag.Validation(err.Error()).WithSpan(span).WithCause(err)
		}
		return nil, err
	}

	return &track, nil
}
