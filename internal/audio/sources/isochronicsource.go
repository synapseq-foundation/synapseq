// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"

type Isochronic struct {
	waveform  efx.WaveformMorph
	amplitude int
}

func NewIsochronic(signal Signal) Isochronic {
	return Isochronic{waveform: signal.Waveform, amplitude: signal.Amplitude[0]}
}

func (source Isochronic) Sample(processor *efx.Processor, carrierPhase int, modulationFactor float64) int {
	carrier := float64(processor.WaveformSampleForMorph(source.waveform, carrierPhase))
	return int(float64(source.amplitude) * carrier * modulationFactor)
}