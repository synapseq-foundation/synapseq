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

package types

const (
	// HubManifestURL is the URL to fetch the Hub manifest
	HubManifestURL = "https://hub.synapseq.org/manifest.json"
)

// HubDependencyType represents the type of a Hub dependency
type HubDependencyType string

const (
	// HubDependencyTypeExtends represents a extend dependency
	HubDependencyTypeExtends HubDependencyType = "extends"
	// HubDependencyTypeAmbiance represents a ambiance dependency
	HubDependencyTypeAmbiance HubDependencyType = "ambiance"
)

// String returns the string representation of the HubDependencyType
func (dt HubDependencyType) String() string {
	return string(dt)
}

// HubDependency represents a dependency for a Hub entry
type HubDependency struct {
	Type        HubDependencyType `json:"type"`
	ID          string            `json:"id"`
	DownloadUrl string            `json:"download_url"`
}

// HubEntry represents an entry in the Hub index
type HubEntry struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Author       string          `json:"author"`
	Category     string          `json:"category"`
	Path         string          `json:"path"`
	DownloadUrl  string          `json:"download_url"`
	UpdatedAt    string          `json:"updated_at"`
	Dependencies []HubDependency `json:"dependencies,omitempty"`
}

// HubManifest represents the manifest of available Hub entries
type HubManifest struct {
	Version     string
	LastUpdated string
	Entries     []HubEntry
}
