// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "strings"

// WaveformName identifies a built-in or document-defined waveform.
type WaveformName string

// Built-in waveform names.
const (
	WaveformSine     WaveformName = KeywordSine
	WaveformSquare   WaveformName = KeywordSquare
	WaveformTriangle WaveformName = KeywordTriangle
	WaveformSawtooth WaveformName = KeywordSawtooth
)

const (
	MinWaveformPoints = 2
	MaxWaveformPoints = SineTableSize
)

// WaveformDefinition contains one normalized periodic waveform cycle.
type WaveformDefinition struct {
	Name   WaveformName
	Points []float64
}

// Effective returns sine for an omitted waveform, preserving the historical default.
func (name WaveformName) Effective() WaveformName {
	if name == "" {
		return WaveformSine
	}
	return name
}

// String returns the waveform name.
func (name WaveformName) String() string {
	return string(name.Effective())
}

// IsBuiltinWaveformName reports whether name conflicts with a built-in waveform.
func IsBuiltinWaveformName(name string) bool {
	switch strings.ToLower(name) {
	case KeywordSine, KeywordSquare, KeywordTriangle, KeywordSawtooth:
		return true
	default:
		return false
	}
}

// BuiltinWaveformNames returns the names accepted without a custom definition.
func BuiltinWaveformNames() []string {
	return []string{KeywordSine, KeywordSquare, KeywordTriangle, KeywordSawtooth}
}

// NormalizeWaveformPoint maps the source-language range 0..100 to -1..1.
func NormalizeWaveformPoint(value float64) float64 {
	return value/50.0 - 1.0
}

// DenormalizeWaveformPoint maps the internal range -1..1 back to 0..100.
func DenormalizeWaveformPoint(value float64) float64 {
	return (value + 1.0) * 50.0
}
