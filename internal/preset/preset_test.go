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

package preset

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestFindPreset(ts *testing.T) {
	var presets []t.Preset
	presets = append(presets, *t.NewBuiltinSilencePreset())

	alpha, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset alpha: %v", err)
	}
	beta, err := t.NewPreset("beta", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset beta: %v", err)
	}
	presets = append(presets, *alpha, *beta)

	got := FindPreset("alpha", presets)
	if got == nil || got != &presets[1] || got.String() != "alpha" {
		ts.Fatalf("FindPreset('alpha') failed, got=%v, want address of presets[1]", got)
	}
	got = FindPreset("silence", presets)
	if got == nil || got != &presets[0] || got.String() != "silence" {
		ts.Fatalf("FindPreset('silence') failed, got=%v, want address of presets[0]", got)
	}
	got = FindPreset("missing", presets)
	if got != nil {
		ts.Fatalf("FindPreset('missing') should be nil, got=%v", got)
	}
}

func TestAllocateTrack(ts *testing.T) {
	p, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	for i := range t.NumberOfChannels {
		idx, err := AllocateTrack(p)
		if err != nil {
			ts.Fatalf("AllocateTrack failed at i=%d: %v", i, err)
		}
		if idx != i {
			ts.Fatalf("AllocateTrack index mismatch: got %d, want %d", idx, i)
		}
		p.Track[idx].Type = t.TrackBinauralBeat
	}

	if _, err := AllocateTrack(p); err == nil {
		ts.Fatalf("AllocateTrack should fail when no free tracks")
	}
}

func TestIsPresetEmpty(ts *testing.T) {
	p, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if !IsPresetEmpty(p) {
		ts.Fatalf("new preset should be empty")
	}

	p.Track[0].Type = t.TrackWhiteNoise
	if IsPresetEmpty(p) {
		ts.Fatalf("preset with one active track should not be empty")
	}

	silencePreset := t.NewBuiltinSilencePreset()
	if IsPresetEmpty(silencePreset) {
		ts.Fatalf("silence preset should not be considered empty")
	}
}