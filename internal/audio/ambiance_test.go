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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	bwav "github.com/gopxl/beep/v2/wav"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

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

	// Stream all frames and convert to interleaved int16 samples to match AmbianceAudio
	var data []int
	const scale = 32768.0 // 2^15 for 16-bit
	buf := make([][2]float64, 4096)
	for {
		n, ok := s.Stream(buf)
		for i := 0; i < n; i++ {
			l := int(buf[i][0] * scale)
			r := int(buf[i][1] * scale)
			if l > audioMaxValue {
				l = audioMaxValue
			}
			if l < audioMinValue {
				l = audioMinValue
			}
			if r > audioMaxValue {
				r = audioMaxValue
			}
			if r < audioMinValue {
				r = audioMinValue
			}
			data = append(data, l, r)
		}
		if !ok {
			break
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	// Return data, sample rate, channels, bit depth
	return data, uint32(fmt.SampleRate), fmt.NumChannels, fmt.Precision * 8
}

func TestAmbianceAudio_LoadReadAndLoop(t *testing.T) {
	path := filepath.Join("testdata", "noise.wav")
	data, sr, chans, depth := mustReadWavAll(t, path)

	aa, err := NewAmbianceAudio([]string{path}, int(sr))
	if err != nil {
		t.Fatalf("NewAmbianceAudio: %v", err)
	}
	defer aa.Close()

	if aa.sampleRate != int(sr) || aa.channels != chans || aa.bitDepth != depth {
		t.Fatalf("mismatched ambiance props sr=%d ch=%d bd=%d vs file sr=%d ch=%d bd=%d", aa.sampleRate, aa.channels, aa.bitDepth, sr, chans, depth)
	}

	// Force looping at least once, reading in chunks up to aa.bufferSize
	target := len(data) + 123
	var buf []int
	chunk := aa.bufferSize
	if chunk <= 0 {
		t.Fatalf("invalid aa.bufferSize: %d", chunk)
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

	// Prefix must match the original file data
	for i := 0; i < len(data) && i < len(buf); i++ {
		if buf[i] != data[i] {
			t.Fatalf("prefix mismatch at %d: got %d want %d", i, buf[i], data[i])
		}
	}

	// After exactly len(data) samples, sequence should restart at beginning
	if total > len(data) {
		if buf[len(data)] != data[0] {
			t.Fatalf("loop restart mismatch: got %d want %d", buf[len(data)], data[0])
		}
	}
}

func TestAmbianceAudio_MultipleIndicesHaveIndependentReadPosition(t *testing.T) {
	path := filepath.Join("testdata", "noise.wav")
	data, sr, _, _ := mustReadWavAll(t, path)

	aa, err := NewAmbianceAudio([]string{path, path}, int(sr))
	if err != nil {
		t.Fatalf("NewAmbianceAudio: %v", err)
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

	// Index 1 must start from beginning, independent from index 0 progress.
	for i := 0; i < len(first1) && i < len(data); i++ {
		if first1[i] != data[i] {
			t.Fatalf("index 1 prefix mismatch at %d: got %d want %d", i, first1[i], data[i])
		}
	}

	// Second read from index 0 must continue from where first read stopped.
	for i := 0; i < len(second0); i++ {
		want := data[chunk+i]
		if second0[i] != want {
			t.Fatalf("index 0 continuation mismatch at %d: got %d want %d", i, second0[i], want)
		}
	}
}

func TestAmbianceAudio_ReadSamplesAt_InvalidIndex(t *testing.T) {
	path := filepath.Join("testdata", "noise.wav")
	_, sr, _, _ := mustReadWavAll(t, path)

	aa, err := NewAmbianceAudio([]string{path, path}, int(sr))
	if err != nil {
		t.Fatalf("NewAmbianceAudio: %v", err)
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

func TestAmbianceAudio_DisabledAndClose(t *testing.T) {
	aa, err := NewAmbianceAudio(nil, 44100)
	if err != nil {
		t.Fatalf("NewAmbianceAudio empty: %v", err)
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

	// Closing should keep it disabled and safe
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

func TestAmbianceAudio_RemoteWAV(t *testing.T) {
	// Read local WAV file to serve
	path := filepath.Join("testdata", "noise.wav")
	_, sr, _, _ := mustReadWavAll(t, path)
	wavData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test WAV: %v", err)
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		w.Write(wavData)
	}))
	defer server.Close()

	// Load from remote URL
	aa, err := NewAmbianceAudio([]string{server.URL}, int(sr))
	if err != nil {
		t.Fatalf("NewAmbianceAudio remote: %v", err)
	}
	defer aa.Close()

	// Verify cache was populated
	if aa.cachedData == nil {
		t.Fatalf("expected cachedData to be populated")
	}
	if len(aa.cachedData) != 1 {
		t.Fatalf("cached tracks mismatch: got %d want 1", len(aa.cachedData))
	}
	if len(aa.cachedData[0]) != len(wavData) {
		t.Fatalf("cached data size mismatch: got %d want %d", len(aa.cachedData[0]), len(wavData))
	}

	// Read some samples to verify it works
	buf := make([]int, 1024)
	n, err := aa.ReadSamplesAt(0, buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamplesAt remote error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamplesAt remote count: got %d want %d", n, len(buf))
	}

	// Verify non-zero samples
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

func TestAmbianceAudio_Remote10MBLimit(ts *testing.T) {
	// Create a server that serves more than the configured WAV max size
	const size = t.MaxWavFileSize + 2*1024*1024
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		// Write a simple WAV header (44 bytes) + data
		header := make([]byte, 44)
		copy(header[0:4], "RIFF")
		copy(header[8:12], "WAVE")
		copy(header[12:16], "fmt ")
		// fmt chunk size
		header[16] = 16
		// PCM format (1)
		header[20] = 1
		// 2 channels
		header[22] = 2
		// 44100 sample rate
		header[24] = 0x44
		header[25] = 0xac
		// byte rate
		header[28] = 0x10
		header[29] = 0xb1
		header[30] = 0x02
		// block align
		header[32] = 4
		// bits per sample
		header[34] = 16
		// data chunk
		copy(header[36:40], "data")
		// data size (size - 44)
		dataSize := size - 44
		header[40] = byte(dataSize)
		header[41] = byte(dataSize >> 8)
		header[42] = byte(dataSize >> 16)
		header[43] = byte(dataSize >> 24)

		w.Write(header)
		// Write more data to exceed 10MB
		chunk := make([]byte, 1024*1024) // 1MB chunks
		for i := 0; i < size-44; i += len(chunk) {
			remaining := size - 44 - i
			if remaining < len(chunk) {
				w.Write(chunk[:remaining])
			} else {
				w.Write(chunk)
			}
		}
	}))
	defer server.Close()

	aa, err := NewAmbianceAudio([]string{server.URL}, 44100)
	if err != nil {
		ts.Fatalf("NewAmbianceAudio 10MB limit: %v", err)
	}
	defer aa.Close()

	// Verify that data was capped at the configured max size
	if len(aa.cachedData) != 1 {
		ts.Fatalf("expected one cached track, got %d", len(aa.cachedData))
	}
	if len(aa.cachedData[0]) != t.MaxWavFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", t.MaxWavFileSize, len(aa.cachedData[0]))
	}
}

func TestAmbianceAudio_Local10MBLimit(ts *testing.T) {
	// Create a temporary WAV file larger than the configured WAV max size
	tmpDir := ts.TempDir()
	path := filepath.Join(tmpDir, "large.wav")

	f, err := os.Create(path)
	if err != nil {
		ts.Fatalf("failed to create temp file: %v", err)
	}

	const size = t.MaxWavFileSize + 2*1024*1024
	// Write WAV header
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

	// Write data
	chunk := make([]byte, 1024*1024) // 1MB chunks
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
	f.Close()

	aa, err := NewAmbianceAudio([]string{path}, 44100)
	if err != nil {
		ts.Fatalf("NewAmbianceAudio local 10MB limit: %v", err)
	}
	defer aa.Close()

	// Verify that data was capped at the configured max size
	if len(aa.cachedData) != 1 {
		ts.Fatalf("expected one cached track, got %d", len(aa.cachedData))
	}
	if len(aa.cachedData[0]) != t.MaxWavFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", t.MaxWavFileSize, len(aa.cachedData[0]))
	}
}

func TestAmbianceAudio_InvalidPath(t *testing.T) {
	if _, err := NewAmbianceAudio([]string{filepath.Join("testdata", "missing.wav")}, 44100); err == nil {
		t.Fatalf("expected error for missing ambiance file")
	}
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
