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

import wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"

func newTestProcessor() *Processor {
	return NewProcessor(44100, wt.Init())
}
