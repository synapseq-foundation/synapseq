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

package hub

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type manifestCatalog struct {
	manifest *t.HubManifest
}

func loadManifestCatalog() (*manifestCatalog, error) {
	manifest, err := GetManifest()
	if err != nil {
		return nil, err
	}

	return &manifestCatalog{manifest: manifest}, nil
}

func (catalog *manifestCatalog) findEntry(sequenceID string) *t.HubEntry {
	if catalog == nil {
		return nil
	}

	return findEntryByID(catalog.manifest, sequenceID)
}