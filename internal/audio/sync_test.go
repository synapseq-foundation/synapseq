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

func TestAudioRendererSync_InterpolatesTrackAndSignal(ts *testing.T) {
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

	renderer := newTestRenderer(ts, []t.Period{p0, p1})
	renderer.sync(500, 0)

	channel := renderer.channels[0]

	assertAlmostEqual(ts, float64(channel.Track.Amplitude), float64(t.AmplitudePercentToRaw(20)), 0.0001)
	assertAlmostEqual(ts, channel.Track.Carrier, 250, 0.0001)
	assertAlmostEqual(ts, channel.Track.Resonance, 10, 0.0001)
	assertAlmostEqual(ts, channel.Track.Effect.Value, 3, 0.0001)
	assertAlmostEqual(ts, float64(channel.Track.Effect.Intensity), float64(t.IntensityPercentToRaw(40)), 0.0001)

	if channel.Track.Type != t.TrackBinauralBeat {
		ts.Fatalf("unexpected track type: got %v", channel.Track.Type)
	}
	if channel.Track.Waveform != t.WaveformSine {
		ts.Fatalf("unexpected waveform: got %v", channel.Track.Waveform)
	}
	if channel.Track.Effect.Type != t.EffectModulation {
		ts.Fatalf("unexpected effect type: got %v", channel.Track.Effect.Type)
	}

	if channel.Amplitude[0] != int(channel.Track.Amplitude) || channel.Amplitude[1] != int(channel.Track.Amplitude) {
		ts.Fatalf("unexpected amplitudes: got %v", channel.Amplitude)
	}
	if channel.Increment[0] != renderer.frequencyToIncrement(255) {
		ts.Fatalf("unexpected high increment: got %d", channel.Increment[0])
	}
	if channel.Increment[1] != renderer.frequencyToIncrement(245) {
		ts.Fatalf("unexpected low increment: got %d", channel.Increment[1])
	}
	if channel.Effect.Increment != renderer.frequencyToIncrement(3) {
		ts.Fatalf("unexpected effect increment: got %d", channel.Effect.Increment)
	}
}

func TestAudioRendererSync_ResetsOffsetsAndClearsResidualState(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackWhiteNoise,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	renderer := newTestRenderer(ts, []t.Period{p0, p1})
	renderer.channels[0].Type = t.TrackBinauralBeat
	renderer.channels[0].Track.Effect.Type = t.EffectPan
	renderer.channels[0].Offset = [2]int{77, 99}
	renderer.channels[0].Amplitude = [2]int{11, 22}
	renderer.channels[0].Increment = [2]int{33, 44}
	renderer.channels[0].Effect.Offset = 55
	renderer.channels[0].Effect.Increment = 66

	renderer.sync(0, 0)

	channel := renderer.channels[0]
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
}

func TestAudioRendererSync_ResetsEffectPhaseWhenEffectChanges(ts *testing.T) {
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

	renderer := newTestRenderer(ts, []t.Period{p0, p1})
	renderer.channels[0].Track.Effect.Type = t.EffectPan
	renderer.channels[0].Effect.Offset = 1234
	renderer.channels[0].Effect.Increment = 4321

	renderer.sync(0, 0)

	channel := renderer.channels[0]
	if channel.Effect.Offset != 0 {
		ts.Fatalf("effect offset was not reset: got %d", channel.Effect.Offset)
	}
	if channel.Effect.Increment != renderer.frequencyToIncrement(5) {
		ts.Fatalf("unexpected effect increment after effect change: got %d", channel.Effect.Increment)
	}
	if channel.Increment[0] != renderer.frequencyToIncrement(220) {
		ts.Fatalf("unexpected carrier increment: got %d", channel.Increment[0])
	}
	if channel.Increment[1] != 0 {
		ts.Fatalf("unexpected secondary increment for pure tone: got %d", channel.Increment[1])
	}
}

func newTestRenderer(ts *testing.T, periods []t.Period) *AudioRenderer {
	ts.Helper()

	renderer, err := NewAudioRenderer(periods, &AudioRendererOptions{
		SampleRate: 44100,
		Volume:     80,
		Ambiance:   map[string]string{},
	})
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	return renderer
}

func assertAlmostEqual(ts *testing.T, got, want, tolerance float64) {
	ts.Helper()

	if math.Abs(got-want) > tolerance {
		ts.Fatalf("unexpected value: got %.6f want %.6f", got, want)
	}
}