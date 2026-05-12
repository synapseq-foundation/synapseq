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
