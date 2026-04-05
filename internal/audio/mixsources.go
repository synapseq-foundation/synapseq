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

func (r *AudioRenderer) mixPureTone(channel *t.Channel) stereoSample {
	inc0 := r.effectProcessor.ApplyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)

	sample := channel.Amplitude[0] * r.effectProcessor.WaveformSample(channel, channel.Offset[0])
	return r.applyEffectToMono(channel, sample)
}

func (r *AudioRenderer) mixBinauralBeat(channel *t.Channel) stereoSample {
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	left := channel.Amplitude[0] * r.effectProcessor.WaveformSample(channel, channel.Offset[0])
	right := channel.Amplitude[1] * r.effectProcessor.WaveformSample(channel, channel.Offset[1])
	return r.applyEffectToStereo(channel, left, right)
}

func (r *AudioRenderer) mixMonauralBeat(channel *t.Channel) stereoSample {
	inc0, inc1 := r.effectProcessor.ApplyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	freqHigh := r.effectProcessor.WaveformSample(channel, channel.Offset[0])
	freqLow := r.effectProcessor.WaveformSample(channel, channel.Offset[1])
	mixed := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

	return r.applyEffectToMono(channel, mixed)
}

func (r *AudioRenderer) mixIsochronicBeat(channel *t.Channel) stereoSample {
	incCarrier := r.effectProcessor.ApplyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], incCarrier)
	channel.Offset[1] = advancePhase(channel.Offset[1], channel.Increment[1])

	modFactor := r.effectProcessor.CalcModulationFactor(channel, channel.Offset[1])
	carrier := float64(r.effectProcessor.WaveformSample(channel, channel.Offset[0]))
	out := int(float64(channel.Amplitude[0]) * carrier * modFactor)

	return r.applyEffectToMono(channel, out)
}

func (r *AudioRenderer) mixNoise(channel *t.Channel) stereoSample {
	sample := channel.Amplitude[0] * r.noiseGenerator.Generate(channel.Track.Type, channel.Track.NoiseSmooth)
	return r.applyEffectToMono(channel, sample)
}

func (r *AudioRenderer) mixAmbiance(channel *t.Channel, ch, frame int) stereoSample {
	const bgScaleFactor = 16

	if r.ambianceState == nil {
		return stereoSample{}
	}

	bgBuf := r.ambianceState.ChannelBuffer(ch)
	if len(bgBuf) < frame*2+2 {
		return stereoSample{}
	}

	left := bgBuf[frame*2] * bgScaleFactor * channel.Amplitude[0]
	right := bgBuf[frame*2+1] * bgScaleFactor * channel.Amplitude[0]

	return r.applyEffectToStereo(channel, left, right)
}