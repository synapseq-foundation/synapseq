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

package output

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	bwav "github.com/gopxl/beep/v2/wav"
)

func TestWAVOutputWrite_EmitsDecodableWAV(ts *testing.T) {
	const (
		sampleRate = 44100
		channels   = 2
		precision  = 2
	)

	render := func(consume func(samples []int) error) error {
		samples := make([]int, 128*channels)
		for i := 0; i < len(samples); i += 2 {
			samples[i] = 1000
			samples[i+1] = -1000
		}
		return consume(samples)
	}

	buffer := &testMemoryWriteSeeker{}
	if err := NewWAVOutput(sampleRate, channels, precision, render).Write(buffer); err != nil {
		ts.Fatalf("Write failed: %v", err)
	}

	decoded, format, err := bwav.Decode(bytesReaderCloser(buffer.Bytes()))
	if err != nil {
		ts.Fatalf("Decode failed: %v", err)
	}
	defer decoded.Close()

	if int(format.SampleRate) != sampleRate {
		ts.Fatalf("unexpected sample rate: got %d", format.SampleRate)
	}
	if format.NumChannels != channels {
		ts.Fatalf("unexpected channel count: got %d", format.NumChannels)
	}

	buf := make([][2]float64, 128)
	n, ok := decoded.Stream(buf)
	if n == 0 && !ok {
		ts.Fatalf("expected decodable audio frames")
	}

	hasSignal := false
	for i := 0; i < n; i++ {
		if buf[i][0] != 0 || buf[i][1] != 0 {
			hasSignal = true
			break
		}
	}
	if !hasSignal {
		ts.Fatalf("expected non-zero decoded samples")
	}
}

type byteReadSeekCloser struct {
	*bytes.Reader
}

func bytesReaderCloser(data []byte) *byteReadSeekCloser {
	return &byteReadSeekCloser{Reader: bytes.NewReader(data)}
}

type testMemoryWriteSeeker struct {
	data []byte
	pos  int
}

func (m *testMemoryWriteSeeker) Write(p []byte) (int, error) {
	end := m.pos + len(p)
	if end > len(m.data) {
		grown := make([]byte, end)
		copy(grown, m.data)
		m.data = grown
	}
	copy(m.data[m.pos:end], p)
	m.pos = end
	return len(p), nil
}

func (m *testMemoryWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var next int
	switch whence {
	case io.SeekStart:
		next = int(offset)
	case io.SeekCurrent:
		next = m.pos + int(offset)
	case io.SeekEnd:
		next = len(m.data) + int(offset)
	default:
		return 0, fmt.Errorf("invalid seek whence: %d", whence)
	}
	if next < 0 {
		return 0, fmt.Errorf("invalid seek position: %d", next)
	}
	m.pos = next
	return int64(next), nil
}

func (m *testMemoryWriteSeeker) Bytes() []byte {
	return append([]byte(nil), m.data...)
}

func (brc *byteReadSeekCloser) Close() error {
	return nil
}