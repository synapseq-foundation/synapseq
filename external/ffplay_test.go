// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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