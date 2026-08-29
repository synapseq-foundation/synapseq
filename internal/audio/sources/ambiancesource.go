// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import "github.com/synapseq-foundation/synapseq/v4/internal/audio/audiosource"

type Ambiance struct {
}

func NewAmbiance(signal Signal) Ambiance {
	return Ambiance{}
}

func (source Ambiance) Sample(runtime *audiosource.Runtime, ch, frame int) (int, int, bool) {
	const bgScaleFactor = 16

	if runtime == nil {
		return 0, 0, false
	}

	bgBuf := runtime.ChannelBuffer(ch)
	if len(bgBuf) < frame*2+2 {
		return 0, 0, false
	}

	left := bgBuf[frame*2] * bgScaleFactor
	right := bgBuf[frame*2+1] * bgScaleFactor
	return left, right, true
}

func (source Ambiance) SampleDoppler(runtime *audiosource.Runtime, ch int, rate float64) (int, int, bool) {
	const bgScaleFactor = 16

	left, right, ok := runtime.SampleDoppler(ch, rate)
	if !ok {
		return 0, 0, false
	}
	return left * bgScaleFactor, right * bgScaleFactor, true
}
