// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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