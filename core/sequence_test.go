// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testSequenceContent = `
@samplerate 48000
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
`

func TestAppContext_LoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "seq.spsq")
	content := strings.TrimSpace(`
@ambiance rain audio/rain
@music meditation audio/meditation
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
`) + "\n"

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp sequence: %v", err)
	}
	ambiancePath := filepath.Join(dir, "audio", "rain.wav")
	if err := os.MkdirAll(filepath.Dir(ambiancePath), 0o755); err != nil {
		t.Fatalf("create temp ambiance dir: %v", err)
	}
	if err := os.WriteFile(ambiancePath, nil, 0o600); err != nil {
		t.Fatalf("write temp ambiance: %v", err)
	}
	musicPath := filepath.Join(dir, "audio", "meditation.mp3")
	if err := os.WriteFile(musicPath, nil, 0o600); err != nil {
		t.Fatalf("write temp music: %v", err)
	}

	loaded, err := NewAppContext().LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile error: %v", err)
	}

	if string(loaded.RawContent()) != content {
		t.Fatalf("LoadFile did not preserve raw content")
	}

	wantAmbiancePath := filepath.Join(dir, "audio", "rain.wav")
	if loaded.Ambiance()["rain"] != wantAmbiancePath {
		t.Fatalf("expected ambiance %q, got %q", wantAmbiancePath, loaded.Ambiance()["rain"])
	}
	wantMusicPath := filepath.Join(dir, "audio", "meditation.mp3")
	if loaded.Music()["meditation"] != wantMusicPath {
		t.Fatalf("expected music %q, got %q", wantMusicPath, loaded.Music()["meditation"])
	}
}

func TestAppContext_LoadFileMissing(t *testing.T) {
	_, err := NewAppContext().LoadFile(filepath.Join(t.TempDir(), "missing.spsq"))
	if err == nil {
		t.Fatalf("expected missing file error")
	}
}

func TestAppContext_LoadContent(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(testSequenceContent)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	if loaded.SampleRate() != 48000 {
		t.Fatalf("expected sample rate 48000, got %d", loaded.SampleRate())
	}
}

func TestLoadedContext_JSON(t *testing.T) {
	loaded, err := NewAppContext().LoadContent(testSequenceContent)
	if err != nil {
		t.Fatalf("LoadContent error: %v", err)
	}

	content, err := loaded.JSON()
	if err != nil {
		t.Fatalf("JSON error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(content, &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, content)
	}

	options := got["options"].(map[string]any)
	if options["sampleRate"] != float64(48000) {
		t.Fatalf("expected sampleRate 48000, got %#v", options["sampleRate"])
	}

	presets := got["presets"].([]any)
	if len(presets) == 0 {
		t.Fatalf("expected presets in JSON: %s", content)
	}

	alpha := presets[0].(map[string]any)
	if alpha["name"] != "alpha" {
		t.Fatalf("expected alpha preset in JSON: %s", content)
	}

	tracks := alpha["tracks"].([]any)
	if len(tracks) == 0 {
		t.Fatalf("expected alpha tracks in JSON: %s", content)
	}

	track := tracks[0].(map[string]any)
	if track["type"] != "binaural" || track["resonance"] != float64(1) {
		t.Fatalf("expected core-style track in JSON: %#v", track)
	}

	timeline := got["timeline"].([]any)
	entry := timeline[0].(map[string]any)
	if entry["presetName"] != "alpha" || entry["timestamp"] != "00:00:00" {
		t.Fatalf("expected core-style timeline in JSON: %#v", entry)
	}
}
