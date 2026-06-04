// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
