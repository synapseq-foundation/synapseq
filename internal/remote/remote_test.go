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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestFindEntryByID(ts *testing.T) {
	index := &t.RemoteIndex{
		Entries: []t.RemoteEntry{{ID: "focus-pack"}, {ID: "sleep-pack"}},
	}

	entry := findEntryByID(index, "sleep-pack")
	if entry == nil {
		ts.Fatalf("expected entry to be found")
	}
	if entry.ID != "sleep-pack" {
		ts.Fatalf("expected sleep-pack, got %q", entry.ID)
	}

	if missing := findEntryByID(index, "missing"); missing != nil {
		ts.Fatalf("expected missing entry lookup to return nil")
	}
}

func TestIndexCatalogFindEntry(ts *testing.T) {
	catalog := &indexCatalog{index: &t.RemoteIndex{Entries: []t.RemoteEntry{{ID: "focus-pack"}}}}
	if entry := catalog.findEntry("focus-pack"); entry == nil || entry.ID != "focus-pack" {
		ts.Fatalf("expected catalog to resolve focus-pack, got %#v", entry)
	}
	if entry := catalog.findEntry("missing"); entry != nil {
		ts.Fatalf("expected missing catalog lookup to return nil")
	}
}

func TestRemoteSequenceURL(ts *testing.T) {
	if got := remoteSequenceURL("/free/focus.spsq"); got != "https://sequence.synapseq.org/free/focus.spsq" {
		ts.Fatalf("unexpected url: %q", got)
	}
	if got := remoteSequenceURL("free/focus.spsq"); got != "https://sequence.synapseq.org/free/focus.spsq" {
		ts.Fatalf("unexpected url: %q", got)
	}
}

func TestDownloadURLRejectsUnexpectedStatus(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "nope", http.StatusBadGateway)
	}))
	defer server.Close()

	_, _, err := downloadURL(server.URL)
	if err == nil {
		ts.Fatalf("expected unexpected status to fail")
	}
}

func TestValidateJSONContentType(ts *testing.T) {
	response := &http.Response{Header: make(http.Header)}
	response.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err := validateJSONContentType(response); err != nil {
		ts.Fatalf("expected json content type to pass, got %v", err)
	}

	response.Header.Set("Content-Type", "text/plain")
	if err := validateJSONContentType(response); err == nil {
		ts.Fatalf("expected text/plain content type to fail")
	}
}

func TestEntryCacheSequencePath(ts *testing.T) {
	cache := entryCache{dir: "/tmp/remote-entry", entry: &t.RemoteEntry{ID: "focus-pack"}}
	if path := cache.sequencePath(); path != filepath.Join(cache.dir, "focus-pack.spsq") {
		ts.Fatalf("expected sequence path to end with focus-pack.spsq, got %q", path)
	}
}

func TestEntryCacheHasSequence(ts *testing.T) {
	tempDir := ts.TempDir()
	cache := entryCache{dir: tempDir, entry: &t.RemoteEntry{ID: "cached"}}
	sequencePath := cache.sequencePath()

	cached, err := cache.hasSequence()
	if err != nil {
		ts.Fatalf("unexpected error checking missing cache file: %v", err)
	}
	if cached {
		ts.Fatalf("expected missing file to not be cached")
	}

	if err := os.WriteFile(sequencePath, []byte("test"), 0644); err != nil {
		ts.Fatalf("failed to seed cached sequence: %v", err)
	}

	cached, err = cache.hasSequence()
	if err != nil {
		ts.Fatalf("unexpected error checking cached file: %v", err)
	}
	if !cached {
		ts.Fatalf("expected existing file to be reported as cached")
	}
}

func TestIndexCacheExists(ts *testing.T) {
	tempDir := ts.TempDir()
	cache := indexCache{path: filepath.Join(tempDir, "index.json")}
	if cache.exists() {
		ts.Fatalf("expected missing index cache to report false")
	}

	if err := os.WriteFile(cache.path, []byte("{}"), 0644); err != nil {
		ts.Fatalf("failed to seed index cache: %v", err)
	}

	if !cache.exists() {
		ts.Fatalf("expected written index cache to report true")
	}
}
