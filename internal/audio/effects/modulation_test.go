// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

func TestCalcModulationFactorForMorphWithSquareKeepsSoftEdges(ts *testing.T) {
	processor := newTestProcessor()
	waveform := WaveformMorph{Start: t.WaveformSquare, End: t.WaveformSine, Alpha: 0.5}

	before := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.5)-1)
	at := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.5))
	after := processor.CalcModulationFactorForMorph(waveform, phaseOffset(0.5)+1)

	if math.Abs(before-at) > 0.01 || math.Abs(at-after) > 0.01 {
		ts.Fatalf("expected square morph modulation to cross edge smoothly: before=%f at=%f after=%f", before, at, after)
	}
	if at <= 0 || at >= 1 {
		ts.Fatalf("expected interpolated square morph modulation factor between bounds, got %f", at)
	}
}

func TestCalcModulationFactorForMorphInterpolatesFactors(ts *testing.T) {
	processor := newTestProcessor()
	offset := phaseOffset(0.25)
	waveform := WaveformMorph{Start: t.WaveformSquare, End: t.WaveformTriangle, Alpha: 0.25}

	start := processor.modulationFactorForWaveform(t.WaveformSquare, offset)
	end := processor.modulationFactorForWaveform(t.WaveformTriangle, offset)
	want := start + (end-start)*0.25

	got := processor.CalcModulationFactorForMorph(waveform, offset)
	if !nearlyEqual(got, want) {
		ts.Fatalf("unexpected interpolated modulation factor: got %f, want %f", got, want)
	}
}

func phaseOffset(cycleFraction float64) int {
	return int(float64(t.SineTableSize*t.PhasePrecision) * cycleFraction)
}

func nearlyEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-12
}
