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
	"os"
	"path/filepath"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func writePresetFile(tst *testing.T, name string, content string) string {
	tst.Helper()
	dir := tst.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content+"\n"), 0o600); err != nil {
		tst.Fatalf("write temp preset file: %v", err)
	}
	return p
}

func eqTrackGotWantLP(got, want t.Track) bool {
	if got.Type != want.Type {
		return false
	}
	if got.Carrier != want.Carrier {
		return false
	}
	if got.Resonance != want.Resonance {
		return false
	}
	if got.Amplitude != want.Amplitude {
		return false
	}
	if got.Waveform != want.Waveform {
		return false
	}
	if got.Effect.Type != want.Effect.Type {
		return false
	}
	if got.Effect.Intensity != want.Effect.Intensity {
		return false
	}
	return true
}

func hasTrackLP(tracks [t.NumberOfChannels]t.Track, want t.Track) bool {
	for ch := range t.NumberOfChannels {
		if eqTrackGotWantLP(tracks[ch], want) {
			return true
		}
	}
	return false
}

func TestLoadPresets_Success(ts *testing.T) {
	content := `
# Presets
preparation
  noise pink amplitude 50
  tone 300 binaural 10 amplitude 10
`
	path := writePresetFile(ts, "presets.spsq", content)

	presets, err := loadPresets(path)
	if err != nil {
		ts.Fatalf("loadPresets error: %v", err)
	}
	if len(presets) != 1 {
		ts.Fatalf("expected 1 preset, got %d", len(presets))
	}
	if presets[0].String() != "preparation" {
		ts.Fatalf("unexpected preset name: %q", presets[0].String())
	}

	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(50),
	}
	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(10),
	}

	if !(hasTrackLP(presets[0].Track, wantNoise) && hasTrackLP(presets[0].Track, wantTone)) {
		ts.Fatalf("missing tracks in loaded preset: %+v", presets[0].Track)
	}
}

func TestLoadPresets_Errors(ts *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "track_before_preset",
			content: `
  tone 100 binaural 1 amplitude 1
alpha
`,
		},
		{
			name: "no_presets_defined",
			content: `
## only comments
`,
		},
		{
			name: "unexpected_content",
			content: `
alpha
00:00:00 alpha
`,
		},
	}

	for _, tt := range tests {
		path := writePresetFile(ts, tt.name+".spsq", tt.content)
		if _, err := loadPresets(path); err == nil {
			ts.Fatalf("%s: expected error, got nil", tt.name)
		}
	}
}
