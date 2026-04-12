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

package sequence

import (
	"fmt"
	"path/filepath"

	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// LoadTextSequence loads a sequence from a text file
func LoadTextSequence(fileName string) (*t.Sequence, error) {
	rawContent, err := r.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	absInputFile, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	baseDir := filepath.Dir(absInputFile)

	return parseSequenceContent(rawContent, absInputFile, baseDir)
}
