// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const (
	shiftTaps   = 63
	shiftCenter = shiftTaps / 2
)

type shiftState struct {
	history       [shiftTaps]float64
	writeIndex    int
	oscillatorSin float64
	oscillatorCos float64
	deltaSin      float64
	deltaCos      float64
	separation    float64
	initialized   bool
	samples       uint32
}

// ApplyShift derives a mono wet signal and shifts it symmetrically between the
// stereo channels. The original stereo signal remains the dry path.
func (p *Processor) ApplyShift(ch int, effect t.Effect, left, right int) (int, int) {
	if ch < 0 || ch >= len(p.shiftStates) {
		return left, right
	}

	state := &p.shiftStates[ch]
	state.updateOscillator(p.sampleRate, effect.Value)

	mono := (float64(left) + float64(right)) * 0.5
	state.history[state.writeIndex] = mono
	delayed := state.history[ringIndex(state.writeIndex-shiftCenter)]

	var quadrature float64
	for tap, coefficient := range p.shiftCoefficients {
		quadrature += coefficient * state.history[ringIndex(state.writeIndex-tap)]
	}
	state.writeIndex = ringIndex(state.writeIndex + 1)

	wetLeft := delayed*state.oscillatorCos - quadrature*state.oscillatorSin
	wetRight := delayed*state.oscillatorCos + quadrature*state.oscillatorSin
	state.advanceOscillator()

	intensity := float64(effect.Intensity)
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}
	if intensity == 0 || effect.Value == 0 {
		return left, right
	}
	dry := 1 - intensity
	return int(math.Round(float64(left)*dry + wetLeft*intensity)),
		int(math.Round(float64(right)*dry + wetRight*intensity))
}

// ResetShift clears the frequency shifter history for one sequencer channel.
func (p *Processor) ResetShift(ch int) {
	if ch < 0 || ch >= len(p.shiftStates) {
		return
	}
	p.shiftStates[ch] = shiftState{}
}

func (state *shiftState) updateOscillator(sampleRate int, separation float64) {
	if state.initialized && state.separation == separation {
		return
	}

	// Each side moves by half of the declared total stereo separation.
	delta := math.Pi * separation / float64(sampleRate)
	state.deltaSin = math.Sin(delta)
	state.deltaCos = math.Cos(delta)
	state.separation = separation
	if !state.initialized {
		state.oscillatorCos = 1
		state.initialized = true
	}
}

func (state *shiftState) advanceOscillator() {
	nextCos := state.oscillatorCos*state.deltaCos - state.oscillatorSin*state.deltaSin
	nextSin := state.oscillatorSin*state.deltaCos + state.oscillatorCos*state.deltaSin
	state.oscillatorCos = nextCos
	state.oscillatorSin = nextSin
	state.samples++

	if state.samples%4096 == 0 {
		norm := math.Hypot(state.oscillatorCos, state.oscillatorSin)
		if norm != 0 {
			state.oscillatorCos /= norm
			state.oscillatorSin /= norm
		}
	}
}

func makeShiftCoefficients() [shiftTaps]float64 {
	var coefficients [shiftTaps]float64
	for tap := range shiftTaps {
		offset := tap - shiftCenter
		if offset == 0 || offset%2 == 0 {
			continue
		}

		window := 0.42 - 0.5*math.Cos(2*math.Pi*float64(tap)/float64(shiftTaps-1)) +
			0.08*math.Cos(4*math.Pi*float64(tap)/float64(shiftTaps-1))
		coefficients[tap] = 2 * window / (math.Pi * float64(offset))
	}
	return coefficients
}

func ringIndex(index int) int {
	index %= shiftTaps
	if index < 0 {
		index += shiftTaps
	}
	return index
}
