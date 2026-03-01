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

package core

import (
	"io"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// AppContext holds the configuration for the application.
// It provides a safe, immutable context for sequence processing.
// Methods that modify the context return a new instance.
type AppContext struct {
	inputFile    string
	outputFile   string
	statusOutput io.Writer
	sequence     *t.Sequence
}

// NewAppContext creates a new AppContext instance.
//
// Parameters:
//   - inputFile: path to the input sequence file (can be local path, stdin "-", or HTTP/HTTPS URL)
//   - outputFile: path to the output WAV file (local path only)
//
// Returns:
//   - *AppContext: a new AppContext instance with the provided configuration
//   - error: an error if the input parameters are invalid (e.g., unsupported output file path)
func NewAppContext(inputFile, outputFile string) (*AppContext, error) {
	return &AppContext{
		inputFile:    inputFile,
		outputFile:   outputFile,
		statusOutput: nil,
	}, nil
}

// InputFile returns the input file path.
func (ac *AppContext) InputFile() string {
	return ac.inputFile
}

// OutputFile returns the output file path.
func (ac *AppContext) OutputFile() string {
	return ac.outputFile
}

// Verbose returns whether verbose mode is enabled.
// When true, status output will be written to the configured writer.
func (ac *AppContext) Verbose() bool {
	return ac.statusOutput != nil
}

// WithVerbose returns a new AppContext with verbose mode enabled.
// Status output will be written to the provided writer (typically os.Stderr).
//
// Example:
//
//	ctx = ctx.WithVerbose(os.Stderr)
func (ac *AppContext) WithVerbose(data io.Writer) *AppContext {
	newCtx := *ac
	newCtx.statusOutput = data
	return &newCtx
}
