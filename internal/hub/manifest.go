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

// GetManifest retrieves and parses the Hub manifest file from the cache
func GetManifest() (*t.HubManifest, error) {
	cache, err := openHubCache()
	if err != nil {
		return nil, err
	}

	return cache.manifest().read()
}

// ManifestExists checks if the Hub manifest file exists in the cache
func ManifestExists() bool {
	cache, err := openHubCache()
	if err != nil {
		return false
	}

	return cache.manifest().exists()
}
