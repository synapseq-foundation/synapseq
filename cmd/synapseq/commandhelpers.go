//go:build !js && !wasm

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

package main

import (
	"path/filepath"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func defaultOutputExtension(opts *cli.CLIOptions) string {
	if opts.Preview {
		return ".html"
	}
	if opts.Mp3 {
		return ".mp3"
	}
	return ".wav"
}

func defaultOutputFileName(baseName string, opts *cli.CLIOptions) string {
	return baseName + defaultOutputExtension(opts)
}

func resolveOutputTarget(defaultBaseName, requestedOutput string, opts *cli.CLIOptions) (string, string) {
	if requestedOutput == "" {
		outputFile := defaultOutputFileName(defaultBaseName, opts)
		return outputFile, defaultOutputExtension(opts)
	}

	return requestedOutput, strings.ToLower(filepath.Ext(requestedOutput))
}

func buildOutputOptions(outputFile, outputFormat string, opts *cli.CLIOptions) *outputOptions {
	return &outputOptions{
		OutputFile: outputFile,
		Quiet:      opts.Quiet,
		Preview:    opts.Preview,
		Play:       opts.Play,
		Mp3:        outputFormat == ".mp3" || opts.Mp3,
		FFplayPath: opts.FFplayPath,
		FFmpegPath: opts.FFmpegPath,
	}
}
