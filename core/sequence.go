//go:build !wasm

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

package core

import (
	"fmt"
	"maps"
	"path/filepath"

	preview "github.com/synapseq-foundation/synapseq/v4/internal/preview"
	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	seq "github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// LoadFile loads the sequence from the input file.
func (ac *AppContext) LoadFile(path string) (*LoadedContext, error) {
	rawContent, err := r.GetFile(path, t.FormatText)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	absInputFile, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	sequence, err := seq.LoadTextSequence(rawContent, absInputFile, filepath.Dir(absInputFile))
	if err != nil {
		return nil, err
	}

	return &LoadedContext{
		appCtx:   ac,
		sequence: sequence,
	}, nil
}

// LoadContent loads the sequence from raw text content.
func (ac *AppContext) LoadContent(content string) (*LoadedContext, error) {
	sequence, err := seq.LoadTextSequence([]byte(content), "", "")
	if err != nil {
		return nil, err
	}

	return &LoadedContext{
		appCtx:   ac,
		sequence: sequence,
	}, nil
}

// Comments returns a defensive copy of sequence comments.
func (lc *LoadedContext) Comments() []string {
	if lc.sequence == nil || len(lc.sequence.Comments) == 0 {
		return []string{}
	}

	comments := make([]string, len(lc.sequence.Comments))
	copy(comments, lc.sequence.Comments)

	return comments
}

// SampleRate returns the sample rate from the loaded sequence options.
func (lc *LoadedContext) SampleRate() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.SampleRate
}

// Extends returns a defensive copy of extends list.
func (lc *LoadedContext) Extends() []string {
	if lc.sequence == nil || lc.sequence.Options == nil || len(lc.sequence.Options.Extends) == 0 {
		return []string{}
	}

	extends := make([]string, len(lc.sequence.Options.Extends))
	copy(extends, lc.sequence.Options.Extends)

	return extends
}

// Volume returns the volume from the loaded sequence options.
func (lc *LoadedContext) Volume() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.Volume
}

// Ambiance returns a defensive copy of ambiance map.
func (lc *LoadedContext) Ambiance() map[string]string {
	if lc.sequence == nil || lc.sequence.Options == nil || len(lc.sequence.Options.Ambiance) == 0 {
		return map[string]string{}
	}

	ambiance := make(map[string]string, len(lc.sequence.Options.Ambiance))
	maps.Copy(ambiance, lc.sequence.Options.Ambiance)

	return ambiance
}

// RawContent returns a defensive copy of raw content.
func (lc *LoadedContext) RawContent() []byte {
	if lc.sequence == nil || len(lc.sequence.RawContent) == 0 {
		return []byte{}
	}

	raw := make([]byte, len(lc.sequence.RawContent))
	copy(raw, lc.sequence.RawContent)

	return raw
}

// Preview renders the loaded sequence as an HTML preview.
func (lc *LoadedContext) Preview() ([]byte, error) {
	if lc.sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}

	return preview.GetPreviewContent(lc.sequence.Periods)
}
