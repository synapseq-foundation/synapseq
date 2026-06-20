// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dump

import (
	"encoding/json"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestJSONSerializesLoadedSequence(ts *testing.T) {
	alpha, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("new preset: %v", err)
	}

	alpha.Track[0] = t.Track{
		Type:        t.TrackPinkNoise,
		Amplitude:   t.AmplitudePercentToRaw(40),
		Waveform:    t.WaveformSine,
		NoiseSmooth: 0,
		Effect:      t.Effect{Type: t.EffectOff},
	}
	alpha.Track[1] = t.Track{
		Type:       t.TrackAmbiance,
		Amplitude:  t.AmplitudePercentToRaw(15),
		SourceName: "rain",
		Waveform:   t.WaveformSine,
		Effect: t.Effect{
			Type:      t.EffectPan,
			Value:     0.3,
			Intensity: t.IntensityPercentToRaw(10),
		},
	}
	alpha.Track[2] = t.Track{
		Type:       t.TrackMusic,
		Amplitude:  t.AmplitudePercentToRaw(50),
		SourceName: "meditation_song",
		Waveform:   t.WaveformSine,
		Effect: t.Effect{
			Type:      t.EffectModulation,
			Value:     10,
			Intensity: t.IntensityPercentToRaw(15),
		},
	}
	alpha.Track[3] = t.Track{
		Type:      t.TrackBinauralBeat,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
		Carrier:   300,
		Resonance: 10,
		Effect:    t.Effect{Type: t.EffectOff},
	}

	seq := &t.Sequence{
		Comments: []string{},
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
			Ambiance:   map[string]string{"rain": "sounds/rain"},
			Music:      map[string]string{"meditation_song": "musics/med_-_45"},
			Extends:    []string{"common/presets", "common/options"},
		},
		Presets: []t.Preset{*t.NewBuiltinSilencePreset(), *alpha},
		Periods: []t.Period{
			{Time: 0, PresetName: "alpha", Transition: t.TransitionSteady},
		},
	}

	content, err := JSON(seq)
	if err != nil {
		ts.Fatalf("JSON error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(content, &got); err != nil {
		ts.Fatalf("invalid JSON: %v\n%s", err, content)
	}

	options := got["options"].(map[string]any)
	if options["sampleRate"] != float64(44100) {
		ts.Fatalf("expected sampleRate 44100, got %#v", options["sampleRate"])
	}
	if got := options["ambiance"].(map[string]any)["rain"]; got != "sounds/rain" {
		ts.Fatalf("expected ambiance rain source, got %#v", got)
	}
	if got := options["music"].(map[string]any)["meditation_song"]; got != "musics/med_-_45" {
		ts.Fatalf("expected music meditation_song source, got %#v", got)
	}
	if _, ok := options["sources"]; ok {
		ts.Fatalf("did not expect legacy sources object: %#v", options["sources"])
	}

	presets := got["presets"].([]any)
	if len(presets) != 1 {
		ts.Fatalf("expected 1 preset, got %d", len(presets))
	}

	alphaJSON := presets[0].(map[string]any)
	if alphaJSON["name"] != "alpha" {
		ts.Fatalf("expected alpha preset, got %#v", alphaJSON["name"])
	}

	alphaTracks := alphaJSON["tracks"].([]any)
	if len(alphaTracks) != 4 {
		ts.Fatalf("expected 4 alpha tracks, got %d", len(alphaTracks))
	}

	musicTrack := alphaTracks[2].(map[string]any)
	if musicTrack["type"] != "music" || musicTrack["sourceName"] != "meditation_song" || musicTrack["amplitude"] != float64(50) {
		ts.Fatalf("unexpected music track: %#v", musicTrack)
	}

	binauralTrack := alphaTracks[3].(map[string]any)
	if binauralTrack["index"] != float64(4) || binauralTrack["type"] != "binaural" || binauralTrack["resonance"] != float64(10) {
		ts.Fatalf("unexpected binaural track: %#v", binauralTrack)
	}

	effect := musicTrack["effect"].(map[string]any)
	if effect["type"] != "modulation" || effect["value"] != float64(10) || effect["intensity"] != float64(15) {
		ts.Fatalf("unexpected music effect: %#v", effect)
	}

	timeline := got["timeline"].([]any)
	entry := timeline[0].(map[string]any)
	if entry["presetName"] != "alpha" || entry["timestamp"] != "00:00:00" {
		ts.Fatalf("unexpected timeline entry: %#v", entry)
	}

	for _, legacy := range []string{`"family"`, `"kind"`, `"beat"`, `"smooth"`, `"source"`, `"loop"`, `"startSeconds"`} {
		if strings.Contains(string(content), legacy) {
			ts.Fatalf("did not expect legacy field %s in JSON:\n%s", legacy, content)
		}
	}
}

func TestJSONRejectsNilSequence(ts *testing.T) {
	if _, err := JSON(nil); err == nil {
		ts.Fatal("expected nil sequence error")
	}
}
