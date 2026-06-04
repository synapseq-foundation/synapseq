// Deprecated: the SynapSeq WebAssembly runtime is kept only for historical
// reference and is no longer recommended for new integrations.
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

package main

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/audio"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type rendererBuilder struct{}

func buildWASMRendererOptions(sequence *t.Sequence) (*audio.AudioRendererOptions, error) {
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}
	if sequence.Options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	return &audio.AudioRendererOptions{
		SampleRate: sequence.Options.SampleRate,
		Volume:     sequence.Options.Volume,
		Ambiance:   sequence.Options.Ambiance,
		Music:      sequence.Options.Music,
		Colors:     false,
	}, nil
}

func (rendererBuilder) Build(sequence *t.Sequence) (renderableAudio, error) {
	options, err := buildWASMRendererOptions(sequence)
	if err != nil {
		return nil, err
	}

	return audio.NewAudioRenderer(sequence.Periods, options)
}
