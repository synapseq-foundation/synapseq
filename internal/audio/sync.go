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

package audio

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// sync synchronizes the audio renderer state with the current time
func (r *AudioRenderer) sync(timeMs int, periodIdx int) {
	if periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[periodIdx]
	nextTime := timeMs + 1000 // Default next time
	if periodIdx+1 < len(r.periods) {
		nextTime = r.periods[periodIdx+1].Time
	}

	// Calculate interpolation factor (0.0 to 1.0)
	progress := float64(timeMs-period.Time) / float64(nextTime-period.Time)
	// Clamp progress between 0 and 1
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	// Update each channel
	for ch := range t.NumberOfChannels {
		if ch >= len(r.channels) || ch >= len(period.TrackStart) {
			return // Bounds protection
		}

		channel := &r.channels[ch]
		tr0 := period.TrackStart[ch]
		tr1 := period.TrackEnd[ch]

		// Interpolate channel parameters
		alpha := progress
		switch period.Transition {
		case t.TransitionEaseOut:
			alpha = math.Log1p(math.Expm1(t.TransitionCurveK)*progress) / t.TransitionCurveK
		case t.TransitionEaseIn:
			alpha = math.Expm1(t.TransitionCurveK*progress) / math.Expm1(t.TransitionCurveK)
		case t.TransitionSmooth:
			// Normalized sigmoid
			raw := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*(progress-0.5)))
			min := 1.0 / (1.0 + math.Exp(t.TransitionCurveK*0.5))
			max := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*0.5))
			alpha = (raw - min) / (max - min)
		}

		prevEffectType := channel.Track.Effect.Type

		channel.Track.Type = tr0.Type
		channel.Track.Effect.Type = tr0.Effect.Type
		channel.Track.Amplitude = t.AmplitudeType(float64(tr0.Amplitude)*(1-alpha) + float64(tr1.Amplitude)*alpha)
		channel.Track.Carrier = tr0.Carrier*(1-alpha) + tr1.Carrier*alpha
		channel.Track.Resonance = tr0.Resonance*(1-alpha) + tr1.Resonance*alpha
		channel.Track.Waveform = tr0.Waveform
		channel.Track.Intensity = t.IntensityType(float64(tr0.Intensity)*(1-alpha) + float64(tr1.Intensity)*alpha)

		// Effects
		if channel.Track.Effect.Type == t.EffectSpin {
			cfg0 := tr0.Effect.Configuration.(t.EffectSpinConfiguration)
			if tr1.Effect.Type == t.EffectSpin {
				cfg1 := tr1.Effect.Configuration.(t.EffectSpinConfiguration)
				channel.Track.Effect.Configuration = t.EffectSpinConfiguration{
					Rate: cfg0.Rate*(1-alpha) + cfg1.Rate*alpha,
				}
			} else {
				channel.Track.Effect.Configuration = t.EffectSpinConfiguration{
					Rate: cfg0.Rate,
				}
			}
		}
		if channel.Track.Effect.Type == t.EffectPulse {
			cfg0 := tr0.Effect.Configuration.(t.EffectPulseConfiguration)
			if tr1.Effect.Type == t.EffectPulse {
				cfg1 := tr1.Effect.Configuration.(t.EffectPulseConfiguration)
				channel.Track.Effect.Configuration = t.EffectPulseConfiguration{
					Pulse: cfg0.Pulse*(1-alpha) + cfg1.Pulse*alpha,
				}
			} else {
				channel.Track.Effect.Configuration = t.EffectPulseConfiguration{
					Pulse: cfg0.Pulse,
				}
			}
		}

		// Reset offsets if track type has changed
		if channel.Type != channel.Track.Type {
			channel.Type = channel.Track.Type
			channel.Offset[0] = 0
			channel.Offset[1] = 0
		}

		// Reset effect phase if effect type changed
		if prevEffectType != channel.Track.Effect.Type {
			channel.Effect.Offset = 0
		}

		switch channel.Track.Effect.Type {
		case t.EffectSpin:
			cfg := channel.Track.Effect.Configuration.(t.EffectSpinConfiguration)
			channel.Effect.Increment = int(cfg.Rate / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.EffectPulse:
			cfg := channel.Track.Effect.Configuration.(t.EffectPulseConfiguration)
			channel.Effect.Increment = int(cfg.Pulse / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		}

		switch channel.Track.Type {
		case t.TrackPureTone:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(channel.Track.Carrier / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackBinauralBeat:
			freq1 := channel.Track.Carrier + channel.Track.Resonance/2
			freq2 := channel.Track.Carrier - channel.Track.Resonance/2
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Amplitude[1] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(freq1 / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freq2 / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackMonauralBeat:
			freqHigh := channel.Track.Carrier + channel.Track.Resonance/2
			freqLow := channel.Track.Carrier - channel.Track.Resonance/2
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(freqHigh / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freqLow / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackIsochronicBeat:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(channel.Track.Carrier / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(channel.Track.Resonance / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
		case t.TrackBackground:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
		}
	}
}
