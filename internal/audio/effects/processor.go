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

const (
	modulationSlewTimeMs = 2.0
)

type Processor struct {
	sampleRate int
	waveTables [4][]int
}

func NewProcessor(sampleRate int, waveTables [4][]int) *Processor {
	return &Processor{sampleRate: sampleRate, waveTables: waveTables}
}