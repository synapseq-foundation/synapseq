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