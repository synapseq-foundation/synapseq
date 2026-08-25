// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type stereoSample struct {
	left  int
	right int
}

func (r *AudioRenderer) mixChannelSample(ch, frame int) stereoSample {
	channel := &r.channels[ch]
	signal := &r.signals[ch]
	if !signal.resolved {
		fallback := r.signalStateFor(ch, channel)
		signal = &fallback
	}

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
	case t.TrackMusic:
		return r.mixMusic(channel, signal, ch, frame)
	default:
		return stereoSample{}
	}
}

func (r *AudioRenderer) applyEffectToMono(channel *t.Channel, signal *channelSignalState, sample int) stereoSample {
	left, right := r.effectProcessor.ApplyEffectToMono(channel, signal.effect, signal.waveform, sample)
	return applyAmplitude(signal.amplitude, left, right)
}

func (r *AudioRenderer) applyEffectToStereo(channel *t.Channel, signal *channelSignalState, left, right int) stereoSample {
	left, right = r.effectProcessor.ApplyEffectToStereo(channel, signal.effect, signal.waveform, left, right)
	return applyAmplitude(signal.amplitude, left, right)
}

func applyAmplitude(amplitude [2]int, left, right int) stereoSample {
	return stereoSample{
		left:  int(int64(left) * int64(amplitude[0])),
		right: int(int64(right) * int64(amplitude[1])),
	}
}
