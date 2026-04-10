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

import (
	"testing"

	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestIsochronicSampleReturnsSilenceWhenModulationFactorIsZero(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewIsochronic(Signal{
		Waveform:  efx.WaveformMorph{Start: t.WaveformSawtooth, End: t.WaveformSawtooth, Alpha: 0},
		Amplitude: [2]int{4096, 0},
	})

	got := source.Sample(processor, t.PhasePrecision, 0)
	if got != 0 {
		ts.Fatalf("unexpected isochronic silence sample: got %d, want 0", got)
	}
}