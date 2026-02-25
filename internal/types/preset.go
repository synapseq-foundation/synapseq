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

package types

import (
	"fmt"
	"strings"
)

const (
	MaxPresets     = 32        // Maximum number of presets
	builtinSilence = "silence" // Represents silence built-in preset
)

// Preset represents a named preset
type Preset struct {
	name       string                  // Name of preset
	IsTemplate bool                    // Whether this preset is a template
	Track      [NumberOfChannels]Track // Track-set for it
	From       *Preset                 // Optional preset to copy from (template)
}

// NewPreset creates a new preset with the given name
func NewPreset(name string, template bool, from *Preset) (*Preset, error) {
	n := strings.ToLower(name)
	if n == builtinSilence {
		return nil, fmt.Errorf("preset name %q is reserved", builtinSilence)
	}

	preset := &Preset{name: n, IsTemplate: template}

	if from != nil {
		preset.From = from
		preset.Track = from.Track
		return preset, nil
	}

	for i := range NumberOfChannels {
		preset.Track[i].Type = TrackOff
		preset.Track[i].Carrier = 0.0
		preset.Track[i].Resonance = 0.0
		preset.Track[i].Amplitude = 0.0
		preset.Track[i].Waveform = WaveformSine
		preset.Track[i].Effect = Effect{Type: EffectOff, Intensity: 0.0}
	}
	return preset, nil
}

// NewBuiltinSilencePreset creates a new silence preset
func NewBuiltinSilencePreset() *Preset {
	preset := &Preset{name: builtinSilence}
	for i := range NumberOfChannels {
		preset.Track[i].Type = TrackSilence
		preset.Track[i].Carrier = 0.0
		preset.Track[i].Resonance = 0.0
		preset.Track[i].Amplitude = 0.0
		preset.Track[i].Waveform = WaveformSine
		preset.Track[i].Effect = Effect{Type: EffectOff, Intensity: 0.0}
	}
	return preset
}

// String returns the name of the preset
func (p *Preset) String() string {
	return p.name
}
