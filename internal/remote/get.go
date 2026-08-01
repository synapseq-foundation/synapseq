// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	"fmt"
	"os"

	"github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// RemoteGet retrieves a sequence by its ID from the Remote index.
func RemoteGet(sequenceID string) (*t.RemoteEntry, error) {
	catalog, err := loadIndexCatalog()
	if err != nil {
		return nil, err
	}

	return catalog.findEntry(sequenceID), nil
}

// RemoteDownload downloads a sequence from Remote.
func RemoteDownload(entry *t.RemoteEntry) (string, error) {
	if entry == nil {
		return "", fmt.Errorf("remote entry is nil")
	}

	cache, err := openRemoteCache()
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
		if err := validateCachedSequence(entryCache); err != nil {
			return "", err
		}
		return entryCache.sequencePath(), nil
	}

	if err := downloadEntrySequence(cache, entryCache, entry); err != nil {
		return "", err
	}

	return entryCache.sequencePath(), nil
}

func prepareEntryDownload(cache *remoteCache, entry *t.RemoteEntry) (entryCache, error) {
	entryCache := cache.entry(entry)
	if err := entryCache.prepare(); err != nil {
		return entryCache, err
	}

	return entryCache, nil
}

func downloadEntrySequence(remoteCache *remoteCache, cache entryCache, entry *t.RemoteEntry) error {
	data, _, err := downloadURL(remoteCache.source.sequenceURL(entry.Sequence))
	if err != nil {
		return fmt.Errorf("error downloading sequence %s: %v", entry.ID, err)
	}
	if err := validateRemoteSequence(data, cache.sequencePath(), cache.dir); err != nil {
		return fmt.Errorf("invalid remote sequence %s: %w", entry.ID, err)
	}
	if err := os.WriteFile(cache.sequencePath(), data, 0644); err != nil {
		return fmt.Errorf("error saving sequence %s: %v", entry.ID, err)
	}

	return nil
}

func validateCachedSequence(cache entryCache) error {
	data, err := os.ReadFile(cache.sequencePath())
	if err != nil {
		return err
	}

	return validateRemoteSequence(data, cache.sequencePath(), cache.dir)
}

func validateRemoteSequence(data []byte, sourceFile, baseRef string) error {
	_, err := sequence.LoadTextSequence(data, sourceFile, baseRef)
	return err
}
