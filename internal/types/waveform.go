// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

// WaveformType represents the waveform shape
type WaveformType int

// Waveform types
const (
	WaveformSine     WaveformType = iota // Sine
	WaveformSquare                       // Square
	WaveformTriangle                     // Triangle
	WaveformSawtooth                     // Sawtooth
)

// String returns the string representation of WaveformType
func (wt WaveformType) String() string {
	switch wt {
	case WaveformSine:
		return KeywordSine
	case WaveformSquare:
		return KeywordSquare
	case WaveformTriangle:
		return KeywordTriangle
	case WaveformSawtooth:
		return KeywordSawtooth
	default:
		return ""
	}
}

// WaveformString returns the WaveformType from a string representation
func WaveformString(s string) WaveformType {
	switch s {
	case KeywordSine:
		return WaveformSine
	case KeywordSquare:
		return WaveformSquare
	case KeywordTriangle:
		return WaveformTriangle
	case KeywordSawtooth:
		return WaveformSawtooth
	default:
		return WaveformSine
	}
}
