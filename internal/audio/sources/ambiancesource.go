// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import "github.com/synapseq-foundation/synapseq/v4/internal/audio/audiosource"

type Ambiance struct {
	amplitude int
}

func NewAmbiance(signal Signal) Ambiance {
	return Ambiance{amplitude: signal.Amplitude[0]}
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

	left := bgBuf[frame*2] * bgScaleFactor * source.amplitude
	right := bgBuf[frame*2+1] * bgScaleFactor * source.amplitude
	return left, right, true
}
