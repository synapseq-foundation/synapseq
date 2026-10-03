// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audio

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestRenderProgressConsumedFrames(tst *testing.T) {
	for _, fail := range []bool{false, true} {
		var status bytes.Buffer
		renderer, err := NewAudioRenderer([]t.Period{{Time: 0}, {Time: 1234}}, &AudioRendererOptions{SampleRate: 8000, Volume: 100, StatusOutput: &status, Progress: true})
		if err != nil {
			tst.Fatal(err)
		}
		var frames int
		sentinel := errors.New("consumer failed")
		err = renderer.Render(func(samples []int) error {
			if fail {
				return sentinel
			}
			frames += len(samples) / 2
			return nil
		})
		if fail {
			if !errors.Is(err, sentinel) || strings.Contains(status.String(), "100%") || !strings.Contains(status.String(), "failed: 0%") {
				tst.Fatalf("%v %s", err, status.String())
			}
		} else if err != nil || frames != 9872 || !strings.Contains(status.String(), "100%  00:00:01 / 00:00:01") {
			tst.Fatalf("frames=%d err=%v status=%s", frames, err, status.String())
		}
	}
}

func TestRenderProgressZeroDuration(tst *testing.T) {
	var status bytes.Buffer
	renderer, err := NewAudioRenderer([]t.Period{{Time: 0}}, &AudioRendererOptions{SampleRate: 8000, Volume: 100, StatusOutput: &status, Progress: true})
	if err != nil {
		tst.Fatal(err)
	}
	if err := renderer.Render(func([]int) error { tst.Fatal("unexpected audio"); return nil }); err != nil {
		tst.Fatal(err)
	}
	if !strings.Contains(status.String(), "Audio rendered: 100%") {
		tst.Fatal(status.String())
	}
}
