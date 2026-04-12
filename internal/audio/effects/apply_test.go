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
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestApplyEffectToStereoUsesSharedModulationPhase(ts *testing.T) {
	processor := newTestProcessor()
	channel := &t.Channel{Effect: t.EffectState{Offset: 0, Increment: 0}}
	waveform := WaveformMorph{Start: t.WaveformSine, End: t.WaveformSine, Alpha: 0}

	left, right := processor.ApplyEffectToStereo(channel, t.Effect{Type: t.EffectModulation, Intensity: 1}, waveform, 1000, 2000)
	if left != 300 || right != 600 {
		ts.Fatalf("unexpected modulation stereo output: got [%d %d], want [300 600]", left, right)
	}
	if !channel.Effect.ModulationInitialized || math.Abs(channel.Effect.ModulationGain-0.3) > 1e-9 {
		ts.Fatalf("unexpected modulation runtime state: initialized=%t gain=%f", channel.Effect.ModulationInitialized, channel.Effect.ModulationGain)
	}
}

func TestApplyEffectToMonoPanRoutesSignal(ts *testing.T) {
	processor := newTestProcessor()
	step := int(t.SineTableSize/4) * t.PhasePrecision
	channel := &t.Channel{Effect: t.EffectState{Offset: 0, Increment: step}}
	waveform := WaveformMorph{Start: t.WaveformSine, End: t.WaveformSine, Alpha: 0}

	left, right := processor.ApplyEffectToMono(channel, t.Effect{Type: t.EffectPan, Intensity: 1}, waveform, 1000)
	if left != 0 || right != 1000 {
		ts.Fatalf("unexpected pan mono routing: got [%d %d], want [0 1000]", left, right)
	}
	if channel.Effect.Offset != step {
		ts.Fatalf("unexpected pan phase advance: got %d, want %d", channel.Effect.Offset, step)
	}
}
