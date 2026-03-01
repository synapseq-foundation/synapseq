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

package external

import (
	"fmt"
	"os"
	"strconv"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

// FFmpeg represents the ffmpeg external tool
type FFmpeg struct{ baseUtility }

// NewFFmpeg creates a new FFmpeg instance with given ffmpeg path
func NewFFmpeg(ffmpegPath string) (*FFmpeg, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	util, err := newUtility(ffmpegPath)
	if err != nil {
		return nil, err
	}

	return &FFmpeg{baseUtility: *util}, nil
}

// NewFFmpegUnsafe creates an FFmpeg instance without validating the path.
// Useful for documentation examples and testing environments.
func NewFFmpegUnsafe(path string) *FFmpeg {
	if path == "" {
		path = "ffmpeg"
	}
	return &FFmpeg{baseUtility: baseUtility{path: path}}
}

// Convert encodes streaming PCM into the specified format using ffmpeg.
func (fm *FFmpeg) Convert(appCtx *synapseq.AppContext, format string) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	// Remove existing output file if it exists
	outputFile := appCtx.OutputFile()
	if _, err := os.Stat(outputFile); err == nil {
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("failed to remove existing output file: %v", err)
		}
	}

	// Base ffmpeg arguments
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
	}

	// Determine format and corresponding options
	switch format {
	case "mp3":
		args = append(args, []string{
			"-c:a", "libmp3lame",
			"-b:a", "320k",
			"-f", "mp3",
		}...)
	// TODO: more formats can be added here
	// BUT for brainwave entrainment, only MP3 is currently relevant
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	args = append(args, []string{
		"-vn",
		outputFile,
	}...)

	ffmpeg := fm.Command(args...)
	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}
