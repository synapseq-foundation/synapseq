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

/*
Package spsq provides a programmatic builder for SynapSeq .spsq sequence text.

# Overview

This package is a Go API for constructing SynapSeq sequences dynamically. It
uses a fluent Builder that records options and timeline entries, plus Preset
builders that record tracks and effects, then renders them as .spsq text.

Build validates the generated text through the core package and returns the
loaded sequence context. The spsq package does not preview, stream, or render
audio itself; those responsibilities remain in core and the internal sequence
and audio packages.

# Example Usage

	package main

	import (
	    "log"
	    "time"

	    "github.com/synapseq-foundation/synapseq/v4/spsq"
	)

	func main() {
	    builder := spsq.New().SampleRate(44100).Volume(100)
	    alpha := builder.NewPreset("alpha")
	    alpha.PinkNoise(0).Amplitude(30)
	    alpha.Tone(300).Binaural(10).Amplitude(15)

	    loaded, err := builder.
	        SilenceAt(0).
	        At(15*time.Second, alpha).
	        SilenceAt(time.Minute).
	        Build()
	    if err != nil {
	        log.Fatal(err)
	    }

	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }
	}

# Builder Flow

Typical builder usage follows the .spsq document shape:
  - add sequence options such as sample rate, volume, or ambiance;
  - create presets and add tracks with track modifiers;
  - add timeline entries that select presets or silence at specific times;
  - call Build to validate and load the generated .spsq content.

Builder methods return the same Builder so calls can be chained. Methods that
modify the last track or timeline entry are no-ops when there is no matching
target. Build returns a core LoadedContext or parser and validation errors
produced by core.

# More Information

For the loading, rendering, streaming, and preview API, see the core package:
https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/core

For complete project documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package spsq
