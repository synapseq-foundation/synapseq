// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type Noise struct {
	kind        t.TrackType
	noiseSmooth float64
	amplitude   int
}

func NewNoise(signal Signal) Noise {
	return Noise{kind: signal.Kind, noiseSmooth: signal.NoiseSmooth, amplitude: signal.Amplitude[0]}
}

func (source Noise) Sample(generator NoiseGenerator) int {
	return source.amplitude * generator.Generate(source.kind, source.noiseSmooth)
}