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
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	// NoiseShift is the bit shift for noise generation
	noiseShift = 12
	// NoiseAmplitude is the amplitude for noise generation
	noiseAmplitude = t.WaveTableAmplitude << noiseShift
	// NoiseBands is the number of bands for noise generation
	noiseBands = 9
	// Random multiplier for noise generation
	randMult = 75
)

// NoiseGenerator handles all noise generation
type NoiseGenerator struct {
	// Pink noise state
	noiseTables    [noiseBands]pinkNoise
	noiseOffset    int
	noiseBuffer    [256]int
	noiseBufferOff uint8

	// Random seed (shared across all noise types)
	seed int

	// Brown noise state
	brownLast int
}

// pinkNoise represents a pink noise generator state
type pinkNoise struct {
	// Current output value
	value int
	// Increment
	increment int
}

// NewNoiseGenerator creates a new noise generator with initial seed
func NewNoiseGenerator() *NoiseGenerator {
	return &NoiseGenerator{
		// Initial seed
		seed:      2,
		brownLast: 0,
	}
}

// Generate generates a noise sample based on the track type
func (ng *NoiseGenerator) Generate(tr t.TrackType) int {
	switch tr {
	case t.TrackWhiteNoise:
		return ng.generateWhiteNoise()
	case t.TrackPinkNoise:
		return ng.generatePinkNoise()
	case t.TrackBrownNoise:
		return ng.generateBrownNoise()
	default:
		return 0
	}
}

// generateWhiteNoise generates white noise sample
func (ng *NoiseGenerator) generateWhiteNoise() int {
	// White noise is simply a random value without filtering
	// return ((seed = seed * RAND_MULT % 131074) - 65535) * (ST_AMP / 65535);
	ng.seed = (ng.seed * randMult) % 131074
	randomValue := ng.seed - 65535
	return randomValue * (t.WaveTableAmplitude / 65535)
}

// generateBrownNoise generates brown noise sample
func (ng *NoiseGenerator) generateBrownNoise() int {
	// Generate a random value
	ng.seed = (ng.seed * randMult) % 131074
	random := ng.seed - 65535

	// Integrate the random value with a decay factor to avoid overflow
	ng.brownLast = int(float64(ng.brownLast+random/16) * 0.9)

	// Limit the value to avoid overflow
	if ng.brownLast > 65535 {
		ng.brownLast = 65535
	}
	if ng.brownLast < -65535 {
		ng.brownLast = -65535
	}

	// Scale to the same level as the wave table
	return ng.brownLast * (t.WaveTableAmplitude / 65535)
}

// generatePinkNoise generates pink noise sample
func (ng *NoiseGenerator) generatePinkNoise() int {
	var tot int
	off := ng.noiseOffset
	ng.noiseOffset++
	cnt := 1
	ns := 0 // index into noiseTables

	// Generate base random value
	ng.seed = (ng.seed * randMult) % 131074
	tot = (ng.seed - 65535) * (noiseAmplitude / 65535 / (noiseBands + 1))

	// Process noise bands
	for (cnt&off) != 0 && ns < noiseBands {
		ng.seed = (ng.seed * randMult) % 131074
		val := (ng.seed - 65535) * (noiseAmplitude / 65535 / (noiseBands + 1))

		cnt += cnt
		ng.noiseTables[ns].increment = (val - ng.noiseTables[ns].value) / cnt
		ng.noiseTables[ns].value += ng.noiseTables[ns].increment
		tot += ng.noiseTables[ns].value
		ns++
	}

	// Add remaining noise bands
	for ns < noiseBands {
		ng.noiseTables[ns].value += ng.noiseTables[ns].increment
		tot += ng.noiseTables[ns].value
		ns++
	}

	// Store in buffer and return scaled value
	ng.noiseBuffer[ng.noiseBufferOff] = tot >> noiseShift
	ng.noiseBufferOff++

	return ng.noiseBuffer[ng.noiseBufferOff-1]
}
