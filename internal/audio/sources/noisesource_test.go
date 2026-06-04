// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type stubNoiseGenerator struct {
	gotKind   t.TrackType
	gotSmooth float64
	value     int
}

func (stub *stubNoiseGenerator) Generate(trackType t.TrackType, smooth float64) int {
	stub.gotKind = trackType
	stub.gotSmooth = smooth
	return stub.value
}

func TestNoiseSampleUsesCompiledSignalState(ts *testing.T) {
	generator := &stubNoiseGenerator{value: 7}
	source := NewNoise(Signal{Kind: t.TrackPinkNoise, NoiseSmooth: 50, Amplitude: [2]int{3, 0}})

	got := source.Sample(generator)
	want := 21
	if got != want {
		ts.Fatalf("unexpected noise sample: got %d, want %d", got, want)
	}
	if generator.gotKind != t.TrackPinkNoise || generator.gotSmooth != 50 {
		ts.Fatalf("unexpected generator inputs: got kind=%v smooth=%.2f", generator.gotKind, generator.gotSmooth)
	}
}