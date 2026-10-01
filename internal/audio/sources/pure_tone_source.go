// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type PureTone struct {
	waveform efx.WaveformMorph
}

func NewPureTone(signal Signal) PureTone {
	return PureTone{waveform: signal.Waveform}
}

func (source PureTone) Sample(processor *efx.Processor, phase int) int {
	return processor.WaveformSampleForMorph(source.waveform, phase)
}
