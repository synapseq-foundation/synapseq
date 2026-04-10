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
	modVal := p.WaveformValueForMorph(waveform, offset)

	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}