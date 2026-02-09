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

	switch kind {
	case t.KeywordTone, t.KeywordPan:
		track := preset.Track[idx]

		if kind == t.KeywordTone && track.Type == t.TrackBackground {
			return fmt.Errorf("background track %d cannot have a tone carrier", trackIdx)
		}
		if kind == t.KeywordPan && track.Type != t.TrackBackground {
			return fmt.Errorf("track %d must be a background track to set pan width, it is %q", trackIdx, track.Type.String())
		}
		if kind == t.KeywordPan && track.Effect.Type != t.EffectPan {
			return fmt.Errorf("pan width can only be set on track %d with pan effect, it is %q", trackIdx, track.Effect.Type.String())
		}

		carrier, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("carrier: %w", err)
		}

		preset.Track[idx].Carrier = carrier
	case t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, t.KeywordModulation:
		track := preset.Track[idx]

		// Validate that the track type matches the keyword being set
		if (kind == t.KeywordBinaural && track.Type != t.TrackBinauralBeat) ||
			(kind == t.KeywordMonaural && track.Type != t.TrackMonauralBeat) ||
			(kind == t.KeywordIsochronic && track.Type != t.TrackIsochronicBeat) ||
			(kind == t.KeywordModulation && track.Type != t.TrackBackground) {
			return fmt.Errorf("cannot change track %d type to %q, it is %q", trackIdx, kind, track.Type.String())
		}

		// Validate that the effect type matches the keyword being set
		if kind == t.KeywordModulation && track.Effect.Type != t.EffectModulation {
			return fmt.Errorf("cannot change track %d effect to %q, it is %q", trackIdx, kind, track.Effect.Type.String())
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

	// Validate the updated track
	if err := preset.Track[idx].Validate(); err != nil {
		return fmt.Errorf("invalid track %d after override: %w", trackIdx, err)
	}

	return nil
}
