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
