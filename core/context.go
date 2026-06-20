// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"io"
)

// NewAppContext creates a new AppContext instance.
func NewAppContext() *AppContext {
	return &AppContext{
		statusOutput: nil,
		statusColors: false,
	}
}

// Verbose returns whether verbose mode is enabled.
// When true, status output will be written to the configured writer.
func (ac *AppContext) Verbose() bool {
	return ac.statusOutput != nil
}

// WithVerbose returns a new AppContext with verbose mode enabled.
// Status output will be written to the provided writer (typically os.Stderr),
// and colors controls whether ANSI color sequences are emitted.
//
// Example:
//
//	ctx = ctx.WithVerbose(os.Stderr, true)
func (ac *AppContext) WithVerbose(data io.Writer, colors bool) *AppContext {
	newCtx := *ac
	newCtx.statusOutput = data
	newCtx.statusColors = colors
	return &newCtx
}
