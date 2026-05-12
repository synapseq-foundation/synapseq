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
	// RemoteIndexURL is the URL to fetch the Remote index.
	RemoteIndexURL = "https://sequence.synapseq.org/free/index.json"
)

// RemoteEntry represents an entry in the Remote index.
type RemoteEntry struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
	Sequence        string `json:"sequence"`
	Category        string `json:"category"`
	CreatedAt       string `json:"createdAt"`
}

// RemoteIndex represents the index of available Remote entries.
type RemoteIndex struct {
	Version     string        `json:"version"`
	LastUpdated string        `json:"lastUpdated"`
	Entries     []RemoteEntry `json:"entries"`
}
