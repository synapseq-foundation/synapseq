// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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