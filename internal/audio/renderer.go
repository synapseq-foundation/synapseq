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
	"fmt"
	"io"

	amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	audiosync "github.com/synapseq-foundation/synapseq/v4/internal/audio/sync"
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
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
	signals         [t.NumberOfChannels]channelSignalState
	plan            renderPlan
	periods         []t.Period
	waveTables      [4][]int
	noiseGenerator  *NoiseGenerator
	syncEngine      *audiosync.Engine
	effectProcessor *efx.Processor
	ambianceState   *amb.Runtime

	// Embedding options
	*AudioRendererOptions
}

// AudioRendererOptions holds options for the audio renderer
type AudioRendererOptions struct {
	SampleRate   int
	Volume       int
	Ambiance     map[string]string
	StatusOutput io.Writer
	Colors       bool
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

	ambianceState, err := amb.NewRuntime(p, ar.Ambiance, ar.SampleRate, func(paths []string, sampleRate int) (amb.SampleAudio, error) {
		return amb.NewAudio(paths, sampleRate)
	})
	if err != nil {
		return nil, err
	}

	renderer := &AudioRenderer{
		plan:                 compileRenderPlan(p, ar.SampleRate),
		periods:              p,
		waveTables:           wt.Init(),
		noiseGenerator:       NewNoiseGenerator(),
		ambianceState:        ambianceState,
		AudioRendererOptions: ar,
	}
	renderer.syncEngine = audiosync.NewEngine(renderer.SampleRate, func(ch int, periodIdx int, trackType t.TrackType) {
		if renderer.ambianceState == nil {
			return
		}

		renderer.ambianceState.UpdateChannelIndex(ch, periodIdx, trackType)
	})
	renderer.effectProcessor = efx.NewProcessor(renderer.SampleRate, renderer.waveTables)

	return renderer, nil
}

// Render generates the audio and passes buffers to the consume function
func (r *AudioRenderer) Render(consume func(samples []int) error) error {
	defer func() {
		if r.ambianceState != nil {
			r.ambianceState.Close()
		}
	}()

	runtime := newRenderRuntime(r, consume)
	return runtime.run()
}
