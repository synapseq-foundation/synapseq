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
	"encoding/json"
	"os"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// GetManifest retrieves and parses the Hub manifest file from the cache
func GetManifest() (*t.HubManifest, error) {
	cache, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	manifestPath := cache + "/manifest.json"
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var hubManifest *t.HubManifest
	if err := json.Unmarshal(manifest, &hubManifest); err != nil {
		return nil, err
	}

	return hubManifest, nil
}

// ManifestExists checks if the Hub manifest file exists in the cache
func ManifestExists() bool {
	cache, err := GetCacheDir()
	if err != nil {
		return false
	}

	manifestPath := cache + "/manifest.json"
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return false
	}

	return true
}
