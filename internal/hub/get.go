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
	"fmt"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HubGet retrieves a sequence by its ID from the Hub
func HubGet(sequenceID string) (*t.HubEntry, error) {
	catalog, err := loadManifestCatalog()
	if err != nil {
		return nil, err
	}

	return catalog.findEntry(sequenceID), nil
}

// HubDownload downloads a sequence and its dependencies from the Hub
func HubDownload(entry *t.HubEntry) (string, error) {
	if entry == nil {
		return "", fmt.Errorf("hub entry is nil")
	}

	cache, err := openHubCache()
	if err != nil {
		return "", err
	}

	entryCache, err := prepareEntryDownload(cache, entry)
	if err != nil {
		return "", err
	}

	cached, err := entryCache.hasSequence()
	if err != nil {
		return "", err
	}
	if cached {
		return entryCache.sequencePath(), nil
	}

	if err := downloadEntryDependencies(entryCache, entry); err != nil {
		return "", err
	}

	if err := downloadEntrySequence(entryCache, entry); err != nil {
		return "", err
	}

	return entryCache.sequencePath(), nil
}

func prepareEntryDownload(cache *hubCache, entry *t.HubEntry) (entryCache, error) {
	entryCache := cache.entry(entry)
	if err := entryCache.prepare(); err != nil {
		return entryCache, err
	}

	return entryCache, nil
}

func downloadEntryDependencies(cache entryCache, entry *t.HubEntry) error {
	for _, dependency := range entry.Dependencies {
		if err := downloadFile(dependency.DownloadUrl, cache.dependencyPath(dependency)); err != nil {
			return fmt.Errorf("error saving dependency %s: %v", dependency.ID, err)
		}
	}

	return nil
}

func downloadEntrySequence(cache entryCache, entry *t.HubEntry) error {
	if err := downloadFile(entry.DownloadUrl, cache.sequencePath()); err != nil {
		return fmt.Errorf("error saving sequence %s: %v", entry.ID, err)
	}

	return nil
}
