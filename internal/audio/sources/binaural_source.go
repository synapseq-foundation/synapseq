// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type Binaural struct {
	waveform efx.WaveformMorph
}

func NewBinaural(signal Signal) Binaural {
	return Binaural{waveform: signal.Waveform}
}

func (source Binaural) Sample(processor *efx.Processor, leftPhase, rightPhase int) (int, int) {
	left := processor.WaveformSampleForMorph(source.waveform, leftPhase)
	right := processor.WaveformSampleForMorph(source.waveform, rightPhase)
	return left, right
}
