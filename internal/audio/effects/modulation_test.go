// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"
	"testing"

	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestCalcModulationFactorForSquareUsesSoftEdges(ts *testing.T) {
	processor := newTestProcessor()
	waveform := WaveformMorph{Start: wt.SquareID, End: wt.SquareID, Alpha: 0}

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
	waveform := WaveformMorph{Start: wt.SquareID, End: wt.SineID, Alpha: 0.5}

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
	waveform := WaveformMorph{Start: wt.SquareID, End: wt.TriangleID, Alpha: 0.25}

	start := processor.modulationFactorForWaveform(wt.SquareID, offset)
	end := processor.modulationFactorForWaveform(wt.TriangleID, offset)
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

func TestIsochronicEnvelopeModes(ts *testing.T) {
	p := newTestProcessor()
	for _, id := range []wt.ID{wt.SineID, wt.SquareID, wt.TriangleID, wt.SawtoothID} {
		waveform := WaveformMorph{Start: id, End: id}
		for _, fraction := range []float64{0, 0.01, 0.125, 0.25, 0.49, 0.5, 0.75, 1.125} {
			offset := phaseOffset(fraction) & (t.SineTableSize*t.PhasePrecision - 1)
			if got, want := p.CalcIsochronicFactor(waveform, offset, t.IsochronicStandard), p.CalcModulationFactorForMorph(waveform, offset); got != want {
				ts.Fatal("standard changed")
			}
			got := p.CalcIsochronicFactor(waveform, offset, t.IsochronicShort)
			normalized := fraction - math.Floor(fraction)
			if normalized >= 0.5 {
				if got != 0 {
					ts.Fatal("second half not silent")
				}
				continue
			}
			position := float64(offset&(t.SineTableSize*t.PhasePrecision-1)) / float64(t.SineTableSize*t.PhasePrecision/2)
			want := 0.75 * smoothstep(position/0.08) * smoothstep((1-position)/0.08) * p.CalcModulationFactorForMorph(waveform, (offset&(t.SineTableSize*t.PhasePrecision-1))*2)
			if !nearlyEqual(got, want) {
				ts.Fatalf("got %f want %f", got, want)
			}
		}
	}
	square := WaveformMorph{Start: wt.SquareID, End: wt.SquareID}
	if p.CalcIsochronicFactor(square, phaseOffset(0.125), t.IsochronicShort) != 0.75 {
		ts.Fatal("wrong peak gain")
	}
	if p.CalcIsochronicFactor(square, 0, t.IsochronicShort) != 0 {
		ts.Fatal("edge not zero")
	}
	if p.CalcIsochronicFactor(square, 1, t.IsochronicShort) > 0.001 {
		ts.Fatal("edge discontinuity")
	}
}

func TestShortIsochronicCustomEnvelopeEdges(ts *testing.T) {
	tables := wt.Init()
	flat := make([]int, t.SineTableSize)
	for i := range flat {
		flat[i] = t.WaveTableAmplitude
	}
	id := wt.ID(len(tables))
	tables = append(tables, flat)
	p := NewProcessor(44100, tables)
	waveform := WaveformMorph{Start: id, End: id}
	if got := p.CalcIsochronicFactor(waveform, phaseOffset(0.25), t.IsochronicShort); got != 0.75 {
		ts.Fatalf("flat peak %f", got)
	}
	for _, offset := range []int{0, phaseOffset(0.5)} {
		if p.CalcIsochronicFactor(waveform, offset, t.IsochronicShort) != 0 {
			ts.Fatal("edge not silent")
		}
	}
	for _, offset := range []int{1, phaseOffset(0.5) - 1} {
		if p.CalcIsochronicFactor(waveform, offset, t.IsochronicShort) > 0.001 {
			ts.Fatal("custom edge discontinuity")
		}
	}
	if p.CalcModulationFactorForMorph(waveform, phaseOffset(0.75)) != 1 {
		ts.Fatal("shared modulation changed")
	}
}
