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