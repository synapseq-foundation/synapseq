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

package preset_test

import (
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	"github.com/synapseq-foundation/synapseq/v4/internal/preset"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func parseOverrideSpec(ts *testing.T, line string) *preset.TrackOverrideSpec {
	ts.Helper()
	decl, err := parser.NewTextParser(line).ParseTrackOverrideDeclaration()
	if err != nil {
		ts.Fatalf("parse override declaration %q: %v", line, err)
	}
	return decl
}

func testTemplatePreset(ts *testing.T) *t.Preset {
	ts.Helper()
	templatePreset, err := t.NewPreset("base", true, nil)
	if err != nil {
		ts.Fatalf("failed to create template: %v", err)
	}
	templatePreset.Track[0] = t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine}
	templatePreset.Track[1] = t.Track{Type: t.TrackAmbiance, AmbianceName: "rain", Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPan, Value: 5, Intensity: t.IntensityPercentToRaw(75)}}
	templatePreset.Track[2] = t.Track{Type: t.TrackMonauralBeat, Carrier: 440, Resonance: 8, Amplitude: t.AmplitudePercentToRaw(15), Waveform: t.WaveformSquare}
	templatePreset.Track[3] = t.Track{Type: t.TrackPinkNoise, NoiseSmooth: 20, Amplitude: t.AmplitudePercentToRaw(25)}
	return templatePreset
}

func testDerivedPreset(ts *testing.T, templatePreset *t.Preset) *t.Preset {
	ts.Helper()
	derivedPreset, err := t.NewPreset("derived", false, templatePreset)
	if err != nil {
		ts.Fatalf("failed to create derived preset: %v", err)
	}
	return derivedPreset
}

func TestApplyTrackOverride_Success(ts *testing.T) {
	templatePreset := testTemplatePreset(ts)
	derivedPreset := testDerivedPreset(ts, templatePreset)

	tests := []struct {
		name      string
		line      string
		checkFunc func(*testing.T, *t.Preset)
	}{
		{name: "override tone carrier", line: "  track 1 tone 350", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[0].Carrier != 350 { tst.Errorf("expected carrier 350, got %v", p.Track[0].Carrier) } }},
		{name: "override tone relative", line: "  track 1 tone +50", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[0].Carrier != 350 { tst.Errorf("expected carrier 350, got %v", p.Track[0].Carrier) } }},
		{name: "override resonance", line: "  track 1 binaural -3", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[0].Resonance != 7 { tst.Errorf("expected resonance 7, got %v", p.Track[0].Resonance) } }},
		{name: "override amplitude", line: "  track 1 amplitude 30", checkFunc: func(tst *testing.T, p *t.Preset) { expected := t.AmplitudePercentToRaw(30); if p.Track[0].Amplitude != expected { tst.Errorf("expected amplitude %v, got %v", expected, p.Track[0].Amplitude) } }},
		{name: "override waveform", line: "  track 1 waveform triangle", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[0].Waveform != t.WaveformTriangle { tst.Errorf("expected waveform %v, got %v", t.WaveformTriangle, p.Track[0].Waveform) } }},
		{name: "override pan", line: "  track 2 pan +3", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[1].Effect.Value != 8 { tst.Errorf("expected pan value 8, got %v", p.Track[1].Effect.Value) } }},
		{name: "override intensity", line: "  track 2 intensity -15", checkFunc: func(tst *testing.T, p *t.Preset) { expected := t.IntensityPercentToRaw(60); if p.Track[1].Effect.Intensity != expected { tst.Errorf("expected intensity %v, got %v", expected, p.Track[1].Effect.Intensity) } }},
		{name: "override ambiance waveform", line: "  track 2 waveform sawtooth", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[1].Waveform != t.WaveformSawtooth { tst.Errorf("expected waveform %v, got %v", t.WaveformSawtooth, p.Track[1].Waveform) } }},
		{name: "override monaural resonance", line: "  track 3 monaural +2", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[2].Resonance != 10 { tst.Errorf("expected resonance 10, got %v", p.Track[2].Resonance) } }},
		{name: "override smooth", line: "  track 4 smooth -5", checkFunc: func(tst *testing.T, p *t.Preset) { if p.Track[3].NoiseSmooth != 15 { tst.Errorf("expected smooth 15, got %v", p.Track[3].NoiseSmooth) } }},
	}

	for _, tt := range tests {
		derivedPreset.Track = templatePreset.Track
		decl := parseOverrideSpec(ts, tt.line)
		if err := preset.ApplyTrackOverride(derivedPreset, decl); err != nil {
			ts.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		tt.checkFunc(ts, derivedPreset)
	}
}

func TestApplyTrackOverride_Errors(ts *testing.T) {
	templatePreset := testTemplatePreset(ts)
	templatePreset.Track[2] = t.Track{Type: t.TrackAmbiance, AmbianceName: "river", Resonance: 2.5, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectModulation, Intensity: t.IntensityPercentToRaw(60)}}
	templatePreset.Track[3] = t.Track{Type: t.TrackBrownNoise, NoiseSmooth: 15, Amplitude: t.AmplitudePercentToRaw(30)}
	derivedPreset := testDerivedPreset(ts, templatePreset)

	tests := []struct {
		name string
		line string
	}{
		{"override off track", "  track 5 amplitude 10"},
		{"tone on ambiance track", "  track 2 tone 300"},
		{"waveform on noise track", "  track 4 waveform square"},
		{"smooth on non-noise track", "  track 1 smooth 30"},
		{"wrong beat type override", "  track 1 monaural 8"},
		{"modulation on pan effect", "  track 2 modulation 3"},
		{"invalid amplitude", "  track 1 amplitude 150"},
		{"invalid intensity", "  track 2 intensity 150"},
		{"invalid smooth", "  track 4 smooth 150"},
	}

	for _, tt := range tests {
		derivedPreset.Track = templatePreset.Track
		decl := parseOverrideSpec(ts, tt.line)
		if err := preset.ApplyTrackOverride(derivedPreset, decl); err == nil {
			ts.Errorf("%s: expected error but got none for line: %q", tt.name, tt.line)
		}
	}
}

func TestApplyTrackOverride_WithoutFromPreset(ts *testing.T) {
	standalonePreset, err := t.NewPreset("standalone", false, nil)
	if err != nil {
		ts.Fatalf("failed to create preset: %v", err)
	}

	decl := parseOverrideSpec(ts, "  track 1 amplitude 10")
	err = preset.ApplyTrackOverride(standalonePreset, decl)
	if err == nil {
		ts.Errorf("expected error when overriding track on preset without 'from' source")
	}
}