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
	"fmt"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const remoteSequenceRootURL = "https://sequence.synapseq.org"

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
		return entryCache.sequencePath(), nil
	}

	if err := downloadEntrySequence(entryCache, entry); err != nil {
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

func downloadEntrySequence(cache entryCache, entry *t.RemoteEntry) error {
	if err := downloadFile(remoteSequenceURL(entry.Sequence), cache.sequencePath()); err != nil {
		return fmt.Errorf("error saving sequence %s: %v", entry.ID, err)
	}

	return nil
}

func remoteSequenceURL(sequencePath string) string {
	return remoteSequenceRootURL + "/" + strings.TrimPrefix(sequencePath, "/")
}
