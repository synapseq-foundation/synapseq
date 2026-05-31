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
	"time"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// SilenceAt adds a silence timeline entry at the given time.
func (b *Builder) SilenceAt(at time.Duration) *Builder {
	b.timeline = append(b.timeline, timelineEntry{
		at:         at,
		presetName: t.KeywordSilence,
		transition: t.KeywordTransitionSteady,
		steps:      0,
	})
	return b
}

// PresetAt adds a preset timeline entry at the given time.
func (b *Builder) PresetAt(at time.Duration, preset *Preset) *Builder {
	if preset == nil {
		return b
	}

	b.timeline = append(b.timeline, timelineEntry{
		at:         at,
		presetName: preset.name,
		transition: t.KeywordTransitionSteady,
		steps:      0,
	})
	return b
}

// Steady sets the transition of the last timeline entry to steady.
func (b *Builder) Steady() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx].transition = t.KeywordTransitionSteady
	return b
}

// EaseIn sets the transition of the last timeline entry to ease-in.
func (b *Builder) EaseIn() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx].transition = t.KeywordTransitionEaseIn
	return b
}

// EaseOut sets the transition of the last timeline entry to ease-out.
func (b *Builder) EaseOut() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx].transition = t.KeywordTransitionEaseOut
	return b
}

// Smooth sets the transition of the last timeline entry to smooth.
func (b *Builder) Smooth() *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx].transition = t.KeywordTransitionSmooth
	return b
}

// Step sets the step count of the last timeline entry.
func (b *Builder) Step(s int) *Builder {
	if len(b.timeline) == 0 {
		return b
	}

	timeIdx := len(b.timeline) - 1
	b.timeline[timeIdx].steps = s
	return b
}

// formatTimelineTime formats a time duration as a string in the format "HH:MM:SS"
func formatTimelineTime(at time.Duration) string {
	totalSeconds := int(at / time.Second)
	hh := totalSeconds / 3600
	mm := (totalSeconds % 3600) / 60
	ss := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
