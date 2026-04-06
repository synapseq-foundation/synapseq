//go:build !wasm

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

package core

import (
	"fmt"
	"io"

	"github.com/synapseq-foundation/synapseq/v4/internal/audio"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// generate generates the audio renderer based on the loaded sequence
func (lc *LoadedContext) generate() (*audio.AudioRenderer, error) {
	sequence, err := lc.renderableSequence()
	if err != nil {
		return nil, err
	}

	renderer, err := audio.NewAudioRenderer(sequence.Periods, lc.buildAudioRendererOptions(sequence))
	if err != nil {
		return nil, err
	}

	return renderer, nil
}

func (lc *LoadedContext) renderableSequence() (*t.Sequence, error) {
	if lc.sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}
	if lc.sequence.Options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	return lc.sequence, nil
}

func (lc *LoadedContext) buildAudioRendererOptions(sequence *t.Sequence) *audio.AudioRendererOptions {
	return &audio.AudioRendererOptions{
		SampleRate:   sequence.Options.SampleRate,
		Volume:       sequence.Options.Volume,
		Ambiance:     sequence.Options.Ambiance,
		StatusOutput: lc.appCtx.statusOutput,
		Colors:       lc.appCtx.statusColors,
	}
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
