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

func TestOutputFormatFromFileName(ts *testing.T) {
	format, err := outputFormatFromFileName("output.mp3")
	if err != nil {
		ts.Fatalf("unexpected error resolving format: %v", err)
	}
	if format != "mp3" {
		ts.Fatalf("expected mp3 format, got %q", format)
	}
}

func TestOutputFormatFromFileNameRejectsMissingExtension(ts *testing.T) {
	format, err := outputFormatFromFileName("output")
	if err == nil {
		ts.Fatalf("expected missing extension to fail, got format %q", format)
	}
}

func TestOutputFormatFromFileNameRejectsUnsupportedExtension(ts *testing.T) {
	format, err := outputFormatFromFileName("output.wav")
	if err == nil {
		ts.Fatalf("expected unsupported extension to fail, got format %q", format)
	}
}

func TestBuildConvertArgs(ts *testing.T) {
	args := buildConvertArgs(44100, "output.mp3", "mp3")
	expected := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", "44100",
		"-i", "pipe:0",
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-f", "mp3",
		"-vn",
		"output.mp3",
	}
	if !reflect.DeepEqual(args, expected) {
		ts.Fatalf("unexpected ffmpeg args:\nwant: %#v\n got: %#v", expected, args)
	}
}