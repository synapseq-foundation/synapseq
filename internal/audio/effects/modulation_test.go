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

func TestCalcModulationFactorForSquareUsesSoftEdges(ts *testing.T) {
	processor := newTestProcessor()
	waveform := WaveformMorph{Start: t.WaveformSquare, End: t.WaveformSquare, Alpha: 0}

	if got := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0)); !nearlyEqual(got, 0.5) {
		ts.Fatalf("expected square modulation to start on soft edge, got %f", got)
	}
	if got := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.25)); !nearlyEqual(got, 1) {
		ts.Fatalf("expected square modulation high plateau, got %f", got)
	}
	if got := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.5)); !nearlyEqual(got, 0.5) {
		ts.Fatalf("expected square modulation falling soft edge, got %f", got)
	}
	if got := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.75)); !nearlyEqual(got, 0) {
		ts.Fatalf("expected square modulation low plateau, got %f", got)
	}
}

func phaseOffset(cycleFraction float64) int {
	return int(float64(t.SineTableSize*t.PhasePrecision) * cycleFraction)
}

func nearlyEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-12
}
