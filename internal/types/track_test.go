// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "testing"

func TestShiftEffectValidation(t *testing.T) {
	tests := []Track{
		{Type: TrackPureTone, Carrier: 300},
		{Type: TrackBinauralBeat, Carrier: 300, Resonance: 8},
		{Type: TrackMonauralBeat, Carrier: 300, Resonance: 8},
		{Type: TrackIsochronicBeat, Carrier: 300, Resonance: 8},
		{Type: TrackWhiteNoise},
		{Type: TrackPinkNoise},
		{Type: TrackBrownNoise},
		{Type: TrackAmbiance},
		{Type: TrackMusic},
	}
	for _, track := range tests {
		track.Effect = Effect{Type: EffectShift, Value: 10, Intensity: 0.25}
		if err := track.Validate(); err != nil {
			t.Fatalf("valid shift on %s: %v", track.Type.String(), err)
		}
	}
}

func TestShiftEffectString(t *testing.T) {
	if EffectShift.String() != KeywordShift || EffectString(KeywordShift) != EffectShift {
		t.Fatal("shift effect string conversion is inconsistent")
	}
}

func TestDopplerEffectTrackCompatibility(t *testing.T) {
	for _, trackType := range []TrackType{TrackPureTone, TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat, TrackAmbiance, TrackMusic} {
		track := Track{Type: trackType, Carrier: 300, Effect: Effect{Type: EffectDoppler, Value: 1, Intensity: 0.5}}
		if trackType == TrackBinauralBeat || trackType == TrackMonauralBeat || trackType == TrackIsochronicBeat {
			track.Resonance = 8
		}
		if err := track.Validate(); err != nil {
			t.Fatalf("valid doppler on %s: %v", trackType.String(), err)
		}
	}

	for _, trackType := range []TrackType{TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise} {
		track := Track{Type: trackType, Effect: Effect{Type: EffectDoppler, Value: 1, Intensity: 0.5}}
		if err := track.Validate(); err == nil {
			t.Fatalf("expected doppler to be rejected on %s", trackType.String())
		}
	}
}

func TestTrackAmplitudeStringsAlwaysIncludeChannels(t *testing.T) {
	track := Track{
		Type:      TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: [2]AmplitudeType{
			AmplitudePercentToRawChannel(50),
			AmplitudePercentToRawChannel(25),
		},
		Waveform: WaveformSine,
	}

	if got, want := track.String(), "waveform sine tone 300.00 binaural 10.00 amplitude left 50.00 right 25.00"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if got, want := track.ShortString(), " (tone:300.00 binaural:10.00 left:50.00 right:25.00)"; got != want {
		t.Fatalf("ShortString() = %q, want %q", got, want)
	}

	track.Amplitude = AmplitudePercentToRaw(50)
	if got, want := track.String(), "waveform sine tone 300.00 binaural 10.00 amplitude left 50.00 right 50.00"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if got, want := track.ShortString(), " (tone:300.00 binaural:10.00 left:50.00 right:50.00)"; got != want {
		t.Fatalf("ShortString() = %q, want %q", got, want)
	}
}
