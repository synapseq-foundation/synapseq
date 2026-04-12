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

type PureTone struct {
	waveform  efx.WaveformMorph
	amplitude int
}

func NewPureTone(signal Signal) PureTone {
	return PureTone{waveform: signal.Waveform, amplitude: signal.Amplitude[0]}
}

func (source PureTone) Sample(processor *efx.Processor, phase int) int {
	return source.amplitude * processor.WaveformSampleForMorph(source.waveform, phase)
}