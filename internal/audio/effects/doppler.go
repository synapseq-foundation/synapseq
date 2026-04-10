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

func (p *Processor) ApplyDoppler(channel *t.Channel, effect t.Effect, increment int) int {
	if effect.Type != t.EffectDoppler {
		return increment
	}

	p.advanceEffectPhase(channel)
	factor := p.calcDopplerFactor(channel.Effect.Offset, effect.Intensity)
	return int(math.Round(float64(increment) * factor))
}

func (p *Processor) ApplyDopplerPair(channel *t.Channel, effect t.Effect, inc0, inc1 int) (int, int) {
	if effect.Type != t.EffectDoppler {
		return inc0, inc1
	}

	p.advanceEffectPhase(channel)
	factor := p.calcDopplerFactor(channel.Effect.Offset, effect.Intensity)
	return int(math.Round(float64(inc0) * factor)), int(math.Round(float64(inc1) * factor))
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