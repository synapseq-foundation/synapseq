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

package types

import "fmt"

// Sequence represents a brainwave sequence
type Sequence struct {
	Periods    []Period
	Options    *SequenceOptions
	Comments   []string
	RawContent []byte
}

// SequenceOptions represents configuration options for a sequence
type SequenceOptions struct {
	// Sample rate (e.g., 44100)
	SampleRate int
	// Volume level (0-100 for 0-100%)
	Volume int
	// List of background audio files
	BackgroundList []string
	// List of preset configuration files
	PresetList []string
	// Gain level (20, 16, 12, 6, 0) for audio processing
	GainLevel GainLevel
}

// Validate checks if the sequence options are valid
func (so *SequenceOptions) Validate() error {
	if so.SampleRate <= 0 {
		return fmt.Errorf("invalid sample rate: %d", so.SampleRate)
	}
	if so.Volume < 0 || so.Volume > 100 {
		return fmt.Errorf("invalid volume: %d", so.Volume)
	}
	return nil
}
