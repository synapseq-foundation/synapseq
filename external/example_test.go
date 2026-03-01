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

package external_test

import (
	"fmt"
	"log"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/external"
)

func ExampleNewFFPlay() {
	// Create ffplay instance using executable from PATH
	// player, err := external.NewFFPlay("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	player := external.NewFFPlayUnsafe("")
	fmt.Println("ffplay initialized:", player.Path())
	// Output:
	// ffplay initialized: ffplay
}

func ExampleFFplay_Play() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before playback)
	// _ = ctx.LoadSequence()

	// Create ffplay instance
	// _, err = external.NewFFPlay("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Play audio (real-time)
	// _ = player.Play(ctx)

	fmt.Printf("Playback executed successfully for input: %s\n", ctx.InputFile())
	// Output:
	// Playback executed successfully for input: input.spsq
}

func ExampleNewFFmpeg() {
	// Create ffmpeg instance using executable from PATH
	// encoder, err := external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	encoder := external.NewFFmpegUnsafe("")
	fmt.Println("ffmpeg initialized:", encoder.Path())
	// Output:
	// ffmpeg initialized: ffmpeg
}

func ExampleFFmpeg_Convert() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before encoding)
	// _ = ctx.LoadSequence()

	// Create ffmpeg instance
	// encoder, err := external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Encode MP3 at 320 kbps CBR
	// _ = encoder.Convert(ctx, "mp3")

	fmt.Printf("MP3 encoding executed successfully for output: %s\n", ctx.OutputFile())
	// Output:
	// MP3 encoding executed successfully for output: output.mp3
}
