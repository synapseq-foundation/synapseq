// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

const (
	modulationSlewTimeMs = 8.0
	panSlewTimeMs        = 2.0
)

type Processor struct {
	sampleRate int
	waveTables [4][]int
}

func NewProcessor(sampleRate int, waveTables [4][]int) *Processor {
	return &Processor{sampleRate: sampleRate, waveTables: waveTables}
}
