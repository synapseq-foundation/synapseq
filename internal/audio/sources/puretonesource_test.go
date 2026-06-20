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

func TestPureToneSampleUsesCompiledWaveformMorph(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewPureTone(Signal{
		Waveform: efx.WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25},
		Amplitude: [2]int{4096, 0},
	})

	got := source.Sample(processor, t.PhasePrecision)
	want := 4096 * processor.WaveformSampleForMorph(efx.WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25}, t.PhasePrecision)
	if got != want {
		ts.Fatalf("unexpected pure tone sample: got %d, want %d", got, want)
	}
}