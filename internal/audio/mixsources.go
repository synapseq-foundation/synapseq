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
	src "github.com/synapseq-foundation/synapseq/v4/internal/audio/sources"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func (r *AudioRenderer) mixPureTone(channel *t.Channel, signal channelSignalState) stereoSample {
	source := src.NewPureTone(signal.sourceSignal())
	inc0 := r.effectProcessor.ApplyDoppler(channel, signal.effect, signal.increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)

	sample := source.Sample(r.effectProcessor, channel.Offset[0])
	return r.applyEffectToMono(channel, signal, sample)
}

func (r *AudioRenderer) mixBinauralBeat(channel *t.Channel, signal channelSignalState) stereoSample {
	source := src.NewBinaural(signal.sourceSignal())
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, signal.effect, signal.increment[0], signal.increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	left, right := source.Sample(r.effectProcessor, channel.Offset[0], channel.Offset[1])
	return r.applyEffectToStereo(channel, signal, left, right)
}

func (r *AudioRenderer) mixMonauralBeat(channel *t.Channel, signal channelSignalState) stereoSample {
	source := src.NewMonaural(signal.sourceSignal())
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, signal.effect, signal.increment[0], signal.increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	mixed := source.Sample(r.effectProcessor, channel.Offset[0], channel.Offset[1])

	return r.applyEffectToMono(channel, signal, mixed)
}

func (r *AudioRenderer) mixIsochronicBeat(channel *t.Channel, signal channelSignalState) stereoSample {
	source := src.NewIsochronic(signal.sourceSignal())
	incCarrier := r.effectProcessor.ApplyDoppler(channel, signal.effect, signal.increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], incCarrier)
	channel.Offset[1] = advancePhase(channel.Offset[1], signal.increment[1])

	modFactor := r.effectProcessor.CalcModulationFactorForMorph(signal.waveform, channel.Offset[1])
	out := source.Sample(r.effectProcessor, channel.Offset[0], modFactor)

	return r.applyEffectToMono(channel, signal, out)
}

func (r *AudioRenderer) mixNoise(channel *t.Channel, signal channelSignalState) stereoSample {
	source := src.NewNoise(signal.sourceSignal())
	sample := source.Sample(r.noiseGenerator)
	return r.applyEffectToMono(channel, signal, sample)
}

func (r *AudioRenderer) mixAmbiance(channel *t.Channel, signal channelSignalState, ch, frame int) stereoSample {
	source := src.NewAmbiance(signal.sourceSignal())
	left, right, ok := source.Sample(r.ambianceState, ch, frame)
	if !ok {
		return stereoSample{}
	}

	return r.applyEffectToStereo(channel, signal, left, right)
}
