// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

import "testing"

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
