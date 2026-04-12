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

package timeline

import (
	"math"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestMaxPeriodSteps(ts *testing.T) {
	tests := []struct {
		durationMs int
		want       int
	}{
		{0, 0},
		{4_999, 0},
		{10_000, 0},
		{15_000, 1},
		{30_000, 2},
		{60_000, 5},
		{180_000, 12},
		{600_000, 12},
	}

	for _, test := range tests {
		if got := MaxPeriodSteps(test.durationMs); got != test.want {
			ts.Fatalf("MaxPeriodSteps(%d) = %d, want %d", test.durationMs, got, test.want)
		}
	}
}

func TestAdjustPeriods_NormalCopy(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	last.TrackEnd[0] = t.Track{Type: t.TrackBinauralBeat, Carrier: 310, Resonance: 11, Amplitude: t.AmplitudePercentToRaw(12), Waveform: t.WaveformSine}
	next.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Carrier: 350, Resonance: 12, Amplitude: t.AmplitudePercentToRaw(15), Waveform: t.WaveformSine}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_FadeInFromSilence(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{Type: t.TrackSilence, Waveform: t.WaveformSquare}
	last.TrackEnd[0] = t.Track{Type: t.TrackSilence, Amplitude: 0}
	next.TrackStart[0] = t.Track{Type: t.TrackMonauralBeat, Carrier: 200, Resonance: 6, Amplitude: t.AmplitudePercentToRaw(25), Waveform: t.WaveformTriangle}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	got := last.TrackStart[0]
	if got.Type != t.TrackMonauralBeat || got.Amplitude != 0 || got.Carrier != 200 || got.Resonance != 6 || got.Waveform != t.WaveformTriangle {
		ts.Fatalf("fade-in not applied as expected: %+v", got)
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch after fade-in: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_FadeOutToSilence(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{Type: t.TrackAmbiance, Carrier: 200, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(50), Waveform: t.WaveformSquare, Effect: t.Effect{Type: t.EffectPan, Intensity: t.IntensityPercentToRaw(70)}}
	last.TrackEnd[0] = last.TrackStart[0]
	next.TrackStart[0] = t.Track{Type: t.TrackSilence, Amplitude: 0}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	if next.TrackStart[0].Type != t.TrackSilence || next.TrackStart[0].Carrier != 200 || next.TrackStart[0].Resonance != 5 || next.TrackStart[0].Effect.Intensity != t.IntensityPercentToRaw(70) {
		ts.Fatalf("fade-out not applied as expected: %+v", next.TrackStart[0])
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch after fade-out: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_AllowsWaveformChangeWhileOn(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	last.TrackEnd[0] = t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	next.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(12), Waveform: t.WaveformTriangle}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error when changing waveform: %v", err)
	}
	if last.TrackEnd[0].Waveform != t.WaveformTriangle {
		ts.Fatalf("waveform was not carried forward: got %v", last.TrackEnd[0].Waveform)
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch after waveform change: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_Errors(ts *testing.T) {
	makePer := func(trackStart, trackEnd, nextStart t.Track) (t.Period, t.Period) {
		var last, next t.Period
		last.TrackStart[0] = trackStart
		last.TrackEnd[0] = trackEnd
		next.TrackStart[0] = nextStart
		return last, next
	}

	tests := []struct {
		name string
		tr0  t.Track
		tr1  t.Track
		tr2  t.Track
	}{
		{"turn off directly", t.Track{}, t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}, t.Track{Type: t.TrackOff}},
		{"turn on directly", t.Track{}, t.Track{Type: t.TrackOff}, t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}},
		{"change type while on", t.Track{}, t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}, t.Track{Type: t.TrackMonauralBeat, Amplitude: t.AmplitudePercentToRaw(12), Waveform: t.WaveformSine}},
		{"change effect type while on (ambiance)", t.Track{}, t.Track{Type: t.TrackAmbiance, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPan}}, t.Track{Type: t.TrackAmbiance, Amplitude: t.AmplitudePercentToRaw(25), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectModulation}}},
	}

	for _, test := range tests {
		last, next := makePer(test.tr0, test.tr1, test.tr2)
		if err := AdjustPeriods(&last, &next); err == nil {
			ts.Fatalf("%s: expected error, got nil", test.name)
		}
	}
}

func TestAdjustPeriods_ValidatesStepsAgainstDuration(ts *testing.T) {
	var last, next t.Period
	last.Time = 120_000
	next.Time = 150_000
	last.Steps = 3
	last.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	last.TrackEnd[0] = last.TrackStart[0]
	next.TrackStart[0] = last.TrackStart[0]

	err := AdjustPeriods(&last, &next)
	if err == nil {
		ts.Fatal("expected steps validation error, got nil")
	}
	if !strings.Contains(err.Error(), "uses 3 steps") {
		ts.Fatalf("expected steps error message, got: %v", err)
	}
}

func TestAdjustPeriods_AllowsStepsWithinDurationLimit(ts *testing.T) {
	var last, next t.Period
	last.Time = 120_000
	next.Time = 180_000
	last.Steps = 5
	last.TrackStart[0] = t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	last.TrackEnd[0] = last.TrackStart[0]
	next.TrackStart[0] = last.TrackStart[0]

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error for valid steps: %v", err)
	}
}

func TestStepAlpha(ts *testing.T) {
	assertAlmostEqual(ts, StepAlpha(0.5, t.TransitionSteady, 0), 0.5, 0.000001)
	assertAlmostEqual(ts, StepAlpha(1.0/3.0, t.TransitionSteady, 1), 1.0, 0.000001)
	assertAlmostEqual(ts, StepAlpha(2.0/3.0, t.TransitionSteady, 1), 0.0, 0.000001)
	assertAlmostEqual(ts, StepAlpha(1.0, t.TransitionSteady, 1), 1.0, 0.000001)
}

func assertAlmostEqual(ts *testing.T, got, want, tolerance float64) {
	ts.Helper()

	if math.Abs(got-want) > tolerance {
		ts.Fatalf("unexpected value: got %.6f want %.6f", got, want)
	}
}