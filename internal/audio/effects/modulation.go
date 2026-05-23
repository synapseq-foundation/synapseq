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
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const modulationSquareEdgeRatio = 0.08

func (p *Processor) ApplyModulation(channel *t.Channel, effect t.Effect, waveform WaveformMorph, sample int) int {
	p.advanceEffectPhase(channel)
	return p.ApplyModulationToCurrentPhase(channel, effect, waveform, sample)
}

func (p *Processor) ApplyModulationToCurrentPhase(channel *t.Channel, effect t.Effect, waveform WaveformMorph, sample int) int {
	modFactor := p.CalcModulationFactorForMorph(waveform, channel.Effect.Offset)
	effectIntensity := float64(effect.Intensity) * 0.7
	targetGain := (1.0 - effectIntensity) + (effectIntensity * modFactor)
	gain := p.smoothedModulationGain(channel, targetGain)
	return int(float64(sample) * gain)
}

func (p *Processor) CalcModulationFactor(channel *t.Channel, offset int) float64 {
	return p.CalcModulationFactorForMorph(WaveformMorphFromChannel(channel), offset)
}

func (p *Processor) CalcModulationFactorForMorph(waveform WaveformMorph, offset int) float64 {
	startWaveform, endWaveform, alpha := normalizedWaveformMorph(waveform)
	start := p.modulationFactorForWaveform(startWaveform, offset)
	if alpha <= 0 || startWaveform == endWaveform {
		return start
	}

	end := p.modulationFactorForWaveform(endWaveform, offset)
	if alpha >= 1 {
		return end
	}

	return lerpFloat64(start, end, alpha)
}

func (p *Processor) modulationFactorForWaveform(waveform t.WaveformType, offset int) float64 {
	if waveform == t.WaveformSquare {
		return softSquareModulationFactor(offset)
	}

	modVal := p.WaveformValueForMorph(WaveformMorph{Start: waveform, End: waveform, Alpha: 0}, offset)
	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}

func softSquareModulationFactor(offset int) float64 {
	cycle := float64(t.SineTableSize * t.PhasePrecision)
	phase := float64(offset&(t.SineTableSize*t.PhasePrecision-1)) / cycle
	edge := modulationSquareEdgeRatio
	halfEdge := edge / 2

	switch {
	case phase < halfEdge:
		return smoothstep((phase + halfEdge) / edge)
	case phase > 1-halfEdge:
		return smoothstep((phase - (1 - halfEdge)) / edge)
	case phase < 0.5-halfEdge:
		return 1
	case phase < 0.5+halfEdge:
		return 1 - smoothstep((phase-(0.5-halfEdge))/edge)
	default:
		return 0
	}
}

func smoothstep(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	return x * x * (3 - 2*x)
}
