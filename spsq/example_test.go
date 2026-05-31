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

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func ExampleNew() {
	sequence := spsq.New().
		AddPreset("alpha").
		AddNoiseTrack().WithPinkNoise(0).WithAmplitude(30).
		AddToneTrack(300).WithBinauralTone(10).WithAmplitude(15).
		SilenceAt(0, 0, 0).
		PresetAt(0, 0, 15).
		SilenceAt(0, 1, 0)

	ctx := synapseq.NewAppContext()
	loaded, err := ctx.LoadContent(sequence.String())
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Println("SPSQ content built and loaded successfully")
	// Output:
	// SPSQ content built and loaded successfully
}
