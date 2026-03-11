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
	"math/bits"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	// NoiseShift is the bit shift for noise generation
	noiseShift = 12
	// NoiseAmplitude is the amplitude for noise generation
	noiseAmplitude = t.WaveTableAmplitude << noiseShift
	// NoiseBands is the number of bands for noise generation
	noiseBands = 9
	// Initial seed for deterministic noise generation
	initialNoiseSeed uint32 = 0x9E3779B9
	// Max centered random value for 16-bit signed samples
	maxCenteredRandom = 1 << 15
	// White/brown noise scale to the wave table amplitude range
	whiteNoiseScale = t.WaveTableAmplitude / maxCenteredRandom
	// Pink noise contribution for the base value and each band
	pinkContributionScale = noiseAmplitude / maxCenteredRandom / (noiseBands + 1)
	// Brown noise input attenuation to keep integration stable
	brownInputDivisor = 16
	// Brown noise decay factor expressed as an integer fraction
	brownDecayNumerator   = 9
	brownDecayDenominator = 10
)

// NoiseGenerator handles all noise generation
type NoiseGenerator struct {
	// Pink noise state
	pinkState pinkNoiseState

	// Random seed (shared across all noise types)
	seed uint32

	// Brown noise state
	brownLast int
}

// pinkNoiseState holds the per-band state for the Voss-McCartney pink noise generator.
type pinkNoiseState struct {
	bands   [noiseBands]pinkNoiseBand
	counter uint32
}

// pinkNoiseBand represents one octave band in the pink noise generator.
type pinkNoiseBand struct {
	value     int
	increment int
}

// NewNoiseGenerator creates a new noise generator with initial seed
func NewNoiseGenerator() *NoiseGenerator {
	ng := &NoiseGenerator{
		seed: initialNoiseSeed,
	}
	ng.initPinkNoise()
	return ng
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

// nextRandom generates the next deterministic pseudo-random value.
func (ng *NoiseGenerator) nextRandom() uint32 {
	if ng.seed == 0 {
		ng.seed = initialNoiseSeed
	}

	ng.seed ^= ng.seed << 13
	ng.seed ^= ng.seed >> 17
	ng.seed ^= ng.seed << 5

	return ng.seed
}

// nextCenteredRandom returns a signed 16-bit value in the range [-32768, 32767].
func (ng *NoiseGenerator) nextCenteredRandom() int {
	return int(ng.nextRandom()>>16) - maxCenteredRandom
}

func (ng *NoiseGenerator) nextPinkContribution() int {
	return ng.nextCenteredRandom() * pinkContributionScale
}

func (ng *NoiseGenerator) initPinkNoise() {
	for i := range ng.pinkState.bands {
		ng.pinkState.bands[i].value = ng.nextPinkContribution()
	}
}

// generateWhiteNoise generates white noise sample
func (ng *NoiseGenerator) generateWhiteNoise() int {
	return ng.nextCenteredRandom() * whiteNoiseScale
}

// generateBrownNoise generates brown noise sample
func (ng *NoiseGenerator) generateBrownNoise() int {
	random := ng.nextCenteredRandom()

	ng.brownLast += random / brownInputDivisor
	ng.brownLast = ng.brownLast * brownDecayNumerator / brownDecayDenominator

	if ng.brownLast > maxCenteredRandom {
		ng.brownLast = maxCenteredRandom
	}
	if ng.brownLast < -maxCenteredRandom {
		ng.brownLast = -maxCenteredRandom
	}

	return ng.brownLast * whiteNoiseScale
}

// generatePinkNoise generates pink noise sample
func (ng *NoiseGenerator) generatePinkNoise() int {
	total := ng.nextPinkContribution()
	updatedBands := bits.TrailingZeros32(^ng.pinkState.counter)
	if updatedBands > noiseBands {
		updatedBands = noiseBands
	}
	ng.pinkState.counter++

	for bandIdx := range ng.pinkState.bands {
		band := &ng.pinkState.bands[bandIdx]
		if bandIdx < updatedBands {
			steps := 1 << (bandIdx + 1)
			target := ng.nextPinkContribution()
			band.increment = (target - band.value) / steps
		}

		band.value += band.increment
		total += band.value
	}

	return total >> noiseShift
}
