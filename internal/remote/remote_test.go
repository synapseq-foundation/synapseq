// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package remote

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

//go:embed testdata/index.json
var customIndexFixture []byte

const validRemoteSPSQ = `entry
  tone 220 binaural 10 amplitude 10
00:00:00 entry
00:00:01 entry
`

type remoteTestServer struct {
	index     []byte
	sequence  []byte
	requested []string
}

func newRemoteTestServer() (*remoteTestServer, *httptest.Server) {
	fixture := &remoteTestServer{index: customIndexFixture, sequence: []byte(validRemoteSPSQ)}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fixture.requested = append(fixture.requested, request.URL.Path)
		switch {
		case request.URL.Path == "/index.json":
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(fixture.index)
		case strings.HasSuffix(request.URL.Path, ".spsq"):
			_, _ = writer.Write(fixture.sequence)
		default:
			http.NotFound(writer, request)
		}
	}))
	return fixture, server
}

func customIndexMutation(ts *testing.T, mutate func(map[string]any)) []byte {
	ts.Helper()
	var index map[string]any
	if err := json.Unmarshal(customIndexFixture, &index); err != nil {
		ts.Fatalf("unmarshal fixture: %v", err)
	}
	mutate(index)
	data, err := json.Marshal(index)
	if err != nil {
		ts.Fatalf("marshal fixture: %v", err)
	}
	return data
}

func fixtureEntry(index map[string]any, position int) map[string]any {
	return index["entries"].([]any)[position].(map[string]any)
}

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

func TestCustomRemoteSource(ts *testing.T) {
	source, err := customRemoteSource("https://my-sequences.com:8443/")
	if err != nil {
		ts.Fatalf("customRemoteSource error: %v", err)
	}
	if source.cacheKey != "my-sequences.com_8443" {
		ts.Fatalf("unexpected cache key: %q", source.cacheKey)
	}
	if got := source.indexURL(); got != "https://my-sequences.com:8443/index.json" {
		ts.Fatalf("unexpected index URL: %q", got)
	}
	if got := source.sequenceURL("focus.spsq"); got != "https://my-sequences.com:8443/focus.spsq" {
		ts.Fatalf("unexpected sequence URL: %q", got)
	}
}

func TestCustomRemoteSourceRejectsNonBaseURL(ts *testing.T) {
	invalidURLs := []string{
		"https://my-sequences.com/sequences",
		"ftp://my-sequences.com",
		"https://user@my-sequences.com",
		"https://my-sequences.com?query=value",
		"https://my-sequences.com#fragment",
	}
	for _, rawURL := range invalidURLs {
		if _, err := customRemoteSource(rawURL); err == nil {
			ts.Fatalf("expected %q to be rejected", rawURL)
		}
	}
}

func TestDefaultRemoteSourceUsesOfficialIndex(ts *testing.T) {
	ts.Setenv(remoteBaseURLEnv, "")
	source, err := defaultRemoteSource()
	if err != nil {
		ts.Fatalf("defaultRemoteSource error: %v", err)
	}
	if got := source.indexURL(); got != t.RemoteIndexURL {
		ts.Fatalf("unexpected index URL: %q", got)
	}
}

func TestParseRemoteIndex(t *testing.T) {
	valid := []byte(`{
  "version": "1",
  "lastUpdated": "2026-08-01T00:00:00Z",
  "entries": [{
    "id": "focus-pack",
    "name": "Focus Pack",
    "description": "",
    "durationMinutes": 15,
    "sequence": "focus/focus-pack.spsq",
    "artwork": "focus/focus-pack.webp",
    "category": "Focus",
    "createdAt": "2026-08-01T00:00:00Z"
  }]
}`)
	if _, err := parseRemoteIndex(valid); err != nil {
		t.Fatalf("parseRemoteIndex error: %v", err)
	}

	invalid := [][]byte{
		[]byte(`{"version":"1","lastUpdated":"invalid","entries":[]}`),
		[]byte(`{"version":"1","lastUpdated":"2026-08-01T00:00:00Z","entries":[{"id":"../escape","name":"Escape","durationMinutes":1,"sequence":"escape.spsq","category":"Focus","createdAt":"2026-08-01T00:00:00Z"}]}`),
		[]byte(`{"version":"1","lastUpdated":"2026-08-01T00:00:00Z","entries":[],"unexpected":true}`),
	}
	for _, data := range invalid {
		if _, err := parseRemoteIndex(data); err == nil {
			t.Fatalf("expected index to be rejected: %s", data)
		}
	}
}

func TestValidateRemoteSequence(t *testing.T) {
	valid := []byte("alpha\n  tone 220 binaural 10 amplitude 10\n00:00:00 alpha\n00:00:01 alpha\n")
	if err := validateRemoteSequence(valid, "focus.spsq", "."); err != nil {
		t.Fatalf("validateRemoteSequence error: %v", err)
	}
	if err := validateRemoteSequence([]byte("not valid"), "focus.spsq", "."); err == nil {
		t.Fatal("expected invalid SPSQ to be rejected")
	}
}

func TestRemoteDownloadRejectsInvalidSequenceWithoutCaching(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/index.json":
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"version":"1","lastUpdated":"2026-08-01T00:00:00Z","entries":[{"id":"focus","name":"Focus","description":"","durationMinutes":1,"sequence":"focus.spsq","category":"Focus","createdAt":"2026-08-01T00:00:00Z"}]}`))
		case "/focus.spsq":
			_, _ = writer.Write([]byte("not valid SPSQ"))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	ts.Setenv(remoteBaseURLEnv, server.URL)
	ts.Setenv("HOME", ts.TempDir())
	if err := RemoteSync(); err != nil {
		ts.Fatalf("RemoteSync error: %v", err)
	}
	entry, err := RemoteGet("focus")
	if err != nil {
		ts.Fatalf("RemoteGet error: %v", err)
	}
	if _, err := RemoteDownload(entry); err == nil {
		ts.Fatal("expected invalid SPSQ to be rejected")
	}

	cache, err := openRemoteCache()
	if err != nil {
		ts.Fatalf("openRemoteCache error: %v", err)
	}
	cached, err := cache.entry(entry).hasSequence()
	if err != nil {
		ts.Fatalf("hasSequence error: %v", err)
	}
	if cached {
		ts.Fatal("invalid SPSQ must not be cached")
	}
}

func TestCustomRemoteServerSyncsOfficialFormatIndexAndDownloadsSequence(ts *testing.T) {
	fixture, server := newRemoteTestServer()
	defer server.Close()

	ts.Setenv("HOME", ts.TempDir())
	if err := RemoteSyncURL(server.URL); err != nil {
		ts.Fatalf("RemoteSyncURL error: %v", err)
	}
	ts.Setenv(remoteBaseURLEnv, server.URL)

	index, err := GetIndex()
	if err != nil {
		ts.Fatalf("GetIndex error: %v", err)
	}
	if len(index.Entries) != 3 {
		ts.Fatalf("expected 3 entries, got %d", len(index.Entries))
	}
	entry, err := RemoteGet("calm-state")
	if err != nil || entry == nil {
		ts.Fatalf("RemoteGet error: %v, entry: %#v", err, entry)
	}
	sequencePath, err := RemoteDownload(entry)
	if err != nil {
		ts.Fatalf("RemoteDownload error: %v", err)
	}
	content, err := os.ReadFile(sequencePath)
	if err != nil {
		ts.Fatalf("read cached sequence: %v", err)
	}
	if !bytes.Equal(content, fixture.sequence) {
		ts.Fatal("cached sequence content differs from server response")
	}
	if !slicesContain(fixture.requested, "/index.json") || !slicesContain(fixture.requested, "/free/relaxation/calm-state/calm-state.spsq") {
		ts.Fatalf("unexpected remote paths: %v", fixture.requested)
	}

	source, err := customRemoteSource(server.URL)
	if err != nil {
		ts.Fatalf("customRemoteSource error: %v", err)
	}
	cacheDir, err := getCacheDir(source)
	if err != nil {
		ts.Fatalf("getCacheDir error: %v", err)
	}
	if !strings.HasSuffix(cacheDir, filepath.Join("custom", source.cacheKey)) {
		ts.Fatalf("custom cache directory not isolated: %q", cacheDir)
	}
}

func TestRemoteSyncRejectsInvalidHostedIndexWithoutReplacingCache(ts *testing.T) {
	fixture, server := newRemoteTestServer()
	defer server.Close()

	ts.Setenv("HOME", ts.TempDir())
	if err := RemoteSyncURL(server.URL); err != nil {
		ts.Fatalf("initial RemoteSyncURL error: %v", err)
	}
	source, err := customRemoteSource(server.URL)
	if err != nil {
		ts.Fatalf("customRemoteSource error: %v", err)
	}
	cache, err := openRemoteCacheForSource(source)
	if err != nil {
		ts.Fatalf("openRemoteCacheForSource error: %v", err)
	}
	original, err := os.ReadFile(cache.index().path)
	if err != nil {
		ts.Fatalf("read cached index: %v", err)
	}
	fixture.index = customIndexMutation(ts, func(index map[string]any) {
		fixtureEntry(index, 0)["sequence"] = "../escape.spsq"
	})

	if err := RemoteSyncURL(server.URL); err == nil {
		ts.Fatal("expected invalid hosted index to be rejected")
	}
	after, err := os.ReadFile(cache.index().path)
	if err != nil {
		ts.Fatalf("read cached index after rejected sync: %v", err)
	}
	if !bytes.Equal(after, original) {
		ts.Fatal("rejected index must not replace the cached index")
	}
}

func TestRemoteSyncRejectsHostedIndexValidationFailures(ts *testing.T) {
	tests := []struct {
		name  string
		index []byte
	}{
		{name: "malformed JSON", index: []byte("{")},
		{name: "unknown field", index: customIndexMutation(ts, func(index map[string]any) { index["unexpected"] = true })},
		{name: "missing version", index: customIndexMutation(ts, func(index map[string]any) { index["version"] = "" })},
		{name: "invalid last updated", index: customIndexMutation(ts, func(index map[string]any) { index["lastUpdated"] = "invalid" })},
		{name: "invalid ID", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["id"] = "../escape" })},
		{name: "duplicate ID", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 1)["id"] = fixtureEntry(index, 0)["id"] })},
		{name: "empty name", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["name"] = "" })},
		{name: "empty category", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["category"] = "" })},
		{name: "non-positive duration", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["durationMinutes"] = 0 })},
		{name: "invalid created at", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["createdAt"] = "invalid" })},
		{name: "absolute sequence", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["sequence"] = "/escape.spsq" })},
		{name: "external sequence", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["sequence"] = "https://example.com/escape.spsq" })},
		{name: "non SPSQ sequence", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["sequence"] = "escape.txt" })},
		{name: "unsafe artwork", index: customIndexMutation(ts, func(index map[string]any) { fixtureEntry(index, 0)["artwork"] = "../escape.webp" })},
	}

	for _, test := range tests {
		ts.Run(test.name, func(ts *testing.T) {
			fixture, server := newRemoteTestServer()
			defer server.Close()
			fixture.index = test.index
			ts.Setenv("HOME", ts.TempDir())

			if err := RemoteSyncURL(server.URL); err == nil {
				ts.Fatal("expected RemoteSyncURL to reject invalid index")
			}
			source, err := customRemoteSource(server.URL)
			if err != nil {
				ts.Fatalf("customRemoteSource error: %v", err)
			}
			cache, err := openRemoteCacheForSource(source)
			if err != nil {
				ts.Fatalf("openRemoteCacheForSource error: %v", err)
			}
			if cache.index().exists() {
				ts.Fatal("invalid index must not be cached")
			}
		})
	}
}

func TestRemoteDownloadRejectsInvalidCachedSequence(ts *testing.T) {
	_, server := newRemoteTestServer()
	defer server.Close()

	ts.Setenv("HOME", ts.TempDir())
	if err := RemoteSyncURL(server.URL); err != nil {
		ts.Fatalf("RemoteSyncURL error: %v", err)
	}
	ts.Setenv(remoteBaseURLEnv, server.URL)
	entry, err := RemoteGet("calm-state")
	if err != nil {
		ts.Fatalf("RemoteGet error: %v", err)
	}
	sequencePath, err := RemoteDownload(entry)
	if err != nil {
		ts.Fatalf("RemoteDownload error: %v", err)
	}
	if err := os.WriteFile(sequencePath, []byte("invalid SPSQ"), 0644); err != nil {
		ts.Fatalf("overwrite cached sequence: %v", err)
	}
	if _, err := RemoteDownload(entry); err == nil {
		ts.Fatal("expected invalid cached SPSQ to be rejected")
	}
}

func slicesContain(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
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
