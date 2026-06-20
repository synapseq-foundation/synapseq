// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preset

import (
	"fmt"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// FindPreset searches for a preset by name in a slice of presets.
func FindPreset(name string, presets []t.Preset) *t.Preset {
	for i := range presets {
		if presets[i].String() == name {
			return &presets[i]
		}
	}
	return nil
}

// AllocateTrack allocates a free track in the preset.
func AllocateTrack(preset *t.Preset) (int, error) {
	for index, track := range preset.Track {
		if track.Type == t.TrackOff {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no available tracks for preset %q", preset.String())
}

// IsPresetEmpty checks if all tracks in the preset are off.
func IsPresetEmpty(preset *t.Preset) bool {
	for _, track := range preset.Track {
		if track.Type != t.TrackOff {
			return false
		}
	}
	return true
}
