// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preset

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type TrackOverrideSpec struct {
	TrackIndex        int
	TrackSpan         diag.Span
	Kind              string
	KindSpan          diag.Span
	Value             float64
	RawValue          string
	ValueSpan         diag.Span
	Relative          bool
	Amplitude         [2]float64
	RawAmplitude      [2]string
	AmplitudeSpan     [2]diag.Span
	RelativeAmplitude [2]bool
	Waveform          t.WaveformName
}

func ApplyTrackOverride(preset *t.Preset, spec *TrackOverrideSpec) error {
	if preset == nil || preset.From == nil {
		return diag.Validation("cannot override tracks on a preset without a 'from' source")
	}

	idx := spec.TrackIndex - 1
	from := preset.From

	if from.Track[idx].Type == t.TrackOff {
		return diag.Validation(fmt.Sprintf("cannot override track %d which is off in the template preset %q", spec.TrackIndex, from.String())).WithSpan(spec.TrackSpan).WithFound(fmt.Sprintf("%d", spec.TrackIndex))
	}

	track := preset.Track[idx]

	switch spec.Kind {
	case t.KeywordTone:
		if track.Type == t.TrackAmbiance ||
			track.Type == t.TrackMusic ||
			track.Type == t.TrackWhiteNoise ||
			track.Type == t.TrackPinkNoise ||
			track.Type == t.TrackBrownNoise {
			return diag.Validation(fmt.Sprintf("cannot set tone frequency on track %d of type %q", spec.TrackIndex, track.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}

		carrier := spec.Value
		if spec.Relative {
			carrier = from.Track[idx].Carrier + carrier
		}
		preset.Track[idx].Carrier = carrier

	case t.KeywordPan, t.KeywordModulation, t.KeywordDoppler, t.KeywordShift:
		if spec.Kind == t.KeywordPan && track.Effect.Type != t.EffectPan {
			return diag.Validation(fmt.Sprintf("pan can only be set on track %d with pan effect, it is %q", spec.TrackIndex, track.Effect.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}
		if spec.Kind == t.KeywordModulation && track.Effect.Type != t.EffectModulation {
			return diag.Validation(fmt.Sprintf("modulation rate can only be set on track %d with modulation effect, it is %q", spec.TrackIndex, track.Effect.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}
		if spec.Kind == t.KeywordDoppler && track.Effect.Type != t.EffectDoppler {
			return diag.Validation(fmt.Sprintf("doppler speed can only be set on track %d with doppler effect, it is %q", spec.TrackIndex, track.Effect.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}
		if spec.Kind == t.KeywordShift && track.Effect.Type != t.EffectShift {
			return diag.Validation(fmt.Sprintf("shift separation can only be set on track %d with shift effect, it is %q", spec.TrackIndex, track.Effect.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}

		effectValue := spec.Value
		if spec.Relative {
			effectValue = from.Track[idx].Effect.Value + effectValue
		}
		preset.Track[idx].Effect.Value = effectValue

	case t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic:
		if (spec.Kind == t.KeywordBinaural && track.Type != t.TrackBinauralBeat) ||
			(spec.Kind == t.KeywordMonaural && track.Type != t.TrackMonauralBeat) ||
			(spec.Kind == t.KeywordIsochronic && track.Type != t.TrackIsochronicBeat) {
			return diag.Validation(fmt.Sprintf("cannot change track %d type to %q, it is %q", spec.TrackIndex, spec.Kind, track.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}

		resonance := spec.Value
		if spec.Relative {
			resonance = from.Track[idx].Resonance + resonance
		}
		preset.Track[idx].Resonance = resonance

	case t.KeywordSmooth:
		if track.Type != t.TrackWhiteNoise &&
			track.Type != t.TrackPinkNoise &&
			track.Type != t.TrackBrownNoise {
			return diag.Validation(fmt.Sprintf("cannot set smooth on track %d of type %q", spec.TrackIndex, track.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}

		smooth := spec.Value
		if spec.Relative {
			smooth = from.Track[idx].NoiseSmooth + smooth
		}
		preset.Track[idx].NoiseSmooth = smooth

	case t.KeywordAmplitude:
		for channel, amplitude := range spec.Amplitude {
			if spec.RelativeAmplitude[channel] {
				amplitude = from.Track[idx].Amplitude[channel].ToPercent() + amplitude
			}
			preset.Track[idx].Amplitude[channel] = t.AmplitudePercentToRawChannel(amplitude)
		}
	case t.KeywordLeft, t.KeywordRight:
		channel := 0
		if spec.Kind == t.KeywordRight {
			channel = 1
		}

		amplitude := spec.Value
		if spec.Relative {
			amplitude = from.Track[idx].Amplitude[channel].ToPercent() + amplitude
		}
		preset.Track[idx].Amplitude[channel] = t.AmplitudePercentToRawChannel(amplitude)

	case t.KeywordIntensity:
		intensity := spec.Value
		if spec.Relative {
			intensity = from.Track[idx].Effect.Intensity.ToPercent() + intensity
		}
		preset.Track[idx].Effect.Intensity = t.IntensityPercentToRaw(intensity)

	case t.KeywordWaveform:
		if track.Type == t.TrackBrownNoise ||
			track.Type == t.TrackPinkNoise ||
			track.Type == t.TrackWhiteNoise {
			return diag.Validation(fmt.Sprintf("cannot set waveform on track %d of type %q", spec.TrackIndex, track.Type.String())).WithSpan(spec.KindSpan).WithFound(spec.Kind)
		}

		preset.Track[idx].Waveform = spec.Waveform

	default:
		return diag.Parse("unexpected keyword").WithSpan(spec.KindSpan).WithFound(spec.Kind)
	}

	span := spec.ValueSpan
	if spec.Kind == t.KeywordAmplitude {
		span = spec.AmplitudeSpan[0]
	}
	if !span.HasLocation() {
		span = spec.KindSpan
	}
	if err := preset.Track[idx].Validate(); err != nil {
		return diag.Validation(fmt.Sprintf("invalid track %d after override: %v", spec.TrackIndex, err)).WithSpan(span).WithCause(err)
	}

	return nil
}
