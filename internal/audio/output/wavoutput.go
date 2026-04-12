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

package output

import (
	"io"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
)

type WAVOutput struct {
	sampleRate int
	channels   int
	precision  int
	render     RenderFunc
}

func NewWAVOutput(sampleRate, channels, precision int, render RenderFunc) *WAVOutput {
	return &WAVOutput{
		sampleRate: sampleRate,
		channels:   channels,
		precision:  precision,
		render:     render,
	}
}

func (wo *WAVOutput) Write(out io.WriteSeeker) error {
	streamer := NewRendererStreamer(wo.render)
	if err := bwav.Encode(out, streamer, wo.format()); err != nil {
		return err
	}

	return streamer.Err()
}

func (wo *WAVOutput) format() beep.Format {
	return beep.Format{
		SampleRate:  beep.SampleRate(wo.sampleRate),
		NumChannels: wo.channels,
		Precision:   wo.precision,
	}
}