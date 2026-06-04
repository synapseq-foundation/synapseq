// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package spsq

import (
	"strconv"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// SampleRate sets the sample rate for the sequence.
func (b *Builder) SampleRate(sampleRate int) *Builder {
	b.options[t.KeywordOptionSampleRate] = strconv.Itoa(sampleRate)
	return b
}

// Volume sets the volume for the sequence.
func (b *Builder) Volume(volume int) *Builder {
	b.options[t.KeywordOptionVolume] = strconv.Itoa(volume)
	return b
}

// Ambiance registers an ambiance source for the sequence.
func (b *Builder) Ambiance(name, path string) *Builder {
	b.ambiance = append(b.ambiance, ambianceOption{name: name, path: path})
	return b
}

// Music registers a music source for the sequence.
func (b *Builder) Music(name, path string) *Builder {
	b.music = append(b.music, musicOption{name: name, path: path})
	return b
}
