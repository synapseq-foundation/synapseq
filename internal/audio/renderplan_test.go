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

package audio

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestCompileRenderPlanBuildsTemporalWindows(ts *testing.T) {
	periods := []t.Period{{Time: 0}, {Time: 500}, {Time: 1250}, {Time: 2000}}

	plan := compileRenderPlan(periods, 44100)

	if plan.totalFrames != 88200 {
		ts.Fatalf("unexpected total frame count: got %d, want %d", plan.totalFrames, 88200)
	}

	if len(plan.windows) != len(periods) {
		ts.Fatalf("unexpected window count: got %d, want %d", len(plan.windows), len(periods))
	}

	expected := []renderWindow{
		{PeriodIndex: 0, StartMs: 0, EndMs: 500},
		{PeriodIndex: 1, StartMs: 500, EndMs: 1250},
		{PeriodIndex: 2, StartMs: 1250, EndMs: 2000},
		{PeriodIndex: 3, StartMs: 2000, EndMs: 2000},
	}

	for index, want := range expected {
		if got := plan.windows[index]; got != want {
			ts.Fatalf("unexpected window %d: got %+v, want %+v", index, got, want)
		}
	}
}

func TestRenderPlanPeriodIndexAtAdvancesAcrossWindows(ts *testing.T) {
	plan := compileRenderPlan([]t.Period{{Time: 0}, {Time: 500}, {Time: 1250}, {Time: 2000}}, 44100)

	tests := []struct {
		name      string
		timeMs    int
		startFrom int
		want      int
	}{
		{name: "initial window", timeMs: 0, startFrom: 0, want: 0},
		{name: "before second window", timeMs: 499, startFrom: 0, want: 0},
		{name: "second window", timeMs: 500, startFrom: 0, want: 1},
		{name: "third window", timeMs: 1250, startFrom: 1, want: 2},
		{name: "near end", timeMs: 1999, startFrom: 2, want: 2},
	}

	for _, tc := range tests {
		if got := plan.periodIndexAt(tc.timeMs, tc.startFrom); got != tc.want {
			ts.Fatalf("%s: unexpected period index: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestRenderPlanCueResolvesInterpolatedTrackState(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(10),
		Carrier:   200,
		Resonance: 8,
		Waveform:  t.WaveformSine,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     2,
			Intensity: t.IntensityPercentToRaw(20),
		},
	}
	p0.TrackEnd[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(30),
		Carrier:   300,
		Resonance: 12,
		Waveform:  t.WaveformTriangle,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     4,
			Intensity: t.IntensityPercentToRaw(60),
		},
	}

	plan := compileRenderPlan([]t.Period{p0, p1}, 44100)
	cue := plan.cue(0, 500)
	channel := cue.Channels[0]

	assertAlmostEqual(ts, channel.Track.Carrier, 250, 0.0001)
	assertAlmostEqual(ts, channel.Track.Resonance, 10, 0.0001)
	assertAlmostEqual(ts, channel.Track.Effect.Value, 3, 0.0001)
	assertAlmostEqual(ts, float64(channel.Track.Effect.Intensity), float64(t.IntensityPercentToRaw(40)), 0.0001)
	assertAlmostEqual(ts, channel.WaveformAlpha, 0.5, 0.0001)
	if channel.Amplitude[0] != int(channel.Track.Amplitude) || channel.Amplitude[1] != int(channel.Track.Amplitude) {
		ts.Fatalf("unexpected amplitude state: got %v", channel.Amplitude)
	}
	if channel.Increment[0] != frequencyToIncrement(44100, 255) || channel.Increment[1] != frequencyToIncrement(44100, 245) {
		ts.Fatalf("unexpected increment state: got %v", channel.Increment)
	}
	if channel.EffectStep != frequencyToIncrement(44100, 3) {
		ts.Fatalf("unexpected effect step: got %d", channel.EffectStep)
	}
	if channel.WaveformStart != t.WaveformSine || channel.WaveformEnd != t.WaveformTriangle {
		ts.Fatalf("unexpected waveform state: got %v -> %v", channel.WaveformStart, channel.WaveformEnd)
	}
}

func assertAlmostEqual(ts *testing.T, got, want, tolerance float64) {
	ts.Helper()

	if math.Abs(got-want) > tolerance {
		ts.Fatalf("unexpected value: got %.6f want %.6f", got, want)
	}
}
