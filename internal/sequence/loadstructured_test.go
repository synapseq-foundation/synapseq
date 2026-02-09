//go:build !wasm

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

package sequence

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

func writeTemp(ts *testing.T, name, content string) string {
	ts.Helper()
	dir := ts.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		ts.Fatalf("write temp file %s: %v", name, err)
	}
	return p
}

func assertIncreasing(ts *testing.T, times []int) {
	ts.Helper()
	if len(times) == 0 {
		ts.Fatalf("empty times")
	}
	if times[0] != 0 {
		ts.Fatalf("first time must be 0, got %d", times[0])
	}
	for i := 1; i < len(times); i++ {
		if !(times[i] > times[i-1]) {
			ts.Fatalf("times not strictly increasing at %d: %v", i, times)
		}
	}
}

func periodTimes(periods []t.Period) []int {
	out := make([]int, len(periods))
	for i, p := range periods {
		out[i] = p.Time
	}
	return out
}

func hasTrackIn(tracks [t.NumberOfChannels]t.Track, want t.Track) bool {
	for i := range tracks {
		got := tracks[i]
		if got.Type == want.Type &&
			got.Carrier == want.Carrier &&
			got.Resonance == want.Resonance &&
			got.Amplitude == want.Amplitude &&
			got.Waveform == want.Waveform {
			return true
		}
	}
	return false
}

func TestLoadStructured_JSON_Standalone(ts *testing.T) {
	json := `{
  "description": ["Standalone structured test"],
  "options": { "samplerate": 44100, "volume": 100 },
  "sequence": [
    {
      "time": 0,
	  "transition": "steady",
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 0, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 0 }
        ]
      }
    },
    {
      "time": 15000,
	  "transition": "steady",
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 15, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 30 }
        ]
      }
    }
  ]
}`
	p := writeTemp(ts, "seq.json", json)

	res, err := LoadStructuredSequence(p, t.FormatJSON)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(json) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}

func TestLoadStructured_XML_Standalone(ts *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SynapSeqInput>
  <description>
    <line>Standalone structured test</line>
  </description>
  <options>
    <samplerate>44100</samplerate>
    <volume>100</volume>
  </options>
  <sequence>
    <entry time="0" transition="steady">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="0" waveform="sine"/>
        <noise mode="pink" amplitude="0"/>
      </track>
    </entry>
    <entry time="15000" transition="steady">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="15" waveform="sine"/>
        <noise mode="pink" amplitude="30"/>
      </track>
    </entry>
  </sequence>
</SynapSeqInput>`
	p := writeTemp(ts, "seq.xml", xml)

	res, err := LoadStructuredSequence(p, t.FormatXML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(xml) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}

func TestLoadStructured_YAML_Standalone(ts *testing.T) {
	yaml := `description:
  - Standalone structured test
options:
  samplerate: 44100
  volume: 100
sequence:
  - time: 0
    transition: steady
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 0
          waveform: sine
      noises:
        - mode: pink
          amplitude: 0
  - time: 15000
    transition: steady
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 15
          waveform: sine
      noises:
        - mode: pink
          amplitude: 30
`
	p := writeTemp(ts, "seq.yaml", yaml)

	res, err := LoadStructuredSequence(p, t.FormatYAML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(yaml) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}

// Helpers for tests
func sampleContentFor(format t.FileFormat) string {
	switch format {
	case t.FormatJSON:
		return `{
  "description": ["Standalone structured test"],
  "options": { "samplerate": 44100, "volume": 100 },
  "sequence": [
    {
      "time": 0,
	  "transition": "steady",
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 0, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 0 }
        ]
      }
    },
    {
      "time": 15000,
	  "transition": "steady",
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 15, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 30 }
        ]
      }
    }
  ]
}`
	case t.FormatXML:
		return `<?xml version="1.0" encoding="UTF-8"?>
<SynapSeqInput>
  <description>
    <line>Standalone structured test</line>
  </description>
  <options>
    <samplerate>44100</samplerate>
    <volume>100</volume>
  </options>
  <sequence>
    <entry time="0" transition="steady">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="0" waveform="sine"/>
        <noise mode="pink" amplitude="0"/>
      </track>
    </entry>
    <entry time="15000" transition="steady">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="15" waveform="sine"/>
        <noise mode="pink" amplitude="30"/>
      </track>
    </entry>
  </sequence>
</SynapSeqInput>`
	case t.FormatYAML:
		return `description:
  - Standalone structured test
options:
  samplerate: 44100
  volume: 100
sequence:
  - time: 0
    transition: steady
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 0
          waveform: sine
      noises:
        - mode: pink
          amplitude: 0
  - time: 15000
    transition: steady
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 15
          waveform: sine
      noises:
        - mode: pink
          amplitude: 30
`
	default:
		return ""
	}
}

func verifyBasicLoadResult(tst *testing.T, res *t.Sequence) {
	tst.Helper()
	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		tst.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		tst.Fatalf("expected non-empty description/comments")
	}
	if len(res.Periods) != 2 {
		tst.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(tst, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		tst.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		tst.Fatalf("missing expected noise in period[1]")
	}
}

// STDIN tests

func TestLoadStructured_JSON_FromStdin(ts *testing.T) {
	content := sampleContentFor(t.FormatJSON)
	r, w, err := os.Pipe()
	if err != nil {
		ts.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = w.Write([]byte(content))
		_ = w.Close()
	}()

	res, err := LoadStructuredSequence("-", t.FormatJSON)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(stdin,json): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

func TestLoadStructured_XML_FromStdin(ts *testing.T) {
	content := sampleContentFor(t.FormatXML)
	r, w, err := os.Pipe()
	if err != nil {
		ts.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = w.Write([]byte(content))
		_ = w.Close()
	}()

	res, err := LoadStructuredSequence("-", t.FormatXML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(stdin,xml): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

func TestLoadStructured_YAML_FromStdin(ts *testing.T) {
	content := sampleContentFor(t.FormatYAML)
	r, w, err := os.Pipe()
	if err != nil {
		ts.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = w.Write([]byte(content))
		_ = w.Close()
	}()

	res, err := LoadStructuredSequence("-", t.FormatYAML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(stdin,yaml): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

// WEB tests (http server local)

func TestLoadStructured_JSON_FromHTTP(ts *testing.T) {
	content := sampleContentFor(t.FormatJSON)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(content))
	}))
	defer srv.Close()

	res, err := LoadStructuredSequence(srv.URL, t.FormatJSON)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(http,json): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

func TestLoadStructured_XML_FromHTTP(ts *testing.T) {
	content := sampleContentFor(t.FormatXML)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(content))
	}))
	defer srv.Close()

	res, err := LoadStructuredSequence(srv.URL, t.FormatXML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(http,xml): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

func TestLoadStructured_YAML_FromHTTP(ts *testing.T) {
	content := sampleContentFor(t.FormatYAML)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// intentionally using x-yaml content type
		// to test content-type flexibility
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write([]byte(content))
	}))
	defer srv.Close()

	res, err := LoadStructuredSequence(srv.URL, t.FormatYAML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(http,yaml): %v", err)
	}
	verifyBasicLoadResult(ts, res)
}

func TestLoadStructured_JSON_FileTooLarge(ts *testing.T) {
	// Generate a JSON file larger than maxStructuredFileSize
	over := t.MaxStructuredFileSize + 1024 // 1KB over the limit
	huge := strings.Repeat("A", over)
	json := fmt.Sprintf(`{"description":["%s"],"options":{"samplerate":44100,"volume":100},"sequence":[{"time":0,"transition":"steady","track":{"tones":[{"mode":"binaural","carrier":250,"resonance":8,"amplitude":0,"waveform":"sine"}]}}]}`, huge)

	p := writeTemp(ts, "too-big.json", json)

	if _, err := LoadStructuredSequence(p, t.FormatJSON); err == nil {
		ts.Fatalf("expected error for file > %d bytes, got nil", t.MaxStructuredFileSize)
	}
}

func TestLoadStructured_JSON_WithBackground(ts *testing.T) {
	json := `{
  "description": ["Background test with pulse effect"],
  "options": {
    "samplerate": 44100,
    "volume": 100,
    "background": "sounds/pink-noise.wav",
    "gainlevel": "high"
  },
  "sequence": [
    {
      "time": 0,
      "transition": "steady",
      "track": {
        "tones": [
          { "mode": "monaural", "carrier": 300, "resonance": 10, "amplitude": 0, "waveform": "sine" }
        ],
        "background": {
          "amplitude": 0,
          "waveform": "sine",
          "effect": {
            "intensity": 50,
            "pulse": {
              "resonance": 10
            }
          }
        }
      }
    },
    {
      "time": 15000,
      "transition": "smooth",
      "track": {
        "tones": [
          { "mode": "monaural", "carrier": 300, "resonance": 10, "amplitude": 15, "waveform": "sine" }
        ],
        "background": {
          "amplitude": 50,
          "waveform": "sine",
          "effect": {
            "intensity": 50,
            "pulse": {
              "resonance": 10
            }
          }
        }
      }
    }
  ]
}`
	p := writeTemp(ts, "bg-pulse.json", json)

	res, err := LoadStructuredSequence(p, t.FormatJSON)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(json with background) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if !strings.HasSuffix(res.Options.BackgroundPath, "sounds/pink-noise.wav") {
		ts.Fatalf("expected background to end with 'sounds/pink-noise.wav', got %q", res.Options.BackgroundPath)
	}
	if res.Options.GainLevel != t.GainLevelHigh {
		ts.Fatalf("expected gainlevel high, got %v", res.Options.GainLevel)
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	// Check transition type
	if res.Periods[1].Transition != t.TransitionSmooth {
		ts.Fatalf("expected smooth transition, got %v", res.Periods[1].Transition)
	}

	// Verify background track with pulse effect
	found := false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Amplitude == t.AmplitudePercentToRaw(50) &&
			track.Effect.Type == t.EffectModulation &&
			track.Resonance == 10 {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track with modulation effect in period[1]")
	}
}

func TestLoadStructured_XML_WithBackground(ts *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SynapSeqInput>
  <description>
    <line>Background test with pan effect</line>
  </description>
  <options>
    <samplerate>44100</samplerate>
    <volume>100</volume>
    <background>sounds/pink-noise.wav</background>
    <gainlevel>high</gainlevel>
  </options>
  <sequence>
    <entry time="0" transition="steady">
      <track>
        <tone mode="binaural" carrier="300" resonance="10" amplitude="0" waveform="sine"></tone>
        <background amplitude="0" waveform="sine">
          <effect intensity="45">
            <spin width="400" rate="10"></spin>
          </effect>
        </background>
      </track>
    </entry>
    <entry time="15000" transition="ease-in">
      <track>
        <tone mode="binaural" carrier="300" resonance="10" amplitude="15" waveform="sine"></tone>
        <background amplitude="50" waveform="sine">
          <effect intensity="45">
            <spin width="400" rate="10"></spin>
          </effect>
        </background>
      </track>
    </entry>
  </sequence>
</SynapSeqInput>`
	p := writeTemp(ts, "bg-spin.xml", xml)

	res, err := LoadStructuredSequence(p, t.FormatXML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(xml with background) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if !strings.HasSuffix(res.Options.BackgroundPath, "sounds/pink-noise.wav") {
		ts.Fatalf("expected background to end with 'sounds/pink-noise.wav', got %q", res.Options.BackgroundPath)
	}
	if res.Options.GainLevel != t.GainLevelHigh {
		ts.Fatalf("expected gainlevel high, got %v", res.Options.GainLevel)
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	// Check transition type
	if res.Periods[1].Transition != t.TransitionEaseIn {
		ts.Fatalf("expected ease-in transition, got %v", res.Periods[1].Transition)
	}

	// Verify background track with spin effect
	found := false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Amplitude == t.AmplitudePercentToRaw(50) &&
			track.Effect.Type == t.EffectPan &&
			track.Carrier == 400 &&
			track.Resonance == 10 {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track with pan effect in period[1]")
	}
}

func TestLoadStructured_YAML_WithBackground(ts *testing.T) {
	yaml := `description:
  - Background test with pulse effect
options:
  samplerate: 44100
  volume: 100
  background: sounds/pink-noise.wav
  gainlevel: high
sequence:
  - time: 0
    transition: steady
    track:
      tones:
        - mode: monaural
          carrier: 300
          resonance: 10
          amplitude: 0
          waveform: sine
      background:
        amplitude: 0
        waveform: sine
        effect:
          intensity: 50
          pulse:
            resonance: 10
  - time: 15000
    transition: ease-out
    track:
      tones:
        - mode: monaural
          carrier: 300
          resonance: 10
          amplitude: 15
          waveform: sine
      background:
        amplitude: 50
        waveform: sine
        effect:
          intensity: 50
          pulse:
            resonance: 10
`
	p := writeTemp(ts, "bg-pulse.yaml", yaml)

	res, err := LoadStructuredSequence(p, t.FormatYAML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(yaml with background) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if !strings.HasSuffix(res.Options.BackgroundPath, "sounds/pink-noise.wav") {
		ts.Fatalf("expected background to end with 'sounds/pink-noise.wav', got %q", res.Options.BackgroundPath)
	}
	if res.Options.GainLevel != t.GainLevelHigh {
		ts.Fatalf("expected gainlevel high, got %v", res.Options.GainLevel)
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	// Check transition type
	if res.Periods[1].Transition != t.TransitionEaseOut {
		ts.Fatalf("expected ease-out transition, got %v", res.Periods[1].Transition)
	}

	// Verify background track with pulse effect
	found := false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Amplitude == t.AmplitudePercentToRaw(50) &&
			track.Effect.Type == t.EffectModulation &&
			track.Resonance == 10 {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track with modulation effect in period[1]")
	}
}

func TestLoadStructured_JSON_BackgroundWithVaryingIntensity(ts *testing.T) {
	json := `{
  "description": ["Background test with varying pulse intensity"],
  "options": {
    "samplerate": 44100,
    "volume": 100,
    "background": "sounds/pink-noise.wav"
  },
  "sequence": [
    {
      "time": 0,
      "transition": "steady",
      "track": {
        "background": {
          "amplitude": 50,
          "waveform": "sine",
          "effect": {
            "intensity": 60,
            "pulse": {
              "resonance": 8
            }
          }
        }
      }
    },
    {
      "time": 30000,
      "transition": "smooth",
      "track": {
        "background": {
          "amplitude": 50,
          "waveform": "sine",
          "effect": {
            "intensity": 30,
            "pulse": {
              "resonance": 8
            }
          }
        }
      }
    }
  ]
}`
	p := writeTemp(ts, "bg-varying-intensity.json", json)

	res, err := LoadStructuredSequence(p, t.FormatJSON)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(json with varying intensity) error: %v", err)
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}

	// Verify first period has higher intensity
	found := false
	for _, track := range res.Periods[0].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Effect.Type == t.EffectModulation &&
			track.Effect.Intensity == t.IntensityPercentToRaw(60) {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track with high intensity in period[0]")
	}

	// Verify second period has lower intensity
	found = false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Effect.Type == t.EffectModulation &&
			track.Effect.Intensity == t.IntensityPercentToRaw(30) {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track with low intensity in period[1]")
	}
}

func TestLoadStructured_XML_BackgroundWithoutEffect(ts *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SynapSeqInput>
  <description>
    <line>Background test without effects</line>
  </description>
  <options>
    <samplerate>44100</samplerate>
    <volume>100</volume>
    <background>sounds/pink-noise.wav</background>
  </options>
  <sequence>
    <entry time="0" transition="steady">
      <track>
        <background amplitude="0" waveform="square"></background>
      </track>
    </entry>
    <entry time="15000" transition="steady">
      <track>
        <background amplitude="40" waveform="square"></background>
      </track>
    </entry>
  </sequence>
</SynapSeqInput>`
	p := writeTemp(ts, "bg-no-effect.xml", xml)

	res, err := LoadStructuredSequence(p, t.FormatXML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(xml background without effect) error: %v", err)
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}

	// Verify background track without effect
	found := false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Amplitude == t.AmplitudePercentToRaw(40) &&
			track.Waveform == t.WaveformSquare &&
			track.Effect.Type == t.EffectOff {
			found = true
			break
		}
	}
	if !found {
		ts.Fatalf("missing expected background track without effects")
	}
}

func TestLoadStructured_YAML_BackgroundWithDifferentIntensities(ts *testing.T) {
	yaml := `description:
  - Background test with varying intensities on same waveform
options:
  samplerate: 44100
  volume: 100
  background: sounds/pink-noise.wav
sequence:
  - time: 0
    transition: steady
    track:
      background:
        amplitude: 0
        waveform: sine
        effect:
          intensity: 40
          pulse:
            resonance: 6
  - time: 15000
    transition: smooth
    track:
      background:
        amplitude: 30
        waveform: sine
        effect:
          intensity: 40
          pulse:
            resonance: 6
  - time: 30000
    transition: smooth
    track:
      background:
        amplitude: 30
        waveform: sine
        effect:
          intensity: 70
          pulse:
            resonance: 6
`
	p := writeTemp(ts, "bg-intensities.yaml", yaml)

	res, err := LoadStructuredSequence(p, t.FormatYAML)
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(yaml with intensities) error: %v", err)
	}

	if len(res.Periods) != 3 {
		ts.Fatalf("expected 3 periods, got %d", len(res.Periods))
	}

	// Check period[1] has lower intensity
	foundLowIntensity := false
	for _, track := range res.Periods[1].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Waveform == t.WaveformSine &&
			track.Effect.Type == t.EffectModulation &&
			track.Effect.Intensity == t.IntensityPercentToRaw(40) &&
			track.Amplitude > 0 {
			foundLowIntensity = true
			break
		}
	}
	if !foundLowIntensity {
		ts.Fatalf("missing expected low intensity background in period[1]")
	}

	// Check period[2] has higher intensity
	foundHighIntensity := false
	for _, track := range res.Periods[2].TrackStart {
		if track.Type == t.TrackBackground &&
			track.Waveform == t.WaveformSine &&
			track.Effect.Type == t.EffectModulation &&
			track.Effect.Intensity == t.IntensityPercentToRaw(70) &&
			track.Amplitude > 0 {
			foundHighIntensity = true
			break
		}
	}
	if !foundHighIntensity {
		ts.Fatalf("missing expected high intensity background in period[2]")
	}
}
