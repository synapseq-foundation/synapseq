// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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