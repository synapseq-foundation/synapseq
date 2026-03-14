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
	"math"
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
		WaveformStart: t.WaveformSquare,
		WaveformEnd:   t.WaveformSquare,
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 0},
		Increment:     [2]int{t.PhasePrecision, 0},
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
		WaveformStart: t.WaveformSquare,
		WaveformEnd:   t.WaveformSquare,
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 0},
		Increment:     [2]int{t.PhasePrecision, 0},
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
			Waveform: t.WaveformSawtooth,
			Effect: t.Effect{
				Type:      t.EffectModulation,
				Intensity: t.IntensityPercentToRaw(100),
			},
		},
		WaveformStart: t.WaveformSawtooth,
		WaveformEnd:   t.WaveformSawtooth,
		Type:          t.TrackBinauralBeat,
		Amplitude:     [2]int{4096, 4096},
		Increment:     [2]int{t.PhasePrecision, t.PhasePrecision},
		Effect: t.EffectState{
			Increment: int(t.SineTableSize/2) * t.PhasePrecision,
		},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	baseSample := renderer.channels[0].Amplitude[0] * renderer.waveTables[int(t.WaveformSawtooth)][1]
	modFactor := 0.5
	gain := 0.3 + 0.7*modFactor
	expected := clampPCM16(int(math.Round(float64(baseSample)*gain)) >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected modulation output: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}

	expectedEffectOffset := (t.BufferSize * int(t.SineTableSize/2) * t.PhasePrecision) & phaseMask
	if renderer.channels[0].Effect.Offset != expectedEffectOffset {
		ts.Fatalf("unexpected modulation phase offset: got %d, want %d", renderer.channels[0].Effect.Offset, expectedEffectOffset)
	}
}

func TestAudioRendererMix_SawtoothModulationUsesLinearRamp(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackIsochronicBeat,
		Waveform: t.WaveformSawtooth,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = t.WaveformSawtooth
	channel.WaveformEnd = t.WaveformSawtooth
	channel.Effect.Offset = int(t.SineTableSize/2) * t.PhasePrecision

	got := renderer.applyModulationToCurrentPhase(channel, 1000)
	expected := int(math.Round(1000 * (0.3 + 0.7*0.5)))
	if got != expected {
		ts.Fatalf("unexpected sawtooth modulation output at mid ramp: got %d, want %d", got, expected)
	}
}

func TestAudioRendererMix_PureToneMorphsBetweenWaveforms(ts *testing.T) {
	renderer := newMixTestRenderer()
	renderer.channels[0] = t.Channel{
		Track: t.Track{
			Type:     t.TrackPureTone,
			Waveform: t.WaveformSine,
		},
		WaveformStart: t.WaveformSine,
		WaveformEnd:   t.WaveformSquare,
		WaveformAlpha: 0.25,
		Type:          t.TrackPureTone,
		Amplitude:     [2]int{4096, 0},
		Increment:     [2]int{t.PhasePrecision, 0},
	}

	samples := renderer.mix(make([]int, t.BufferSize*audioChannels))

	sine := float64(renderer.waveTables[int(t.WaveformSine)][1])
	square := float64(renderer.waveTables[int(t.WaveformSquare)][1])
	blended := lerpFloat64(sine, square, 0.25)
	expectedRaw := int(math.Round(float64(renderer.channels[0].Amplitude[0]) * blended))
	expected := clampPCM16(expectedRaw >> audioBitShift)

	if samples[0] != expected || samples[1] != expected {
		ts.Fatalf("unexpected morphed sample: got [%d %d], want [%d %d]", samples[0], samples[1], expected, expected)
	}
}

func TestAudioRendererMix_ModulationSlewsAbruptSquareGainChanges(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackAmbiance,
		Waveform: t.WaveformSquare,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = t.WaveformSquare
	channel.WaveformEnd = t.WaveformSquare
	channel.Effect.ModulationGain = 1
	channel.Effect.ModulationInitialized = true

	got := renderer.applyModulationToCurrentPhase(channel, 1000)
	expectedFloor := int(1000 * 0.3)
	if got <= expectedFloor || got >= 1000 {
		ts.Fatalf("unexpected slewed modulation output: got %d, want between %d and 1000", got, expectedFloor)
	}
	if channel.Effect.ModulationGain >= 1 || channel.Effect.ModulationGain <= 0.3 {
		ts.Fatalf("unexpected slewed modulation gain: got %f", channel.Effect.ModulationGain)
	}
}

func TestAudioRendererMix_PanUsesWaveformAndSlewsSquareSwitches(ts *testing.T) {
	renderer := newMixTestRenderer()
	channel := &renderer.channels[0]
	channel.Track = t.Track{
		Type:     t.TrackAmbiance,
		Waveform: t.WaveformSquare,
		Effect: t.Effect{
			Type:      t.EffectPan,
			Intensity: t.IntensityPercentToRaw(100),
		},
	}
	channel.WaveformStart = t.WaveformSquare
	channel.WaveformEnd = t.WaveformSquare
	channel.Effect.PanPosition = -1
	channel.Effect.PanInitialized = true
	channel.Effect.Offset = int(t.SineTableSize/4) * t.PhasePrecision

	left, right := renderer.applyPan(channel, 1000, 1000)
	if left <= 0 || right <= 0 {
		ts.Fatalf("unexpected hard-switched pan output: got [%d %d]", left, right)
	}
	if channel.Effect.PanPosition <= -1 || channel.Effect.PanPosition >= 1 {
		ts.Fatalf("unexpected slewed pan position: got %f", channel.Effect.PanPosition)
	}
	if channel.Effect.PanPosition >= -0.9 {
		ts.Fatalf("expected a small initial slew step, got pan position %f", channel.Effect.PanPosition)
	}
	if right <= 0 || left >= 1000 {
		ts.Fatalf("expected pan to start moving toward right channel without hard switch: got [%d %d]", left, right)
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
