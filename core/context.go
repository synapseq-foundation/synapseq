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
	statusOutput io.Writer
}

// LoadedContext holds a loaded sequence and execution settings.
type LoadedContext struct {
	appCtx    *AppContext
	sequence  *t.Sequence
}

// NewAppContext creates a new AppContext instance.
func NewAppContext() *AppContext {
	return &AppContext{
		statusOutput: nil,
	}
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
