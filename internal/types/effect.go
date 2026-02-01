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

// EffectConfiguration represents the configuration for an effect
type EffectConfiguration interface {
	// effectType is a marker method to indicate this type is an EffectConfiguration
	effectType() EffectType
}

// EffectType represents the type of effect applied to a background track
type EffectType int

const (
	// Effect is off
	EffectOff EffectType = iota
	// Effect is spin
	EffectSpin
	// Effect is pulse
	EffectPulse
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
	default:
		return "unknown"
	}
}

// Effect represents a effect configuration
type Effect struct {
	// Effect type
	Type EffectType
	// Effect configuration
	Configuration EffectConfiguration
	// Intensity (0-1.0 for 0-100%)
	Intensity IntensityType
}

// EffectSpinConfiguration represents the configuration for a spin effect
type EffectSpinConfiguration struct {
	Width float64
	Rate  float64
}

// EffectPulseConfiguration represents the configuration for a pulse effect
type EffectPulseConfiguration struct {
	Pulse float64
}

// Marker methods
func (EffectSpinConfiguration) effectType() EffectType  { return EffectSpin }
func (EffectPulseConfiguration) effectType() EffectType { return EffectPulse }
