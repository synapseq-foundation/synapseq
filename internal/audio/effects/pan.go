// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func (p *Processor) ApplyPanEffect(channel *t.Channel, effect t.Effect, waveform WaveformMorph, left, right int) (int, int) {
	p.advanceEffectPhase(channel)
	return p.ApplyPanForMorph(channel, effect, waveform, left, right)
}

func (p *Processor) ApplyPan(channel *t.Channel, inL, inR int) (outL, outR int) {
	return p.ApplyPanForMorph(channel, channel.Track.Effect, WaveformMorphFromChannel(channel), inL, inR)
}

func (p *Processor) ApplyPanForMorph(channel *t.Channel, effect t.Effect, waveform WaveformMorph, inL, inR int) (outL, outR int) {
	intensity := float64(effect.Intensity)
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}

	targetX := (p.WaveformValueForMorph(waveform, channel.Effect.Offset) / float64(t.WaveTableAmplitude)) * intensity
	x := p.smoothedPanPosition(channel, targetX)

	pos := int(math.Round((x + 1.0) * 32768.0))
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