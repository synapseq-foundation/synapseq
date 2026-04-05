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

package sync

import (
	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const defaultSyncWindowMs = 1000

type Engine struct {
	SampleRate          int
	UpdateAmbianceIndex func(ch int, periodIdx int, trackType t.TrackType)
}

func NewEngine(sampleRate int, updateAmbianceIndex func(ch int, periodIdx int, trackType t.TrackType)) *Engine {
	return &Engine{
		SampleRate:          sampleRate,
		UpdateAmbianceIndex: updateAmbianceIndex,
	}
}

func (e *Engine) Sync(channels []t.Channel, periods []t.Period, timeMs int, periodIdx int) {
	if periodIdx >= len(periods) {
		return
	}

	period := periods[periodIdx]
	alpha := tl.StepAlpha(e.interpolationProgress(periods, timeMs, periodIdx), period.Transition, period.Steps)

	for ch := range channels {
		e.syncChannel(ch, channels, periodIdx, period, alpha)
	}
}

func FrequencyToIncrement(sampleRate int, frequency float64) int {
	return int(frequency / float64(sampleRate) * t.SineTableSize * t.PhasePrecision)
}

func (e *Engine) interpolationProgress(periods []t.Period, timeMs int, periodIdx int) float64 {
	period := periods[periodIdx]
	return clampUnit(float64(timeMs-period.Time) / float64(e.nextPeriodTime(periods, periodIdx, timeMs)-period.Time))
}

func (e *Engine) nextPeriodTime(periods []t.Period, periodIdx int, timeMs int) int {
	if periodIdx+1 < len(periods) {
		return periods[periodIdx+1].Time
	}

	return timeMs + defaultSyncWindowMs
}

func (e *Engine) syncChannel(ch int, channels []t.Channel, periodIdx int, period t.Period, alpha float64) {
	channel := &channels[ch]
	track := interpolateTrack(period.TrackStart[ch], period.TrackEnd[ch], alpha)
	previousTrackType := channel.Type
	previousEffectType := channel.Track.Effect.Type

	channel.Track = track
	channel.WaveformStart = period.TrackStart[ch].Waveform
	channel.WaveformEnd = period.TrackEnd[ch].Waveform
	channel.WaveformAlpha = alpha
	e.updateAmbianceIndex(ch, periodIdx, track.Type)
	e.resetRuntimeState(channel, previousTrackType, previousEffectType)
	e.configureEffectState(channel)
	e.configureTrackSignal(channel)
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

func (e *Engine) updateAmbianceIndex(ch int, periodIdx int, trackType t.TrackType) {
	if e.UpdateAmbianceIndex == nil {
		return
	}

	e.UpdateAmbianceIndex(ch, periodIdx, trackType)
}

func (e *Engine) resetRuntimeState(channel *t.Channel, previousTrackType t.TrackType, previousEffectType t.EffectType) {
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

func (e *Engine) configureEffectState(channel *t.Channel) {
	channel.Effect.Increment = 0
	if channel.Track.Effect.Type == t.EffectOff {
		return
	}

	channel.Effect.Increment = FrequencyToIncrement(e.SampleRate, channel.Track.Effect.Value)
}

func (e *Engine) configureTrackSignal(channel *t.Channel) {
	channel.Amplitude = [2]int{}
	channel.Increment = [2]int{}

	amplitude := int(channel.Track.Amplitude)

	switch channel.Track.Type {
	case t.TrackPureTone:
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = FrequencyToIncrement(e.SampleRate, channel.Track.Carrier)
	case t.TrackBinauralBeat:
		freq1 := channel.Track.Carrier + channel.Track.Resonance/2
		freq2 := channel.Track.Carrier - channel.Track.Resonance/2
		channel.Amplitude[0] = amplitude
		channel.Amplitude[1] = amplitude
		channel.Increment[0] = FrequencyToIncrement(e.SampleRate, freq1)
		channel.Increment[1] = FrequencyToIncrement(e.SampleRate, freq2)
	case t.TrackMonauralBeat:
		freqHigh := channel.Track.Carrier + channel.Track.Resonance/2
		freqLow := channel.Track.Carrier - channel.Track.Resonance/2
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = FrequencyToIncrement(e.SampleRate, freqHigh)
		channel.Increment[1] = FrequencyToIncrement(e.SampleRate, freqLow)
	case t.TrackIsochronicBeat:
		channel.Amplitude[0] = amplitude
		channel.Increment[0] = FrequencyToIncrement(e.SampleRate, channel.Track.Carrier)
		channel.Increment[1] = FrequencyToIncrement(e.SampleRate, channel.Track.Resonance)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance:
		channel.Amplitude[0] = amplitude
	}
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

func lerpFloat64(start, end, alpha float64) float64 {
	return start + (end-start)*alpha
}
