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

package parser

import (
	"fmt"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestHasTrack(ts *testing.T) {
	trLnTone := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()

	trLnNoise := (&t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}).String()

	trLnAmbiance := (&t.Track{
		Type:         t.TrackAmbiance,
		AmbianceName: "rain",
		Amplitude:    t.AmplitudePercentToRaw(50),
	}).String()

	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("  %s", trLnTone), true},
		{fmt.Sprintf("  %s", trLnNoise), true},
		{fmt.Sprintf("  %s", trLnAmbiance), true},
		{fmt.Sprintf(" %s", trLnTone), false},
		{fmt.Sprintf("   %s", trLnTone), false},
		{trLnTone, false},
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		got := ctx.HasTrack()
		if got != test.expected {
			ts.Errorf("For line '%s', expected HasTrack()=%v, got %v", test.line, test.expected, got)
		}
	}
}

func TestParseTrack_Tones(ts *testing.T) {
	trs := []*t.Track{
		{
			Type:      t.TrackBinauralBeat,
			Carrier:   300,
			Resonance: 10,
			Amplitude: t.AmplitudePercentToRaw(15),
		},
		{
			Type:      t.TrackMonauralBeat,
			Carrier:   440,
			Resonance: 11,
			Amplitude: t.AmplitudePercentToRaw(20),
		},
		{
			Type:      t.TrackIsochronicBeat,
			Carrier:   220,
			Resonance: 8,
			Amplitude: t.AmplitudePercentToRaw(5),
		},
		{
			Type:      t.TrackBinauralBeat,
			Carrier:   300,
			Resonance: 10,
			Amplitude: t.AmplitudePercentToRaw(15),
			Waveform:  t.WaveformTriangle,
		},
		{
			Type:      t.TrackPureTone,
			Carrier:   350,
			Amplitude: t.AmplitudePercentToRaw(10),
			Waveform:  t.WaveformSquare,
		},
	}

	// Helper to format track without extra waveform
	fmtLine := func(tr *t.Track) string {
		return strings.Join(strings.Fields(tr.String())[2:], " ")
	}

	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{fmtLine(trs[0]), *trs[0]},
		{fmtLine(trs[1]), *trs[1]},
		{fmtLine(trs[2]), *trs[2]},
		{trs[3].String(), *trs[3]},
		{trs[4].String(), *trs[4]},
	}

	for i, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("Test %d: For line '%s', expected track %+v but got %+v", i, tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Noise(ts *testing.T) {
	trs := []*t.Track{
		{
			Type:      t.TrackWhiteNoise,
			Amplitude: t.AmplitudePercentToRaw(5),
		},
		{
			Type:      t.TrackPinkNoise,
			Amplitude: t.AmplitudePercentToRaw(40),
		},
		{
			Type:      t.TrackBrownNoise,
			Amplitude: t.AmplitudePercentToRaw(15),
		},
	}

	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{trs[0].String(), *trs[0]},
		{trs[1].String(), *trs[1]},
		{trs[2].String(), *trs[2]},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Ambiance(ts *testing.T) {
	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{
			line: "ambiance rain amplitude 50",
			wantTrack: t.Track{
				Type:         t.TrackAmbiance,
				AmbianceName: "rain",
				Amplitude:    t.AmplitudePercentToRaw(50),
				Waveform:     t.WaveformSine,
			},
		},
		{
			line: "ambiance beach effect pan 10 intensity 75 amplitude 50",
			wantTrack: t.Track{
				Type:         t.TrackAmbiance,
				AmbianceName: "beach",
				Effect: t.Effect{
					Type:      t.EffectPan,
					Value:     10,
					Intensity: t.IntensityPercentToRaw(75),
				},
				Amplitude: t.AmplitudePercentToRaw(50),
				Waveform:  t.WaveformSine,
			},
		},
		{
			line: "ambiance music effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: t.Track{
				Type:         t.TrackAmbiance,
				AmbianceName: "music",
				Effect: t.Effect{
					Type:      t.EffectModulation,
					Value:     2.5,
					Intensity: t.IntensityPercentToRaw(60),
				},
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSine,
			},
		},
		{
			line: "waveform square ambiance river effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: t.Track{
				Type:         t.TrackAmbiance,
				AmbianceName: "river",
				Effect: t.Effect{
					Type:      t.EffectModulation,
					Value:     2.5,
					Intensity: t.IntensityPercentToRaw(60),
				},
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSquare,
			},
		},
		{
			line: "ambiance stream_01 amplitude 33",
			wantTrack: t.Track{
				Type:         t.TrackAmbiance,
				AmbianceName: "stream_01",
				Amplitude:    t.AmplitudePercentToRaw(33),
				Waveform:     t.WaveformSine,
			},
		},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Errors(ts *testing.T) {
	tests := []string{
		"  tone 300 binaural amplitude 10",
		"  tone 300 unknown 10 amplitude 10",
		"  noise white amplitude",
		"  ambiance spin 200 rate five intensity 75 amplitude 50",
		"  ambiance pulse effect modulation 2.5 intensity sixty amplitude 40",
		"  ambiance amplitude 50 extra",
		"  tone 300 binaural 10 amplitude 120",
		"  ambiance pulse effect modulation 2.5 intensity 150 amplitude 40",
		"  unknown something",
	}

	for _, line := range tests {
		ctx := NewTextParser(line)
		_, err := ctx.ParseTrack()
		if err == nil {
			ts.Errorf("For line '%s', expected error but got none", line)
		}
	}
}
