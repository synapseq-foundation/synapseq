// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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