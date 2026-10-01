// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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