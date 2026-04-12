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

package external

import (
	"reflect"
	"testing"
)

func TestNewFFPlayUnsafeUsesDefaultPath(ts *testing.T) {
	player := NewFFPlayUnsafe("")
	if player.Path() != "ffplay" {
		ts.Fatalf("expected default ffplay path, got %q", player.Path())
	}
}

func TestBuildPlayArgs(ts *testing.T) {
	args := buildPlayArgs(44100)
	expected := []string{
		"-nodisp",
		"-hide_banner",
		"-loglevel", "error",
		"-autoexit",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", "44100",
		"-i", "pipe:0",
	}
	if !reflect.DeepEqual(args, expected) {
		ts.Fatalf("unexpected ffplay args:\nwant: %#v\n got: %#v", expected, args)
	}
}