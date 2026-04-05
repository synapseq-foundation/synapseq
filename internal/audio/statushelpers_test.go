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
	"testing"

	audiostatus "github.com/synapseq-foundation/synapseq/v4/internal/audio/status"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestCountActiveChannels(ts *testing.T) {
	tests := []struct {
		name     string
		channels []t.Channel
		expected int
	}{
		{"empty slice -> at least 1", []t.Channel{}, 1},
		{"all off -> 1", make([]t.Channel, 5), 1},
		{"single active at 0 -> 1", func() []t.Channel { channels := make([]t.Channel, 4); channels[0].Track.Type = t.TrackBinauralBeat; return channels }(), 1},
		{"last active at end -> len", func() []t.Channel { channels := make([]t.Channel, 4); channels[3].Track.Type = t.TrackPinkNoise; return channels }(), 4},
		{"last active in the middle -> index+1", func() []t.Channel { channels := make([]t.Channel, 5); channels[2].Track.Type = t.TrackBrownNoise; return channels }(), 3},
		{"multiple actives -> last index+1", func() []t.Channel { channels := make([]t.Channel, 8); channels[1].Track.Type = t.TrackBinauralBeat; channels[6].Track.Type = t.TrackAmbiance; return channels }(), 7},
		{"all active -> len", func() []t.Channel { channels := make([]t.Channel, 7); for i := range channels { channels[i].Track.Type = t.TrackPinkNoise }; return channels }(), 7},
		{"last off but previous active", func() []t.Channel { channels := make([]t.Channel, 6); channels[4].Track.Type = t.TrackAmbiance; return channels }(), 5},
	}

	for _, test := range tests {
		got := audiostatus.CountActiveChannels(test.channels)
		if got != test.expected {
			ts.Errorf("%s: expected %d, got %d", test.name, test.expected, got)
		}
	}
}

func TestIsTrackEqual(ts *testing.T) {
	base := &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}

	tests := []struct {
		name string
		a    *t.Track
		b    *t.Track
		eq   bool
	}{
		{"identical", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, true},
		{"different amplitude", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(30), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different carrier", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 320, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different resonance", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 12, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different waveform", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformTriangle, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different intensity", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(50)}}, false},
		{"different type", base, &t.Track{Type: t.TrackMonauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"ambiance effect type ignored", &t.Track{Type: t.TrackAmbiance, Carrier: 200, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPan, Intensity: t.IntensityPercentToRaw(60)}}, &t.Track{Type: t.TrackAmbiance, Carrier: 200, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectModulation, Intensity: t.IntensityPercentToRaw(60)}}, true},
	}

	for _, test := range tests {
		got := audiostatus.IsTrackEqual(test.a, test.b)
		if got != test.eq {
			ts.Errorf("%s: expected %v, got %v\nA=%+v\nB=%+v", test.name, test.eq, got, *test.a, *test.b)
		}
	}
}