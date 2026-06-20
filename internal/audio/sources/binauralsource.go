// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type Binaural struct {
	waveform  efx.WaveformMorph
	amplitude [2]int
}

func NewBinaural(signal Signal) Binaural {
	return Binaural{waveform: signal.Waveform, amplitude: signal.Amplitude}
}

func (source Binaural) Sample(processor *efx.Processor, leftPhase, rightPhase int) (int, int) {
	left := source.amplitude[0] * processor.WaveformSampleForMorph(source.waveform, leftPhase)
	right := source.amplitude[1] * processor.WaveformSampleForMorph(source.waveform, rightPhase)
	return left, right
}