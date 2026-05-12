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

// GetIndex retrieves and parses the Remote index file from the cache.
func GetIndex() (*t.RemoteIndex, error) {
	cache, err := openRemoteCache()
	if err != nil {
		return nil, err
	}

	return cache.index().read()
}

// IndexExists checks if the Remote index file exists in the cache.
func IndexExists() bool {
	cache, err := openRemoteCache()
	if err != nil {
		return false
	}

	return cache.index().exists()
}
