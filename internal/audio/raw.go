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
	"bufio"
	"io"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// RenderRaw renders the audio to a raw PCM stream (16-bit little-endian)
func (r *AudioRenderer) RenderRaw(w io.Writer) error {
	bw := bufio.NewWriter(w)
	// 2 bytes per sample (16-bit)
	out := make([]byte, t.BufferSize*audioChannels*2)

	err := r.Render(func(samples []int) error {
		need := len(samples) * 2
		if cap(out) < need {
			out = make([]byte, need)
		}
		b := out[:need]

		j := 0
		for _, s := range samples {
			if s > audioMaxValue {
				s = audioMaxValue
			} else if s < audioMinValue {
				s = audioMinValue
			}
			v := int16(s)
			b[j] = byte(v)        // LSB
			b[j+1] = byte(v >> 8) // MSB
			j += 2
		}
		_, err := bw.Write(b)
		return err
	})
	if err != nil {
		return err
	}
	return bw.Flush()
}
