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
	"fmt"
	"io"

	"github.com/synapseq-foundation/synapseq/v3/internal/audio"
	"github.com/synapseq-foundation/synapseq/v3/internal/info"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// generate generates the audio renderer based on the loaded sequence
func (ac *AppContext) generate() (*audio.AudioRenderer, error) {
	sequence := ac.sequence
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}

	options := sequence.Options
	if options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	renderer, err := audio.NewAudioRenderer(sequence.Periods, &audio.AudioRendererOptions{
		SampleRate:     options.SampleRate,
		Volume:         options.Volume,
		GainLevel:      options.GainLevel,
		BackgroundList: options.BackgroundList,
		StatusOutput:   ac.statusOutput,
	})
	if err != nil {
		return nil, err
	}

	return renderer, nil
}

// WAV generates the WAV file from the loaded sequence
func (ac *AppContext) WAV() error {
	renderer, err := ac.generate()
	if err != nil {
		return err
	}

	if err = renderer.RenderWav(ac.outputFile); err != nil {
		return err
	}

	presetList := ac.sequence.Options.PresetList
	if ac.format == t.FormatText && len(presetList) == 0 && !ac.unsafeNoMetadata {
		metadata, err := info.NewMetadata(ac.sequence.RawContent)
		if err != nil {
			return err
		}

		if err = audio.WriteICMTChunkFromTextFile(ac.outputFile, metadata); err != nil {
			return err
		}

	}

	return nil
}

// Stream generates the raw audio stream from the loaded sequence
func (ac *AppContext) Stream(data io.Writer) error {
	renderer, err := ac.generate()
	if err != nil {
		return err
	}

	err = renderer.RenderRaw(data)
	if err != nil {
		return err
	}

	return nil
}
