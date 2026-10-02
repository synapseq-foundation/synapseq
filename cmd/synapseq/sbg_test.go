// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/sbg"
)

func TestRunSBGConversionUsesDefaultOutputAndWarns(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var status bytes.Buffer
	if err := runSBGConversion([]string{inputPath}, &cli.CLIOptions{}, &status, nil); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !strings.Contains(status.String(), "conversion is approximate") || !strings.Contains(status.String(), "Converted:") || !strings.Contains(status.String(), inputPath[:len(inputPath)-len(".sbg")]+".spsq") {
		t.Fatalf("status = %q", status.String())
	}
	outputPath := filepath.Join(dir, "session.spsq")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(content), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected output:\n%s", content)
	}
}

func TestRunSBGConversionOverwritesOutput(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	outputPath := filepath.Join(dir, "custom.spsq")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("old"), 0o600); err != nil {
		t.Fatalf("write old output: %v", err)
	}
	if err := runSBGConversion([]string{inputPath, outputPath}, &cli.CLIOptions{}, nil, nil); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(content) == "old" {
		t.Fatal("output was not overwritten")
	}
}

func TestRunSBGConversionWritesToStandardOutput(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var output bytes.Buffer
	if err := runSBGConversion([]string{inputPath, "-"}, &cli.CLIOptions{}, nil, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !strings.Contains(output.String(), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected standard output:\n%s", output.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "session.spsq")); !os.IsNotExist(err) {
		t.Fatalf("default output file was created: %v", err)
	}
}

func TestRunSBGConversionQuietSuppressesStatus(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var status bytes.Buffer
	var output bytes.Buffer
	if err := runSBGConversion([]string{inputPath, "-"}, &cli.CLIOptions{Quiet: true}, &status, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if status.Len() != 0 {
		t.Fatalf("quiet status = %q", status.String())
	}
	if !strings.Contains(output.String(), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected standard output:\n%s", output.String())
	}
}

func TestRunSBGConversionExportsMP3(t *testing.T) {
	tests := []struct {
		name               string
		outputName         string
		useMP3Flag         bool
		expectedOutputName string
	}{
		{name: "explicit MP3 target", outputName: "custom.mp3", expectedOutputName: "custom.mp3"},
		{name: "default MP3 target", useMP3Flag: true, expectedOutputName: "session.mp3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			inputPath := filepath.Join(dir, "session.sbg")
			if err := os.WriteFile(inputPath, []byte("alpha: 300+10/20\noff: -\n00:00:00 alpha\n+00:00:01 off\n"), 0o600); err != nil {
				t.Fatalf("write input: %v", err)
			}

			audioCapture := filepath.Join(dir, "audio.raw")
			outputCapture := filepath.Join(dir, "output-path")
			t.Setenv("SYNAPSEQ_TEST_FFMPEG_CAPTURE", audioCapture)
			t.Setenv("SYNAPSEQ_TEST_FFMPEG_OUTPUT", outputCapture)

			ffmpegPath := filepath.Join(dir, "ffmpeg")
			ffmpegScript := "#!/bin/sh\nset -e\noutput=\nfor argument do output=\"$argument\"; done\ncat > \"$SYNAPSEQ_TEST_FFMPEG_CAPTURE\"\nprintf '%s' \"$output\" > \"$SYNAPSEQ_TEST_FFMPEG_OUTPUT\"\n: > \"$output\"\n"
			if err := os.WriteFile(ffmpegPath, []byte(ffmpegScript), 0o700); err != nil {
				t.Fatalf("write fake ffmpeg: %v", err)
			}

			opts := &cli.CLIOptions{Mp3: test.useMP3Flag, FFmpegPath: ffmpegPath}
			args := []string{inputPath}
			if test.outputName != "" {
				args = append(args, filepath.Join(dir, test.outputName))
			}
			var status bytes.Buffer
			if err := runSBGConversion(args, opts, &status, nil); err != nil {
				t.Fatalf("runSBGConversion error: %v", err)
			}

			expectedOutputPath := filepath.Join(dir, test.expectedOutputName)
			if !strings.Contains(status.String(), "review the converted output") || !strings.Contains(status.String(), "Converted:") || !strings.Contains(status.String(), expectedOutputPath) {
				t.Fatalf("unexpected status: %q", status.String())
			}
			outputPath, err := os.ReadFile(outputCapture)
			if err != nil {
				t.Fatalf("read ffmpeg output path: %v", err)
			}
			if string(outputPath) != expectedOutputPath {
				t.Fatalf("ffmpeg output path = %q, want %q", outputPath, expectedOutputPath)
			}
			audio, err := os.ReadFile(audioCapture)
			if err != nil {
				t.Fatalf("read captured audio: %v", err)
			}
			if len(audio) == 0 {
				t.Fatal("no audio was streamed to ffmpeg")
			}
			if _, err := os.Stat(expectedOutputPath); err != nil {
				t.Fatalf("MP3 output was not created: %v", err)
			}
			if _, err := os.Stat(filepath.Join(dir, "session.spsq")); !os.IsNotExist(err) {
				t.Fatalf("intermediate SPSQ file was created: %v", err)
			}
		})
	}
}

func TestRunSBGConversionStreamsPCMWhenMP3TargetsStandardOutput(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	if err := os.WriteFile(inputPath, []byte("alpha: 300+10/20\noff: -\n00:00:00 alpha\n+00:00:01 off\n"), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	converter, err := sbg.New(synapseq.NewAppContext())
	if err != nil {
		t.Fatalf("create SBG converter: %v", err)
	}
	loaded, err := converter.LoadFile(inputPath)
	if err != nil {
		t.Fatalf("load expected sequence: %v", err)
	}
	var expectedPCM bytes.Buffer
	if err := loaded.Stream(&expectedPCM); err != nil {
		t.Fatalf("stream expected PCM: %v", err)
	}

	var output bytes.Buffer
	if err := runSBGConversion([]string{inputPath, "-"}, &cli.CLIOptions{Mp3: true}, nil, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !bytes.Equal(output.Bytes(), expectedPCM.Bytes()) {
		t.Fatal("standard output did not contain the unencoded PCM stream")
	}
}

func TestRunSBGConversionPlaysWithoutWritingSequence(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	if err := os.WriteFile(inputPath, []byte("alpha: 300+10/20\noff: -\n00:00:00 alpha\n+00:00:01 off\n"), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	playbackMarker := filepath.Join(dir, "played")
	audioCapture := filepath.Join(dir, "audio.raw")
	for name, value := range map[string]string{
		"SYNAPSEQ_TEST_FFPLAY_MARKER":  playbackMarker,
		"SYNAPSEQ_TEST_FFPLAY_CAPTURE": audioCapture,
	} {
		original, wasSet := os.LookupEnv(name)
		if err := os.Setenv(name, value); err != nil {
			t.Fatalf("set %s: %v", name, err)
		}
		name, original := name, original
		t.Cleanup(func() {
			if wasSet {
				_ = os.Setenv(name, original)
			} else {
				_ = os.Unsetenv(name)
			}
		})
	}

	ffplayPath := filepath.Join(dir, "ffplay")
	ffplayScript := "#!/bin/sh\ncat > \"$SYNAPSEQ_TEST_FFPLAY_CAPTURE\"\nprintf played > \"$SYNAPSEQ_TEST_FFPLAY_MARKER\"\n"
	if err := os.WriteFile(ffplayPath, []byte(ffplayScript), 0o700); err != nil {
		t.Fatalf("write fake ffplay: %v", err)
	}

	var status bytes.Buffer
	var output bytes.Buffer
	opts := &cli.CLIOptions{Play: true, Mp3: true, FFplayPath: ffplayPath}
	if err := runSBGConversion([]string{inputPath}, opts, &status, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !strings.Contains(status.String(), "conversion is approximate") {
		t.Fatalf("expected conversion warning, got %q", status.String())
	}
	if strings.Contains(status.String(), "Converted:") {
		t.Fatalf("unexpected conversion success message during playback: %q", status.String())
	}
	if output.Len() != 0 {
		t.Fatalf("unexpected sequence output: %q", output.String())
	}
	if _, err := os.Stat(playbackMarker); err != nil {
		t.Fatalf("fake ffplay was not invoked: %v", err)
	}
	audio, err := os.ReadFile(audioCapture)
	if err != nil {
		t.Fatalf("read captured audio: %v", err)
	}
	if len(audio) == 0 {
		t.Fatal("no audio was streamed to ffplay")
	}
	if _, err := os.Stat(filepath.Join(dir, "session.spsq")); !os.IsNotExist(err) {
		t.Fatalf("converted SPSQ file was created during playback: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "session.mp3")); !os.IsNotExist(err) {
		t.Fatalf("MP3 file was created during playback: %v", err)
	}
}

func TestRunSBGConversionRejectsOutputPathDuringPlayback(t *testing.T) {
	var status bytes.Buffer
	err := runSBGConversion(
		[]string{"session.sbg", "session.spsq"},
		&cli.CLIOptions{Play: true},
		&status,
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "output destination cannot be provided with -play") {
		t.Fatalf("expected playback output path error, got %v", err)
	}
	if status.Len() != 0 {
		t.Fatalf("unexpected status before validation error: %q", status.String())
	}
}

func TestRunSBGConversionRejectsIncompatiblePlaybackFlags(t *testing.T) {
	tests := []struct {
		name string
		opts *cli.CLIOptions
	}{
		{name: "dump", opts: &cli.CLIOptions{Play: true, Dump: true}},
		{name: "test", opts: &cli.CLIOptions{Play: true, Test: true}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := runSBGConversion([]string{"session.sbg"}, test.opts, nil, nil)
			if err == nil || !strings.Contains(err.Error(), "cannot be combined with -dump or -test") {
				t.Fatalf("expected incompatible playback flag error, got %v", err)
			}
		})
	}
}

func TestSBGConversionWarningRespectsColorSetting(t *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	cli.SetColorEnabled(true)
	var colored bytes.Buffer
	if err := writeSBGConversionWarning(&colored); err != nil {
		t.Fatalf("write colored warning: %v", err)
	}
	if !strings.Contains(colored.String(), "\x1b[") {
		t.Fatalf("expected ANSI warning, got %q", colored.String())
	}

	cli.SetColorEnabled(false)
	var plain bytes.Buffer
	if err := writeSBGConversionWarning(&plain); err != nil {
		t.Fatalf("write plain warning: %v", err)
	}
	if strings.Contains(plain.String(), "\x1b[") {
		t.Fatalf("unexpected ANSI warning: %q", plain.String())
	}
}

func TestSBGConversionSuccessRespectsColorSetting(t *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	cli.SetColorEnabled(true)
	var colored bytes.Buffer
	if err := writeSBGConversionSuccess(&colored, "session.spsq"); err != nil {
		t.Fatalf("write colored success: %v", err)
	}
	if !strings.Contains(colored.String(), "\x1b[") {
		t.Fatalf("expected ANSI success, got %q", colored.String())
	}

	cli.SetColorEnabled(false)
	var plain bytes.Buffer
	if err := writeSBGConversionSuccess(&plain, "session.spsq"); err != nil {
		t.Fatalf("write plain success: %v", err)
	}
	if strings.Contains(plain.String(), "\x1b[") {
		t.Fatalf("unexpected ANSI success: %q", plain.String())
	}
}
