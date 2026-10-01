// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"io"
	"os"
	"strings"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

// Converter converts SBaGen content through the provided SynapSeq application context.
type Converter struct {
	ctx *synapseq.AppContext
}

// New creates an SBaGen converter using ctx.
func New(ctx *synapseq.AppContext) (*Converter, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	return &Converter{ctx: ctx}, nil
}

// LoadFile converts the SBaGen sequence at path into a validated SynapSeq sequence.
func (c *Converter) LoadFile(path string) (*synapseq.LoadedContext, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open SBaGen file %q: %w", path, err)
	}
	defer file.Close()

	return c.load(path, file)
}

// LoadContent converts in-memory SBaGen content into a validated SynapSeq sequence.
func (c *Converter) LoadContent(content string) (*synapseq.LoadedContext, error) {
	return c.load("<content>", strings.NewReader(content))
}

func (c *Converter) load(source string, reader io.Reader) (*synapseq.LoadedContext, error) {
	parsedSequence, err := parse(source, reader)
	if err != nil {
		return nil, err
	}
	builder, err := build(c.ctx, parsedSequence)
	if err != nil {
		return nil, err
	}
	loaded, err := builder.Load()
	if err != nil {
		return nil, fmt.Errorf("validate converted sequence: %w", err)
	}
	return loaded, nil
}
