// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"fmt"
	"io"

	amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"
	efx "github.com/synapseq-foundation/synapseq/v4/internal/audio/effects"
	mus "github.com/synapseq-foundation/synapseq/v4/internal/audio/music"
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
	waveTables      [][]int
	noiseGenerator  *NoiseGenerator
	syncEngine      *audiosync.Engine
	effectProcessor *efx.Processor
	ambianceState   *amb.Runtime
	musicState      *mus.Runtime

	// Embedding options
	*AudioRendererOptions
}

// AudioRendererOptions holds options for the audio renderer
type AudioRendererOptions struct {
	SampleRate   int
	Volume       int
	Ambiance     map[string]string
	Music        map[string]string
	StatusOutput io.Writer
	Colors       bool
	Waveforms    []t.WaveformDefinition
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

	waveforms, err := wt.Compile(ar.Waveforms)
	if err != nil {
		return nil, err
	}
	if err := validateRendererWaveforms(p, waveforms); err != nil {
		return nil, err
	}

	ambianceState, err := amb.NewRuntime(p, ar.Ambiance, ar.SampleRate, func(paths []string, sampleRate int) (amb.SampleAudio, error) {
		return amb.NewAudio(paths, sampleRate)
	})
	if err != nil {
		return nil, err
	}
	musicState, err := mus.NewRuntime(p, ar.Music, ar.SampleRate, func(paths []string, sampleRate int) (mus.SampleAudio, error) {
		return mus.NewAudio(paths, sampleRate)
	})
	if err != nil {
		if ambianceState != nil {
			ambianceState.Close()
		}
		return nil, err
	}

	renderer := &AudioRenderer{
		plan:                 compileRenderPlanWithWaveforms(p, ar.SampleRate, waveforms),
		periods:              p,
		waveTables:           waveforms.Tables,
		noiseGenerator:       NewNoiseGenerator(),
		ambianceState:        ambianceState,
		musicState:           musicState,
		AudioRendererOptions: ar,
	}
	renderer.syncEngine = audiosync.NewEngine(renderer.SampleRate, func(ch int, periodIdx int, trackType t.TrackType) {
		if renderer.ambianceState != nil {
			renderer.ambianceState.UpdateChannelIndex(ch, periodIdx, trackType)
		}
		if renderer.musicState != nil {
			renderer.musicState.UpdateChannelIndex(ch, periodIdx, trackType)
		}
	})
	renderer.effectProcessor = efx.NewProcessor(renderer.SampleRate, renderer.waveTables)

	return renderer, nil
}

func validateRendererWaveforms(periods []t.Period, registry *wt.Registry) error {
	for periodIndex := range periods {
		for channel := range t.NumberOfChannels {
			tracks := []t.Track{
				periods[periodIndex].TrackStart[channel],
				periods[periodIndex].TrackEnd[channel],
				periods[periodIndex].CrossfadeOut[channel].Track,
				periods[periodIndex].CrossfadeIn[channel].Track,
			}
			for _, track := range tracks {
				if _, ok := registry.Lookup(track.Waveform); !ok {
					return fmt.Errorf("unknown waveform %q in rendering period %d channel %d", track.Waveform, periodIndex, channel+1)
				}
			}
		}
	}
	return nil
}

// Render generates the audio and passes buffers to the consume function
func (r *AudioRenderer) Render(consume func(samples []int) error) error {
	defer func() {
		if r.ambianceState != nil {
			r.ambianceState.Close()
		}
		if r.musicState != nil {
			r.musicState.Close()
		}
	}()

	runtime := newRenderRuntime(r, consume)
	return runtime.run()
}
