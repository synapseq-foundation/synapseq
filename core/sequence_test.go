//go:build !wasm

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

package core

import (
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
