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

func TestMonauralSampleAveragesDistinctPhases(ts *testing.T) {
	processor := efx.NewProcessor(44100, wt.Init())
	source := NewMonaural(Signal{
		Waveform:  efx.WaveformMorph{Start: t.WaveformSquare, End: t.WaveformSquare, Alpha: 0},
		Amplitude: [2]int{7, 0},
	})

	highPhase := t.PhasePrecision
	lowPhase := int(t.SineTableSize/2) * t.PhasePrecision
	got := source.Sample(processor, highPhase, lowPhase)
	want := (7 * (processor.WaveformSampleForMorph(source.waveform, highPhase) + processor.WaveformSampleForMorph(source.waveform, lowPhase))) >> 1
	if got != want {
		ts.Fatalf("unexpected monaural sample: got %d, want %d", got, want)
	}
}
