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
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

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

// calcSpinIncrement calculates the spin effect modulation increment for a channel
func (r *AudioRenderer) calcSpinIncrement(channel *t.Channel) float64 {
	spinCarrierMax := 127.0 / 1e-6 / float64(r.SampleRate)
	clampedWidth := channel.Track.Effect.Configuration.(t.EffectSpinConfiguration).Width

	if clampedWidth > spinCarrierMax {
		clampedWidth = spinCarrierMax
	}
	if clampedWidth < -spinCarrierMax {
		clampedWidth = -spinCarrierMax
	}

	return clampedWidth * 1e-6 * float64(r.SampleRate) * float64(1<<24) / float64(t.WaveTableAmplitude)
}

// calcSpinPan returns a pan position in [-128..127] based on spinPos and intensity.
// spinPos is expected to be roughly in [-128..127] as in the current code.
func calcSpinPan(spinPos int, intensity float64) int {
	spinGain := 0.5 + (intensity*0.7)*3.5

	ampSpin := int(float64(spinPos) * spinGain)
	if ampSpin > 127 {
		ampSpin = 127
	}
	if ampSpin < -128 {
		ampSpin = -128
	}
	return ampSpin
}

// applySpinCrossMix applies the same cross-mix logic you currently use.
// Inputs must already be scaled (i.e., amplitude already applied).
func applySpinCrossMix(inL, inR int, ampSpin int) (outL, outR int) {
	posVal := ampSpin
	if posVal < 0 {
		posVal = -posVal
	}
	if posVal > 128 {
		posVal = 128
	}

	if ampSpin >= 0 {
		outL = (inL * (128 - posVal)) >> 7
		outR = inR + ((inL * posVal) >> 7)
	} else {
		outL = inL + ((inR * posVal) >> 7)
		outR = (inR * (128 - posVal)) >> 7
	}

	return outL, outR
}
