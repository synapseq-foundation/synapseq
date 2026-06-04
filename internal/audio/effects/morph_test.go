// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestWaveformMorphFromChannelFallsBackToTrackWaveform(ts *testing.T) {
	channel := &t.Channel{Track: t.Track{Waveform: t.WaveformTriangle}}

	morph := WaveformMorphFromChannel(channel)
	if morph.Start != t.WaveformTriangle || morph.End != t.WaveformTriangle || morph.Alpha != 0 {
		ts.Fatalf("unexpected fallback morph: %+v", morph)
	}
}
