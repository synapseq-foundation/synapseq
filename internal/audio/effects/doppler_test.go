// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"
	"testing"

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
	factor := processor.calcDopplerFactor(step, 1)
	if left != int(math.Round(100*factor)) || right != int(math.Round(200*factor)) {
		ts.Fatalf("unexpected doppler pair output: got [%d %d], want [%d %d]", left, right, int(math.Round(100*factor)), int(math.Round(200*factor)))
	}
}
