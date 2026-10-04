// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package spsq_test

import (
	"fmt"
	"io"
	"strings"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func ExampleNew() {
	ctx := synapseq.NewAppContext()
	builder, err := spsq.New(ctx)
	if err != nil {
		panic(err)
	}
	builder.SampleRate(44100).Volume(100)
	alpha := builder.NewPreset("alpha")
	alpha.Pink(0).Amplitude(30)
	alpha.Tone(300).Binaural(10).Amplitude(15)

	loaded, err := builder.
		SilenceAt(0).
		PresetAt(15*time.Second, alpha).
		SilenceAt(time.Minute).
		Load()
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Println("SPSQ content built and loaded successfully")
	// Output:
	// SPSQ content built and loaded successfully
}

func ExampleBuilder_Load_verbose() {
	ctx := synapseq.NewAppContext().WithVerbose(io.Discard, false)
	builder, err := spsq.New(ctx)
	if err != nil {
		panic(err)
	}
	alpha := builder.NewPreset("alpha")
	alpha.Pink(0).Amplitude(30)

	loaded, err := builder.
		SilenceAt(0).
		PresetAt(15*time.Second, alpha).
		SilenceAt(time.Minute).
		Load()
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Println("SPSQ content built with verbose context")
	// Output:
	// SPSQ content built with verbose context
}

func ExamplePreset_RightLeft() {
	builder, err := spsq.New(synapseq.NewAppContext())
	if err != nil {
		panic(err)
	}
	focus := builder.NewPreset("focus")
	focus.Tone(220).Binaural(10).RightLeft().Amplitude(10)
	loaded, err := builder.PresetAt(0, focus).PresetAt(time.Minute, focus).Load()
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(loaded.RawContent()), "\n") {
		if strings.Contains(line, "binaural") {
			fmt.Println(strings.TrimSpace(line))
		}
	}
	// Output:
	// waveform sine tone 220.00 binaural right-left 10.00 amplitude left 10.00 right 10.00
}

func ExamplePreset_Short() {
	builder, err := spsq.New(synapseq.NewAppContext())
	if err != nil {
		panic(err)
	}
	focus := builder.NewPreset("focus")
	focus.Tone(220).Isochronic(10).Short().Amplitude(15)
	loaded, err := builder.PresetAt(0, focus).PresetAt(time.Minute, focus).Load()
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(loaded.RawContent()), "\n") {
		if strings.Contains(line, "isochronic") {
			fmt.Println(strings.TrimSpace(line))
		}
	}
	// Output:
	// waveform sine tone 220.00 isochronic short 10.00 amplitude left 15.00 right 15.00
}

func ExamplePreset_Wide() {
	builder, err := spsq.New(synapseq.NewAppContext())
	if err != nil {
		panic(err)
	}
	focus := builder.NewPreset("focus")
	focus.Tone(220).Isochronic(10).Wide().Amplitude(15)
	loaded, err := builder.PresetAt(0, focus).PresetAt(time.Minute, focus).Load()
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(loaded.RawContent()), "\n") {
		if strings.Contains(line, "isochronic") {
			fmt.Println(strings.TrimSpace(line))
		}
	}
	// Output:
	// waveform sine tone 220.00 isochronic wide 10.00 amplitude left 15.00 right 15.00
}
