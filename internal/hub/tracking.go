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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/synapseq-foundation/synapseq/v4/internal/info"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// TrackDownload sends an anonymous download event to the SynapSeq Hub analytics endpoint.
// It only sends technical metadata, no personal or identifying information.
func TrackDownload(sequenceID string, action t.HubActionTracking) error {
	if info.VERSION == "development" {
		// Do not track in development mode
		return nil
	}

	if sequenceID == "" {
		return fmt.Errorf("empty sequence ID")
	}

	payload := map[string]string{
		"id": sequenceID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", t.HubTrackEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SYNAPSEQ-SOURCE", "CLI")
	req.Header.Set("X-SYNAPSEQ-VERSION", info.VERSION)
	req.Header.Set("X-SYNAPSEQ-PLATFORM", runtime.GOOS)
	req.Header.Set("X-SYNAPSEQ-ARCH", runtime.GOARCH)
	req.Header.Set("X-SYNAPSEQ-ACTION", strings.ToUpper(action.String()))

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Fail silently. Tracking must never break CLI functionality
		return nil
	}
	resp.Body.Close()

	return nil
}
