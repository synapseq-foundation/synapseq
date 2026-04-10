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

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type Engine struct {
	SampleRate          int
	UpdateAmbianceIndex func(ch int, periodIdx int, trackType t.TrackType)
}

type Cue struct {
	PeriodIndex int
	Channels    [t.NumberOfChannels]ChannelCue
}

type ChannelCue struct {
	Track         t.Track
	WaveformStart t.WaveformType
	WaveformEnd   t.WaveformType
	WaveformAlpha float64
	Amplitude     [2]int
	Increment     [2]int
	EffectStep    int
}

func NewEngine(sampleRate int, updateAmbianceIndex func(ch int, periodIdx int, trackType t.TrackType)) *Engine {
	return &Engine{
		SampleRate:          sampleRate,
		UpdateAmbianceIndex: updateAmbianceIndex,
	}
}

func (e *Engine) Sync(channels []t.Channel, cue Cue) {
	for ch := range channels {
		e.syncChannel(ch, channels, cue.PeriodIndex, cue.Channels[ch])
	}
}

func FrequencyToIncrement(sampleRate int, frequency float64) int {
	return int(frequency / float64(sampleRate) * t.SineTableSize * t.PhasePrecision)
}

func (e *Engine) syncChannel(ch int, channels []t.Channel, periodIdx int, cue ChannelCue) {
	channel := &channels[ch]
	previousTrackType := channel.Type
	previousEffectType := channel.Track.Effect.Type

	channel.Track = cue.Track
	channel.WaveformStart = cue.WaveformStart
	channel.WaveformEnd = cue.WaveformEnd
	channel.WaveformAlpha = cue.WaveformAlpha
	channel.Amplitude = cue.Amplitude
	channel.Increment = cue.Increment
	e.updateAmbianceIndex(ch, periodIdx, cue.Track.Type)
	e.resetRuntimeState(channel, previousTrackType, previousEffectType)
	e.applyEffectState(channel, cue)
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


func (e *Engine) applyEffectState(channel *t.Channel, cue ChannelCue) {
	channel.Effect.Increment = cue.EffectStep
}

