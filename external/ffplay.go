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
	"strconv"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

// FFplay represents the ffplay external tool
type FFplay struct{ baseUtility }

// NewFFPlay creates a new FFplay instance with given ffplay path
func NewFFPlay(ffplayPath string) (*FFplay, error) {
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	util, err := newUtility(ffplayPath)
	if err != nil {
		return nil, err
	}

	return &FFplay{baseUtility: *util}, nil
}

// NewFFPlayUnsafe creates an FFplay instance without validating the path.
// Useful for documentation examples and testing environments.
func NewFFPlayUnsafe(path string) *FFplay {
	if path == "" {
		path = "ffplay"
	}
	return &FFplay{baseUtility: baseUtility{path: path}}
}

// Play invokes ffplay to play from streaming audio input.
func (fp *FFplay) Play(loadedCtx *synapseq.LoadedContext) error {
	if loadedCtx == nil {
		return fmt.Errorf("loaded context cannot be nil")
	}

	ffplay := fp.Command(
		"-nodisp",
		"-hide_banner",
		"-loglevel", "error",
		"-autoexit",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(loadedCtx.SampleRate()),
		"-i", "pipe:0",
	)

	if err := startPipeCmd(ffplay, loadedCtx); err != nil {
		return err
	}

	return nil
}
