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

const (
	BufferSize         = 1024    // Buffer size for audio processing
	SineTableSize      = 16384   // Number of elements in sine-table (power of 2)
	WaveTableAmplitude = 0x7FFFF // Amplitude of wave in wave-table
	PhasePrecision     = 65536   // Phase precision (1/65536 of a cycle)
)

type AmplitudeType float64 // Amplitude level (0-4096 for 0-100%)

// ToPercent converts a raw amplitude value to a float64 percentage
func (a AmplitudeType) ToPercent() float64 {
	return float64(a / 40.96)
}

// AmplitudePercentToRaw converts a float64 value to a raw amplitude value
func AmplitudePercentToRaw(v float64) AmplitudeType {
	return AmplitudeType(v * 40.96)
}

type IntensityType float64 // Intensity level (0-1.0 for 0-100%)

// ToPercent converts a raw intensity value to a float64 percentage
func (i IntensityType) ToPercent() float64 {
	return float64(i * 100)
}

// IntensityPercentToRaw converts a float64 value to a raw intensity value
func IntensityPercentToRaw(v float64) IntensityType {
	return IntensityType(v / 100)
}
