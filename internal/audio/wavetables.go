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

package audio

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// InitWaveformTables initializes the waveform tables
func InitWaveformTables() [4][]int {
	var waveTables [4][]int
	for i := range waveTables {
		waveformTable := make([]int, t.SineTableSize)

		for j := range t.SineTableSize {
			phase := float64(j) * 2.0 * float64(math.Pi) / float64(t.SineTableSize)
			var val float64

			switch i { // i is the waveform type
			case int(t.WaveformSine):
				val = math.Sin(phase)
			case int(t.WaveformSquare):
				if math.Sin(phase) > 0 {
					val = 1.0
				} else {
					val = -1.0
				}
			case int(t.WaveformTriangle):
				if phase < math.Pi {
					val = (2.0 * phase / math.Pi) - 1.0
				} else {
					val = 3.0 - (2.0 * phase / math.Pi)
				}
			case int(t.WaveformSawtooth):
				// Center the discontinuity at pi so waveform morphs stay phase-aligned with
				// square/triangle and do not sound like a fade during transitions.
				val = 2.0 * (phase/(2.0*math.Pi) - math.Floor(phase/(2.0*math.Pi)+0.5))
			default:
				val = math.Sin(phase)
			}

			waveformTable[j] = int(t.WaveTableAmplitude * val)
		}

		waveTables[i] = waveformTable
	}
	return waveTables
}
