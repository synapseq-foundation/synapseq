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