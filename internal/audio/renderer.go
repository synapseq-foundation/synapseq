/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"fmt"
	"io"
	"math"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

const (
	audioChannels = 2      // Stereo
	audioBitDepth = 16     // 16-bit audio
	audioBitShift = 16     // 16 Bit shift
	audioMaxValue = 32767  // 2^15 - 1
	audioMinValue = -32768 // -2^15
)

// AudioRenderer handle audio generation
type AudioRenderer struct {
	channels        [t.NumberOfChannels]t.Channel
	periods         []t.Period
	waveTables      [4][]int
	noiseGenerator  *NoiseGenerator
	backgroundAudio *BackgroundAudio

	// Reusable buffer to avoid allocating every mix() call
	backgroundSamplesByIndex [][]int
	// Track indices that have active background audio (for optimization)
	activeBGIndices []int
	// Mask to track which background audio tracks are currently active
	activeBGMask []bool

	// Cache for the current background index of each channel to optimize lookups during sync
	channelBGIndex [t.NumberOfChannels]int

	// cache for the current period's background names for each channel to optimize lookups during sync
	periodBGStart [][]int

	// Embedding options
	*AudioRendererOptions
}

// AudioRendererOptions holds options for the audio renderer
type AudioRendererOptions struct {
	SampleRate     int
	Volume         int
	GainLevel      t.GainLevel
	BackgroundList map[string]string
	StatusOutput   io.Writer
}

// NewAudioRenderer creates a new AudioRenderer instance
func NewAudioRenderer(p []t.Period, ar *AudioRendererOptions) (*AudioRenderer, error) {
	if ar == nil {
		return nil, fmt.Errorf("audio renderer options cannot be nil")
	}

	if ar.SampleRate <= 0 {
		return nil, fmt.Errorf("invalid sample rate: %d", ar.SampleRate)
	}

	if ar.Volume < 0 || ar.Volume > 100 {
		return nil, fmt.Errorf("volume must be between 0 and 100, got %d", ar.Volume)
	}

	if len(p) == 0 {
		return nil, fmt.Errorf("no periods defined in the sequence")
	}

	bgPaths, bgNameToIndex, err := buildBackgroundIndex(ar.BackgroundList)
	if err != nil {
		return nil, err
	}

	periodBGStart, err := precomputePeriodBGStart(p, bgNameToIndex)
	if err != nil {
		return nil, err
	}

	backgroundAudio, err := NewBackgroundAudio(bgPaths, ar.SampleRate)
	if err != nil {
		return nil, err
	}

	renderer := &AudioRenderer{
		periods:                  p,
		waveTables:               InitWaveformTables(),
		noiseGenerator:           NewNoiseGenerator(),
		backgroundAudio:          backgroundAudio,
		backgroundSamplesByIndex: make([][]int, len(bgPaths)),
		activeBGIndices:          make([]int, 0, t.NumberOfChannels),
		activeBGMask:             make([]bool, len(bgPaths)),
		periodBGStart:            periodBGStart,
		AudioRendererOptions:     ar,
	}

	for i := range renderer.channelBGIndex {
		renderer.channelBGIndex[i] = -1
	}

	return renderer, nil
}

// Render generates the audio and passes buffers to the consume function
func (r *AudioRenderer) Render(consume func(samples []int) error) error {
	// Ensure background audio file is closed if opened
	defer func() {
		if r.backgroundAudio != nil {
			r.backgroundAudio.Close()
		}
	}()

	endMs := r.periods[len(r.periods)-1].Time
	totalFrames := int64(math.Round(float64(endMs) * float64(r.SampleRate) / 1000.0))
	chunkFrames := int64(t.BufferSize)
	framesWritten := int64(0)

	var statusReporter *StatusReporter
	if r.StatusOutput != nil {
		statusReporter = NewStatusReporter(r.StatusOutput)
		defer statusReporter.FinalStatus()
	}

	// Stereo: left + right
	samples := make([]int, t.BufferSize*audioChannels)
	periodIdx := 0

	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.SampleRate))
		// Find the correct period for the current time
		for periodIdx+1 < len(r.periods) && currentTimeMs >= r.periods[periodIdx+1].Time {
			periodIdx++
		}

		r.sync(currentTimeMs, periodIdx)
		r.collectActiveBackgroundIndices()
		r.prepareBackgroundBuffers()

		if statusReporter != nil {
			statusReporter.CheckPeriodChange(r, periodIdx)
		}

		data := r.mix(samples)

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			// stereo interleaved
			data = data[:remain*audioChannels]
		}

		if consume != nil {
			if err := consume(data); err != nil {
				return fmt.Errorf("failed to consume audio buffer: %w", err)
			}
		}

		framesWritten += framesToWrite

		if statusReporter != nil && statusReporter.ShouldUpdateStatus() {
			statusReporter.DisplayStatus(r, currentTimeMs)
		}
	}

	return nil
}
