/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const phaseMask = (t.SineTableSize << 16) - 1

type stereoSample struct {
	left  int
	right int
}

// mix generates a stereo audio sample by mixing all channels
func (r *AudioRenderer) mix(samples []int) []int {
	for i := range t.BufferSize {
		var left, right int

		for ch := range t.NumberOfChannels {
			mixed := r.mixChannelSample(ch, i)
			left += mixed.left
			right += mixed.right
		}

		final := r.finalizeMixedSample(left, right)
		samples[i*2] = final.left
		samples[i*2+1] = final.right
	}

	return samples
}

func (r *AudioRenderer) mixChannelSample(ch, frame int) stereoSample {
	channel := &r.channels[ch]
	waveIdx := int(channel.Track.Waveform)

	switch channel.Track.Type {
	case t.TrackPureTone:
		return r.mixPureTone(channel, waveIdx)
	case t.TrackBinauralBeat:
		return r.mixBinauralBeat(channel, waveIdx)
	case t.TrackMonauralBeat:
		return r.mixMonauralBeat(channel, waveIdx)
	case t.TrackIsochronicBeat:
		return r.mixIsochronicBeat(channel, waveIdx)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return r.mixNoise(channel)
	case t.TrackAmbiance:
		return r.mixAmbiance(channel, ch, frame)
	default:
		return stereoSample{}
	}
}

func (r *AudioRenderer) mixPureTone(channel *t.Channel, waveIdx int) stereoSample {
	inc0 := r.applyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)

	sample := channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
	return r.applyEffectToMono(channel, sample)
}

func (r *AudioRenderer) mixBinauralBeat(channel *t.Channel, waveIdx int) stereoSample {
	inc0, inc1 := r.applyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	left := channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
	right := channel.Amplitude[1] * r.waveTables[waveIdx][channel.Offset[1]>>16]
	return r.applyEffectToStereo(channel, left, right)
}

func (r *AudioRenderer) mixMonauralBeat(channel *t.Channel, waveIdx int) stereoSample {
	inc0, inc1 := r.applyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	freqHigh := r.waveTables[waveIdx][channel.Offset[0]>>16]
	freqLow := r.waveTables[waveIdx][channel.Offset[1]>>16]
	mixed := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

	return r.applyEffectToMono(channel, mixed)
}

func (r *AudioRenderer) mixIsochronicBeat(channel *t.Channel, waveIdx int) stereoSample {
	incCarrier := r.applyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], incCarrier)
	channel.Offset[1] = advancePhase(channel.Offset[1], channel.Increment[1])

	modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Offset[1])
	carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
	out := int(float64(channel.Amplitude[0]) * carrier * modFactor)

	return r.applyEffectToMono(channel, out)
}

func (r *AudioRenderer) mixNoise(channel *t.Channel) stereoSample {
	sample := channel.Amplitude[0] * r.noiseGenerator.Generate(channel.Track.Type, channel.Track.NoiseSmooth)
	return r.applyEffectToMono(channel, sample)
}

func (r *AudioRenderer) mixAmbiance(channel *t.Channel, ch, frame int) stereoSample {
	const bgScaleFactor = 16

	idx := r.channelAmbianceIndex[ch]
	if idx < 0 || idx >= len(r.ambianceSamplesByIndex) {
		return stereoSample{}
	}

	bgBuf := r.ambianceSamplesByIndex[idx]
	if len(bgBuf) < frame*2+2 {
		return stereoSample{}
	}

	left := bgBuf[frame*2] * bgScaleFactor * channel.Amplitude[0]
	right := bgBuf[frame*2+1] * bgScaleFactor * channel.Amplitude[0]

	return r.applyEffectToStereo(channel, left, right)
}

func (r *AudioRenderer) applyDoppler(channel *t.Channel, increment int) int {
	if channel.Track.Effect.Type != t.EffectDoppler {
		return increment
	}

	r.advanceEffectPhase(channel)
	factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
	return int(math.Round(float64(increment) * factor))
}

func (r *AudioRenderer) applyDopplerPair(channel *t.Channel, inc0, inc1 int) (int, int) {
	if channel.Track.Effect.Type != t.EffectDoppler {
		return inc0, inc1
	}

	r.advanceEffectPhase(channel)
	factor := r.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
	return int(math.Round(float64(inc0) * factor)), int(math.Round(float64(inc1) * factor))
}

func (r *AudioRenderer) applyEffectToMono(channel *t.Channel, sample int) stereoSample {
	if channel.Track.Effect.Type == t.EffectModulation {
		sample = r.applyModulation(channel, sample)
	}

	if channel.Track.Effect.Type == t.EffectPan {
		return r.applyPanEffect(channel, sample, sample)
	}

	return stereoSample{left: sample, right: sample}
}

func (r *AudioRenderer) applyEffectToStereo(channel *t.Channel, left, right int) stereoSample {
	if channel.Track.Effect.Type == t.EffectModulation {
		left = r.applyModulation(channel, left)
		right = r.applyModulationToCurrentPhase(channel, right)
		return stereoSample{left: left, right: right}
	}

	if channel.Track.Effect.Type == t.EffectPan {
		return r.applyPanEffect(channel, left, right)
	}

	return stereoSample{left: left, right: right}
}

func (r *AudioRenderer) applyModulation(channel *t.Channel, sample int) int {
	r.advanceEffectPhase(channel)
	return r.applyModulationToCurrentPhase(channel, sample)
}

func (r *AudioRenderer) applyModulationToCurrentPhase(channel *t.Channel, sample int) int {
	modFactor := r.calcModulationFactor(channel.Track.Waveform, channel.Effect.Offset)
	effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
	gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
	return int(float64(sample) * gain)
}

func (r *AudioRenderer) applyPanEffect(channel *t.Channel, left, right int) stereoSample {
	r.advanceEffectPhase(channel)
	outLeft, outRight := r.applyPan(channel, left, right)
	return stereoSample{left: outLeft, right: outRight}
}

func (r *AudioRenderer) advanceEffectPhase(channel *t.Channel) {
	channel.Effect.Offset = advancePhase(channel.Effect.Offset, channel.Effect.Increment)
}

func advancePhase(offset, increment int) int {
	return (offset + increment) & phaseMask
}

func (r *AudioRenderer) finalizeMixedSample(left, right int) stereoSample {
	if r.Volume != 100 {
		left = left * r.Volume / 100
		right = right * r.Volume / 100
	}

	left >>= audioBitShift
	right >>= audioBitShift

	return stereoSample{
		left:  clampPCM16(left),
		right: clampPCM16(right),
	}
}

func clampPCM16(sample int) int {
	if sample > audioMaxValue {
		return audioMaxValue
	}
	if sample < audioMinValue {
		return audioMinValue
	}

	return sample
}
