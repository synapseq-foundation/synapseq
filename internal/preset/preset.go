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
