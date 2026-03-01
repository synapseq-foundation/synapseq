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

	"github.com/synapseq-foundation/synapseq/v4/internal/audio"
)

// generate generates the audio renderer based on the loaded sequence
func (lc *LoadedContext) generate() (*audio.AudioRenderer, error) {
	sequence := lc.sequence
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}

	options := sequence.Options
	if options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	renderer, err := audio.NewAudioRenderer(sequence.Periods, &audio.AudioRendererOptions{
		SampleRate:   options.SampleRate,
		Volume:       options.Volume,
		AmbianceList: options.AmbianceList,
		StatusOutput: lc.appCtx.statusOutput,
	})
	if err != nil {
		return nil, err
	}

	return renderer, nil
}

// WAV generates the WAV file from the loaded sequence.
func (lc *LoadedContext) WAV(outputFile string) error {
	if outputFile == "" {
		return fmt.Errorf("output file cannot be empty")
	}

	renderer, err := lc.generate()
	if err != nil {
		return err
	}

	if err = renderer.RenderWav(outputFile); err != nil {
		return err
	}

	return nil
}

// Stream generates the raw audio stream from the loaded sequence.
func (lc *LoadedContext) Stream(data io.Writer) error {
	renderer, err := lc.generate()
	if err != nil {
		return err
	}

	err = renderer.RenderRaw(data)
	if err != nil {
		return err
	}

	return nil
}
