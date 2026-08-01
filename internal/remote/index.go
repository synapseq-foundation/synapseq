// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// GetIndex retrieves and parses the Remote index file from the cache.
func GetIndex() (*t.RemoteIndex, error) {
	cache, err := openRemoteCache()
	if err != nil {
		return nil, err
	}

	return cache.index().read()
}

// GetIndexURL retrieves and parses the index cached for a custom Remote base URL.
func GetIndexURL(baseURL string) (*t.RemoteIndex, error) {
	source, err := customRemoteSource(baseURL)
	if err != nil {
		return nil, err
	}
	cache, err := openRemoteCacheForSource(source)
	if err != nil {
		return nil, err
	}

	return cache.index().read()
}

// IndexExists checks if the Remote index file exists in the cache.
func IndexExists() bool {
	cache, err := openRemoteCache()
	if err != nil {
		return false
	}

	return cache.index().exists()
}

func parseRemoteIndex(data []byte) (*t.RemoteIndex, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var index t.RemoteIndex
	if err := decoder.Decode(&index); err != nil {
		return nil, fmt.Errorf("invalid remote index: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("invalid remote index: multiple JSON values")
	}

	if strings.TrimSpace(index.Version) == "" {
		return nil, fmt.Errorf("invalid remote index: version is required")
	}
	if _, err := time.Parse(time.RFC3339, index.LastUpdated); err != nil {
		return nil, fmt.Errorf("invalid remote index: lastUpdated must be RFC 3339")
	}

	seenIDs := make(map[string]struct{}, len(index.Entries))
	for _, entry := range index.Entries {
		if err := validateRemoteEntry(entry); err != nil {
			return nil, err
		}
		if _, exists := seenIDs[entry.ID]; exists {
			return nil, fmt.Errorf("invalid remote index: duplicate entry ID %q", entry.ID)
		}
		seenIDs[entry.ID] = struct{}{}
	}

	return &index, nil
}

func validateRemoteEntry(entry t.RemoteEntry) error {
	if !validRemoteID(entry.ID) {
		return fmt.Errorf("invalid remote index: invalid entry ID %q", entry.ID)
	}
	if strings.TrimSpace(entry.Name) == "" || strings.TrimSpace(entry.Category) == "" {
		return fmt.Errorf("invalid remote index: name and category are required for %q", entry.ID)
	}
	if entry.DurationMinutes <= 0 {
		return fmt.Errorf("invalid remote index: durationMinutes must be positive for %q", entry.ID)
	}
	if _, err := time.Parse(time.RFC3339, entry.CreatedAt); err != nil {
		return fmt.Errorf("invalid remote index: createdAt must be RFC 3339 for %q", entry.ID)
	}
	if !validRemoteSequencePath(entry.Sequence) {
		return fmt.Errorf("invalid remote index: invalid sequence path for %q", entry.ID)
	}
	if entry.Artwork != "" && !validRemoteRelativePath(entry.Artwork) {
		return fmt.Errorf("invalid remote index: invalid artwork path for %q", entry.ID)
	}

	return nil
}

func validRemoteID(id string) bool {
	if id == "" || strings.HasPrefix(id, "-") || strings.HasSuffix(id, "-") {
		return false
	}
	for _, char := range id {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '-' {
			continue
		}
		return false
	}
	return true
}

func validRemoteSequencePath(value string) bool {
	return strings.HasSuffix(value, ".spsq") && validRemoteRelativePath(value)
}

func validRemoteRelativePath(value string) bool {
	if value == "" || strings.HasPrefix(value, "/") || strings.Contains(value, "\\") || strings.Contains(value, "://") || strings.ContainsAny(value, "?#") {
		return false
	}

	return path.Clean(value) == value && !strings.HasPrefix(value, "../") && value != ".."
}
