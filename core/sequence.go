//go:build !wasm

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

package core

import (
	seq "github.com/synapseq-foundation/synapseq/v4/internal/sequence"
)

// Load loads the sequence from the input file.
func (ac *AppContext) Load(path string) (*LoadedContext, error) {
	sequence, err := seq.LoadTextSequence(path)
	if err != nil {
		return nil, err
	}

	return &LoadedContext{
		appCtx:   ac,
		sequence: sequence,
	}, nil
}

// Comments returns the comments from the loaded sequence.
func (lc *LoadedContext) Comments() []string {
	if lc.sequence == nil {
		return nil
	}
	return lc.sequence.Comments
}

// SampleRate returns the sample rate from the loaded sequence options.
func (lc *LoadedContext) SampleRate() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.SampleRate
}

// PresetList returns the preset list from the loaded sequence options.
func (lc *LoadedContext) PresetList() []string {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return []string{}
	}

	return lc.sequence.Options.PresetList
}

// Volume returns the volume from the loaded sequence options.
func (lc *LoadedContext) Volume() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.Volume
}

// AmbianceList returns the ambiance audio list from the loaded sequence options.
func (lc *LoadedContext) AmbianceList() map[string]string {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return map[string]string{}
	}

	return lc.sequence.Options.AmbianceList
}

// RawContent returns the raw content of the loaded sequence.
func (lc *LoadedContext) RawContent() []byte {
	if lc.sequence == nil {
		return nil
	}

	return lc.sequence.RawContent
}
