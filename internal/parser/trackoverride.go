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
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
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

// ParseTrackOverride applies track overrides to the given preset
func (ctx *TextParser) ParseTrackOverride(preset *t.Preset) error {
	if preset == nil || preset.From == nil {
		return diag.Validation("cannot override tracks on a preset without a 'from' source")
	}

	_, ok := ctx.Line.NextToken()
	if !ok {
		return diag.UnexpectedEOF(ctx.Line.EOFSpan(), t.KeywordTrack)
	}

	trackIdx, err := ctx.Line.NextIntStrict()
	if err != nil {
		return err
	}
	trackSpan, _ := ctx.Line.LastTokenSpan()

	if trackIdx <= 0 || trackIdx >= t.NumberOfChannels {
		return diag.Validation(fmt.Sprintf("track index out of range (1-%d): %d", t.NumberOfChannels-1, trackIdx)).WithSpan(trackSpan).WithFound(fmt.Sprintf("%d", trackIdx))
	}

	idx := trackIdx - 1 // Convert to 0-based index
	from := preset.From

	if from.Track[idx].Type == t.TrackOff {
		return diag.Validation(fmt.Sprintf("cannot override track %d which is off in the template preset %q", trackIdx, from.String())).WithSpan(trackSpan).WithFound(fmt.Sprintf("%d", trackIdx))
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
		return err
	}
	kindSpan, _ := ctx.Line.LastTokenSpan()

	track := preset.Track[idx]

	switch kind {
	case t.KeywordTone:
		if track.Type == t.TrackAmbiance ||
			track.Type == t.TrackWhiteNoise ||
			track.Type == t.TrackPinkNoise ||
			track.Type == t.TrackBrownNoise {
			return diag.Validation(fmt.Sprintf("cannot set tone frequency on track %d of type %q", trackIdx, track.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}

		rawValue, _ := ctx.Line.Peek()

		carrier, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			carrier = from.Track[idx].Carrier + carrier
		}

		preset.Track[idx].Carrier = carrier
	case t.KeywordPan, t.KeywordModulation, t.KeywordDoppler:
		if kind == t.KeywordPan && track.Effect.Type != t.EffectPan {
			return diag.Validation(fmt.Sprintf("pan can only be set on track %d with pan effect, it is %q", trackIdx, track.Effect.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}
		if kind == t.KeywordModulation && track.Effect.Type != t.EffectModulation {
			return diag.Validation(fmt.Sprintf("modulation rate can only be set on track %d with modulation effect, it is %q", trackIdx, track.Effect.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}
		if kind == t.KeywordDoppler && track.Effect.Type != t.EffectDoppler {
			return diag.Validation(fmt.Sprintf("doppler speed can only be set on track %d with doppler effect, it is %q", trackIdx, track.Effect.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}

		rawValue, _ := ctx.Line.Peek()

		effectValue, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			effectValue = from.Track[idx].Effect.Value + effectValue
		}

		preset.Track[idx].Effect.Value = effectValue
	case t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic:
		if (kind == t.KeywordBinaural && track.Type != t.TrackBinauralBeat) ||
			(kind == t.KeywordMonaural && track.Type != t.TrackMonauralBeat) ||
			(kind == t.KeywordIsochronic && track.Type != t.TrackIsochronicBeat) {
			return diag.Validation(fmt.Sprintf("cannot change track %d type to %q, it is %q", trackIdx, kind, track.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}

		rawValue, _ := ctx.Line.Peek()

		resonance, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			resonance = from.Track[idx].Resonance + resonance
		}

		preset.Track[idx].Resonance = resonance
	case t.KeywordSmooth:
		if track.Type != t.TrackWhiteNoise &&
			track.Type != t.TrackPinkNoise &&
			track.Type != t.TrackBrownNoise {
			return diag.Validation(fmt.Sprintf("cannot set smooth on track %d of type %q", trackIdx, track.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}

		rawValue, _ := ctx.Line.Peek()

		smooth, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			smooth = from.Track[idx].NoiseSmooth + smooth
		}

		preset.Track[idx].NoiseSmooth = smooth
	case t.KeywordAmplitude:
		rawValue, _ := ctx.Line.Peek()

		amplitude, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			amplitude = from.Track[idx].Amplitude.ToPercent() + amplitude
		}

		preset.Track[idx].Amplitude = t.AmplitudePercentToRaw(amplitude)
	case t.KeywordIntensity:
		rawValue, _ := ctx.Line.Peek()

		intensity, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return err
		}

		if strings.HasPrefix(rawValue, "+") || strings.HasPrefix(rawValue, "-") {
			intensity = from.Track[idx].Effect.Intensity.ToPercent() + intensity
		}

		preset.Track[idx].Effect.Intensity = t.IntensityPercentToRaw(intensity)
	case t.KeywordWaveform:
		if track.Type == t.TrackBrownNoise ||
			track.Type == t.TrackPinkNoise ||
			track.Type == t.TrackWhiteNoise {
			return diag.Validation(fmt.Sprintf("cannot set waveform on track %d of type %q", trackIdx, track.Type.String())).WithSpan(kindSpan).WithFound(kind)
		}

		waveform, err := ctx.Line.NextExpectOneOf(
			t.KeywordSine,
			t.KeywordSquare,
			t.KeywordTriangle,
			t.KeywordSawtooth)

		if err != nil {
			return err
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
			return diag.Parse("unexpected waveform type").WithSpan(kindSpan).WithFound(waveform)
		}

		preset.Track[idx].Waveform = waveformType
	default:
		return diag.Parse("unexpected keyword").WithSpan(kindSpan).WithFound(kind)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return diag.Parse("unexpected token after track override definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	if err := preset.Track[idx].Validate(); err != nil {
		if span, ok := ctx.Line.LastTokenSpan(); ok {
			return diag.Validation(fmt.Sprintf("invalid track %d after override: %v", trackIdx, err)).WithSpan(span).WithCause(err)
		}
		return err
	}

	return nil
}
