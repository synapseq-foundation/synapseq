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

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

func (p *Processor) ApplyEffectToMono(channel *t.Channel, effect t.Effect, waveform WaveformMorph, sample int) (int, int) {
	if effect.Type == t.EffectModulation {
		sample = p.ApplyModulation(channel, effect, waveform, sample)
	}

	if effect.Type == t.EffectPan {
		return p.ApplyPanEffect(channel, effect, waveform, sample, sample)
	}

	return sample, sample
}

func (p *Processor) ApplyEffectToStereo(channel *t.Channel, effect t.Effect, waveform WaveformMorph, left, right int) (int, int) {
	if effect.Type == t.EffectModulation {
		left = p.ApplyModulation(channel, effect, waveform, left)
		right = p.ApplyModulationToCurrentPhase(channel, effect, waveform, right)
		return left, right
	}

	if effect.Type == t.EffectPan {
		return p.ApplyPanEffect(channel, effect, waveform, left, right)
	}

	return left, right
}