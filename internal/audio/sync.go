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

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const defaultSyncWindowMs = 1000

// sync synchronizes the audio renderer state with the current time
func (r *AudioRenderer) sync(timeMs int, periodIdx int) {
	if periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[periodIdx]
	alpha := r.transitionAlpha(r.interpolationProgress(timeMs, periodIdx), period.Transition)

	for ch := range r.channels {
		r.syncChannel(ch, periodIdx, period, alpha)
	}
}

func (r *AudioRenderer) interpolationProgress(timeMs int, periodIdx int) float64 {
	period := r.periods[periodIdx]
	return clampUnit(float64(timeMs-period.Time) / float64(r.nextPeriodTime(periodIdx, timeMs)-period.Time))
}

func (r *AudioRenderer) nextPeriodTime(periodIdx int, timeMs int) int {
	if periodIdx+1 < len(r.periods) {
		return r.periods[periodIdx+1].Time
	}

	return timeMs + defaultSyncWindowMs
}

func (r *AudioRenderer) transitionAlpha(progress float64, transition t.TransitionType) float64 {
	alpha := progress

	switch transition {
	case t.TransitionEaseOut:
		alpha = math.Log1p(math.Expm1(t.TransitionCurveK)*progress) / t.TransitionCurveK
	case t.TransitionEaseIn:
		alpha = math.Expm1(t.TransitionCurveK*progress) / math.Expm1(t.TransitionCurveK)
	case t.TransitionSmooth:
		raw := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*(progress-0.5)))
		min := 1.0 / (1.0 + math.Exp(t.TransitionCurveK*0.5))
		max := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*0.5))
		alpha = (raw - min) / (max - min)
	}

	return alpha
}

func (r *AudioRenderer) syncChannel(ch int, periodIdx int, period t.Period, alpha float64) {
	channel := &r.channels[ch]
	track := interpolateTrack(period.TrackStart[ch], period.TrackEnd[ch], alpha)
	previousTrackType := channel.Type
	previousEffectType := channel.Track.Effect.Type

	channel.Track = track
	channel.WaveformStart = period.TrackStart[ch].Waveform
	channel.WaveformEnd = period.TrackEnd[ch].Waveform
	channel.WaveformAlpha = alpha
	r.updateAmbianceIndex(ch, periodIdx, track.Type)
	r.resetRuntimeState(channel, previousTrackType, previousEffectType)
	r.configureEffectState(channel)
	r.configureTrackSignal(channel)
}

func interpolateTrack(start, end t.Track, alpha float64) t.Track {
	return t.Track{
		Type:         start.Type,
		Amplitude:    t.AmplitudeType(lerpFloat64(float64(start.Amplitude), float64(end.Amplitude), alpha)),
		Carrier:      lerpFloat64(start.Carrier, end.Carrier, alpha),
		Resonance:    lerpFloat64(start.Resonance, end.Resonance, alpha),
		NoiseSmooth:  lerpFloat64(start.NoiseSmooth, end.NoiseSmooth, alpha),
		Waveform:     start.Waveform,
		AmbianceName: start.AmbianceName,
		Effect: t.Effect{
			Type:      start.Effect.Type,
			Value:     lerpFloat64(start.Effect.Value, end.Effect.Value, alpha),
			Intensity: t.IntensityType(lerpFloat64(float64(start.Effect.Intensity), float64(end.Effect.Intensity), alpha)),
		},
	}
}

func (r *AudioRenderer) updateAmbianceIndex(ch int, periodIdx int, trackType t.TrackType) {
	if trackType == t.TrackAmbiance {
		r.channelAmbianceIndex[ch] = r.periodAmbianceStart[periodIdx][ch]
		return
	}

	r.channelAmbianceIndex[ch] = -1
}

func (r *AudioRenderer) resetRuntimeState(channel *t.Channel, previousTrackType t.TrackType, previousEffectType t.EffectType) {
	if previousTrackType != channel.Track.Type {
		channel.Offset = [2]int{}
	}
	channel.Type = channel.Track.Type

	if previousEffectType != channel.Track.Effect.Type {
		channel.Effect.Offset = 0
		channel.Effect.ModulationGain = 0
		channel.Effect.ModulationInitialized = false
		channel.Effect.PanPosition = 0
		channel.Effect.PanInitialized = false
	}
}

func (r *AudioRenderer) configureEffectState(channel *t.Channel) {
	channel.Effect.Increment = 0
	if channel.Track.Effect.Type == t.EffectOff {
		return
	}

	channel.Effect.Increment = r.frequencyToIncrement(channel.Track.Effect.Value)
}

func (r *AudioRenderer) configureTrackSignal(channel *t.Channel) {
	channel.Amplitude = [2]int{}
	channel.Increment = [2]int{}

	amplitude := int(channel.Track.Amplitude)

	switch channel.Track.Type {
	case t.TrackPureTone:
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = r.frequencyToIncrement(channel.Track.Carrier)
	case t.TrackBinauralBeat:
		freq1 := channel.Track.Carrier + channel.Track.Resonance/2
		freq2 := channel.Track.Carrier - channel.Track.Resonance/2
		channel.Amplitude[0] = amplitude
		channel.Amplitude[1] = amplitude
		channel.Increment[0] = r.frequencyToIncrement(freq1)
		channel.Increment[1] = r.frequencyToIncrement(freq2)
	case t.TrackMonauralBeat:
		freqHigh := channel.Track.Carrier + channel.Track.Resonance/2
		freqLow := channel.Track.Carrier - channel.Track.Resonance/2
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = r.frequencyToIncrement(freqHigh)
		channel.Increment[1] = r.frequencyToIncrement(freqLow)
	case t.TrackIsochronicBeat:
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = r.frequencyToIncrement(channel.Track.Carrier)
		channel.Increment[1] = r.frequencyToIncrement(channel.Track.Resonance)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance:
		channel.Amplitude[0] = amplitude
	}
}

func (r *AudioRenderer) frequencyToIncrement(frequency float64) int {
	return int(frequency / float64(r.SampleRate) * t.SineTableSize * t.PhasePrecision)
}

func lerpFloat64(start, end, alpha float64) float64 {
	return start + (end-start)*alpha
}

func clampUnit(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}

	return value
}
