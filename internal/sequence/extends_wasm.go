//go:build wasm

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
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// extends loads preset and option definitions from a remote .spsc file.
func extends(fileName string) (*t.Extends, error) {
	rawContent, err := s.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, err
	}

	return parseExtendsContent(rawContent, "", "")
}
