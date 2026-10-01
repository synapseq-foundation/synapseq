// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"math"

	audiosync "github.com/synapseq-foundation/synapseq/v4/internal/audio/sync"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	defaultPlanWindowMs = 1000
)

type renderPlan struct {
	periods     []t.Period
	windows     []renderWindow
	sampleRate  int
	totalFrames int64
	waveforms   []periodWaveforms
	active      [][t.NumberOfChannels]bool
	transitions map[string][]float64
}

type periodWaveforms struct {
	Start        [t.NumberOfChannels]wt.ID
	End          [t.NumberOfChannels]wt.ID
	CrossfadeOut [t.NumberOfChannels]wt.ID
	CrossfadeIn  [t.NumberOfChannels]wt.ID
}

type renderWindow struct {
	PeriodIndex int
	StartMs     int
	EndMs       int
}

func compileRenderPlan(periods []t.Period, sampleRate int) renderPlan {
	waveforms, _ := wt.Compile(nil)
	return compileRenderPlanWithWaveformsAndTransitions(periods, sampleRate, waveforms, nil)
}

func compileRenderPlanWithWaveforms(periods []t.Period, sampleRate int, registry *wt.Registry) renderPlan {
	return compileRenderPlanWithWaveformsAndTransitions(periods, sampleRate, registry, nil)
}

func compileRenderPlanWithWaveformsAndTransitions(periods []t.Period, sampleRate int, registry *wt.Registry, definitions []t.TransitionDefinition) renderPlan {
	transitions := make(map[string][]float64, len(definitions))
	for _, definition := range definitions {
		transitions[definition.Name] = definition.Points
	}
	plan := renderPlan{
		periods:     periods,
		windows:     make([]renderWindow, len(periods)),
		sampleRate:  sampleRate,
		totalFrames: totalFramesFromDuration(durationMs(periods), sampleRate),
		waveforms:   make([]periodWaveforms, len(periods)),
		active:      make([][t.NumberOfChannels]bool, len(periods)),
		transitions: transitions,
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
		for channel := range t.NumberOfChannels {
			plan.active[index][channel] = periodCanProduceAudio(periods[index], channel)
			plan.waveforms[index].Start[channel], _ = registry.Lookup(periods[index].TrackStart[channel].Waveform)
			plan.waveforms[index].End[channel], _ = registry.Lookup(periods[index].TrackEnd[channel].Waveform)
			plan.waveforms[index].CrossfadeOut[channel], _ = registry.Lookup(periods[index].CrossfadeOut[channel].Track.Waveform)
			plan.waveforms[index].CrossfadeIn[channel], _ = registry.Lookup(periods[index].CrossfadeIn[channel].Track.Waveform)
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
	alpha := tl.StepAlphaWithPoints(rp.interpolationProgress(window, currentTimeMs), period.Transition, rp.transitions[period.TransitionName], period.Steps)
	cue := audiosync.Cue{
		PeriodIndex: window.PeriodIndex,
		Channels:    [t.NumberOfChannels]audiosync.ChannelCue{},
	}

	for index := 0; index < t.NumberOfChannels; index++ {
		if !rp.active[periodIdx][index] {
			continue
		}
		track, waveformStart, waveformEnd, waveformAlpha, crossfade := rp.trackStateAt(window, period, index, alpha, currentTimeMs)
		signal := compileSignalState(planTrackState{
			track:      track,
			sampleRate: rp.sampleRate,
		})
		cue.Channels[index] = audiosync.ChannelCue{
			Track:         signal.Track,
			WaveformStart: waveformStart,
			WaveformEnd:   waveformEnd,
			WaveformAlpha: waveformAlpha,
			Amplitude:     signal.Amplitude,
			Increment:     signal.Increment,
			EffectStep:    signal.EffectStep,
			Crossfade:     crossfade,
		}
	}

	return cue
}

func periodCanProduceAudio(period t.Period, channel int) bool {
	tracks := [...]t.Track{
		period.TrackStart[channel],
		period.TrackEnd[channel],
		period.CrossfadeOut[channel].Track,
		period.CrossfadeIn[channel].Track,
	}
	for _, track := range tracks {
		if track.Type != t.TrackOff && track.Type != t.TrackSilence {
			return true
		}
	}
	return false
}

func (rp renderPlan) trackStateAt(window renderWindow, period t.Period, ch int, alpha float64, currentTimeMs int) (t.Track, wt.ID, wt.ID, float64, audiosync.CrossfadeCue) {
	waveforms := rp.waveforms[window.PeriodIndex]
	if crossfade := period.CrossfadeOut[ch]; crossfade.Active {
		duration := tl.CrossfadeDuration(window.EndMs - window.StartMs)
		if duration > 0 && currentTimeMs >= window.EndMs-duration {
			progress := clampUnit(float64(currentTimeMs-(window.EndMs-duration)) / float64(duration))
			track := scaleTrackAmplitude(crossfade.Track, 1-progress)
			return track, waveforms.CrossfadeOut[ch], waveforms.CrossfadeOut[ch], 0, audiosync.CrossfadeCue{Active: true, Direction: audiosync.CrossfadeOut, Alpha: progress}
		}
		track := crossfade.Track
		return track, waveforms.CrossfadeOut[ch], waveforms.CrossfadeOut[ch], 0, audiosync.CrossfadeCue{}
	}

	if crossfade := period.CrossfadeIn[ch]; crossfade.Active {
		duration := tl.CrossfadeDuration(window.EndMs - window.StartMs)
		if duration > 0 && currentTimeMs <= window.StartMs+duration {
			progress := clampUnit(float64(currentTimeMs-window.StartMs) / float64(duration))
			track := scaleTrackAmplitude(crossfade.Track, progress)
			return track, waveforms.CrossfadeIn[ch], waveforms.CrossfadeIn[ch], 0, audiosync.CrossfadeCue{Active: true, Direction: audiosync.CrossfadeIn, Alpha: progress}
		}
		track := crossfade.Track
		return track, waveforms.CrossfadeIn[ch], waveforms.CrossfadeIn[ch], 0, audiosync.CrossfadeCue{}
	}

	track := interpolateTrack(period.TrackStart[ch], period.TrackEnd[ch], alpha)
	return track, waveforms.Start[ch], waveforms.End[ch], alpha, audiosync.CrossfadeCue{}
}

func scaleTrackAmplitude(track t.Track, scale float64) t.Track {
	scale = clampUnit(scale)
	for channel, amplitude := range track.Amplitude {
		track.Amplitude[channel] = t.AmplitudeType(float64(amplitude) * scale)
	}
	return track
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
		Type: start.Type,
		Amplitude: [2]t.AmplitudeType{
			t.AmplitudeType(lerpFloat64(float64(start.Amplitude[0]), float64(end.Amplitude[0]), alpha)),
			t.AmplitudeType(lerpFloat64(float64(start.Amplitude[1]), float64(end.Amplitude[1]), alpha)),
		},
		Carrier:     lerpFloat64(start.Carrier, end.Carrier, alpha),
		Resonance:   lerpFloat64(start.Resonance, end.Resonance, alpha),
		NoiseSmooth: lerpFloat64(start.NoiseSmooth, end.NoiseSmooth, alpha),
		Waveform:    start.Waveform,
		SourceName:  start.SourceName,
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

	amplitude := [2]int{int(state.track.Amplitude[0]), int(state.track.Amplitude[1])}
	switch state.track.Type {
	case t.TrackPureTone:
		compiled.Amplitude = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, state.track.Carrier)
	case t.TrackBinauralBeat:
		freq1 := state.track.Carrier + state.track.Resonance/2
		freq2 := state.track.Carrier - state.track.Resonance/2
		compiled.Amplitude = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, freq1)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, freq2)
	case t.TrackMonauralBeat:
		freqHigh := state.track.Carrier + state.track.Resonance/2
		freqLow := state.track.Carrier - state.track.Resonance/2
		compiled.Amplitude = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, freqHigh)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, freqLow)
	case t.TrackIsochronicBeat:
		compiled.Amplitude = amplitude
		compiled.Increment[0] = frequencyToIncrement(state.sampleRate, state.track.Carrier)
		compiled.Increment[1] = frequencyToIncrement(state.sampleRate, state.track.Resonance)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance, t.TrackMusic:
		compiled.Amplitude = amplitude
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
