// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

const (
	BufferSize         = 1024
	SineTableSize      = 16384
	WaveTableAmplitude = 0x7FFFF
	PhasePrecision     = 65536
)

// AmplitudeType represents the amplitude level of an audio signal
type AmplitudeType float64

// ToPercent converts a raw amplitude value to a float64 percentage
func (a AmplitudeType) ToPercent() float64 {
	return float64(a / 40.96)
}

// AmplitudePercentToRaw converts a float64 value to a raw amplitude value
func AmplitudePercentToRaw(v float64) AmplitudeType {
	return AmplitudeType(v * 40.96)
}

// IntensityType represents the intensity level of an audio signal
type IntensityType float64

// ToPercent converts a raw intensity value to a float64 percentage
func (i IntensityType) ToPercent() float64 {
	return float64(i * 100)
}

// IntensityPercentToRaw converts a float64 value to a raw intensity value
func IntensityPercentToRaw(v float64) IntensityType {
	return IntensityType(v / 100)
}
