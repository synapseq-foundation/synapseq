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
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

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

func TestAudioRenderer_RenderWav_Integration(ts *testing.T) {
	// Create test periods (2 seconds total) with different track types
	var p0, p1, p2 t.Period

	// Period 0: 0-500ms binaural beat
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	// Period 1: 500-1000ms monaural beat with interpolation
	p1.Time = 500
	p1.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p1.TrackEnd[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   280,
		Resonance: 12,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformTriangle,
	}

	// Period 2: 1000-2000ms with multiple tracks (noise + isochronic)
	p2.Time = 1000
	p2.TrackStart[0] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p2.TrackStart[1] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   40,
		Resonance: 2.5,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSawtooth,
	}
	p2.TrackEnd[0] = p2.TrackStart[0]
	p2.TrackEnd[1] = p2.TrackStart[1]

	// End period at 2s
	var pEnd t.Period
	pEnd.Time = 2000

	periods := []t.Period{p0, p1, p2, pEnd}

	options := &AudioRendererOptions{
		SampleRate: 44100,
		Volume:     80,
		Ambiance:   map[string]string{},
	}

	renderer, err := NewAudioRenderer(periods, options)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	// Create temp directory and output file
	tempDir := ts.TempDir()
	outPath := filepath.Join(tempDir, "test_output.wav")

	// Render to WAV
	if err := renderer.RenderWav(outPath); err != nil {
		ts.Fatalf("RenderWav failed: %v", err)
	}

	// Validate the generated WAV file
	file, err := os.Open(outPath)
	if err != nil {
		ts.Fatalf("Failed to open generated WAV: %v", err)
	}
	defer file.Close()

	s, f, err := bwav.Decode(file)
	if err != nil {
		ts.Fatalf("Decode failed: %v", err)
	}
	defer s.Close()

	if int(f.SampleRate) != options.SampleRate {
		ts.Fatalf("Sample rate mismatch: got %d, want %d", f.SampleRate, options.SampleRate)
	}
	if f.NumChannels != audioChannels {
		ts.Fatalf("Channel count mismatch: got %d, want %d", f.NumChannels, audioChannels)
	}
	if f.Precision*8 != audioBitDepth {
		ts.Fatalf("Bit depth mismatch: got %d, want %d", f.Precision*8, audioBitDepth)
	}

	// Verify file size is reasonable for 2 seconds of audio
	stat, err := file.Stat()
	if err != nil {
		ts.Fatalf("Failed to stat file: %v", err)
	}
	expectedMinSize := int64(2 * options.SampleRate * audioChannels * audioBitDepth / 8)
	if stat.Size() < expectedMinSize/2 {
		ts.Fatalf("Generated file too small: got %d bytes, expected at least %d", stat.Size(), expectedMinSize/2)
	}

	// Read and verify some audio data exists (non-zero samples)
	if err := s.Seek(0); err != nil {
		ts.Fatalf("Seek to start failed: %v", err)
	}
	foundNonZero := false
	buf := make([][2]float64, 4096)
	for i := 0; i < 8 && !foundNonZero; i++ {
		n, ok := s.Stream(buf)
		for j := 0; j < n; j++ {
			if buf[j][0] != 0 || buf[j][1] != 0 {
				foundNonZero = true
				break
			}
		}
		if !ok {
			break
		}
	}
	if err := s.Err(); err != nil {
		ts.Fatalf("Stream error: %v", err)
	}
	if !foundNonZero {
		ts.Fatalf("All samples are zero - audio generation may be broken")
	}
}

func TestAudioRenderer_RenderWav_WithAmbiance(ts *testing.T) {
	// Create a simple test WAV file as ambiance
	tempDir := ts.TempDir()
	bgPath := filepath.Join(tempDir, "ambiance.wav")

	bgFile, err := os.Create(bgPath)
	if err != nil {
		ts.Fatalf("Failed to create ambiance file: %v", err)
	}
	defer bgFile.Close()

	const sr = 44100
	format := beep.Format{SampleRate: beep.SampleRate(sr), NumChannels: audioChannels, Precision: audioBitDepth / 8}
	val := float64(1000) / 32768.0
	cs := &constStreamer{framesLeft: sr, val: val}
	if err := bwav.Encode(bgFile, cs, format); err != nil {
		ts.Fatalf("Failed to write ambiance: %v", err)
	}
	if _, err := bgFile.Seek(0, 0); err != nil {
		ts.Fatalf("Failed to rewind ambiance file: %v", err)
	}

	// Create test period with ambiance track
	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:         t.TrackAmbiance,
		AmbianceName: "bg",
		Amplitude:    t.AmplitudePercentToRaw(30),
		Waveform:     t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	pEnd.Time = 1000 // 1 second
	periods := []t.Period{p0, pEnd}

	options := &AudioRendererOptions{
		SampleRate: 44100,
		Volume:     100,
		Ambiance: map[string]string{
			"bg": bgPath,
		},
	}

	renderer, err := NewAudioRenderer(periods, options)
	if err != nil {
		ts.Fatalf("NewAudioRenderer with ambiance failed: %v", err)
	}

	outPath := filepath.Join(tempDir, "test_with_bg.wav")
	if err := renderer.RenderWav(outPath); err != nil {
		ts.Fatalf("RenderWav with ambiance failed: %v", err)
	}

	// Basic validation
	if _, err := os.Stat(outPath); err != nil {
		ts.Fatalf("Output file not created: %v", err)
	}
}

func TestAudioRenderer_Render_CallbacksAndSizes(ts *testing.T) {
	sr := 44100

	endMs := 1234
	totalFrames := int64(math.Round(float64(endMs) * float64(sr) / 1000.0))

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   220,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs

	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     80,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	var lens []int
	calls := 0

	consume := func(data []int) error {
		lens = append(lens, len(data))
		calls++
		return nil
	}

	if err := r.Render(consume); err != nil {
		ts.Fatalf("Render failed: %v", err)
	}

	chunk := int64(t.BufferSize)
	full := int(totalFrames / chunk)
	rem := int(totalFrames % chunk)

	expected := make([]int, 0, full+1)
	for i := 0; i < full; i++ {
		expected = append(expected, t.BufferSize*audioChannels)
	}
	if rem > 0 {
		expected = append(expected, rem*audioChannels)
	}

	if calls != len(expected) {
		ts.Fatalf("Expected %d callbacks, got %d", len(expected), calls)
	}
	for i := range expected {
		if lens[i] != expected[i] {
			ts.Fatalf("Chunk %d size mismatch: got %d, want %d", i, lens[i], expected[i])
		}
	}
}

func TestAudioRenderer_Render_PropagatesError(ts *testing.T) {
	sr := 44100
	endMs := 2000

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   200,
		Resonance: 7,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs
	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     100,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	targetErr := errors.New("sink failure")
	consume := func(_ []int) error {
		return targetErr
	}

	err = r.Render(consume)
	if err == nil {
		ts.Fatalf("Expected error from consumer, got nil")
	}
	if !errors.Is(err, targetErr) {
		ts.Fatalf("Expected wrapped target error, got: %v", err)
	}
}

func TestAudioRenderer_Render_NilConsumer(ts *testing.T) {
	sr := 44100
	endMs := 2000

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   10,
		Resonance: 2,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs
	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     90,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	if err := r.Render(nil); err != nil {
		ts.Fatalf("Render with nil consumer failed: %v", err)
	}
}

type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("sink failure")
}

func TestRenderRaw_PropagatesWriteError(ts *testing.T) {
	sr := 44100
	endMs := 50

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   10,
		Resonance: 2,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs

	opts := &AudioRendererOptions{
		SampleRate:   sr,
		Volume:       80,
		StatusOutput: os.Stderr,
	}

	r, err := NewAudioRenderer([]t.Period{p0, pEnd}, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	err = r.RenderRaw(&failingWriter{})
	if err == nil {
		ts.Fatalf("expected error from writer, got nil")
	}
}
