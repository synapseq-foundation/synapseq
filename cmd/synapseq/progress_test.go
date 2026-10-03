// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func TestRenderStatusModes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		opts   cli.CLIOptions
		target string
		want   string
	}{
		{"default", cli.CLIOptions{}, "out.wav", "Audio rendered: 100%"},
		{"verbose", cli.CLIOptions{Verbose: true}, "out.wav", "->"},
		{"quiet", cli.CLIOptions{Quiet: true, Verbose: true}, "out.wav", ""},
		{"stdout", cli.CLIOptions{}, "-", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var status, pcm bytes.Buffer
			loaded, err := newAppContext(tc.target, &status, &tc.opts).LoadContent("alpha\n  tone 100 binaural 1 amplitude 1\n00:00:00 alpha\n00:00:01 alpha\n")
			if err != nil {
				t.Fatal(err)
			}
			if err := loaded.Stream(&pcm); err != nil {
				t.Fatal(err)
			}
			if tc.want == "" {
				if status.Len() != 0 {
					t.Fatal(status.String())
				}
			} else if !strings.Contains(status.String(), tc.want) {
				t.Fatal(status.String())
			}
			if tc.opts.Verbose && strings.Contains(status.String(), "Audio rendered") {
				t.Fatal("verbose contains progress")
			}
		})
	}
}
