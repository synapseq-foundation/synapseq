// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

const phaseMask = (t.SineTableSize << 16) - 1

func (p *Processor) advanceEffectPhase(channel *t.Channel) {
	channel.Effect.Offset = advancePhase(channel.Effect.Offset, channel.Effect.Increment)
}

func advancePhase(offset, increment int) int {
	return (offset + increment) & phaseMask
}
