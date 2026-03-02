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

package audio

import (
	"fmt"
	"os"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
)

// RenderWav renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderWav(outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	streamer := newRendererStreamer(r)
	format := beep.Format{
		SampleRate:  beep.SampleRate(r.SampleRate),
		NumChannels: audioChannels,
		Precision:   audioBitDepth / 8,
	}

	if err := bwav.Encode(out, streamer, format); err != nil {
		return err
	}
	if streamer.err != nil {
		return streamer.err
	}

	return nil
}
