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

package spsq_test

import (
	"fmt"
	"io"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func ExampleNew() {
	builder := spsq.New().SampleRate(44100).Volume(100)
	alpha := builder.NewPreset("alpha")
	alpha.PinkNoise(0).Amplitude(30)
	alpha.Tone(300).Binaural(10).Amplitude(15)

	ctx := synapseq.NewAppContext()
	loaded, err := builder.
		SilenceAt(0).
		PresetAt(15*time.Second, alpha).
		SilenceAt(time.Minute).
		Load(ctx)
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Println("SPSQ content built and loaded successfully")
	// Output:
	// SPSQ content built and loaded successfully
}

func ExampleBuilder_Load_verbose() {
	builder := spsq.New()
	alpha := builder.NewPreset("alpha")
	alpha.PinkNoise(0).Amplitude(30)

	ctx := synapseq.NewAppContext().WithVerbose(io.Discard, false)
	loaded, err := builder.
		SilenceAt(0).
		PresetAt(15*time.Second, alpha).
		SilenceAt(time.Minute).
		Load(ctx)
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Println("SPSQ content built with verbose context")
	// Output:
	// SPSQ content built with verbose context
}
