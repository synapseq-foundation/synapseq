// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"bytes"
	"testing"
)

func TestStatusModes(t *testing.T) {
	var out bytes.Buffer
	base := NewAppContext()
	progress := base.WithProgress(&out, true)
	verbose := progress.WithVerbose(&out, false)
	if base.statusOutput != nil || progress.Verbose() || !progress.statusProgress || !verbose.Verbose() || verbose.statusProgress {
		t.Fatal("incorrect status configuration")
	}
	if verbose.WithProgress(&out, false).Verbose() || verbose.WithVerbose(nil, false).Verbose() {
		t.Fatal("mode selection or nil writer")
	}
	loaded, err := progress.LoadContent("alpha\n  tone 100 binaural 1 amplitude 1\n00:00:00 alpha\n00:00:01 alpha\n")
	if err != nil {
		t.Fatal(err)
	}
	options := loaded.buildAudioRendererOptions(loaded.sequence)
	if !options.Progress || !options.Colors || options.StatusOutput != &out {
		t.Fatal("options not propagated")
	}
}
