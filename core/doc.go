// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package core provides the application context and core functionality
for the SynapSeq text-driven audio sequencer for brainwave entrainment.

# Overview

This package is the public Go API for loading .spsq sequences,
inspecting their metadata, rendering WAV output, streaming raw PCM,
and producing JSON dumps.

# Supported Format

SynapSeq currently supports text input in .spsq format.

# Example Usage

	package main

	import (
	    "log"
	    "os"

	    synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	)

	func main() {
	    // Create application context
	    ctx := synapseq.NewAppContext()

	    // Enable verbose output (optional)
	    ctx = ctx.WithVerbose(os.Stderr, true)

		// Load sequence (required before generating WAV or streaming)
		loaded, err := ctx.LoadFile("input.spsq")
		if err != nil {
			log.Fatal(err)
		}

	    // Generate WAV file
	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }

	    // Or generate a JSON dump
	    jsonDump, err := loaded.JSON()
	    if err != nil {
	        log.Fatal(err)
	    }
	    _ = jsonDump
	}

# File Paths

Input paths support:
  - Local file paths: "path/to/file.spsq"
  - Standard input: "-" (input only)
  - HTTP/HTTPS URLs: "https://example.com/sequence.spsq"

# Raw Content

Use LoadContent when the .spsq source is already available as a string,
such as content received from a web form, embedded resource, or another
in-memory source:

	content := `
	alpha
	  tone 100 binaural 1 amplitude 1
	00:00:00 alpha
	00:01:00 alpha
	`

	loaded, err := ctx.LoadContent(content)
	if err != nil {
		log.Fatal(err)
	}
	_ = loaded

Output methods:
  - loaded.WAV("output.wav") writes a WAV file
  - loaded.Stream(writer) writes raw PCM to an io.Writer
  - loaded.JSON() returns JSON dump bytes

# Sequence Inspection

Use Duration to retrieve the total duration of the loaded sequence as a
time.Duration:

	duration := loaded.Duration()

# Thread Safety

AppContext methods are safe for concurrent use because configuration methods
return new instances rather than mutating the original context.

# More Information

For complete documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package core
