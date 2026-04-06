/*
Package core provides the application context and core functionality
for the SynapSeq text-driven audio sequencer for brainwave entrainment.

# Overview

This package is the public Go API for loading .spsq sequences,
inspecting their metadata, rendering WAV output, streaming raw PCM,
and generating HTML previews.

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
		loaded, err := ctx.Load("input.spsq")
		if err != nil {
			log.Fatal(err)
		}

	    // Generate WAV file
	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }

	    // Or render an HTML preview
	    previewHTML, err := loaded.Preview()
	    if err != nil {
	        log.Fatal(err)
	    }
	    _ = previewHTML
	}

# File Paths

Input paths support:
  - Local file paths: "path/to/file.spsq"
  - Standard input: "-" (input only)
  - HTTP/HTTPS URLs: "https://example.com/sequence.spsq"

Output methods:
  - loaded.WAV("output.wav") writes a WAV file
  - loaded.Stream(writer) writes raw PCM to an io.Writer
  - loaded.Preview() returns HTML preview bytes

# Thread Safety

AppContext methods are safe for concurrent use because configuration methods
return new instances rather than mutating the original context.

# More Information

For complete documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package core
