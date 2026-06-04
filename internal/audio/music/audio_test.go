// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package music

import (
	"os"
	"path/filepath"
	"testing"

	bmp3 "github.com/gopxl/beep/v2/mp3"
	p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"
)

func mustReadMP3All(t *testing.T, path string) ([]int, uint32) {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open mp3: %v", err)
	}
	defer f.Close()

	stream, format, err := bmp3.Decode(f)
	if err != nil {
		t.Fatalf("decode mp3: %v", err)
	}
	defer stream.Close()

	var data []int
	buf := make([][2]float64, 4096)
	for {
		n, ok := stream.Stream(buf)
		for i := 0; i < n; i++ {
			data = append(data, p.FloatToSample16(buf[i][0]), p.FloatToSample16(buf[i][1]))
		}
		if !ok {
			break
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	return data, uint32(format.SampleRate)
}

func TestAudioReturnsSilenceAfterEOF(t *testing.T) {
	path := filepath.Join("..", "testdata", "short.mp3")
	data, sampleRate := mustReadMP3All(t, path)

	audio, err := NewAudio([]string{path}, int(sampleRate))
	if err != nil {
		t.Fatalf("NewAudio short mp3: %v", err)
	}
	defer audio.Close()

	buf := make([]int, len(data)+128)
	n, err := audio.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt short music mp3: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamplesAt short music mp3 count: got %d want %d", n, len(buf))
	}
	for i := len(data); i < len(buf); i++ {
		if buf[i] != 0 {
			t.Fatalf("expected silence after music EOF at %d, got %d", i, buf[i])
		}
	}
}
