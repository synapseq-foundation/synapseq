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
	"path/filepath"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type hubCache struct {
	root string
}

type manifestCache struct {
	path string
}

type entryCache struct {
	dir   string
	entry *t.HubEntry
}

func openHubCache() (*hubCache, error) {
	root, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	return &hubCache{root: root}, nil
}

func (cache *hubCache) manifest() manifestCache {
	return manifestCache{path: filepath.Join(cache.root, "manifest.json")}
}

func (cache *hubCache) entry(entry *t.HubEntry) entryCache {
	return entryCache{
		dir:   filepath.Join(cache.root, strings.TrimSuffix(entry.Path, ".spsq")),
		entry: entry,
	}
}

func (cache manifestCache) read() (*t.HubManifest, error) {
	return readManifestFile(cache.path)
}

func (cache manifestCache) write(data []byte) error {
	return writeManifestFile(cache.path, data)
}

func (cache manifestCache) exists() bool {
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

func (cache entryCache) dependencyPath(dependency t.HubDependency) string {
	extension := ".spsc"
	if dependency.Type == t.HubDependencyTypeAmbiance {
		extension = ".wav"
	}

	return filepath.Join(cache.dir, dependency.ID+extension)
}

func readManifestFile(path string) (*t.HubManifest, error) {
	manifestData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest *t.HubManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func writeManifestFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func findEntryByID(manifest *t.HubManifest, sequenceID string) *t.HubEntry {
	if manifest == nil {
		return nil
	}

	for _, entry := range manifest.Entries {
		if entry.ID == sequenceID {
			match := entry
			return &match
		}
	}

	return nil
}
