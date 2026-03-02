/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package shared

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestCountActiveChannels(ts *testing.T) {
	tests := []struct {
		name     string
		chs      []t.Channel
		expected int
	}{
		{
			name:     "empty slice -> at least 1",
			chs:      []t.Channel{},
			expected: 1,
		},
		{
			name:     "all off -> 1",
			chs:      make([]t.Channel, 5),
			expected: 1,
		},
		{
			name: "single active at 0 -> 1",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 4)
				chs[0].Track.Type = t.TrackBinauralBeat
				return chs
			}(),
			expected: 1,
		},
		{
			name: "last active at end -> len",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 4)
				chs[3].Track.Type = t.TrackPinkNoise
				return chs
			}(),
			expected: 4,
		},
		{
			name: "last active in the middle -> index+1",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 5)
				chs[2].Track.Type = t.TrackBrownNoise
				return chs
			}(),
			expected: 3,
		},
		{
			name: "multiple actives -> last index+1",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 8)
				chs[1].Track.Type = t.TrackBinauralBeat
				chs[6].Track.Type = t.TrackAmbiance
				return chs
			}(),
			expected: 7,
		},
		{
			name: "all active -> len",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 7)
				for i := range chs {
					chs[i].Track.Type = t.TrackPinkNoise
				}
				return chs
			}(),
			expected: 7,
		},
		{
			name: "last off but previous active",
			chs: func() []t.Channel {
				chs := make([]t.Channel, 6)
				chs[4].Track.Type = t.TrackAmbiance
				return chs
			}(),
			expected: 5,
		},
	}

	for _, tt := range tests {
		got := CountActiveChannels(tt.chs)
		if got != tt.expected {
			ts.Errorf("%s: expected %d, got %d", tt.name, tt.expected, got)
		}
	}
}
