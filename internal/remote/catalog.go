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

package remote

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type indexCatalog struct {
	index *t.RemoteIndex
}

func loadIndexCatalog() (*indexCatalog, error) {
	index, err := GetIndex()
	if err != nil {
		return nil, err
	}

	return &indexCatalog{index: index}, nil
}

func (catalog *indexCatalog) findEntry(sequenceID string) *t.RemoteEntry {
	if catalog == nil {
		return nil
	}

	return findEntryByID(catalog.index, sequenceID)
}
