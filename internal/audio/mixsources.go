// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	src "github.com/synapseq-foundation/synapseq/v4/internal/audio/sources"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func (r *AudioRenderer) mixPureTone(ch int, channel *t.Channel, signal *channelSignalState) stereoSample {
	source := src.NewPureTone(signal.sourceSignal())
	inc0 := r.effectProcessor.ApplyDoppler(channel, signal.effect, signal.increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)

	sample := source.Sample(r.effectProcessor, channel.Offset[0])
	return r.applyEffectToMono(ch, channel, signal, sample)
}

func (r *AudioRenderer) mixBinauralBeat(ch int, channel *t.Channel, signal *channelSignalState) stereoSample {
	source := src.NewBinaural(signal.sourceSignal())
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, signal.effect, signal.increment[0], signal.increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	left, right := source.Sample(r.effectProcessor, channel.Offset[0], channel.Offset[1])
	return r.applyEffectToStereo(ch, channel, signal, left, right)
}

func (r *AudioRenderer) mixMonauralBeat(ch int, channel *t.Channel, signal *channelSignalState) stereoSample {
	source := src.NewMonaural(signal.sourceSignal())
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, signal.effect, signal.increment[0], signal.increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	mixed := source.Sample(r.effectProcessor, channel.Offset[0], channel.Offset[1])

	return r.applyEffectToMono(ch, channel, signal, mixed)
}

func (r *AudioRenderer) mixIsochronicBeat(ch int, channel *t.Channel, signal *channelSignalState) stereoSample {
	source := src.NewIsochronic(signal.sourceSignal())
	incCarrier := r.effectProcessor.ApplyDoppler(channel, signal.effect, signal.increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], incCarrier)
	channel.Offset[1] = advancePhase(channel.Offset[1], signal.increment[1])

	modFactor := r.effectProcessor.CalcModulationFactorForMorph(signal.waveform, channel.Offset[1])
	out := source.Sample(r.effectProcessor, channel.Offset[0], modFactor)

	return r.applyEffectToMono(ch, channel, signal, out)
}

func (r *AudioRenderer) mixNoise(ch int, channel *t.Channel, signal *channelSignalState) stereoSample {
	source := src.NewNoise(signal.sourceSignal())
	sample := source.Sample(r.noiseGenerator)
	return r.applyEffectToMono(ch, channel, signal, sample)
}

func (r *AudioRenderer) mixAmbiance(channel *t.Channel, signal *channelSignalState, ch, frame int) stereoSample {
	source := src.NewAmbiance(signal.sourceSignal())
	left, right, ok := source.Sample(r.ambianceState, ch, frame)
	if !ok {
		return stereoSample{}
	}

	return r.applyEffectToStereo(ch, channel, signal, left, right)
}

func (r *AudioRenderer) mixMusic(channel *t.Channel, signal *channelSignalState, ch, frame int) stereoSample {
	source := src.NewAmbiance(signal.sourceSignal())
	left, right, ok := source.Sample(r.musicState, ch, frame)
	if !ok {
		return stereoSample{}
	}

	return r.applyEffectToStereo(ch, channel, signal, left, right)
}
