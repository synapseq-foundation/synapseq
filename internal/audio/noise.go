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

package audio

import (
	"math/bits"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	noiseShift                          = 12
	noiseAmplitude                      = t.WaveTableAmplitude << noiseShift
	noiseBands                          = 9
	initialNoiseSeed             uint32 = 0x9E3779B9
	maxCenteredRandom                   = 1 << 15
	whiteNoiseScale                     = t.WaveTableAmplitude / maxCenteredRandom
	pinkContributionScale               = noiseAmplitude / maxCenteredRandom / (noiseBands + 1)
	brownInputDivisor                   = 16
	brownDecayNumerator                 = 9
	brownDecayDenominator               = 10
	noiseSmoothnessStages               = 3
	noiseSmoothnessLinearCeiling        = 60.0
	noiseSmoothnessMinAlpha             = 0.02
)

// NoiseGenerator handles all noise generation
type NoiseGenerator struct {
	pinkState   pinkNoiseState
	seed        uint32
	brownLast   int
	smoothState map[t.TrackType]noiseSmoothnessState
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

// noiseSmoothnessState holds the per-track state for noise smoothing.
type noiseSmoothnessState struct {
	initialized bool
	stages      [noiseSmoothnessStages]float64
}

// noiseSmoothnessProfile holds the linear target and max effective alpha for a noise smoothness stage.
type noiseSmoothnessProfile struct {
	linearTarget float64
	maxEffective float64
}

// NewNoiseGenerator creates a new noise generator with initial seed
func NewNoiseGenerator() *NoiseGenerator {
	ng := &NoiseGenerator{
		seed:        initialNoiseSeed,
		smoothState: make(map[t.TrackType]noiseSmoothnessState),
	}
	ng.initPinkNoise()
	return ng
}

// Generate generates a noise sample based on the track type
func (ng *NoiseGenerator) Generate(tr t.TrackType, smooth float64) int {
	var sample int

	switch tr {
	case t.TrackWhiteNoise:
		sample = ng.generateWhiteNoise()
	case t.TrackPinkNoise:
		sample = ng.generatePinkNoise()
	case t.TrackBrownNoise:
		sample = ng.generateBrownNoise()
	default:
		return 0
	}

	return ng.applySmoothness(tr, smooth, sample)
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

// nextPinkContribution returns a scaled pink noise contribution.
func (ng *NoiseGenerator) nextPinkContribution() int {
	return ng.nextCenteredRandom() * pinkContributionScale
}

// initPinkNoise initializes the pink noise state with random contributions.
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

// applySmoothness applies noise smoothing to a sample using the per-track state.
func (ng *NoiseGenerator) applySmoothness(tr t.TrackType, smooth float64, sample int) int {
	if smooth <= 0 {
		return sample
	}
	if smooth > 100 {
		smooth = 100
	}

	state := ng.smoothState[tr]

	if !state.initialized {
		for idx := range state.stages {
			state.stages[idx] = float64(sample)
		}
		state.initialized = true
		ng.smoothState[tr] = state
		return sample
	}

	alpha := noiseSmoothnessAlpha(tr, smooth)
	value := float64(sample)
	for idx := range state.stages {
		state.stages[idx] += (value - state.stages[idx]) * alpha
		value = state.stages[idx]
	}

	ng.smoothState[tr] = state
	return clampNoiseSample(int(value))
}

// noiseSmoothnessAlpha returns the noise smoothing alpha for a given track and smoothness level.
func noiseSmoothnessAlpha(tr t.TrackType, smoothness float64) float64 {
	normalized := effectiveNoiseSmoothness(tr, smoothness) / 100.0
	inverse := 1.0 - normalized
	return noiseSmoothnessMinAlpha + (1.0-noiseSmoothnessMinAlpha)*inverse*inverse
}

// effectiveNoiseSmoothness returns the effective noise smoothness for a given track and smoothness level.
func effectiveNoiseSmoothness(tr t.TrackType, smoothness float64) float64 {
	profile := noiseSmoothnessProfileForTrack(tr)

	if smoothness < 0 {
		return 0
	}
	if smoothness > 100 {
		smoothness = 100
	}
	if smoothness <= noiseSmoothnessLinearCeiling {
		return smoothness * profile.linearTarget / noiseSmoothnessLinearCeiling
	}

	progress := (smoothness - noiseSmoothnessLinearCeiling) / (100.0 - noiseSmoothnessLinearCeiling)
	return profile.linearTarget + (profile.maxEffective-profile.linearTarget)*progress*progress
}

// noiseSmoothnessProfileForTrack returns the noise smoothness profile for a given track type.
func noiseSmoothnessProfileForTrack(tr t.TrackType) noiseSmoothnessProfile {
	switch tr {
	case t.TrackWhiteNoise:
		return noiseSmoothnessProfile{linearTarget: 30, maxEffective: 42}
	case t.TrackPinkNoise:
		return noiseSmoothnessProfile{linearTarget: 24, maxEffective: 34}
	case t.TrackBrownNoise:
		return noiseSmoothnessProfile{linearTarget: 18, maxEffective: 28}
	default:
		return noiseSmoothnessProfile{linearTarget: 24, maxEffective: 34}
	}
}

// clampNoiseSample clamps a noise sample to the wave table amplitude range.
func clampNoiseSample(sample int) int {
	if sample > int(t.WaveTableAmplitude) {
		return int(t.WaveTableAmplitude)
	}
	if sample < -int(t.WaveTableAmplitude) {
		return -int(t.WaveTableAmplitude)
	}
	return sample
}
