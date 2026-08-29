// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

const (
	modulationSlewTimeMs = 8.0
	panSlewTimeMs        = 2.0
)

type Processor struct {
	sampleRate         int
	waveTables         [][]int
	modulationMaxDelta float64
	panMaxDelta        float64
	shiftCoefficients  [shiftTaps]float64
	shiftStates        [t.NumberOfChannels]shiftState
}

func NewProcessor(sampleRate int, waveTables [][]int) *Processor {
	processor := &Processor{sampleRate: sampleRate, waveTables: waveTables}
	processor.modulationMaxDelta = processor.effectSlewMaxDelta(modulationSlewTimeMs)
	processor.panMaxDelta = 2 * processor.effectSlewMaxDelta(panSlewTimeMs)
	processor.shiftCoefficients = makeShiftCoefficients()
	return processor
}
