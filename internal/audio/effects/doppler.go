// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
	factor := p.calcDopplerFactor(WaveformMorphFromChannel(channel), channel.Effect.Offset, effect.Intensity)
	return int(math.Round(float64(increment) * factor))
}

func (p *Processor) ApplyDopplerPair(channel *t.Channel, effect t.Effect, inc0, inc1 int) (int, int) {
	if effect.Type != t.EffectDoppler {
		return inc0, inc1
	}

	p.advanceEffectPhase(channel)
	factor := p.calcDopplerFactor(WaveformMorphFromChannel(channel), channel.Effect.Offset, effect.Intensity)
	return int(math.Round(float64(inc0) * factor)), int(math.Round(float64(inc1) * factor))
}

func (p *Processor) calcDopplerFactor(waveform WaveformMorph, offset int, intensity t.IntensityType) float64 {
	inten := float64(intensity)
	if inten < 0 {
		inten = 0
	}
	if inten > 1 {
		inten = 1
	}

	lfoNorm := p.WaveformValueForMorph(waveform, offset) / float64(t.WaveTableAmplitude)

	depth := 0.05 * inten
	return 1.0 + (depth * lfoNorm)
}
