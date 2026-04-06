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

import (
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HubUpdate updates the local Hub manifest cache
func HubUpdate() error {
	cache, err := openHubCache()
	if err != nil {
		return err
	}

	data, response, err := downloadURL(t.HubManifestURL)
	if err != nil {
		return err
	}
	if err := validateJSONContentType(response); err != nil {
		return err
	}

	if err = cache.manifest().write(data); err != nil {
		return err
	}

	return nil
}
