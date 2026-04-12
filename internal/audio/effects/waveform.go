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

package effects

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func (p *Processor) WaveformSample(channel *t.Channel, offset int) int {
	return p.WaveformSampleForMorph(WaveformMorphFromChannel(channel), offset)
}

func (p *Processor) WaveformValue(channel *t.Channel, offset int) float64 {
	return p.WaveformValueForMorph(WaveformMorphFromChannel(channel), offset)
}

func (p *Processor) WaveformSampleForMorph(waveform WaveformMorph, offset int) int {
	return int(math.Round(p.WaveformValueForMorph(waveform, offset)))
}

func (p *Processor) WaveformValueForMorph(waveform WaveformMorph, offset int) float64 {
	startWaveform, endWaveform, alpha := normalizedWaveformMorph(waveform)
	idx := offset >> 16
	start := float64(p.waveTables[int(startWaveform)][idx])
	if alpha <= 0 || startWaveform == endWaveform {
		return start
	}

	end := float64(p.waveTables[int(endWaveform)][idx])
	if alpha >= 1 {
		return end
	}

	return lerpFloat64(start, end, alpha)
}

func lerpFloat64(start, end, alpha float64) float64 {
	return start + (end-start)*alpha
}