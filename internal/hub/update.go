//go:build !nohub

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

package hub

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HubUpdate updates the local Hub manifest cache
func HubUpdate() error {
	cache, err := GetCacheDir()
	if err != nil {
		return err
	}

	resp, err := http.Get(t.HubManifestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("invalid content-type for manifest file: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	manifestPath := cache + "/manifest.json"
	if err = os.WriteFile(manifestPath, data, 0644); err != nil {
		return err
	}

	return nil
}
