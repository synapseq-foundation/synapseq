/*
Package core provides the application context and core functionality
for the SynapSeq brainwave generator.

# Overview

This package is designed to be used as a library by other Go projects
that want to integrate SynapSeq audio generation capabilities.

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
	    ctx = ctx.WithVerbose(os.Stderr)

		// Load sequence (required before generating WAV or streaming)
		loaded, err := ctx.Load("input.spsq")
		if err != nil {
			log.Fatal(err)
		}

	    // Generate WAV file
	    if err := loaded.WAV("output.wav"); err != nil {
	        log.Fatal(err)
	    }
	}

# File Paths

Input and output files support:
  - Local file paths: "path/to/file.spsq"
  - Standard input: "-" (only for input files)
  - HTTP/HTTPS URLs: "https://example.com/sequence.spsq"

# Thread Safety

AppContext methods are safe for concurrent use as they return new instances
rather than modifying the original context.

# More Information

For complete documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package core
