// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/remote"
)

const invalidRemoteSequenceIndex = `{
  "version": "3.0.0",
  "lastUpdated": "2026-05-13T22:30:42Z",
  "entries": [{
    "id": "calm-state",
    "name": "Calm State",
    "description": "A test sequence.",
    "durationMinutes": 15,
    "sequence": "free/relaxation/calm-state/calm-state.spsq",
    "artwork": "free/relaxation/calm-state/calm-state.webp",
    "category": "Relaxation",
    "createdAt": "2026-03-21T00:00:00Z"
  }]
}`

func TestRemoteCommandsRejectInvalidDownloadedSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/index.json":
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(invalidRemoteSequenceIndex))
		case "/free/relaxation/calm-state/calm-state.spsq":
			_, _ = writer.Write([]byte("not valid SPSQ"))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	t.Setenv("HOME", t.TempDir())
	t.Setenv("SYNAPSEQ_REMOTE_BASE_URL", server.URL)
	if err := remote.RemoteSync(); err != nil {
		t.Fatalf("RemoteSync error: %v", err)
	}

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "download",
			run: func() error {
				return remoteRunDownload("calm-state", t.TempDir(), true)
			},
		},
		{
			name: "get",
			run: func() error {
				return remoteRunGet("calm-state", filepath.Join(t.TempDir(), "calm-state.wav"), &cli.CLIOptions{Quiet: true, Test: true})
			},
		},
		{
			name: "info",
			run: func() error {
				return remoteRunInfo("calm-state")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); err == nil {
				t.Fatal("expected invalid downloaded SPSQ to be rejected")
			}
		})
	}
}
