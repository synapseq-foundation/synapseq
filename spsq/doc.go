// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package spsq provides a programmatic builder for SynapSeq .spsq sequence text.

# Overview

This package is a Go API for constructing SynapSeq sequences dynamically. It
uses a fluent Builder that records options and timeline entries, plus Preset
builders that record tracks and effects, then renders them as .spsq text.

Load validates the generated text through the core package and returns the
loaded sequence context. The spsq package does not dump, stream, or render audio
itself; those responsibilities remain in core and the internal sequence and
audio packages.

# Example Usage

	package main

	import (
	    "log"
	    "time"

	    synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	    "github.com/synapseq-foundation/synapseq/v4/spsq"
	)

	func main() {
	    builder := spsq.New().SampleRate(44100).Volume(100)
	    alpha := builder.NewPreset("alpha")
	    alpha.Pink(0).Amplitude(30)
	    alpha.Tone(300).Binaural(10).Amplitude(15)

	    ctx := synapseq.NewAppContext()
	    loaded, err := builder.
	        SilenceAt(0).
	        PresetAt(15*time.Second, alpha).
	        SilenceAt(time.Minute).
	        Load(ctx)
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }
	}

# Verbose Output

Load receives a core AppContext. Configure that context with WithVerbose when
you want progress output from later operations such as WAV, MP3, or Stream:

	ctx := synapseq.NewAppContext().WithVerbose(os.Stderr, true)
	loaded, err := builder.Load(ctx)
	if err != nil {
	    log.Fatal(err)
	}
	if err := loaded.WAV("output.wav"); err != nil {
	    log.Fatal(err)
	}

# Builder Flow

Typical builder usage follows the .spsq document shape:
  - add sequence options such as sample rate, volume, ambiance, or music;
  - create presets and add tracks with track modifiers;
  - add timeline entries that select presets or silence at specific times;
  - call Load with a core AppContext to validate and load the generated .spsq content.

Noise tracks are added with White, Pink, or Brown. Each method receives the
noise smoothness percentage and returns the preset so modifiers such as
Amplitude can be chained:

	preset.White(0).Amplitude(20)
	preset.Pink(10).Amplitude(20)
	preset.Brown(15).Amplitude(20)

Builder methods return the same Builder so calls can be chained. Methods that
modify the last track or timeline entry are no-ops when there is no matching
target. Load requires a non-nil core AppContext and returns a core
LoadedContext or parser and validation errors produced by core.

# More Information

For the loading, rendering, streaming, and dump API, see the core package:
https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/core

For complete project documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package spsq
