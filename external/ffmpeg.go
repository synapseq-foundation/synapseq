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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

// Convert encodes streaming PCM into the output format inferred from the file extension.
func (fm *FFmpeg) Convert(loadedCtx *synapseq.LoadedContext, outputFile string) error {
	if loadedCtx == nil {
		return fmt.Errorf("loaded context cannot be nil")
	}

	if outputFile == "" {
		return fmt.Errorf("output file cannot be empty")
	}

	format, err := outputFormatFromFileName(outputFile)
	if err != nil {
		return err
	}

	// Remove existing output file if it exists
	if _, err := os.Stat(outputFile); err == nil {
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("failed to remove existing output file: %v", err)
		}
	}

	args := buildConvertArgs(loadedCtx.SampleRate(), outputFile, format)

	ffmpeg := fm.Command(args...)
	if err := startPipeCmd(ffmpeg, loadedCtx); err != nil {
		return err
	}

	return nil
}

func outputFormatFromFileName(outputFile string) (string, error) {
	extension := strings.ToLower(filepath.Ext(outputFile))
	if extension == "" {
		return "", fmt.Errorf("output file must include an extension to determine format")
	}

	switch extension {
	case ".mp3":
		return "mp3", nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", extension)
	}
}

func buildConvertArgs(sampleRate int, outputFile, format string) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(sampleRate),
		"-i", "pipe:0",
	}

	switch format {
	case "mp3":
		args = append(args, []string{
			"-c:a", "libmp3lame",
			"-b:a", "320k",
			"-f", "mp3",
		}...)
	}

	args = append(args, []string{
		"-vn",
		outputFile,
	}...)

	return args
}
