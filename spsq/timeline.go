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

package spsq

import (
	"fmt"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// SilenceAt adds a silence timeline entry at the given time.
func (b *Builder) SilenceAt(hh, mm, ss int) *Builder {
	time := fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	b.timeline = append(b.timeline, [4]string{time, t.KeywordSilence, t.KeywordTransitionSteady, "0"})
	return b
}

// At adds a preset timeline entry at the given time.
func (b *Builder) At(hh, mm, ss int, preset *Preset) *Builder {
	if preset == nil {
		return b
	}

	time := fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	b.timeline = append(b.timeline, [4]string{time, preset.name, t.KeywordTransitionSteady, "0"})
	return b
}

// WithSteadyTransition sets the transition of the last timeline entry to steady
func (b *Builder) WithSteadyTransition() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx][timelineTransition] = t.KeywordTransitionSteady
	return b
}

// WithEaseInTransition sets the transition of the last timeline entry to ease-in
func (b *Builder) WithEaseInTransition() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx][timelineTransition] = t.KeywordTransitionEaseIn
	return b
}

// WithEaseOutTransition sets the transition of the last timeline entry to ease-out
func (b *Builder) WithEaseOutTransition() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx][timelineTransition] = t.KeywordTransitionEaseOut
	return b
}

// WithSmoothTransition sets the transition of the last timeline entry to ease-out
func (b *Builder) WithSmoothTransition() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx][timelineTransition] = t.KeywordTransitionSmooth
	return b
}

// WithStep sets the step of the last timeline entry
func (b *Builder) WithStep(s int) *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx][timelineStep] = fmt.Sprintf("%d", s)
	return b
}
