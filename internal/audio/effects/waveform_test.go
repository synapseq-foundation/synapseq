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

func TestWaveformValueForMorphInterpolatesBetweenTables(ts *testing.T) {
	processor := newTestProcessor()
	offset := t.PhasePrecision
	waveform := WaveformMorph{Start: wt.SineID, End: wt.SquareID, Alpha: 0.25}

	got := processor.WaveformValueForMorph(waveform, offset)
	start := float64(processor.waveTables[int(wt.SineID)][1])
	end := float64(processor.waveTables[int(wt.SquareID)][1])
	want := start + (end-start)*0.25
	if got != want {
		ts.Fatalf("unexpected morphed waveform value: got %f, want %f", got, want)
	}

	sample := processor.WaveformSampleForMorph(waveform, offset)
	if sample != int(math.Round(want)) {
		ts.Fatalf("unexpected morphed waveform sample: got %d, want %d", sample, int(math.Round(want)))
	}
}

func TestCustomWaveformDrivesSamplingAndEffects(ts *testing.T) {
	registry, err := wt.Compile([]t.WaveformDefinition{{Name: "pulse", Points: []float64{-1, 1}}})
	if err != nil {
		ts.Fatalf("Compile error: %v", err)
	}
	id, _ := registry.Lookup("pulse")
	processor := NewProcessor(44100, registry.Tables)
	waveform := WaveformMorph{Start: id, End: id}
	highPhase := int(t.SineTableSize/2) * t.PhasePrecision

	if got := processor.WaveformSampleForMorph(waveform, highPhase); got != int(t.WaveTableAmplitude) {
		ts.Fatalf("unexpected custom waveform sample: got %d", got)
	}
	if got := processor.CalcModulationFactorForMorph(waveform, highPhase); got != 1 {
		ts.Fatalf("unexpected custom modulation factor: got %f", got)
	}
	wantBlend := int(math.Round(float64(t.WaveTableAmplitude) * 0.5))
	for _, blend := range []WaveformMorph{
		{Start: id, End: wt.SineID, Alpha: 0.5},
		{Start: wt.SineID, End: id, Alpha: 0.5},
	} {
		if got := processor.WaveformSampleForMorph(blend, highPhase); got != wantBlend {
			ts.Fatalf("unexpected custom/built-in morphed sample: got %d, want %d", got, wantBlend)
		}
	}

	channel := &t.Channel{Effect: t.EffectState{Offset: highPhase}}
	left, right := processor.ApplyPanForMorph(
		channel,
		t.Effect{Type: t.EffectPan, Intensity: 1},
		waveform,
		1000,
		1000,
	)
	if left != 0 || right != 1000 {
		ts.Fatalf("unexpected custom waveform pan output: got [%d %d]", left, right)
	}

	channel.WaveformStart = int(id)
	channel.WaveformEnd = int(id)
	channel.Effect = t.EffectState{Increment: highPhase}
	if got := processor.ApplyDoppler(channel, t.Effect{Type: t.EffectDoppler, Intensity: 1}, 100); got != 105 {
		ts.Fatalf("unexpected custom waveform doppler increment: got %d, want 105", got)
	}
}
