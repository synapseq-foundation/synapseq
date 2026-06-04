// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type Monaural struct {
	waveform  efx.WaveformMorph
	amplitude int
}

func NewMonaural(signal Signal) Monaural {
	return Monaural{waveform: signal.Waveform, amplitude: signal.Amplitude[0]}
}

func (source Monaural) Sample(processor *efx.Processor, highPhase, lowPhase int) int {
	high := processor.WaveformSampleForMorph(source.waveform, highPhase)
	low := processor.WaveformSampleForMorph(source.waveform, lowPhase)
	return (source.amplitude * (high + low)) >> 1
}