// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
