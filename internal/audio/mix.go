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

const modulationSlewTimeMs = 2.0

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

	switch channel.Track.Type {
	case t.TrackPureTone:
		return r.mixPureTone(channel)
	case t.TrackBinauralBeat:
		return r.mixBinauralBeat(channel)
	case t.TrackMonauralBeat:
		return r.mixMonauralBeat(channel)
	case t.TrackIsochronicBeat:
		return r.mixIsochronicBeat(channel)
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return r.mixNoise(channel)
	case t.TrackAmbiance:
		return r.mixAmbiance(channel, ch, frame)
	default:
		return stereoSample{}
	}
}

func (r *AudioRenderer) mixPureTone(channel *t.Channel) stereoSample {
	inc0 := r.applyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)

	sample := channel.Amplitude[0] * r.waveformSample(channel, channel.Offset[0])
	return r.applyEffectToMono(channel, sample)
}

func (r *AudioRenderer) mixBinauralBeat(channel *t.Channel) stereoSample {
	inc0, inc1 := r.applyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	left := channel.Amplitude[0] * r.waveformSample(channel, channel.Offset[0])
	right := channel.Amplitude[1] * r.waveformSample(channel, channel.Offset[1])
	return r.applyEffectToStereo(channel, left, right)
}

func (r *AudioRenderer) mixMonauralBeat(channel *t.Channel) stereoSample {
	inc0, inc1 := r.applyDopplerPair(channel, channel.Increment[0], channel.Increment[1])
	channel.Offset[0] = advancePhase(channel.Offset[0], inc0)
	channel.Offset[1] = advancePhase(channel.Offset[1], inc1)

	freqHigh := r.waveformSample(channel, channel.Offset[0])
	freqLow := r.waveformSample(channel, channel.Offset[1])
	mixed := (channel.Amplitude[0] * (freqHigh + freqLow)) >> 1

	return r.applyEffectToMono(channel, mixed)
}

func (r *AudioRenderer) mixIsochronicBeat(channel *t.Channel) stereoSample {
	incCarrier := r.applyDoppler(channel, channel.Increment[0])
	channel.Offset[0] = advancePhase(channel.Offset[0], incCarrier)
	channel.Offset[1] = advancePhase(channel.Offset[1], channel.Increment[1])

	modFactor := r.calcModulationFactor(channel, channel.Offset[1])
	carrier := float64(r.waveformSample(channel, channel.Offset[0]))
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
	modFactor := r.calcModulationFactor(channel, channel.Effect.Offset)
	effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
	targetGain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
	gain := r.smoothedModulationGain(channel, targetGain)
	return int(float64(sample) * gain)
}

func (r *AudioRenderer) smoothedModulationGain(channel *t.Channel, targetGain float64) float64 {
	if !channel.Effect.ModulationInitialized {
		channel.Effect.ModulationGain = targetGain
		channel.Effect.ModulationInitialized = true
		return targetGain
	}

	maxDelta := r.effectSlewMaxDelta()
	delta := targetGain - channel.Effect.ModulationGain
	if delta > maxDelta {
		delta = maxDelta
	} else if delta < -maxDelta {
		delta = -maxDelta
	}

	channel.Effect.ModulationGain += delta
	return channel.Effect.ModulationGain
}

func (r *AudioRenderer) effectSlewMaxDelta() float64 {
	if r.SampleRate <= 0 {
		return 1
	}

	rampSamples := float64(r.SampleRate) * modulationSlewTimeMs / 1000.0
	if rampSamples < 1 {
		return 1
	}

	return 1 / rampSamples
}

func (r *AudioRenderer) waveformSample(channel *t.Channel, offset int) int {
	return int(math.Round(r.waveformValue(channel, offset)))
}

func (r *AudioRenderer) waveformValue(channel *t.Channel, offset int) float64 {
	startWaveform, endWaveform, alpha := channelWaveformMorph(channel)
	idx := offset >> 16
	start := float64(r.waveTables[int(startWaveform)][idx])
	if alpha <= 0 || startWaveform == endWaveform {
		return start
	}

	end := float64(r.waveTables[int(endWaveform)][idx])
	if alpha >= 1 {
		return end
	}

	return lerpFloat64(start, end, alpha)
}

func channelWaveformMorph(channel *t.Channel) (t.WaveformType, t.WaveformType, float64) {
	if channel.WaveformStart == 0 && channel.WaveformEnd == 0 && channel.WaveformAlpha == 0 {
		return channel.Track.Waveform, channel.Track.Waveform, 0
	}

	return channel.WaveformStart, channel.WaveformEnd, channel.WaveformAlpha
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
