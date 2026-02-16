/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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

	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
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
	ln := ctx.Line.Raw

	if tok, ok := ctx.Line.Peek(); ok && tok == t.KeywordWaveform {
		ctx.Line.NextToken() // skip "waveform"

		wfTok, err := ctx.Line.NextExpectOneOf(t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, %q, or %q after waveform: %s", t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth, ln)
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

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordBackground); err != nil {
			return nil, fmt.Errorf("expected %q or %q after waveform type: %s", t.KeywordTone, t.KeywordBackground, ln)
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected %q, %s or %q: %s", t.KeywordTone, t.KeywordNoise, t.KeywordBackground, ln)
	}

	var (
		carrier, resonance, amplitude float64
		trackType                     t.TrackType
		backgroundName                string
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
			return nil, fmt.Errorf("carrier: %w", err)
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q after carrier: %s", t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordEffect, t.KeywordAmplitude, ln)
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
				return nil, fmt.Errorf("resonance: %w", err)
			}

			kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
			if err != nil {
				return nil, fmt.Errorf("expected %q or %q after resonance: %s", t.KeywordEffect, t.KeywordAmplitude, ln)
			}
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation, t.KeywordDoppler)
			if err != nil {
				return nil, fmt.Errorf("expected %q, %q or %q after effect: %s", t.KeywordPan, t.KeywordModulation, t.KeywordDoppler, ln)
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("effect value: %w", err)
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
				return nil, fmt.Errorf("expected %q after effect value: %s", t.KeywordIntensity, ln)
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, fmt.Errorf("expected %q after intensity: %s", t.KeywordAmplitude, ln)
			}
		}
	case t.KeywordNoise:
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordWhite, t.KeywordPink, t.KeywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after noise: %s", t.KeywordWhite, t.KeywordPink, t.KeywordBrown, ln)
		}

		switch kind {
		case t.KeywordWhite:
			trackType = t.TrackWhiteNoise
		case t.KeywordPink:
			trackType = t.TrackPinkNoise
		case t.KeywordBrown:
			trackType = t.TrackBrownNoise
		}

		kind, err = ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, fmt.Errorf("expected %q or %q after noise type: %s", t.KeywordEffect, t.KeywordAmplitude, ln)
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation)
			if err != nil {
				return nil, fmt.Errorf("expected %q or %q after noise effect: %s", t.KeywordPan, t.KeywordModulation, ln)
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("effect value: %w", err)
			}
			effect.Value = effectValue

			switch effectKind {
			case t.KeywordPan:
				effect.Type = t.EffectPan
			case t.KeywordModulation:
				effect.Type = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after effect value: %s", t.KeywordIntensity, ln)
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, fmt.Errorf("expected %q after intensity: %s", t.KeywordAmplitude, ln)
			}
		}
	case t.KeywordBackground:
		trackType = t.TrackBackground

		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("background name cannot be empty: %s", ln)
		}

		if err := s.IsValidNamedRef(name); err != nil {
			return nil, err
		}

		if name == "" {
			return nil, fmt.Errorf("background name cannot be empty: %s", ln)
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordEffect, t.KeywordAmplitude)
		if err != nil {
			return nil, fmt.Errorf("expected %q or %q after background: %s", t.KeywordEffect, t.KeywordAmplitude, ln)
		}

		if kind == t.KeywordEffect {
			effectKind, err := ctx.Line.NextExpectOneOf(t.KeywordPan, t.KeywordModulation)
			if err != nil {
				return nil, fmt.Errorf("expected %q or %q after noise effect: %s", t.KeywordPan, t.KeywordModulation, ln)
			}

			effectValue, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("effect value: %w", err)
			}
			effect.Value = effectValue

			switch effectKind {
			case t.KeywordPan:
				effect.Type = t.EffectPan
			case t.KeywordModulation:
				effect.Type = t.EffectModulation
			}

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after effect value: %s", t.KeywordIntensity, ln)
			}

			intensity, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
			effect.Intensity = t.IntensityPercentToRaw(intensity)

			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, fmt.Errorf("expected %q after intensity: %s", t.KeywordAmplitude, ln)
			}
		}

		// convert to 0-based index
		backgroundName = name
	default:
		return nil, fmt.Errorf("expected %q, %q, %q or %q. Received: %s", t.KeywordTone, t.KeywordNoise, t.KeywordBackground, t.KeywordTrack, first)
	}

	var err error
	if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
		return nil, fmt.Errorf("amplitude: %w", err)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after track definition: %q", unknown)
	}

	track := t.Track{
		Type:           trackType,
		Carrier:        carrier,
		Resonance:      resonance,
		Amplitude:      t.AmplitudePercentToRaw(amplitude),
		BackgroundName: backgroundName,
		Waveform:       waveform,
		Effect:         effect,
	}
	if err := track.Validate(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &track, nil
}
