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

package main

import (
	"fmt"
	"os"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

// OutputOptions defines options for processing sequence output
type outputOptions struct {
	OutputFile       string
	Quiet            bool
	Preview          bool
	Play             bool
	Mp3              bool
	UnsafeNoMetadata bool
	FFplayPath       string
	FFmpegPath       string
}

// processSequenceOutput processes the output of a loaded sequence
func processSequenceOutput(loadedCtx *synapseq.LoadedContext, opts *outputOptions) error {
	if opts.Preview {
		content, err := loadedCtx.Preview()
		if err != nil {
			return err
		}

		if opts.OutputFile == "-" {
			fmt.Println(string(content))
			return nil
		}

		if err := os.WriteFile(opts.OutputFile, content, 0644); err != nil {
			return fmt.Errorf("failed to write preview HTML: %v", err)
		}

		if !opts.Quiet {
			fmt.Printf("Preview has been generated as %q\n", opts.OutputFile)
			fmt.Printf("You can open the file in a web browser to view the sequence preview.\n")
		}

		return nil
	}

	// --- Handle Stream mode (output = "-")
	if opts.OutputFile == "-" {
		return loadedCtx.Stream(os.Stdout)
	}

	// --- Print comments
	if !opts.Quiet {
		for _, c := range loadedCtx.Comments() {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	// --- Handle Play using external ffplay
	if opts.Play {
		return externalPlay(opts.FFplayPath, loadedCtx)
	}

	// --- Handle MP3 output using external ffmpeg
	if opts.Mp3 {
		return externalMp3(opts.FFmpegPath, loadedCtx, opts.OutputFile)
	}

	// Default: Render to WAV
	return loadedCtx.WAV(opts.OutputFile)
}
