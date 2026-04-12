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
