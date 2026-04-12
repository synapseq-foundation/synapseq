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