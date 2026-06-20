// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"fmt"
	"os"

	audiooutput "github.com/synapseq-foundation/synapseq/v4/internal/audio/output"
)

// RenderWav renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderWav(outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	return audiooutput.NewWAVOutput(r.SampleRate, audioChannels, audioBitDepth/8, r.Render).Write(out)
}
