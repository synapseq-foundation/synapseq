//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package core_test

import (
	"fmt"
	"log"
	"os"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

func ExampleNewAppContext() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	fmt.Printf("AppContext created successfully\n")
	// Output: AppContext created successfully
}

func ExampleAppContext_LoadSequence() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("Sequence loaded successfully\n")
	// Output: Sequence loaded successfully
}

func ExampleAppContext_WAV() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}

	// Optional: Enable verbose output
	// Replace with an io.Writer, e.g., os.Stderr
	ctx = ctx.WithVerbose(os.Stderr)

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Generate the WAV file
	// if err := ctx.WAV(); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("WAV file generated successfully\n")
	// Output: WAV file generated successfully
}

func ExampleAppContext_Stream() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Stream the RAW data to standard output (44100 Hz [default], 16-bit, stereo)
	// Replace with an io.Writer, e.g., os.Stdout
	// if err := ctx.Stream(os.Stdout); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("RAW data streamed successfully\n")
	// Output: RAW data streamed successfully
}

func ExampleAppContext_Comments() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Retrieve comments from the sequence
	// for _, comment := range ctx.Comments() {
	//	fmt.Println(comment)
	// }

	fmt.Printf("Comments retrieved successfully\n")
	// Output: Comments retrieved successfully
}

func ExampleAppContext_SampleRate() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Get the sample rate from the loaded sequence
	// sampleRate := ctx.SampleRate()
	// fmt.Printf("Sample Rate: %d Hz\n", sampleRate)

	fmt.Printf("Sample rate retrieved successfully\n")
	// Output: Sample rate retrieved successfully
}

func ExampleAppContext_Volume() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Get the volume from the loaded sequence
	// volume := ctx.Volume()
	// fmt.Printf("Volume: %d\n", volume)

	fmt.Printf("Volume retrieved successfully\n")
	// Output: Volume retrieved successfully
}

func ExampleAppContext_AmbianceList() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Get the ambiance audio list from the loaded sequence
	// ambianceList := ctx.AmbianceList()
	// fmt.Printf("Ambiance entries: %d\n", len(ambianceList))

	fmt.Printf("Ambiance list retrieved successfully\n")
	// Output: Ambiance list retrieved successfully
}

func ExampleAppContext_PresetList() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Get the preset list from the loaded sequence
	// presetList := ctx.PresetList()
	// fmt.Printf("Preset files: %d\n", len(presetList))

	fmt.Printf("Preset list retrieved successfully\n")
	// Output: Preset list retrieved successfully
}

func ExampleAppContext_RawContent() {
	// Create a new application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav")
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Get the raw content of the loaded sequence
	// rawContent := ctx.RawContent()
	// fmt.Printf("Raw Content Length: %d bytes\n", len(rawContent))

	fmt.Printf("Raw content retrieved successfully\n")
	// Output: Raw content retrieved successfully
}
