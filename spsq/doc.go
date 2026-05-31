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

This package is the public Go API for constructing SynapSeq sequences
dynamically. It uses a fluent Builder that records options, presets, tracks,
effects, and timeline entries, then renders them as .spsq text.

The generated text is intended to be consumed by the core package, usually
through AppContext.LoadContent. The spsq package does not parse, validate,
preview, stream, or render audio itself; those responsibilities remain in
core and the internal sequence and audio packages.

# Example Usage

	package main

	import (
	    "log"

	    synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	    "github.com/synapseq-foundation/synapseq/v4/spsq"
	)

	func main() {
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
	        log.Fatal(err)
	    }

	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }
	}

# Builder Flow

Typical builder usage follows the .spsq document shape:
  - add sequence options such as sample rate, volume, ambiance, or extends;
  - add presets with tracks and track modifiers;
  - add timeline entries that select presets or silence at specific times;
  - call String to produce the final .spsq content.

Builder methods return the same Builder so calls can be chained. Methods that
modify the "last" preset, track, or timeline entry are no-ops when there is no
matching target.

# More Information

For the loading, rendering, streaming, and preview API, see the core package:
https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/core

For complete project documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package spsq
