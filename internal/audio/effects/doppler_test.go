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

func TestApplyDopplerIgnoresOtherEffects(ts *testing.T) {
	processor := newTestProcessor()
	channel := &t.Channel{Effect: t.EffectState{Offset: 123, Increment: t.PhasePrecision}}

	got := processor.ApplyDoppler(channel, t.Effect{Type: t.EffectPan, Intensity: 1}, 456)
	if got != 456 {
		ts.Fatalf("unexpected increment for non-doppler effect: got %d, want 456", got)
	}
	if channel.Effect.Offset != 123 {
		ts.Fatalf("unexpected phase advance for non-doppler effect: got %d, want 123", channel.Effect.Offset)
	}
}

func TestApplyDopplerPairAdvancesPhaseAndScalesBothChannels(ts *testing.T) {
	processor := newTestProcessor()
	step := int(t.SineTableSize/4) * t.PhasePrecision
	channel := &t.Channel{Effect: t.EffectState{Increment: step}}

	left, right := processor.ApplyDopplerPair(channel, t.Effect{Type: t.EffectDoppler, Intensity: 1}, 100, 200)
	if channel.Effect.Offset != step {
		ts.Fatalf("unexpected doppler phase advance: got %d, want %d", channel.Effect.Offset, step)
	}
	factor := processor.calcDopplerFactor(WaveformMorph{}, step, 1)
	if left != int(math.Round(100*factor)) || right != int(math.Round(200*factor)) {
		ts.Fatalf("unexpected doppler pair output: got [%d %d], want [%d %d]", left, right, int(math.Round(100*factor)), int(math.Round(200*factor)))
	}
}

func TestApplyDopplerUsesTrackWaveform(ts *testing.T) {
	processor := newTestProcessor()
	channel := &t.Channel{
		Track:  t.Track{Waveform: t.WaveformSquare},
		Effect: t.EffectState{Increment: t.PhasePrecision},
	}

	got := processor.ApplyDoppler(channel, t.Effect{Type: t.EffectDoppler, Intensity: 1}, 100)
	if got != 105 {
		ts.Fatalf("unexpected square-wave doppler increment: got %d, want 105", got)
	}

	sineFactor := processor.calcDopplerFactor(
		WaveformMorph{Start: wt.SineID, End: wt.SineID},
		t.PhasePrecision,
		1,
	)
	if got == int(math.Round(100*sineFactor)) {
		ts.Fatalf("doppler ignored the track waveform: got %d", got)
	}
}

func TestApplyDopplerMorphsTrackWaveforms(ts *testing.T) {
	processor := newTestProcessor()
	channel := &t.Channel{
		WaveformStart: int(wt.SineID),
		WaveformEnd:   int(wt.SquareID),
		WaveformAlpha: 0.5,
		Effect:        t.EffectState{Increment: t.PhasePrecision},
	}

	got := processor.ApplyDoppler(channel, t.Effect{Type: t.EffectDoppler, Intensity: 1}, 100)
	sine := float64(processor.waveTables[int(wt.SineID)][1])
	square := float64(processor.waveTables[int(wt.SquareID)][1])
	waveformValue := sine + (square-sine)*0.5
	want := int(math.Round(100 * (1 + 0.05*waveformValue/float64(t.WaveTableAmplitude))))
	if got != want {
		ts.Fatalf("unexpected morphed doppler increment: got %d, want %d", got, want)
	}
}
