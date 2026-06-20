// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
