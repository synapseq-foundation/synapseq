//go:build wasm

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

package sequence

import (
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// LoadTextSequence loads a sequence from a file content
func LoadTextSequence(rawContent []byte) (*t.Sequence, error) {
	return parseSequenceContent(rawContent, "")
}
