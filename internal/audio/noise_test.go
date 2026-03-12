/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func genWhiteExpected(n int) []int {
	// Replicates the xorshift generator and the pink-state warm-up in NewNoiseGenerator.
	s := initialNoiseSeed
	for range noiseBands {
		s = nextExpectedRandom(s)
	}
	out := make([]int, n)
	for i := 0; i < n; i++ {
		s = nextExpectedRandom(s)
		out[i] = (int(s>>16) - maxCenteredRandom) * whiteNoiseScale
	}
	return out
}

func nextExpectedRandom(seed uint32) uint32 {
	seed ^= seed << 13
	seed ^= seed >> 17
	seed ^= seed << 5
	return seed
}

func meanAbsDelta(vals []int) float64 {
	if len(vals) < 2 {
		return 0
	}
	var sum float64
	for i := 1; i < len(vals); i++ {
		d := float64(vals[i] - vals[i-1])
		if d < 0 {
			d = -d
		}
		sum += d
	}
	return sum / float64(len(vals)-1)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func TestNoise_White_DeterministicFirstValues(ts *testing.T) {
	ng := NewNoiseGenerator()
	exp := genWhiteExpected(8)
	got := make([]int, len(exp))
	for i := range got {
		got[i] = ng.Generate(t.TrackWhiteNoise, 0)
	}
	for i := range exp {
		if got[i] != exp[i] {
			ts.Fatalf("white[%d]: got %d want %d", i, got[i], exp[i])
		}
	}
}

func TestNoise_EffectiveSmoothnessCurve(ts *testing.T) {
	tests := []struct {
		track t.TrackType
		input float64
		want  float64
	}{
		{track: t.TrackWhiteNoise, input: 0, want: 0},
		{track: t.TrackWhiteNoise, input: 60, want: 30},
		{track: t.TrackWhiteNoise, input: 100, want: 42},
		{track: t.TrackPinkNoise, input: 60, want: 24},
		{track: t.TrackPinkNoise, input: 80, want: 26.5},
		{track: t.TrackPinkNoise, input: 100, want: 34},
		{track: t.TrackBrownNoise, input: 60, want: 18},
		{track: t.TrackBrownNoise, input: 80, want: 20.5},
		{track: t.TrackBrownNoise, input: 100, want: 28},
	}

	for _, test := range tests {
		got := effectiveNoiseSmoothness(test.track, test.input)
		if math.Abs(got-test.want) > 0.0001 {
			ts.Fatalf("effectiveNoiseSmoothness(%v, %.2f): got %.4f want %.4f", test.track, test.input, got, test.want)
		}
	}
}

func TestNoise_BoundsAndSign(ts *testing.T) {
	ng := NewNoiseGenerator()
	N := 16384
	types := []t.TrackType{t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise}
	for _, typ := range types {
		pos, neg := 0, 0
		for i := 0; i < N; i++ {
			v := ng.Generate(typ, 50)
			if v > int(t.WaveTableAmplitude) || v < -int(t.WaveTableAmplitude) {
				ts.Fatalf("%v: out of bounds value %d at i=%d", typ, v, i)
			}
			if v > 0 {
				pos++
			} else if v < 0 {
				neg++
			}
		}
		// Ensure both signs appear reasonably
		if pos == 0 || neg == 0 {
			ts.Fatalf("%v: expected both positive and negative values (pos=%d neg=%d)", typ, pos, neg)
		}
	}
}

func TestNoise_SmoothnessComparison(ts *testing.T) {
	N := 16384

	// White
	ngW := NewNoiseGenerator()
	white := make([]int, N)
	for i := 0; i < N; i++ {
		white[i] = ngW.Generate(t.TrackWhiteNoise, 0)
	}

	// Pink
	ngP := NewNoiseGenerator()
	pink := make([]int, N)
	for i := 0; i < N; i++ {
		pink[i] = ngP.Generate(t.TrackPinkNoise, 0)
	}

	// Brown
	ngB := NewNoiseGenerator()
	brown := make([]int, N)
	for i := 0; i < N; i++ {
		brown[i] = ngB.Generate(t.TrackBrownNoise, 0)
	}

	dW := meanAbsDelta(white)
	dP := meanAbsDelta(pink)
	dB := meanAbsDelta(brown)

	// Brown should be smoothest (lowest delta), then Pink, then White highest
	if !(dB < dP && dP < dW) {
		ts.Fatalf("unexpected smooth order: dB=%.2f dP=%.2f dW=%.2f", dB, dP, dW)
	}

	// Also ensure deltas are non-trivial
	if math.IsNaN(dW) || math.IsNaN(dP) || math.IsNaN(dB) {
		ts.Fatalf("NaN deltas: dB=%.2f dP=%.2f dW=%.2f", dB, dP, dW)
	}
}

func TestNoise_SmoothnessReducesVariation(ts *testing.T) {
	N := 16384
	types := []t.TrackType{t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise}

	for _, typ := range types {
		ng0 := NewNoiseGenerator()
		ng50 := NewNoiseGenerator()
		ng100 := NewNoiseGenerator()

		seq0 := make([]int, N)
		seq50 := make([]int, N)
		seq100 := make([]int, N)

		for i := 0; i < N; i++ {
			seq0[i] = ng0.Generate(typ, 0)
			seq50[i] = ng50.Generate(typ, 50)
			seq100[i] = ng100.Generate(typ, 100)
		}

		d0 := meanAbsDelta(seq0)
		d50 := meanAbsDelta(seq50)
		d100 := meanAbsDelta(seq100)

		if !(d100 < d50 && d50 < d0) {
			ts.Fatalf("%v: expected variation to decrease with smooth, got d0=%.2f d50=%.2f d100=%.2f", typ, d0, d50, d100)
		}
	}
}

func TestNoise_DeterminismAcrossInstances(ts *testing.T) {
	N := 2048
	types := []t.TrackType{t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise}
	smoothnessValues := []float64{0, 50, 100}
	for _, typ := range types {
		for _, smooth := range smoothnessValues {
			ng1 := NewNoiseGenerator()
			ng2 := NewNoiseGenerator()
			for i := 0; i < N; i++ {
				v1 := ng1.Generate(typ, smooth)
				v2 := ng2.Generate(typ, smooth)
				if v1 != v2 {
					ts.Fatalf("%v smooth %.2f: non-deterministic at i=%d: %d vs %d", typ, smooth, i, v1, v2)
				}
			}
		}
	}
}

func TestNoise_SmoothnessRampPreservesFilterState(ts *testing.T) {
	const warmupSamples = 512

	ngPreserved := NewNoiseGenerator()
	ngReset := NewNoiseGenerator()

	prevPreserved := 0
	prevReset := 0
	for i := 0; i < warmupSamples; i++ {
		prevPreserved = ngPreserved.Generate(t.TrackPinkNoise, 100)
		prevReset = ngReset.Generate(t.TrackPinkNoise, 100)
	}

	ngReset.smoothState[t.TrackPinkNoise] = noiseSmoothnessState{}

	outPreserved := ngPreserved.Generate(t.TrackPinkNoise, 50)
	outReset := ngReset.Generate(t.TrackPinkNoise, 50)

	deltaPreserved := absInt(outPreserved - prevPreserved)
	deltaReset := absInt(outReset - prevReset)

	if len(ngPreserved.smoothState) != 1 {
		ts.Fatalf("expected a single smooth state entry for ramped pink noise, got %d", len(ngPreserved.smoothState))
	}
	if deltaPreserved >= deltaReset {
		ts.Fatalf("expected preserved filter state to reduce transition discontinuity, got preserved=%d reset=%d", deltaPreserved, deltaReset)
	}
}
