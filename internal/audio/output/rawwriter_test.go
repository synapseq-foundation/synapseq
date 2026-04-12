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

package output

import (
	"bytes"
	"testing"
)

func TestRawPCMWriter_WriteSamples_EncodesLittleEndianPCM16(ts *testing.T) {
	var out bytes.Buffer
	writer := NewRawPCMWriter(&out, 4)

	if err := writer.WriteSamples([]int{40000, -40000, 1, -1}); err != nil {
		ts.Fatalf("WriteSamples failed: %v", err)
	}
	if err := writer.Flush(); err != nil {
		ts.Fatalf("Flush failed: %v", err)
	}

	expected := []byte{
		0xff, 0x7f,
		0x00, 0x80,
		0x01, 0x00,
		0xff, 0xff,
	}

	if !bytes.Equal(out.Bytes(), expected) {
		ts.Fatalf("unexpected PCM output: got %v want %v", out.Bytes(), expected)
	}
}