// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// LoadTextSequence loads a sequence from text content.
func LoadTextSequence(rawContent []byte, sourceFile string, baseRef string) (*t.Sequence, error) {
	return parseSequenceContent(rawContent, sourceFile, baseRef)
}
