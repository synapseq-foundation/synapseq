// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"testing"

	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestSmoothedPanPositionClampsToBounds(ts *testing.T) {
	processor := NewProcessor(0, wt.Init())
	channel := &t.Channel{Effect: t.EffectState{PanPosition: 0.75, PanInitialized: true}}

	got := processor.smoothedPanPosition(channel, 10)
	if got != 1 {
		ts.Fatalf("unexpected clamped pan position: got %f, want 1", got)
	}
	if channel.Effect.PanPosition != 1 {
		ts.Fatalf("unexpected stored pan position: got %f, want 1", channel.Effect.PanPosition)
	}
}

func TestSmoothedModulationGainUsesLongerDeclickRamp(ts *testing.T) {
	processor := NewProcessor(44100, wt.Init())
	channel := &t.Channel{Effect: t.EffectState{ModulationGain: 1, ModulationInitialized: true}}

	got := processor.smoothedModulationGain(channel, 0.3)
	want := 1 - 1/(44100*modulationSlewTimeMs/1000.0)
	if got != want {
		ts.Fatalf("unexpected modulation gain step: got %f, want %f", got, want)
	}
}

func TestSmoothedPanPositionKeepsShortRamp(ts *testing.T) {
	processor := NewProcessor(44100, wt.Init())
	channel := &t.Channel{Effect: t.EffectState{PanPosition: -1, PanInitialized: true}}

	got := processor.smoothedPanPosition(channel, 1)
	want := -1 + 2/(44100*panSlewTimeMs/1000.0)
	if got != want {
		ts.Fatalf("unexpected pan step: got %f, want %f", got, want)
	}
}
