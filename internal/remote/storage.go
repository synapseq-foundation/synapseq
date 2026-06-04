// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	"encoding/json"
	"os"
	"path/filepath"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type remoteCache struct {
	root string
}

type indexCache struct {
	path string
}

type entryCache struct {
	dir   string
	entry *t.RemoteEntry
}

func openRemoteCache() (*remoteCache, error) {
	root, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	return &remoteCache{root: root}, nil
}

func (cache *remoteCache) index() indexCache {
	return indexCache{path: filepath.Join(cache.root, "index.json")}
}

func (cache *remoteCache) entry(entry *t.RemoteEntry) entryCache {
	return entryCache{
		dir:   filepath.Join(cache.root, entry.ID),
		entry: entry,
	}
}

func (cache indexCache) read() (*t.RemoteIndex, error) {
	return readIndexFile(cache.path)
}

func (cache indexCache) write(data []byte) error {
	return writeIndexFile(cache.path, data)
}

func (cache indexCache) exists() bool {
	if _, err := os.Stat(cache.path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (cache entryCache) prepare() error {
	return os.MkdirAll(cache.dir, 0755)
}

func (cache entryCache) sequencePath() string {
	return filepath.Join(cache.dir, cache.entry.ID+".spsq")
}

func (cache entryCache) hasSequence() (bool, error) {
	sequencePath := cache.sequencePath()
	if _, err := os.Stat(sequencePath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func readIndexFile(path string) (*t.RemoteIndex, error) {
	indexData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var index *t.RemoteIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		return nil, err
	}

	return index, nil
}

func writeIndexFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func findEntryByID(index *t.RemoteIndex, sequenceID string) *t.RemoteEntry {
	if index == nil {
		return nil
	}

	for _, entry := range index.Entries {
		if entry.ID == sequenceID {
			match := entry
			return &match
		}
	}

	return nil
}
