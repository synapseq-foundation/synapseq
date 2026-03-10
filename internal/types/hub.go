/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	// HubTrackEndpoint is the endpoint for tracking downloads
	HubTrackEndpoint = "https://us-central1-synapseq-hub.cloudfunctions.net/trackDownload"
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

// HubActionTracking represents the type of action being tracked
type HubActionTracking string

const (
	// HubActionTrackingGet represents the "get" action being tracked
	HubActionTrackingGet HubActionTracking = "get"
	// HubActionTrackingInfo represents the "info" action being tracked
	HubActionTrackingInfo HubActionTracking = "info"
	// HubActionTrackingDownload represents the "download" action being tracked
	HubActionTrackingDownload HubActionTracking = "download"
)

// String returns the string representation of the HubActionTracking
func (at HubActionTracking) String() string {
	return string(at)
}

// HubDependency represents a dependency for a Hub entry
type HubDependency struct {
	Type        HubDependencyType `json:"type"`
	Name        string            `json:"name"`
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
