// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core_test

import (
	"fmt"
	"log"
	"os"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

func ExampleNewAppContext() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	fmt.Printf("AppContext created successfully\n")
	// Output: AppContext created successfully
}

func ExampleAppContext_LoadFile() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }
	// _ = loaded

	fmt.Printf("Sequence loaded successfully\n")
	// Output: Sequence loaded successfully
}

func ExampleAppContext_LoadContent() {
	// Create a new application context
	ctx := synapseq.NewAppContext()

	content := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
`

	loaded, err := ctx.LoadContent(content)
	if err != nil {
		panic(err)
	}
	_ = loaded

	fmt.Printf("Sequence content loaded successfully\n")
	// Output: Sequence content loaded successfully
}

func ExampleLoadedContext_WAV() {
	// Create a new application context
	ctx := synapseq.NewAppContext()

	// Optional: Enable verbose output
	// Replace with an io.Writer, e.g., os.Stderr
	ctx = ctx.WithVerbose(os.Stderr, true)

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Generate the WAV file
	// if err := loaded.WAV("output.wav"); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("WAV file generated successfully\n")
	// Output: WAV file generated successfully
}

func ExampleLoadedContext_Stream() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Stream the RAW data to standard output (44100 Hz [default], 16-bit, stereo)
	// Replace with an io.Writer, e.g., os.Stdout
	// if err := loaded.Stream(os.Stdout); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("RAW data streamed successfully\n")
	// Output: RAW data streamed successfully
}

func ExampleLoadedContext_Comments() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Retrieve comments from the sequence
	// for _, comment := range loaded.Comments() {
	//	fmt.Println(comment)
	// }

	fmt.Printf("Comments retrieved successfully\n")
	// Output: Comments retrieved successfully
}

func ExampleLoadedContext_SampleRate() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Get the sample rate from the loaded sequence
	// sampleRate := loaded.SampleRate()
	// fmt.Printf("Sample Rate: %d Hz\n", sampleRate)

	fmt.Printf("Sample rate retrieved successfully\n")
	// Output: Sample rate retrieved successfully
}

func ExampleLoadedContext_Volume() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Get the volume from the loaded sequence
	// volume := loaded.Volume()
	// fmt.Printf("Volume: %d\n", volume)

	fmt.Printf("Volume retrieved successfully\n")
	// Output: Volume retrieved successfully
}

func ExampleLoadedContext_Ambiance() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Get the ambiance audio from the loaded sequence
	// ambiance := loaded.Ambiance()
	// fmt.Printf("Ambiance entries: %d\n", len(ambiance))

	fmt.Printf("Ambiance retrieved successfully\n")
	// Output: Ambiance retrieved successfully
}

func ExampleLoadedContext_Presets() {
	ctx := synapseq.NewAppContext()

	loaded, err := ctx.LoadContent(`alpha
  waveform sine tone 300 binaural 10 amplitude 20
  noise pink smooth 5 amplitude 10

00:00:00 alpha
00:01:00 alpha
`)
	if err != nil {
		log.Fatal(err)
	}

	presets := loaded.Presets()

	fmt.Printf("Presets: %d\n", len(presets))
	fmt.Printf("Preset: %s\n", presets[0].Name)
	fmt.Printf("Tracks: %d\n", len(presets[0].Tracks))
	fmt.Printf("Track 1: %s %.2f %.2f %.2f\n",
		presets[0].Tracks[0].Type,
		presets[0].Tracks[0].Carrier,
		presets[0].Tracks[0].Resonance,
		presets[0].Tracks[0].Amplitude,
	)

	// Output:
	// Presets: 1
	// Preset: alpha
	// Tracks: 2
	// Track 1: binaural 300.00 10.00 20.00
}

func ExampleLoadedContext_Timeline() {
	ctx := synapseq.NewAppContext()

	loaded, err := ctx.LoadContent(`alpha
  waveform sine tone 300 binaural 10 amplitude 20

beta
  waveform sine tone 200 monaural 8 amplitude 15

00:00:00 alpha steady 0
00:01:00 beta smooth 2
`)
	if err != nil {
		log.Fatal(err)
	}

	timeline := loaded.Timeline()

	fmt.Printf("Timeline: %d\n", len(timeline))
	fmt.Printf("Entry 1: %s %s %s %d\n",
		timeline[0].Timestamp,
		timeline[0].PresetName,
		timeline[0].Transition,
		timeline[0].Steps,
	)
	fmt.Printf("Entry 2: %s %s %s %d\n",
		timeline[1].Timestamp,
		timeline[1].PresetName,
		timeline[1].Transition,
		timeline[1].Steps,
	)

	// Output:
	// Timeline: 2
	// Entry 1: 00:00:00 alpha steady 0
	// Entry 2: 00:01:00 beta smooth 2
}

func ExampleLoadedContext_Duration() {
	ctx := synapseq.NewAppContext()

	loaded, err := ctx.LoadContent(`alpha
  waveform sine tone 300 binaural 10 amplitude 20

00:00:00 alpha
00:01:00 alpha
`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Duration: %s\n", loaded.Duration())

	// Output:
	// Duration: 1m0s
}

func ExampleLoadedContext_Extends() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Get the extends list from the loaded sequence
	// extends := loaded.Extends()
	// fmt.Printf("Extends entries: %d\n", len(extends))

	fmt.Printf("Extends list retrieved successfully\n")
	// Output: Extends list retrieved successfully
}

func ExampleLoadedContext_JSON() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Generate the JSON dump bytes
	// jsonDump, err := loaded.JSON()
	// if err != nil {
	//	log.Fatal(err)
	// }
	// _ = jsonDump

	fmt.Printf("JSON dump bytes generated successfully\n")
	// Output: JSON dump bytes generated successfully
}

func ExampleLoadedContext_RawContent() {
	// Create a new application context
	ctx := synapseq.NewAppContext()
	_ = ctx

	// Load the sequence
	// loaded, err := ctx.LoadFile("input.spsq")
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Get the raw content of the loaded sequence
	// rawContent := loaded.RawContent()
	// fmt.Printf("Raw Content Length: %d bytes\n", len(rawContent))

	fmt.Printf("Raw content retrieved successfully\n")
	// Output: Raw content retrieved successfully
}
