// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestWaveformValueForMorphInterpolatesBetweenTables(ts *testing.T) {
	processor := newTestProcessor()
	offset := t.PhasePrecision
	waveform := WaveformMorph{Start: t.WaveformSine, End: t.WaveformSquare, Alpha: 0.25}

	got := processor.WaveformValueForMorph(waveform, offset)
	start := float64(processor.waveTables[int(t.WaveformSine)][1])
	end := float64(processor.waveTables[int(t.WaveformSquare)][1])
	want := start + (end-start)*0.25
	if got != want {
		ts.Fatalf("unexpected morphed waveform value: got %f, want %f", got, want)
	}

	sample := processor.WaveformSampleForMorph(waveform, offset)
	if sample != int(math.Round(want)) {
		ts.Fatalf("unexpected morphed waveform sample: got %d, want %d", sample, int(math.Round(want)))
	}
}
