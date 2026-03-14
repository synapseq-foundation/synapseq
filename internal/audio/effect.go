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

// calcModulationFactor calculates the pulse effect modulation factor for a channel
func (r *AudioRenderer) calcModulationFactor(channel *t.Channel, offset int) float64 {
	modVal := r.waveformValue(channel, offset)

	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}

// applyPan applies the pan effect to the given input samples for a channel.
func (r *AudioRenderer) applyPan(channel *t.Channel, inL, inR int) (outL, outR int) {
	// Intensity is already 0..1 (see types.IntensityType)
	intensity := float64(channel.Track.Effect.Intensity)
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}

	// Selected waveform in [-A..A] => normalize to [-1..1]
	targetX := (r.waveformValue(channel, channel.Effect.Offset) / float64(t.WaveTableAmplitude)) * intensity
	x := r.smoothedPanPosition(channel, targetX)

	// High-resolution linear pan gains (0..65536), avoids 8-bit pan stepping artifacts
	// x=-1 => left=1.0 right=0.0 ; x=+1 => left=0.0 right=1.0
	pos := int(math.Round((x + 1.0) * 32768.0)) // 0..65536
	if pos < 0 {
		pos = 0
	}
	if pos > 65536 {
		pos = 65536
	}

	lGain := 65536 - pos
	rGain := pos

	outL = int((int64(inL) * int64(lGain)) >> 16)
	outR = int((int64(inR) * int64(rGain)) >> 16)
	return outL, outR
}

func (r *AudioRenderer) smoothedPanPosition(channel *t.Channel, targetX float64) float64 {
	if !channel.Effect.PanInitialized {
		channel.Effect.PanPosition = targetX
		channel.Effect.PanInitialized = true
		return targetX
	}

	maxDelta := 2 * r.effectSlewMaxDelta()
	delta := targetX - channel.Effect.PanPosition
	if delta > maxDelta {
		delta = maxDelta
	} else if delta < -maxDelta {
		delta = -maxDelta
	}

	channel.Effect.PanPosition += delta
	if channel.Effect.PanPosition > 1 {
		channel.Effect.PanPosition = 1
	}
	if channel.Effect.PanPosition < -1 {
		channel.Effect.PanPosition = -1
	}

	return channel.Effect.PanPosition
}
