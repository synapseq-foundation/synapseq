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

package audio

import (
	"io"

	audiooutput "github.com/synapseq-foundation/synapseq/v4/internal/audio/output"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// RenderRaw renders the audio to a raw PCM stream (16-bit little-endian)
func (r *AudioRenderer) RenderRaw(w io.Writer) error {
	rawWriter := audiooutput.NewRawPCMWriter(w, t.BufferSize*audioChannels)

	err := r.Render(func(samples []int) error {
		return rawWriter.WriteSamples(samples)
	})
	if err != nil {
		return err
	}
	return rawWriter.Flush()
}
