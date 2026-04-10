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

package sources

import (
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type Signal struct {
	Kind        t.TrackType
	NoiseSmooth float64
	Waveform    efx.WaveformMorph
	Amplitude   [2]int
}

type NoiseGenerator interface {
	Generate(trackType t.TrackType, smooth float64) int
}