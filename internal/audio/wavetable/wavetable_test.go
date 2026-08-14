// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package wavetable

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestInit_LengthAndBounds(ts *testing.T) {
	wts := Init()
	for i := range wts {
		tab := wts[i]
		if len(tab) != t.SineTableSize {
			ts.Fatalf("waveform %d: expected len=%d, got %d", i, t.SineTableSize, len(tab))
		}
		for j := 0; j < len(tab); j++ {
			if tab[j] > int(t.WaveTableAmplitude) || tab[j] < -int(t.WaveTableAmplitude) {
				ts.Fatalf("waveform %d: out of bounds at j=%d: %d", i, j, tab[j])
			}
		}
	}
}

func TestInit_Sine(ts *testing.T) {
	wts := Init()
	tab := wts[int(SineID)]
	amp := int(t.WaveTableAmplitude)

	j0 := 0
	j1 := t.SineTableSize / 4
	j2 := t.SineTableSize / 2
	j3 := 3 * t.SineTableSize / 4

	if tab[j0] != 0 {
		ts.Fatalf("sine j=0: want 0, got %d", tab[j0])
	}
	if tab[j1] != amp {
		ts.Fatalf("sine j=pi/2: want %d, got %d", amp, tab[j1])
	}
	if tab[j2] != 0 {
		ts.Fatalf("sine j=pi: want 0, got %d", tab[j2])
	}
	if tab[j3] != -amp {
		ts.Fatalf("sine j=3pi/2: want %d, got %d", -amp, tab[j3])
	}
}

func TestInit_Square(ts *testing.T) {
	wts := Init()
	tab := wts[int(SquareID)]
	amp := int(t.WaveTableAmplitude)

	j0 := 0
	j1 := t.SineTableSize / 4
	j2 := t.SineTableSize / 2
	j3 := 3 * t.SineTableSize / 4

	if tab[j0] != -amp {
		ts.Fatalf("square j=0: want %d, got %d", -amp, tab[j0])
	}
	if tab[j1] != amp {
		ts.Fatalf("square j=pi/2: want %d, got %d", amp, tab[j1])
	}
	if !(tab[j2-1] == amp && tab[j2+1] == -amp) {
		ts.Fatalf("square around pi: expected transition +amp -> -amp, got %d and %d", tab[j2-1], tab[j2+1])
	}
	if tab[j3] != -amp {
		ts.Fatalf("square j=3pi/2: want %d, got %d", -amp, tab[j3])
	}

	for i, v := range tab {
		if v != amp && v != -amp {
			ts.Fatalf("square value must be +/-amp at i=%d: got %d", i, v)
		}
	}
}

func TestInit_Triangle(ts *testing.T) {
	wts := Init()
	tab := wts[int(TriangleID)]
	amp := int(t.WaveTableAmplitude)

	j0 := 0
	j1 := t.SineTableSize / 4
	j2 := t.SineTableSize / 2
	j3 := 3 * t.SineTableSize / 4

	if tab[j0] != -amp {
		ts.Fatalf("triangle j=0: want %d, got %d", -amp, tab[j0])
	}
	if tab[j1] != 0 {
		ts.Fatalf("triangle j=pi/2: want 0, got %d", tab[j1])
	}
	if tab[j2] != amp {
		ts.Fatalf("triangle j=pi: want %d, got %d", amp, tab[j2])
	}
	if tab[j3] != 0 {
		ts.Fatalf("triangle j=3pi/2: want 0, got %d", tab[j3])
	}
}

func TestInit_Sawtooth(ts *testing.T) {
	wts := Init()
	tab := wts[int(SawtoothID)]
	amp := int(t.WaveTableAmplitude)

	exp := func(j int) int {
		phase := float64(j) * 2.0 * math.Pi / float64(t.SineTableSize)
		val := 2.0 * (phase/(2.0*math.Pi) - math.Floor(phase/(2.0*math.Pi)+0.5))
		return int(float64(amp) * val)
	}

	j0 := 0
	j1 := t.SineTableSize / 4
	j2 := t.SineTableSize / 2
	j3 := 3 * t.SineTableSize / 4

	if tab[j0] != exp(j0) {
		ts.Fatalf("sawtooth j=0: want %d, got %d", exp(j0), tab[j0])
	}
	if tab[j1] != exp(j1) {
		ts.Fatalf("sawtooth j=pi/2: want %d, got %d", exp(j1), tab[j1])
	}
	if tab[j2] != exp(j2) {
		ts.Fatalf("sawtooth j=pi: want %d, got %d", exp(j2), tab[j2])
	}
	if tab[j3] != exp(j3) {
		ts.Fatalf("sawtooth j=3pi/2: want %d, got %d", exp(j3), tab[j3])
	}
}

func TestCompileCustomWaveformNormalizesCycleAndWraps(ts *testing.T) {
	registry, err := Compile([]t.WaveformDefinition{{
		Name:   "pulse",
		Points: []float64{-1, 1},
	}})
	if err != nil {
		ts.Fatalf("Compile error: %v", err)
	}

	id, ok := registry.Lookup("pulse")
	if !ok {
		ts.Fatal("expected custom waveform ID")
	}
	if id != 4 {
		ts.Fatalf("expected first custom waveform ID 4, got %d", id)
	}

	table := registry.Tables[id]
	amp := int(t.WaveTableAmplitude)
	if table[0] != -amp || table[t.SineTableSize/2] != amp {
		ts.Fatalf("unexpected normalized extrema: got %d and %d", table[0], table[t.SineTableSize/2])
	}
	if table[t.SineTableSize/4] != 0 || table[3*t.SineTableSize/4] != 0 {
		ts.Fatalf("expected linear midpoint crossings, got %d and %d", table[t.SineTableSize/4], table[3*t.SineTableSize/4])
	}
	if table[len(table)-1] >= table[len(table)-2] || table[len(table)-1] <= -amp {
		ts.Fatalf("expected final segment to approach the first point continuously, got %d", table[len(table)-1])
	}
}

func TestCompilePreservesBuiltInTables(ts *testing.T) {
	want := Init()
	registry, err := Compile([]t.WaveformDefinition{{Name: "custom", Points: []float64{-1, 1}}})
	if err != nil {
		ts.Fatalf("Compile error: %v", err)
	}
	for waveform := range want {
		for index := range want[waveform] {
			if registry.Tables[waveform][index] != want[waveform][index] {
				ts.Fatalf("built-in waveform %d changed at sample %d", waveform, index)
			}
		}
	}
}
