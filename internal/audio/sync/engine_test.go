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

package sync

import (
	"math"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const testSampleRate = 44100

func TestEngineSync_InterpolatesTrackAndSignal(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(10),
		Carrier:   200,
		Resonance: 8,
		Waveform:  t.WaveformSine,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     2,
			Intensity: t.IntensityPercentToRaw(20),
		},
	}
	p0.TrackEnd[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(30),
		Carrier:   300,
		Resonance: 12,
		Waveform:  t.WaveformTriangle,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     4,
			Intensity: t.IntensityPercentToRaw(60),
		},
	}

	channels := make([]t.Channel, t.NumberOfChannels)
	engine := NewEngine(testSampleRate, nil)
	engine.Sync(channels, []t.Period{p0, p1}, 500, 0)

	channel := channels[0]

	assertAlmostEqual(ts, float64(channel.Track.Amplitude), float64(t.AmplitudePercentToRaw(20)), 0.0001)
	assertAlmostEqual(ts, channel.Track.Carrier, 250, 0.0001)
	assertAlmostEqual(ts, channel.Track.Resonance, 10, 0.0001)
	assertAlmostEqual(ts, channel.Track.Effect.Value, 3, 0.0001)
	assertAlmostEqual(ts, float64(channel.Track.Effect.Intensity), float64(t.IntensityPercentToRaw(40)), 0.0001)

	if channel.Track.Type != t.TrackBinauralBeat {
		ts.Fatalf("unexpected track type: got %v", channel.Track.Type)
	}
	if channel.WaveformStart != t.WaveformSine || channel.WaveformEnd != t.WaveformTriangle {
		ts.Fatalf("unexpected waveform morph state: got %v -> %v", channel.WaveformStart, channel.WaveformEnd)
	}
	assertAlmostEqual(ts, channel.WaveformAlpha, 0.5, 0.0001)
	if channel.Track.Waveform != t.WaveformSine {
		ts.Fatalf("unexpected displayed waveform: got %v", channel.Track.Waveform)
	}
	if channel.Track.Effect.Type != t.EffectModulation {
		ts.Fatalf("unexpected effect type: got %v", channel.Track.Effect.Type)
	}

	if channel.Amplitude[0] != int(channel.Track.Amplitude) || channel.Amplitude[1] != int(channel.Track.Amplitude) {
		ts.Fatalf("unexpected amplitudes: got %v", channel.Amplitude)
	}
	if channel.Increment[0] != FrequencyToIncrement(testSampleRate, 255) {
		ts.Fatalf("unexpected high increment: got %d", channel.Increment[0])
	}
	if channel.Increment[1] != FrequencyToIncrement(testSampleRate, 245) {
		ts.Fatalf("unexpected low increment: got %d", channel.Increment[1])
	}
	if channel.Effect.Increment != FrequencyToIncrement(testSampleRate, 3) {
		ts.Fatalf("unexpected effect increment: got %d", channel.Effect.Increment)
	}
}

func TestEngineSync_ResetsOffsetsAndClearsResidualState(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackWhiteNoise,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	channels := make([]t.Channel, t.NumberOfChannels)
	channels[0].Type = t.TrackBinauralBeat
	channels[0].Track.Effect.Type = t.EffectPan
	channels[0].Offset = [2]int{77, 99}
	channels[0].Amplitude = [2]int{11, 22}
	channels[0].Increment = [2]int{33, 44}
	channels[0].Effect.Offset = 55
	channels[0].Effect.Increment = 66
	channels[0].Effect.ModulationGain = 0.75
	channels[0].Effect.ModulationInitialized = true
	channels[0].Effect.PanPosition = 0.25
	channels[0].Effect.PanInitialized = true

	engine := NewEngine(testSampleRate, nil)
	engine.Sync(channels, []t.Period{p0, p1}, 0, 0)

	channel := channels[0]
	if channel.Type != t.TrackWhiteNoise {
		ts.Fatalf("unexpected runtime type: got %v", channel.Type)
	}
	if channel.Offset != [2]int{} {
		ts.Fatalf("offsets were not reset: got %v", channel.Offset)
	}
	if channel.Amplitude[0] != int(channel.Track.Amplitude) || channel.Amplitude[1] != 0 {
		ts.Fatalf("unexpected amplitudes after cleanup: got %v", channel.Amplitude)
	}
	if channel.Increment != [2]int{} {
		ts.Fatalf("increments were not cleared: got %v", channel.Increment)
	}
	if channel.Effect.Increment != 0 {
		ts.Fatalf("effect increment was not cleared: got %d", channel.Effect.Increment)
	}
	if channel.Effect.Offset != 0 {
		ts.Fatalf("effect offset was not reset on effect change: got %d", channel.Effect.Offset)
	}
	if channel.Effect.ModulationGain != 0 || channel.Effect.ModulationInitialized {
		ts.Fatalf("modulation smoothing state was not reset: got gain=%f initialized=%v", channel.Effect.ModulationGain, channel.Effect.ModulationInitialized)
	}
	if channel.Effect.PanPosition != 0 || channel.Effect.PanInitialized {
		ts.Fatalf("pan smoothing state was not reset: got pos=%f initialized=%v", channel.Effect.PanPosition, channel.Effect.PanInitialized)
	}
}

func TestEngineSync_ResetsEffectPhaseWhenEffectChanges(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackPureTone,
		Amplitude: t.AmplitudePercentToRaw(15),
		Carrier:   220,
		Waveform:  t.WaveformSawtooth,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     5,
			Intensity: t.IntensityPercentToRaw(30),
		},
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	channels := make([]t.Channel, t.NumberOfChannels)
	channels[0].Track.Effect.Type = t.EffectPan
	channels[0].Effect.Offset = 1234
	channels[0].Effect.Increment = 4321
	channels[0].Effect.ModulationGain = 0.4
	channels[0].Effect.ModulationInitialized = true
	channels[0].Effect.PanPosition = -0.5
	channels[0].Effect.PanInitialized = true

	engine := NewEngine(testSampleRate, nil)
	engine.Sync(channels, []t.Period{p0, p1}, 0, 0)

	channel := channels[0]
	if channel.Effect.Offset != 0 {
		ts.Fatalf("effect offset was not reset: got %d", channel.Effect.Offset)
	}
	if channel.Effect.Increment != FrequencyToIncrement(testSampleRate, 5) {
		ts.Fatalf("unexpected effect increment after effect change: got %d", channel.Effect.Increment)
	}
	if channel.Effect.ModulationInitialized {
		ts.Fatalf("modulation smoothing state should be reset on effect change")
	}
	if channel.Effect.PanInitialized {
		ts.Fatalf("pan smoothing state should be reset on effect change")
	}
	if channel.Increment[0] != FrequencyToIncrement(testSampleRate, 220) {
		ts.Fatalf("unexpected carrier increment: got %d", channel.Increment[0])
	}
	if channel.Increment[1] != 0 {
		ts.Fatalf("unexpected secondary increment for pure tone: got %d", channel.Increment[1])
	}
}

func TestEngineSync_AppliesStepsTrajectory(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p0.Steps = 1
	p1.Time = 9000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(10),
		Carrier:   200,
		Resonance: 8,
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(20),
		Carrier:   300,
		Resonance: 12,
		Waveform:  t.WaveformTriangle,
	}

	channels := make([]t.Channel, t.NumberOfChannels)
	engine := NewEngine(testSampleRate, nil)

	engine.Sync(channels, []t.Period{p0, p1}, 3000, 0)
	first := channels[0].Track
	assertAlmostEqual(ts, first.Carrier, 300, 0.0001)
	assertAlmostEqual(ts, first.Resonance, 12, 0.0001)

	engine.Sync(channels, []t.Period{p0, p1}, 6000, 0)
	second := channels[0].Track
	assertAlmostEqual(ts, second.Carrier, 200, 0.0001)
	assertAlmostEqual(ts, second.Resonance, 8, 0.0001)

	engine.Sync(channels, []t.Period{p0, p1}, 9000, 0)
	third := channels[0].Track
	assertAlmostEqual(ts, third.Carrier, 300, 0.0001)
	assertAlmostEqual(ts, third.Resonance, 12, 0.0001)
	assertAlmostEqual(ts, channels[0].WaveformAlpha, 1, 0.0001)
}

func assertAlmostEqual(ts *testing.T, got, want, tolerance float64) {
	ts.Helper()

	if math.Abs(got-want) > tolerance {
		ts.Fatalf("unexpected value: got %.6f want %.6f", got, want)
	}
}
