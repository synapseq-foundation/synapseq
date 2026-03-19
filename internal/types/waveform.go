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
