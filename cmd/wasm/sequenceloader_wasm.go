//go:build wasm

// Deprecated: the SynapSeq WebAssembly runtime is kept only for historical
// reference and is no longer recommended for new integrations.
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
	sequencepkg "github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type sequenceLoader struct{}

func (sequenceLoader) Load(rawContent []byte) (*t.Sequence, error) {
	return sequencepkg.LoadTextSequence(rawContent, "", "")
}
