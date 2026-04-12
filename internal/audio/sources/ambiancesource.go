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

package sources

import amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"

type Ambiance struct {
	amplitude int
}

func NewAmbiance(signal Signal) Ambiance {
	return Ambiance{amplitude: signal.Amplitude[0]}
}

func (source Ambiance) Sample(runtime *amb.Runtime, ch, frame int) (int, int, bool) {
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