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
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	src "github.com/synapseq-foundation/synapseq/v4/internal/audio/sources"
	audiosync "github.com/synapseq-foundation/synapseq/v4/internal/audio/sync"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type channelSignalState struct {
	resolved    bool
	kind        t.TrackType
	noiseSmooth float64
	effect      t.Effect
	waveform    efx.WaveformMorph
	amplitude   [2]int
	increment   [2]int
}

func (r *AudioRenderer) applyCueSignalState(cue audiosync.Cue) {
	for index := range cue.Channels {
		channelCue := cue.Channels[index]
		r.signals[index] = channelSignalState{
			resolved:    true,
			kind:        channelCue.Track.Type,
			noiseSmooth: channelCue.Track.NoiseSmooth,
			effect:      channelCue.Track.Effect,
			waveform: efx.WaveformMorph{
				Start: channelCue.WaveformStart,
				End:   channelCue.WaveformEnd,
				Alpha: channelCue.WaveformAlpha,
			},
			amplitude:   channelCue.Amplitude,
			increment:   channelCue.Increment,
		}
	}
}

func (r *AudioRenderer) signalStateFor(ch int, channel *t.Channel) channelSignalState {
	if r.signals[ch].resolved {
		return r.signals[ch]
	}

	return channelSignalState{
		resolved:    false,
		kind:        channel.Track.Type,
		noiseSmooth: channel.Track.NoiseSmooth,
		effect:      channel.Track.Effect,
		waveform:    efx.WaveformMorphFromChannel(channel),
		amplitude:   channel.Amplitude,
		increment:   channel.Increment,
	}
}

func (state channelSignalState) sourceSignal() src.Signal {
	return src.Signal{
		Kind:        state.kind,
		NoiseSmooth: state.noiseSmooth,
		Waveform:    state.waveform,
		Amplitude:   state.amplitude,
	}
}