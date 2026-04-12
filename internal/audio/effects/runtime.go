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

const phaseMask = (t.SineTableSize << 16) - 1

func (p *Processor) advanceEffectPhase(channel *t.Channel) {
	channel.Effect.Offset = advancePhase(channel.Effect.Offset, channel.Effect.Increment)
}

func advancePhase(offset, increment int) int {
	return (offset + increment) & phaseMask
}
