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
	"strconv"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// AddSampleRateOption sets the sample rate for the sequence
func (b *Builder) AddSampleRateOption(sampleRate int) *Builder {
	b.options[t.KeywordOptionSampleRate] = strconv.Itoa(sampleRate)
	return b
}

// AddVolumeOption sets the volume for the sequence
func (b *Builder) AddVolumeOption(volume int) *Builder {
	b.options[t.KeywordOptionVolume] = strconv.Itoa(volume)
	return b
}

// AddAmbianceOption sets the ambiance for the sequence
func (b *Builder) AddAmbianceOption(name, path string) *Builder {
	b.ambiance = append(b.ambiance, []string{name, path})
	return b
}

// AddExtendsOption sets the extends for the sequence
func (b *Builder) AddExtendsOption(extends string) *Builder {
	b.extends = append(b.extends, extends)
	return b
}
