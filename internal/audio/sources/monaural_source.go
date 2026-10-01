// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type Monaural struct {
	waveform efx.WaveformMorph
}

func NewMonaural(signal Signal) Monaural {
	return Monaural{waveform: signal.Waveform}
}

func (source Monaural) Sample(processor *efx.Processor, highPhase, lowPhase int) int {
	high := processor.WaveformSampleForMorph(source.waveform, highPhase)
	low := processor.WaveformSampleForMorph(source.waveform, lowPhase)
	return (high + low) >> 1
}
