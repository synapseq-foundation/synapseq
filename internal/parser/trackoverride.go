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

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
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

// ParseTrackOverride applies track overrides to the given preset
func (ctx *TextParser) ParseTrackOverride(preset *t.Preset) error {
	if preset == nil || preset.From == nil {
		return fmt.Errorf("cannot override tracks on a preset without a 'from' source")
	}

	ln := ctx.Line.Raw
	_, ok := ctx.Line.NextToken()
	if !ok {
		return fmt.Errorf("expected 'track' keyword, got EOF: %s", ln)
	}

	trackIdx, err := ctx.Line.NextIntStrict()
	if err != nil {
		return fmt.Errorf("expected track index after 'track': %s", ln)
	}

	if trackIdx <= 0 || trackIdx >= t.NumberOfChannels {
		return fmt.Errorf("track index out of range (1-%d): %d", t.NumberOfChannels-1, trackIdx)
	}

	idx := trackIdx - 1 // Convert to 0-based index
	from := preset.From

	if from.Track[idx].Type == t.TrackOff {
		return fmt.Errorf("cannot override track %d which is off in the template preset %q", trackIdx, from.String())
	}

	kind, err := ctx.Line.NextExpectOneOf(
		t.KeywordTone,
		t.KeywordBinaural,
		t.KeywordMonaural,
		t.KeywordIsochronic,
		t.KeywordPan,
		t.KeywordModulation,
		t.KeywordDoppler,
		t.KeywordAmplitude,
		t.KeywordIntensity)
	if err != nil {
		return fmt.Errorf(
			"expected one of %q, %q, %q, %q, %q, %q, %q, %q: %s",
			t.KeywordTone,
			t.KeywordBinaural,
			t.KeywordMonaural,
			t.KeywordIsochronic,
			t.KeywordPan,
			t.KeywordModulation,
			t.KeywordDoppler,
			t.KeywordAmplitude,
			ln)
	}

	track := preset.Track[idx]

	switch kind {
	case t.KeywordTone:
		if track.Type == t.TrackBackground ||
			track.Type == t.TrackWhiteNoise ||
			track.Type == t.TrackPinkNoise ||
			track.Type == t.TrackBrownNoise {
			return fmt.Errorf("cannot set tone frequency on track %d of type %q", trackIdx, track.Type.String())
		}

		carrier, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("tone frequency: %w", err)
		}

		preset.Track[idx].Carrier = carrier
	case t.KeywordPan, t.KeywordModulation, t.KeywordDoppler:
		if kind == t.KeywordPan && track.Effect.Type != t.EffectPan {
			return fmt.Errorf("pan can only be set on track %d with pan effect, it is %q", trackIdx, track.Effect.Type.String())
		}
		if kind == t.KeywordModulation && track.Effect.Type != t.EffectModulation {
			return fmt.Errorf("modulation rate can only be set on track %d with modulation effect, it is %q", trackIdx, track.Effect.Type.String())
		}
		if kind == t.KeywordDoppler && track.Effect.Type != t.EffectDoppler {
			return fmt.Errorf("doppler speed can only be set on track %d with doppler effect, it is %q", trackIdx, track.Effect.Type.String())
		}

		effectValue, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("effect value: %w", err)
		}

		preset.Track[idx].Effect.Value = effectValue
	case t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic:
		if (kind == t.KeywordBinaural && track.Type != t.TrackBinauralBeat) ||
			(kind == t.KeywordMonaural && track.Type != t.TrackMonauralBeat) ||
			(kind == t.KeywordIsochronic && track.Type != t.TrackIsochronicBeat) {
			return fmt.Errorf("cannot change track %d type to %q, it is %q", trackIdx, kind, track.Type.String())
		}

		resonance, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("resonance: %w", err)
		}

		preset.Track[idx].Resonance = resonance
	case t.KeywordAmplitude:
		amplitude, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("amplitude: %w", err)
		}

		preset.Track[idx].Amplitude = t.AmplitudePercentToRaw(amplitude)
	case t.KeywordIntensity:
		intensity, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("intensity: %w", err)
		}

		preset.Track[idx].Effect.Intensity = t.IntensityPercentToRaw(intensity)
	default:
		return fmt.Errorf("unexpected keyword: %s", kind)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return fmt.Errorf("unexpected token after track override definition: %q", unknown)
	}

	if err := preset.Track[idx].Validate(); err != nil {
		return fmt.Errorf("invalid track %d after override: %w", trackIdx, err)
	}

	return nil
}
