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

package shared

import (
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestIsTrackEqual(ts *testing.T) {
	base := &t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
		Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
	}

	tests := []struct {
		name string
		a, b *t.Track
		eq   bool
	}{
		{
			name: "identical",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: true,
		},
		{
			name: "different amplitude",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(30),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different carrier",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   320,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different resonance",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 12,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different waveform",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformTriangle,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different intensity",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(50)},
			},
			eq: false,
		},
		{
			name: "different type",
			a:    base,
			b: &t.Track{
				Type:      t.TrackMonauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "ambiance effect type ignored",
			a: &t.Track{
				Type:      t.TrackAmbiance,
				Carrier:   200,
				Resonance: 5,
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectPan, Intensity: t.IntensityPercentToRaw(60)},
			},
			b: &t.Track{
				Type:      t.TrackAmbiance,
				Carrier:   200,
				Resonance: 5,
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectModulation, Intensity: t.IntensityPercentToRaw(60)},
			},
			eq: true,
		},
	}

	for _, tc := range tests {
		got := IsTrackEqual(tc.a, tc.b)
		if got != tc.eq {
			ts.Errorf("%s: expected %v, got %v\nA=%+v\nB=%+v", tc.name, tc.eq, got, *tc.a, *tc.b)
		}
	}
}
