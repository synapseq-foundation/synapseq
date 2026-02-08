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

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// calcDopplerFactor returns a pitch multiplier in [1-depth .. 1+depth] based on a sine LFO.
func (r *AudioRenderer) calcDopplerFactor(offset int, intensity t.IntensityType) float64 {
	inten := float64(intensity) // expected 0..1
	if inten < 0 {
		inten = 0
	}
	if inten > 1 {
		inten = 1
	}

	lfo := r.waveTables[int(t.WaveformSine)][offset>>16] // [-A..A]
	lfoNorm := float64(lfo) / float64(t.WaveTableAmplitude)

	depth := 0.05 * inten
	return 1.0 + (depth * lfoNorm)
}

// calcPulseFactor calculates the pulse effect modulation factor for a channel
func (r *AudioRenderer) calcPulseFactor(waveform t.WaveformType, offset int) float64 {
	modVal := float64(r.waveTables[int(waveform)][offset>>16])

	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}

// applySpin applies the spin effect to the given input samples for a channel.
func (r *AudioRenderer) applySpin(channel *t.Channel, inL, inR int) (outL, outR int) {
	intensity := float64(channel.Track.Intensity) // 0..1

	// Use a sine LFO for autopan (independent from the audio waveform)
	lfo := r.waveTables[int(t.WaveformSine)][channel.Effect.Offset>>16] // [-A..A]

	// Normalize to [-1..1] and map to pan [-128..128] with rounding (less "sticking" at center)
	panF := (float64(lfo) / float64(t.WaveTableAmplitude)) * 128.0 * intensity
	pan := int(math.Round(panF))

	if pan < -128 {
		pan = -128
	}
	if pan > 128 {
		pan = 128
	}

	pos := pan + 128 // 0..256
	lGain := 256 - pos
	rGain := pos

	outL = int((int64(inL) * int64(lGain)) >> 8) // /256
	outR = int((int64(inR) * int64(rGain)) >> 8)

	return outL, outR
}
