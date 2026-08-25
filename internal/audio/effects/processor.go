// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

const (
	modulationSlewTimeMs = 8.0
	panSlewTimeMs        = 2.0
)

type Processor struct {
	sampleRate         int
	waveTables         [][]int
	modulationMaxDelta float64
	panMaxDelta        float64
}

func NewProcessor(sampleRate int, waveTables [][]int) *Processor {
	processor := &Processor{sampleRate: sampleRate, waveTables: waveTables}
	processor.modulationMaxDelta = processor.effectSlewMaxDelta(modulationSlewTimeMs)
	processor.panMaxDelta = 2 * processor.effectSlewMaxDelta(panSlewTimeMs)
	return processor
}
