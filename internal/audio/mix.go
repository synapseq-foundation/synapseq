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

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

const phaseMask = (t.SineTableSize << 16) - 1

func (r *AudioRenderer) mix(samples []int) []int {
	for i := range t.BufferSize {
		var left, right int

		for ch := range t.NumberOfChannels {
			mixed := r.mixChannelSample(ch, i)
			left += mixed.left
			right += mixed.right
		}

		final := r.finalizeMixedSample(left, right)
		samples[i*2] = final.left
		samples[i*2+1] = final.right
	}

	return samples

}

func advancePhase(offset, increment int) int {
	return (offset + increment) & phaseMask
}

func (r *AudioRenderer) finalizeMixedSample(left, right int) stereoSample {
	if r.Volume != 100 {
		left = left * r.Volume / 100
		right = right * r.Volume / 100
	}

	left >>= audioBitShift
	right >>= audioBitShift

	return stereoSample{
		left:  clampPCM16(left),
		right: clampPCM16(right),
	}

}

func clampPCM16(sample int) int {
	if sample > audioMaxValue {
		return audioMaxValue
	}
	if sample < audioMinValue {
		return audioMinValue
	}

	return sample
}
