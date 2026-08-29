// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "testing"

func TestShiftEffectValidation(t *testing.T) {
	valid := Track{Type: TrackAmbiance, Effect: Effect{Type: EffectShift, Value: 10, Intensity: 0.25}}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid ambiance shift: %v", err)
	}

	invalid := Track{Type: TrackPureTone, Carrier: 300, Effect: Effect{Type: EffectShift, Value: 10, Intensity: 0.25}}
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected shift on tone to be rejected")
	}
}

func TestShiftEffectString(t *testing.T) {
	if EffectShift.String() != KeywordShift || EffectString(KeywordShift) != EffectShift {
		t.Fatal("shift effect string conversion is inconsistent")
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
