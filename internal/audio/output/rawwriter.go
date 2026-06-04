// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package output

import (
	"bufio"
	"io"

	p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"
)

type RawPCMWriter struct {
	writer *bufio.Writer
	buffer []byte
}

func NewRawPCMWriter(w io.Writer, initialSampleCapacity int) *RawPCMWriter {
	return &RawPCMWriter{
		writer: bufio.NewWriter(w),
		buffer: make([]byte, initialSampleCapacity*2),
	}
}

func (rw *RawPCMWriter) WriteSamples(samples []int) error {
	rw.buffer = p.EncodePCM16LE(rw.buffer, samples)
	_, err := rw.writer.Write(rw.buffer)
	return err
}

func (rw *RawPCMWriter) Flush() error {
	return rw.writer.Flush()
}