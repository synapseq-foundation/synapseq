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

package audio

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	bwav "github.com/gopxl/beep/v2/wav"
	audiooutput "github.com/synapseq-foundation/synapseq/v4/internal/audio/output"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestWAVOutput_Write_EmitsDecodableWAV(ts *testing.T) {
	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackPureTone,
		Carrier:   220,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = 100

	renderer, err := NewAudioRenderer([]t.Period{p0, pEnd}, &AudioRendererOptions{
		SampleRate: 44100,
		Volume:     100,
		Ambiance:   map[string]string{},
	})
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	buffer := &testMemoryWriteSeeker{}
	if err := audiooutput.NewWAVOutput(renderer.SampleRate, audioChannels, audioBitDepth/8, renderer.Render).Write(buffer); err != nil {
		ts.Fatalf("Write failed: %v", err)
	}

	decoded, format, err := bwav.Decode(bytesReaderCloser(buffer.Bytes()))
	if err != nil {
		ts.Fatalf("Decode failed: %v", err)
	}
	defer decoded.Close()

	if int(format.SampleRate) != 44100 {
		ts.Fatalf("unexpected sample rate: got %d", format.SampleRate)
	}
	if format.NumChannels != audioChannels {
		ts.Fatalf("unexpected channel count: got %d", format.NumChannels)
	}

	buf := make([][2]float64, 128)
	n, ok := decoded.Stream(buf)
	if n == 0 && !ok {
		ts.Fatalf("expected decodable audio frames")
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