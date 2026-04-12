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

package audio

import (
	"math"

	audiosync "github.com/synapseq-foundation/synapseq/v4/internal/audio/sync"
	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const defaultPlanWindowMs = 1000

type renderPlan struct {
	periods     []t.Period
	windows     []renderWindow
	sampleRate  int
	totalFrames int64
}

type renderWindow struct {
	PeriodIndex int
	StartMs     int
	EndMs       int
}

func compileRenderPlan(periods []t.Period, sampleRate int) renderPlan {
	plan := renderPlan{
		periods:     periods,
		windows:     make([]renderWindow, len(periods)),
		sampleRate:  sampleRate,
		totalFrames: totalFramesFromDuration(durationMs(periods), sampleRate),
	}

	for index := range periods {
		endMs := periods[index].Time
		if index+1 < len(periods) {
			endMs = periods[index+1].Time
		}

		plan.windows[index] = renderWindow{
			PeriodIndex: index,
			StartMs:     periods[index].Time,
			EndMs:       endMs,
		}
	}

	return plan
}

func (rp renderPlan) periodIndexAt(currentTimeMs int, currentPeriodIdx int) int {
	for currentPeriodIdx+1 < len(rp.windows) && currentTimeMs >= rp.windows[currentPeriodIdx+1].StartMs {
		currentPeriodIdx++
	}

	return currentPeriodIdx
}

func (rp renderPlan) cue(periodIdx int, currentTimeMs int) audiosync.Cue {
	window := rp.windows[periodIdx]
	period := rp.periods[periodIdx]
	alpha := tl.StepAlpha(rp.interpolationProgress(window, currentTimeMs), period.Transition, period.Steps)
	cue := audiosync.Cue{
		PeriodIndex: window.PeriodIndex,
		Channels:    [t.NumberOfChannels]audiosync.ChannelCue{},
	}

	for index := 0; index < t.NumberOfChannels; index++ {
		signal := compileSignalState(planTrackState{
			track:      interpolateTrack(period.TrackStart[index], period.TrackEnd[index], alpha),
			sampleRate: rp.sampleRate,
		})
		cue.Channels[index] = audiosync.ChannelCue{
			Track:         signal.Track,
			WaveformStart: period.TrackStart[index].Waveform,
			WaveformEnd:   period.TrackEnd[index].Waveform,
			WaveformAlpha: alpha,
			Amplitude:     signal.Amplitude,
			Increment:     signal.Increment,
			EffectStep:    signal.EffectStep,
		}
	}

	return cue
}

func (rp renderPlan) interpolationProgress(window renderWindow, currentTimeMs int) float64 {
	endMs := window.EndMs
	if endMs <= window.StartMs {
		endMs = currentTimeMs + defaultPlanWindowMs
	}

	return clampUnit(float64(currentTimeMs-window.StartMs) / float64(endMs-window.StartMs))
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

type planTrackState struct {
	track      t.Track
	sampleRate int
}

type compiledSignalState struct {
	Track      t.Track
	Amplitude  [2]int
	Increment  [2]int
	EffectStep int
}

func compileSignalState(state planTrackState) compiledSignalState {
	compiled := compiledSignalState{Track: state.track}
	if state.track.Effect.Type != t.EffectOff {
		compiled.EffectStep = frequencyToIncrement(state.sampleRate, state.track.Effect.Value)
	}

	amplitude := int(state.track.Amplitude)
	switch state.track.Type {
	case t.TrackPureTone:
		compiled.Amplitude[0] = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, state.track.Carrier)
	case t.TrackBinauralBeat:
		freq1 := state.track.Carrier + state.track.Resonance/2
		freq2 := state.track.Carrier - state.track.Resonance/2
		compiled.Amplitude[0] = amplitude
		compiled.Amplitude[1] = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, freq1)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, freq2)
	case t.TrackMonauralBeat:
		freqHigh := state.track.Carrier + state.track.Resonance/2
		freqLow := state.track.Carrier - state.track.Resonance/2
		compiled.Amplitude[0] = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, freqHigh)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, freqLow)
	case t.TrackIsochronicBeat:
		compiled.Amplitude[0] = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, state.track.Carrier)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, state.track.Resonance)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance:
		compiled.Amplitude[0] = amplitude
	}

	return compiled
}

func frequencyToIncrement(sampleRate int, frequency float64) int {
	return int(frequency / float64(sampleRate) * t.SineTableSize * t.PhasePrecision)
}

func durationMs(periods []t.Period) int {
	if len(periods) == 0 {
		return 0
	}

	return periods[len(periods)-1].Time
}

func totalFramesFromDuration(durationMs int, sampleRate int) int64 {
	return int64(math.Round(float64(durationMs) * float64(sampleRate) / 1000.0))
}