// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// RemoteSync updates the local Remote index cache.
func RemoteSync() error {
	cache, err := openRemoteCache()
	if err != nil {
		return err
	}

	data, response, err := downloadURL(t.RemoteIndexURL)
	if err != nil {
		return err
	}
	if err := validateJSONContentType(response); err != nil {
		return err
	}

	if err = cache.index().write(data); err != nil {
		return err
	}

	return nil
}
