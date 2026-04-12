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

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type stereoSample struct {
	left  int
	right int
}

func (r *AudioRenderer) mixChannelSample(ch, frame int) stereoSample {
	channel := &r.channels[ch]
	signal := r.signalStateFor(ch, channel)

	switch signal.kind {
	case t.TrackPureTone:
		return r.mixPureTone(channel, signal)
	case t.TrackBinauralBeat:
		return r.mixBinauralBeat(channel, signal)
	case t.TrackMonauralBeat:
		return r.mixMonauralBeat(channel, signal)
	case t.TrackIsochronicBeat:
		return r.mixIsochronicBeat(channel, signal)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return r.mixNoise(channel, signal)
	case t.TrackAmbiance:
		return r.mixAmbiance(channel, signal, ch, frame)
	default:
		return stereoSample{}
	}
}

func (r *AudioRenderer) applyEffectToMono(channel *t.Channel, signal channelSignalState, sample int) stereoSample {
	left, right := r.effectProcessor.ApplyEffectToMono(channel, signal.effect, signal.waveform, sample)
	return stereoSample{left: left, right: right}
}

func (r *AudioRenderer) applyEffectToStereo(channel *t.Channel, signal channelSignalState, left, right int) stereoSample {
	left, right = r.effectProcessor.ApplyEffectToStereo(channel, signal.effect, signal.waveform, left, right)
	return stereoSample{left: left, right: right}
}