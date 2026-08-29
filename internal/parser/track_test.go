// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
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
		Type:       t.TrackAmbiance,
		SourceName: "rain",
		Amplitude:  t.AmplitudePercentToRaw(50),
	}).String()

	trLnMusic := (&t.Track{
		Type:       t.TrackMusic,
		SourceName: "meditation",
		Amplitude:  t.AmplitudePercentToRaw(50),
	}).String()

	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("  %s", trLnTone), true},
		{fmt.Sprintf("  %s", trLnNoise), true},
		{fmt.Sprintf("  %s", trLnAmbiance), true},
		{fmt.Sprintf("  %s", trLnMusic), true},
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
	trs := []*ParsedTrackDeclaration{
		{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, AmplitudePercent: 15, Waveform: t.WaveformSine},
		{Type: t.TrackMonauralBeat, Carrier: 440, Resonance: 11, AmplitudePercent: 20, Waveform: t.WaveformSine},
		{Type: t.TrackIsochronicBeat, Carrier: 220, Resonance: 8, AmplitudePercent: 5, Waveform: t.WaveformSine},
		{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, AmplitudePercent: 15, Waveform: t.WaveformTriangle},
		{Type: t.TrackPureTone, Carrier: 350, AmplitudePercent: 10, Waveform: t.WaveformSquare},
	}

	// Helper to format track without extra waveform
	fmtLine := func(tr *t.Track) string {
		return strings.Join(strings.Fields(tr.String())[2:], " ")
	}
	toTrack := func(decl *ParsedTrackDeclaration) *t.Track {
		return &t.Track{
			Type:      decl.Type,
			Carrier:   decl.Carrier,
			Resonance: decl.Resonance,
			Amplitude: t.AmplitudePercentToRaw(decl.AmplitudePercent),
			Waveform:  decl.Waveform,
		}
	}

	tests := []struct {
		line      string
		wantTrack ParsedTrackDeclaration
	}{
		{fmtLine(toTrack(trs[0])), *trs[0]},
		{fmtLine(toTrack(trs[1])), *trs[1]},
		{fmtLine(toTrack(trs[2])), *trs[2]},
		{toTrack(trs[3]).String(), *trs[3]},
		{toTrack(trs[4]).String(), *trs[4]},
		{"tone 300 effect shift 10 intensity 25 amplitude 20", ParsedTrackDeclaration{Type: t.TrackPureTone, Carrier: 300, EffectType: t.EffectShift, EffectValue: 10, EffectIntensityPercent: 25, AmplitudePercent: 20, Waveform: t.WaveformSine}},
		{"tone 300 binaural 8 effect shift 4 intensity 20 amplitude 20", ParsedTrackDeclaration{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 8, EffectType: t.EffectShift, EffectValue: 4, EffectIntensityPercent: 20, AmplitudePercent: 20, Waveform: t.WaveformSine}},
		{"tone 300 monaural 8 effect shift 4 intensity 20 amplitude 20", ParsedTrackDeclaration{Type: t.TrackMonauralBeat, Carrier: 300, Resonance: 8, EffectType: t.EffectShift, EffectValue: 4, EffectIntensityPercent: 20, AmplitudePercent: 20, Waveform: t.WaveformSine}},
		{"tone 300 isochronic 8 effect shift 4 intensity 20 amplitude 20", ParsedTrackDeclaration{Type: t.TrackIsochronicBeat, Carrier: 300, Resonance: 8, EffectType: t.EffectShift, EffectValue: 4, EffectIntensityPercent: 20, AmplitudePercent: 20, Waveform: t.WaveformSine}},
	}

	for i, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrackDeclaration()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		tr.WaveformSpan = diag.Span{}
		if *tr != tt.wantTrack {
			ts.Errorf("Test %d: For line '%s', expected track %+v but got %+v", i, tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Noise(ts *testing.T) {
	trs := []*ParsedTrackDeclaration{
		{Type: t.TrackWhiteNoise, AmplitudePercent: 5, Waveform: t.WaveformSine},
		{Type: t.TrackPinkNoise, AmplitudePercent: 40, Waveform: t.WaveformSine},
		{Type: t.TrackBrownNoise, AmplitudePercent: 15, Waveform: t.WaveformSine},
	}
	toTrack := func(decl *ParsedTrackDeclaration) *t.Track {
		return &t.Track{Type: decl.Type, Amplitude: t.AmplitudePercentToRaw(decl.AmplitudePercent), Waveform: decl.Waveform}
	}

	tests := []struct {
		line      string
		wantTrack ParsedTrackDeclaration
	}{
		{toTrack(trs[0]).String(), *trs[0]},
		{toTrack(trs[1]).String(), *trs[1]},
		{toTrack(trs[2]).String(), *trs[2]},
		{"noise white effect shift 10 intensity 25 amplitude 20", ParsedTrackDeclaration{Type: t.TrackWhiteNoise, EffectType: t.EffectShift, EffectValue: 10, EffectIntensityPercent: 25, AmplitudePercent: 20, Waveform: t.WaveformSine}},
		{"noise pink smooth 30 effect shift 8 intensity 20 amplitude 20", ParsedTrackDeclaration{Type: t.TrackPinkNoise, NoiseSmooth: 30, EffectType: t.EffectShift, EffectValue: 8, EffectIntensityPercent: 20, AmplitudePercent: 20, Waveform: t.WaveformSine}},
		{"noise brown effect shift 6 intensity 15 amplitude 20", ParsedTrackDeclaration{Type: t.TrackBrownNoise, EffectType: t.EffectShift, EffectValue: 6, EffectIntensityPercent: 15, AmplitudePercent: 20, Waveform: t.WaveformSine}},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrackDeclaration()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		tr.WaveformSpan = diag.Span{}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_NoiseRejectsDoppler(ts *testing.T) {
	ctx := NewTextParser("noise pink effect doppler 1 intensity 40 amplitude 20")
	if _, err := ctx.ParseTrackDeclaration(); err == nil {
		ts.Fatal("expected noise doppler to be rejected")
	}
}

func TestParseTrack_Ambiance(ts *testing.T) {
	tests := []struct {
		line      string
		wantTrack ParsedTrackDeclaration
	}{
		{
			line: "ambiance rain amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:             t.TrackAmbiance,
				SourceName:       "rain",
				AmplitudePercent: 50,
				Waveform:         t.WaveformSine,
			},
		},
		{
			line: "ambiance beach effect pan 10 intensity 75 amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackAmbiance,
				SourceName:             "beach",
				EffectType:             t.EffectPan,
				EffectValue:            10,
				EffectIntensityPercent: 75,
				AmplitudePercent:       50,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "ambiance music effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackAmbiance,
				SourceName:             "music",
				EffectType:             t.EffectModulation,
				EffectValue:            2.5,
				EffectIntensityPercent: 60,
				AmplitudePercent:       40,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "ambiance rain effect shift 10 intensity 25 amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackAmbiance,
				SourceName:             "rain",
				EffectType:             t.EffectShift,
				EffectValue:            10,
				EffectIntensityPercent: 25,
				AmplitudePercent:       50,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "ambiance rain effect doppler 0.8 intensity 40 amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackAmbiance,
				SourceName:             "rain",
				EffectType:             t.EffectDoppler,
				EffectValue:            0.8,
				EffectIntensityPercent: 40,
				AmplitudePercent:       50,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "waveform square ambiance river effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackAmbiance,
				SourceName:             "river",
				EffectType:             t.EffectModulation,
				EffectValue:            2.5,
				EffectIntensityPercent: 60,
				AmplitudePercent:       40,
				Waveform:               t.WaveformSquare,
			},
		},
		{
			line: "ambiance stream_01 amplitude 33",
			wantTrack: ParsedTrackDeclaration{
				Type:             t.TrackAmbiance,
				SourceName:       "stream_01",
				AmplitudePercent: 33,
				Waveform:         t.WaveformSine,
			},
		},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrackDeclaration()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		tr.WaveformSpan = diag.Span{}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Music(ts *testing.T) {
	tests := []struct {
		line      string
		wantTrack ParsedTrackDeclaration
	}{
		{
			line: "music meditation amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:             t.TrackMusic,
				SourceName:       "meditation",
				AmplitudePercent: 50,
				Waveform:         t.WaveformSine,
			},
		},
		{
			line: "music bells effect pan 10 intensity 75 amplitude 50",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackMusic,
				SourceName:             "bells",
				EffectType:             t.EffectPan,
				EffectValue:            10,
				EffectIntensityPercent: 75,
				AmplitudePercent:       50,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "music drones effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackMusic,
				SourceName:             "drones",
				EffectType:             t.EffectModulation,
				EffectValue:            2.5,
				EffectIntensityPercent: 60,
				AmplitudePercent:       40,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "music meditation effect shift 8 intensity 10 amplitude 45",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackMusic,
				SourceName:             "meditation",
				EffectType:             t.EffectShift,
				EffectValue:            8,
				EffectIntensityPercent: 10,
				AmplitudePercent:       45,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "music meditation effect doppler 0.8 intensity 40 amplitude 45",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackMusic,
				SourceName:             "meditation",
				EffectType:             t.EffectDoppler,
				EffectValue:            0.8,
				EffectIntensityPercent: 40,
				AmplitudePercent:       45,
				Waveform:               t.WaveformSine,
			},
		},
		{
			line: "waveform square music bells effect modulation 2.5 intensity 60 amplitude 40",
			wantTrack: ParsedTrackDeclaration{
				Type:                   t.TrackMusic,
				SourceName:             "bells",
				EffectType:             t.EffectModulation,
				EffectValue:            2.5,
				EffectIntensityPercent: 60,
				AmplitudePercent:       40,
				Waveform:               t.WaveformSquare,
			},
		},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrackDeclaration()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		tr.WaveformSpan = diag.Span{}
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
		"  tone 300 amplitude left 10",
		"  tone 300 amplitude left 10 right",
		"  ambiance spin 200 rate five intensity 75 amplitude 50",
		"  ambiance pulse effect modulation 2.5 intensity sixty amplitude 40",
		"  ambiance amplitude 50 extra",
		"  unknown something",
	}

	for _, line := range tests {
		ctx := NewTextParser(line)
		_, err := ctx.ParseTrackDeclaration()
		if err == nil {
			ts.Errorf("For line '%s', expected error but got none", line)
		}
	}
}

func TestParseTrack_TypoDiagnostic(ts *testing.T) {
	ctx := NewTextParser("tone 300 binaual 10 amplitude 10")

	_, err := ctx.ParseTrackDeclaration()
	if err == nil {
		ts.Fatal("expected typo diagnostic")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diag.Diagnostic, got %T", err)
	}
	if diagnostic.Found != "binaual" {
		ts.Fatalf("expected found token binaual, got %q", diagnostic.Found)
	}
	if diagnostic.Suggestion != "did you mean \"binaural\"?" {
		ts.Fatalf("expected binaural suggestion, got %q", diagnostic.Suggestion)
	}
	if diagnostic.Span.Column != 10 || diagnostic.Span.EndColumn != 17 {
		ts.Fatalf("expected typo at 10..17, got %d..%d", diagnostic.Span.Column, diagnostic.Span.EndColumn)
	}
}
