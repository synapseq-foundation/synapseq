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

// EffectType represents the type of effect applied to a background track
type EffectType int

const (
	// Effect is off
	EffectOff EffectType = iota
	// Effect is spin
	EffectSpin
	// Effect is pulse
	EffectPulse
	// Effect is doppler
	EffectDoppler
)

// String returns the string representation of the EffectType
func (et EffectType) String() string {
	switch et {
	case EffectOff:
		return KeywordOff
	case EffectSpin:
		return KeywordSpin
	case EffectPulse:
		return KeywordPulse
	case EffectDoppler:
		return KeywordDoppler
	default:
		return "unknown"
	}
}

// Effect represents a effect configuration
type Effect struct {
	// Effect type
	Type EffectType
	// Effect value
	Value float64
	// Intensity (0-1.0 for 0-100%)
	Intensity IntensityType
}

// EffectState stores runtime (per-channel) effect state without per-effect structs/maps.
type EffectState struct {
	// LFO state for the currently active effect (phase + step)
	Increment int
	Offset    int
}
