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

package effects

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	modulationSlewTimeMs = 2.0
	phaseMask            = (t.SineTableSize << 16) - 1
)

type Processor struct {
	sampleRate int
	waveTables [4][]int
}

func NewProcessor(sampleRate int, waveTables [4][]int) *Processor {
	return &Processor{sampleRate: sampleRate, waveTables: waveTables}
}

func (p *Processor) ApplyDoppler(channel *t.Channel, increment int) int {
	if channel.Track.Effect.Type != t.EffectDoppler {
		return increment
	}

	p.advanceEffectPhase(channel)
	factor := p.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
	return int(math.Round(float64(increment) * factor))
}

func (p *Processor) ApplyDopplerPair(channel *t.Channel, inc0, inc1 int) (int, int) {
	if channel.Track.Effect.Type != t.EffectDoppler {
		return inc0, inc1
	}

	p.advanceEffectPhase(channel)
	factor := p.calcDopplerFactor(channel.Effect.Offset, channel.Track.Effect.Intensity)
	return int(math.Round(float64(inc0) * factor)), int(math.Round(float64(inc1) * factor))
}

func (p *Processor) ApplyEffectToMono(channel *t.Channel, sample int) (int, int) {
	if channel.Track.Effect.Type == t.EffectModulation {
		sample = p.ApplyModulation(channel, sample)
	}

	if channel.Track.Effect.Type == t.EffectPan {
		return p.ApplyPanEffect(channel, sample, sample)
	}

	return sample, sample
}

func (p *Processor) ApplyEffectToStereo(channel *t.Channel, left, right int) (int, int) {
	if channel.Track.Effect.Type == t.EffectModulation {
		left = p.ApplyModulation(channel, left)
		right = p.ApplyModulationToCurrentPhase(channel, right)
		return left, right
	}

	if channel.Track.Effect.Type == t.EffectPan {
		return p.ApplyPanEffect(channel, left, right)
	}

	return left, right
}

func (p *Processor) ApplyModulation(channel *t.Channel, sample int) int {
	p.advanceEffectPhase(channel)
	return p.ApplyModulationToCurrentPhase(channel, sample)
}

func (p *Processor) ApplyModulationToCurrentPhase(channel *t.Channel, sample int) int {
	modFactor := p.CalcModulationFactor(channel, channel.Effect.Offset)
	effectIntensity := float64(channel.Track.Effect.Intensity) * 0.7
	targetGain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
	gain := p.smoothedModulationGain(channel, targetGain)
	return int(float64(sample) * gain)
}

func (p *Processor) ApplyPanEffect(channel *t.Channel, left, right int) (int, int) {
	p.advanceEffectPhase(channel)
	return p.ApplyPan(channel, left, right)
}

func (p *Processor) ApplyPan(channel *t.Channel, inL, inR int) (outL, outR int) {
	intensity := float64(channel.Track.Effect.Intensity)
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}

	targetX := (p.WaveformValue(channel, channel.Effect.Offset) / float64(t.WaveTableAmplitude)) * intensity
	x := p.smoothedPanPosition(channel, targetX)

	pos := int(math.Round((x + 1.0) * 32768.0))
	if pos < 0 {
		pos = 0
	}
	if pos > 65536 {
		pos = 65536
	}

	lGain := 65536 - pos
	rGain := pos

	outL = int((int64(inL) * int64(lGain)) >> 16)
	outR = int((int64(inR) * int64(rGain)) >> 16)
	return outL, outR
}

func (p *Processor) CalcModulationFactor(channel *t.Channel, offset int) float64 {
	modVal := p.WaveformValue(channel, offset)

	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}

func (p *Processor) WaveformSample(channel *t.Channel, offset int) int {
	return int(math.Round(p.WaveformValue(channel, offset)))
}

func (p *Processor) WaveformValue(channel *t.Channel, offset int) float64 {
	startWaveform, endWaveform, alpha := channelWaveformMorph(channel)
	idx := offset >> 16
	start := float64(p.waveTables[int(startWaveform)][idx])
	if alpha <= 0 || startWaveform == endWaveform {
		return start
	}

	end := float64(p.waveTables[int(endWaveform)][idx])
	if alpha >= 1 {
		return end
	}

	return lerpFloat64(start, end, alpha)
}

func (p *Processor) calcDopplerFactor(offset int, intensity t.IntensityType) float64 {
	inten := float64(intensity)
	if inten < 0 {
		inten = 0
	}
	if inten > 1 {
		inten = 1
	}

	lfo := p.waveTables[int(t.WaveformSine)][offset>>16]
	lfoNorm := float64(lfo) / float64(t.WaveTableAmplitude)

	depth := 0.05 * inten
	return 1.0 + (depth * lfoNorm)
}

func (p *Processor) smoothedModulationGain(channel *t.Channel, targetGain float64) float64 {
	if !channel.Effect.ModulationInitialized {
		channel.Effect.ModulationGain = targetGain
		channel.Effect.ModulationInitialized = true
		return targetGain
	}

	maxDelta := p.effectSlewMaxDelta()
	delta := targetGain - channel.Effect.ModulationGain
	if delta > maxDelta {
		delta = maxDelta
	} else if delta < -maxDelta {
		delta = -maxDelta
	}

	channel.Effect.ModulationGain += delta
	return channel.Effect.ModulationGain
}

func (p *Processor) effectSlewMaxDelta() float64 {
	if p.sampleRate <= 0 {
		return 1
	}

	rampSamples := float64(p.sampleRate) * modulationSlewTimeMs / 1000.0
	if rampSamples < 1 {
		return 1
	}

	return 1 / rampSamples
}

func (p *Processor) smoothedPanPosition(channel *t.Channel, targetX float64) float64 {
	if !channel.Effect.PanInitialized {
		channel.Effect.PanPosition = targetX
		channel.Effect.PanInitialized = true
		return targetX
	}

	maxDelta := 2 * p.effectSlewMaxDelta()
	delta := targetX - channel.Effect.PanPosition
	if delta > maxDelta {
		delta = maxDelta
	} else if delta < -maxDelta {
		delta = -maxDelta
	}

	channel.Effect.PanPosition += delta
	if channel.Effect.PanPosition > 1 {
		channel.Effect.PanPosition = 1
	}
	if channel.Effect.PanPosition < -1 {
		channel.Effect.PanPosition = -1
	}

	return channel.Effect.PanPosition
}

func (p *Processor) advanceEffectPhase(channel *t.Channel) {
	channel.Effect.Offset = advancePhase(channel.Effect.Offset, channel.Effect.Increment)
}

func advancePhase(offset, increment int) int {
	return (offset + increment) & phaseMask
}

func channelWaveformMorph(channel *t.Channel) (t.WaveformType, t.WaveformType, float64) {
	if channel.WaveformStart == 0 && channel.WaveformEnd == 0 && channel.WaveformAlpha == 0 {
		return channel.Track.Waveform, channel.Track.Waveform, 0
	}

	return channel.WaveformStart, channel.WaveformEnd, channel.WaveformAlpha
}

func lerpFloat64(start, end, alpha float64) float64 {
	return start + (end-start)*alpha
}