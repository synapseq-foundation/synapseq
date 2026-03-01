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

// LoadSequence loads the sequence from the input file based on the specified format
func (ac *AppContext) LoadSequence() error {
	var err error
	ac.sequence, err = seq.LoadTextSequence(ac.inputFile)
	if err != nil {
		return err
	}

	return nil
}

// Comments returns the comments from the loaded sequence
func (ac *AppContext) Comments() []string {
	if ac.sequence == nil {
		return nil
	}
	return ac.sequence.Comments
}

// SampleRate returns the sample rate from the loaded sequence options
func (ac *AppContext) SampleRate() int {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return 0
	}

	return ac.sequence.Options.SampleRate
}

// PresetList returns the preset list from the loaded sequence options
func (ac *AppContext) PresetList() []string {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return []string{}
	}

	return ac.sequence.Options.PresetList
}

// Volume returns the volume from the loaded sequence options
func (ac *AppContext) Volume() int {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return 0
	}

	return ac.sequence.Options.Volume
}

// AmbianceList returns the ambiance audio list from the loaded sequence options
func (ac *AppContext) AmbianceList() map[string]string {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return map[string]string{}
	}

	return ac.sequence.Options.AmbianceList
}

// RawContent returns the raw content of the loaded sequence
func (ac *AppContext) RawContent() []byte {
	if ac.sequence == nil {
		return nil
	}

	return ac.sequence.RawContent
}
