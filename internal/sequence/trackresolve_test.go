// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func parseTrackDecl(ts *testing.T, line string) *parser.ParsedTrackDeclaration {
	ts.Helper()
	decl, err := parser.NewTextParser(line).ParseTrackDeclaration()
	if err != nil {
		ts.Fatalf("ParseTrackDeclaration(%q): %v", line, err)
	}
	return decl
}

func TestBuildTrackFromDeclaration_Success(ts *testing.T) {
	tests := []struct {
		line string
		want t.Track
	}{
		{line: "tone 300 binaural 10 amplitude 15", want: t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(15), Waveform: t.WaveformSine}},
		{line: "waveform square tone 350 amplitude 10", want: t.Track{Type: t.TrackPureTone, Carrier: 350, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSquare}},
		{line: "noise pink amplitude 40", want: t.Track{Type: t.TrackPinkNoise, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine}},
		{line: "ambiance beach effect pan 10 intensity 75 amplitude 50", want: t.Track{Type: t.TrackAmbiance, SourceName: "beach", Amplitude: t.AmplitudePercentToRaw(50), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPan, Value: 10, Intensity: t.IntensityPercentToRaw(75)}}},
		{line: "music meditation effect modulation 2.5 intensity 60 amplitude 40", want: t.Track{Type: t.TrackMusic, SourceName: "meditation", Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectModulation, Value: 2.5, Intensity: t.IntensityPercentToRaw(60)}}},
		{line: "ambiance rain effect shift 10 intensity 25 amplitude 50", want: t.Track{Type: t.TrackAmbiance, SourceName: "rain", Amplitude: t.AmplitudePercentToRaw(50), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectShift, Value: 10, Intensity: t.IntensityPercentToRaw(25)}}},
		{line: "music meditation effect shift 8 intensity 10 amplitude 45", want: t.Track{Type: t.TrackMusic, SourceName: "meditation", Amplitude: t.AmplitudePercentToRaw(45), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectShift, Value: 8, Intensity: t.IntensityPercentToRaw(10)}}},
		{line: "tone 300 binaural 10 amplitude left 10 right 5", want: t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: [2]t.AmplitudeType{t.AmplitudePercentToRawChannel(10), t.AmplitudePercentToRawChannel(5)}, Waveform: t.WaveformSine}},
		{line: "noise pink amplitude left 30 right 25", want: t.Track{Type: t.TrackPinkNoise, Amplitude: [2]t.AmplitudeType{t.AmplitudePercentToRawChannel(30), t.AmplitudePercentToRawChannel(25)}, Waveform: t.WaveformSine}},
		{line: "ambiance beach amplitude left 10 right 15", want: t.Track{Type: t.TrackAmbiance, SourceName: "beach", Amplitude: [2]t.AmplitudeType{t.AmplitudePercentToRawChannel(10), t.AmplitudePercentToRawChannel(15)}, Waveform: t.WaveformSine}},
		{line: "music meditation amplitude left 5 right 10", want: t.Track{Type: t.TrackMusic, SourceName: "meditation", Amplitude: [2]t.AmplitudeType{t.AmplitudePercentToRawChannel(5), t.AmplitudePercentToRawChannel(10)}, Waveform: t.WaveformSine}},
	}

	for _, tt := range tests {
		decl := parseTrackDecl(ts, tt.line)
		track, err := buildTrackFromDeclaration("test.spsq", 1, tt.line, decl)
		if err != nil {
			ts.Fatalf("buildTrackFromDeclaration(%q): %v", tt.line, err)
		}
		if *track != tt.want {
			ts.Fatalf("line %q: expected %+v, got %+v", tt.line, tt.want, *track)
		}
	}
}

func TestBuildTrackFromDeclaration_Errors(ts *testing.T) {
	tests := []string{
		"tone 300 binaural 10 amplitude 120",
		"ambiance pulse effect modulation 2.5 intensity 150 amplitude 40",
		"noise brown smooth 150 amplitude 30",
		"tone 300 amplitude left 10 right 120",
	}

	for _, line := range tests {
		decl := parseTrackDecl(ts, line)
		if _, err := buildTrackFromDeclaration("test.spsq", 1, line, decl); err == nil {
			ts.Fatalf("expected semantic error for %q", line)
		}
	}
}
