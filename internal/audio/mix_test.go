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
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestAudioRendererMix_PureToneFirstSampleAndPhaseAdvance(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSquare,
		},
		Type:      t.TrackPureTone,
		Amplitude: [2]int{4096, 0},
		Increment: [2]int{t.PhasePrecision, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedRaw := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(t.WaveformSquare)][1]
	expected := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected first stereo sample: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}

	expectedOffset := (t.BufferSize * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Offset[0] != expectedOffset {
		ts.Fatalf("unexpected phase offset: got %d, want %d", renderer.channels[0].Offset[0], expectedOffset)
	}
}

func TestAudioRendererMix_PanEffectRoutesMonoSignal(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSquare,
			Effect: t.Effect{
				Type:      t.EffectPan,
				Intensity: t.IntensityPercentToRaw(100),
			},
		},
		Type:      t.TrackPureTone,
		Amplitude: [2]int{4096, 0},
		Increment: [2]int{t.PhasePrecision, 0},
		Effect: t.EffectState{
			Increment: int(t.SineTableSize/4) * t.PhasePrecision,
		},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedRaw := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(t.WaveformSquare)][1]
	expectedRight := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != 0 || samples[1] != expectedRight {
		ts.Fatalf("unexpected pan output: got [%d %d], want [0 %d]", samples[0], samples[1], expectedRight)
	}

	expectedEffectOffset := (t.BufferSize * int(t.SineTableSize/4) * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Effect.Offset != expectedEffectOffset {
		ts.Fatalf("unexpected effect phase offset: got %d, want %d", renderer.channels[0].Effect.Offset, expectedEffectOffset)
	}
}

func TestAudioRendererMix_ModulationAffectsStereoWithSharedPhase(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackBinauralBeat,
			Waveform: t.WaveformSquare,
			Effect: t.Effect{
				Type:      t.EffectModulation,
				Intensity: t.IntensityPercentToRaw(100),
			},
		},
		Type:      t.TrackBinauralBeat,
		Amplitude: [2]int{4096, 4096},
		Increment: [2]int{t.PhasePrecision, t.PhasePrecision},
		Effect: t.EffectState{
			Increment: int(t.SineTableSize/2) * t.PhasePrecision,
		},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	baseSample := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(t.WaveformSquare)][1]
	expected := clampPCM16(int(float64(baseSample)*0.3) >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected modulation output: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}

	expectedEffectOffset := (t.BufferSize * int(t.SineTableSize/2) * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Effect.Offset != expectedEffectOffset {
		ts.Fatalf("unexpected modulation phase offset: got %d, want %d", renderer.channels[0].Effect.Offset, expectedEffectOffset)
	}
}

func TestAudioRendererMix_AmbianceUsesPreparedStereoBuffer(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.ambianceSamplesByIndex = [][]int{{20000, -10000}}
	renderer.channelAmbianceIndex[0] = 0
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type: t.TrackAmbiance,
		},
		Type:      t.TrackAmbiance,
		Amplitude: [2]int{3, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	expectedLeft := clampPCM16((20000 * 16 * 3) >> audioBitShift)
	expectedRight := clampPCM16((-10000 * 16 * 3) >> audioBitShift)

	if samples[0] != expectedLeft || samples[1] != expectedRight {
		ts.Fatalf("unexpected ambiance output: got [%d %d], want [%d %d]", samples[0], samples[1], expectedLeft, expectedRight)
	}
}

func newMixTestRenderer() *AudioRenderer {
	renderer := &AudioRenderer{
		waveTables:             InitWaveformTables(),
		noiseGenerator:         NewNoiseGenerator(),
		ambianceSamplesByIndex: [][]int{},
		AudioRendererOptions: &AudioRendererOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	for i := range renderer.channelAmbianceIndex {
		renderer.channelAmbianceIndex[i] = -1
	}

	return renderer
}
