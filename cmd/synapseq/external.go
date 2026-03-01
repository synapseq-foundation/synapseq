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

package main

import (
	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/external"
)

// externalPlay invokes utility tool to play from streaming audio input
func externalPlay(ffplayPath string, appCtx *synapseq.AppContext) error {
	ffplay, err := external.NewFFPlay(ffplayPath)
	if err != nil {
		return err
	}

	if err := ffplay.Play(appCtx); err != nil {
		return err
	}

	return nil
}

// externalMp3 encodes streaming PCM into an MP3 file using external utility
func externalMp3(ffmpegPath string, appCtx *synapseq.AppContext) error {
	ffmpeg, err := external.NewFFmpeg(ffmpegPath)
	if err != nil {
		return err
	}

	if err := ffmpeg.Convert(appCtx, "mp3"); err != nil {
		return err
	}

	return nil
}
