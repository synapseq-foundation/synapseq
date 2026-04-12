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

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestWaveformMorphFromChannelFallsBackToTrackWaveform(ts *testing.T) {
	channel := &t.Channel{Track: t.Track{Waveform: t.WaveformTriangle}}

	morph := WaveformMorphFromChannel(channel)
	if morph.Start != t.WaveformTriangle || morph.End != t.WaveformTriangle || morph.Alpha != 0 {
		ts.Fatalf("unexpected fallback morph: %+v", morph)
	}
}
