// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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