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

package ambiance

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
	p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const testBitDepth = 16

type constStreamer struct {
	framesLeft int
	val        float64
}

func (cs *constStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n = len(samples)
	if n > cs.framesLeft {
		n = cs.framesLeft
	}
	for i := 0; i < n; i++ {
		samples[i][0] = cs.val
		samples[i][1] = cs.val
	}
	cs.framesLeft -= n
	ok = cs.framesLeft > 0
	return
}

func (cs *constStreamer) Err() error { return nil }

func writeConstWav(t *testing.T, path string, sampleRate int) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create wav: %v", err)
	}
	defer f.Close()

	format := beep.Format{
		SampleRate:  beep.SampleRate(sampleRate),
		NumChannels: stereoChannels,
		Precision:   testBitDepth / 8,
	}
	cs := &constStreamer{framesLeft: sampleRate, val: float64(1000) / 32768.0}
	if err := bwav.Encode(f, cs, format); err != nil {
		t.Fatalf("encode wav: %v", err)
	}
}

func mustReadWavAll(t *testing.T, path string) ([]int, uint32, int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open wav: %v", err)
	}
	defer f.Close()

	s, fmt, err := bwav.Decode(f)
	if err != nil {
		t.Fatalf("decode wav: %v", err)
	}
	defer s.Close()

	var data []int
	buf := make([][2]float64, 4096)
	for {
		n, ok := s.Stream(buf)
		for i := 0; i < n; i++ {
			l := p.FloatToSample16(buf[i][0])
			r := p.FloatToSample16(buf[i][1])
			data = append(data, l, r)
		}
		if !ok {
			break
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	return data, uint32(fmt.SampleRate), fmt.NumChannels, fmt.Precision * 8
}

func TestAudio_LoadReadAndLoop(t *testing.T) {
	path := filepath.Join("..", "testdata", "noise.wav")
	data, sr, chans, depth := mustReadWavAll(t, path)

	aa, err := NewAudio([]string{path}, int(sr))
	if err != nil {
		t.Fatalf("NewAudio: %v", err)
	}
	defer aa.Close()

	if aa.SampleRate() != int(sr) || aa.Channels() != chans || aa.BitDepth() != depth {
		t.Fatalf("mismatched ambiance props sr=%d ch=%d bd=%d vs file sr=%d ch=%d bd=%d", aa.SampleRate(), aa.Channels(), aa.BitDepth(), sr, chans, depth)
	}

	target := len(data) + 123
	var buf []int
	chunk := aa.BufferSize()
	if chunk <= 0 {
		t.Fatalf("invalid buffer size: %d", chunk)
	}
	tmp := make([]int, chunk)
	total := 0
	for total < target {
		need := target - total
		if need > chunk {
			need = chunk
		}
		n, err := aa.ReadSamplesAt(0, tmp[:need], need)
		if err != nil {
			t.Fatalf("ReadSamplesAt error: %v", err)
		}
		if n != need {
			t.Fatalf("ReadSamplesAt short read: got %d want %d", n, need)
		}
		buf = append(buf, tmp[:need]...)
		total += n
	}

	for i := 0; i < len(data) && i < len(buf); i++ {
		if buf[i] != data[i] {
			t.Fatalf("prefix mismatch at %d: got %d want %d", i, buf[i], data[i])
		}
	}

	if total > len(data) && buf[len(data)] != data[0] {
		t.Fatalf("loop restart mismatch: got %d want %d", buf[len(data)], data[0])
	}
}

func TestAudio_MultipleIndicesHaveIndependentReadPosition(t *testing.T) {
	path := filepath.Join("..", "testdata", "noise.wav")
	data, sr, _, _ := mustReadWavAll(t, path)

	aa, err := NewAudio([]string{path, path}, int(sr))
	if err != nil {
		t.Fatalf("NewAudio: %v", err)
	}
	defer aa.Close()

	const chunk = 1024
	first0 := make([]int, chunk)
	second0 := make([]int, chunk)
	first1 := make([]int, chunk)

	n, err := aa.ReadSamplesAt(0, first0, len(first0))
	if err != nil {
		t.Fatalf("ReadSamplesAt(0) first error: %v", err)
	}
	if n != len(first0) {
		t.Fatalf("ReadSamplesAt(0) first short read: got %d want %d", n, len(first0))
	}

	n, err = aa.ReadSamplesAt(0, second0, len(second0))
	if err != nil {
		t.Fatalf("ReadSamplesAt(0) second error: %v", err)
	}
	if n != len(second0) {
		t.Fatalf("ReadSamplesAt(0) second short read: got %d want %d", n, len(second0))
	}

	n, err = aa.ReadSamplesAt(1, first1, len(first1))
	if err != nil {
		t.Fatalf("ReadSamplesAt(1) first error: %v", err)
	}
	if n != len(first1) {
		t.Fatalf("ReadSamplesAt(1) first short read: got %d want %d", n, len(first1))
	}

	for i := 0; i < len(first1) && i < len(data); i++ {
		if first1[i] != data[i] {
			t.Fatalf("index 1 prefix mismatch at %d: got %d want %d", i, first1[i], data[i])
		}
	}

	for i := 0; i < len(second0); i++ {
		want := data[chunk+i]
		if second0[i] != want {
			t.Fatalf("index 0 continuation mismatch at %d: got %d want %d", i, second0[i], want)
		}
	}
}

func TestAudio_ReadSamplesAt_InvalidIndex(t *testing.T) {
	path := filepath.Join("..", "testdata", "noise.wav")
	_, sr, _, _ := mustReadWavAll(t, path)

	aa, err := NewAudio([]string{path, path}, int(sr))
	if err != nil {
		t.Fatalf("NewAudio: %v", err)
	}
	defer aa.Close()

	buf := make([]int, 128)

	if _, err := aa.ReadSamplesAt(-1, buf, len(buf)); err == nil {
		t.Fatalf("expected error for negative ambiance index")
	}

	if _, err := aa.ReadSamplesAt(2, buf, len(buf)); err == nil {
		t.Fatalf("expected error for out-of-range ambiance index")
	}
}

func TestAudio_DisabledAndClose(t *testing.T) {
	aa, err := NewAudio(nil, 44100)
	if err != nil {
		t.Fatalf("NewAudio empty: %v", err)
	}
	buf := make([]int, 256)
	n, err := aa.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt disabled error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamplesAt disabled count: got %d want %d", n, len(buf))
	}
	for i, v := range buf {
		if v != 0 {
			t.Fatalf("disabled should fill zeros at %d: %d", i, v)
		}
	}

	if err := aa.Close(); err != nil {
		t.Fatalf("Close disabled: %v", err)
	}
	_, err = aa.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt after close error: %v", err)
	}
	for i, v := range buf {
		if v != 0 {
			t.Fatalf("after close should fill zeros at %d: %d", i, v)
		}
	}
}

func TestAudio_RemoteWAV(t *testing.T) {
	path := filepath.Join("..", "testdata", "noise.wav")
	_, sr, _, _ := mustReadWavAll(t, path)
	wavData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test WAV: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		_, _ = w.Write(wavData)
	}))
	defer server.Close()

	aa, err := NewAudio([]string{server.URL}, int(sr))
	if err != nil {
		t.Fatalf("NewAudio remote: %v", err)
	}
	defer aa.Close()

	if aa.CachedData() == nil {
		t.Fatalf("expected cachedData to be populated")
	}
	if len(aa.CachedData()) != 1 {
		t.Fatalf("cached tracks mismatch: got %d want 1", len(aa.CachedData()))
	}
	if len(aa.CachedData()[0]) != len(wavData) {
		t.Fatalf("cached data size mismatch: got %d want %d", len(aa.CachedData()[0]), len(wavData))
	}

	buf := make([]int, 1024)
	n, err := aa.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt remote error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamplesAt remote count: got %d want %d", n, len(buf))
	}

	hasNonZero := false
	for _, v := range buf {
		if v != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Fatalf("expected non-zero samples from remote WAV")
	}
}

func TestAudio_Remote10MBLimit(ts *testing.T) {
	const size = t.MaxWavFileSize + 2*1024*1024
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		header := make([]byte, 44)
		copy(header[0:4], "RIFF")
		copy(header[8:12], "WAVE")
		copy(header[12:16], "fmt ")
		header[16] = 16
		header[20] = 1
		header[22] = 2
		header[24] = 0x44
		header[25] = 0xac
		header[28] = 0x10
		header[29] = 0xb1
		header[30] = 0x02
		header[32] = 4
		header[34] = 16
		copy(header[36:40], "data")
		dataSize := size - 44
		header[40] = byte(dataSize)
		header[41] = byte(dataSize >> 8)
		header[42] = byte(dataSize >> 16)
		header[43] = byte(dataSize >> 24)

		_, _ = w.Write(header)
		chunk := make([]byte, 1024*1024)
		for i := 0; i < size-44; i += len(chunk) {
			remaining := size - 44 - i
			if remaining < len(chunk) {
				_, _ = w.Write(chunk[:remaining])
			} else {
				_, _ = w.Write(chunk)
			}
		}
	}))
	defer server.Close()

	aa, err := NewAudio([]string{server.URL}, 44100)
	if err != nil {
		ts.Fatalf("NewAudio 10MB limit: %v", err)
	}
	defer aa.Close()

	if len(aa.CachedData()) != 1 {
		ts.Fatalf("expected one cached track, got %d", len(aa.CachedData()))
	}
	if len(aa.CachedData()[0]) != t.MaxWavFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", t.MaxWavFileSize, len(aa.CachedData()[0]))
	}
}

func TestAudio_Local10MBLimit(ts *testing.T) {
	tmpDir := ts.TempDir()
	path := filepath.Join(tmpDir, "large.wav")

	f, err := os.Create(path)
	if err != nil {
		ts.Fatalf("failed to create temp file: %v", err)
	}

	const size = t.MaxWavFileSize + 2*1024*1024
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	header[16] = 16
	header[20] = 1
	header[22] = 2
	header[24] = 0x44
	header[25] = 0xac
	header[28] = 0x10
	header[29] = 0xb1
	header[30] = 0x02
	header[32] = 4
	header[34] = 16
	copy(header[36:40], "data")
	dataSize := size - 44
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	if _, err := f.Write(header); err != nil {
		ts.Fatalf("failed to write header: %v", err)
	}

	chunk := make([]byte, 1024*1024)
	for i := 0; i < size-44; i += len(chunk) {
		remaining := size - 44 - i
		if remaining < len(chunk) {
			if _, err := f.Write(chunk[:remaining]); err != nil {
				ts.Fatalf("failed to write data: %v", err)
			}
		} else {
			if _, err := f.Write(chunk); err != nil {
				ts.Fatalf("failed to write data: %v", err)
			}
		}
	}
	_ = f.Close()

	aa, err := NewAudio([]string{path}, 44100)
	if err != nil {
		ts.Fatalf("NewAudio local 10MB limit: %v", err)
	}
	defer aa.Close()

	if len(aa.CachedData()) != 1 {
		ts.Fatalf("expected one cached track, got %d", len(aa.CachedData()))
	}
	if len(aa.CachedData()[0]) != t.MaxWavFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", t.MaxWavFileSize, len(aa.CachedData()[0]))
	}
}

func TestAudio_InvalidPath(t *testing.T) {
	if _, err := NewAudio([]string{filepath.Join("..", "testdata", "missing.wav")}, 44100); err == nil {
		t.Fatalf("expected error for missing ambiance file")
	}
}

func TestAudio_ResamplesMismatchedSampleRate(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "mismatch.wav")
	writeConstWav(t, path, 48000)

	aa, err := NewAudio([]string{path}, 44100)
	if err != nil {
		t.Fatalf("NewAudio resample: %v", err)
	}
	defer aa.Close()

	if aa.SampleRate() != 44100 {
		t.Fatalf("expected resampled sample rate 44100, got %d", aa.SampleRate())
	}
	if aa.Channels() != stereoChannels {
		t.Fatalf("expected %d channels, got %d", stereoChannels, aa.Channels())
	}

	buf := make([]int, 1024)
	n, err := aa.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt resampled wav: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamplesAt short read after resample: got %d want %d", n, len(buf))
	}

	hasNonZero := false
	for _, sample := range buf {
		if sample != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Fatalf("expected non-zero samples from resampled wav")
	}

	reader := bytes.NewReader(aa.CachedData()[0])
	stream, format, err := bwav.Decode(reader)
	if err != nil {
		t.Fatalf("decode resampled cache: %v", err)
	}
	defer stream.Close()

	if int(format.SampleRate) != 44100 {
		t.Fatalf("expected cached wav sample rate 44100, got %d", format.SampleRate)
	}
}