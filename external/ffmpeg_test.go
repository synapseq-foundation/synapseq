// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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