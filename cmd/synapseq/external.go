// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/external"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

// externalPlay invokes utility tool to play from streaming audio input
func externalPlay(ffplayPath string, loadedCtx *synapseq.LoadedContext) error {
	ffplay, err := external.NewFFPlay(ffplayPath)
	if err != nil {
		return diag.Validation("something went wrong while calling ffplay").WithHint("run `synapseq -doctor` to diagnose issues")
	}

	if err := ffplay.Play(loadedCtx); err != nil {
		return err
	}

	return nil
}

// externalMp3 encodes streaming PCM into an MP3 file using external utility
func externalMp3(ffmpegPath string, loadedCtx *synapseq.LoadedContext, outputFile string) error {
	ffmpeg, err := external.NewFFmpeg(ffmpegPath)
	if err != nil {
		return diag.Validation("something went wrong while calling ffmpeg").WithHint("run `synapseq -doctor` to diagnose issues")
	}

	if err := ffmpeg.Convert(loadedCtx, outputFile); err != nil {
		return err
	}

	return nil
}
