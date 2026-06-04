// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sources

import (
	"testing"

	amb "github.com/synapseq-foundation/synapseq/v4/internal/audio/ambiance"
)

func TestAmbianceSampleUsesPreparedStereoBuffer(ts *testing.T) {
	runtime := amb.NewTestRuntime(1)
	runtime.SetChannelBuffer(0, []int{20000, -10000})
	runtime.SetChannelIndex(0, 0)
	source := NewAmbiance(Signal{Amplitude: [2]int{3, 0}})

	left, right, ok := source.Sample(runtime, 0, 0)
	if !ok {
		ts.Fatalf("expected ambiance sample to be available")
	}

	wantLeft := 20000 * 16 * 3
	wantRight := -10000 * 16 * 3
	if left != wantLeft || right != wantRight {
		ts.Fatalf("unexpected ambiance sample: got [%d %d], want [%d %d]", left, right, wantLeft, wantRight)
	}
}

func TestAmbianceSampleReturnsFalseForUnavailableData(ts *testing.T) {
	source := NewAmbiance(Signal{Amplitude: [2]int{3, 0}})

	left, right, ok := source.Sample(nil, 0, 0)
	if ok || left != 0 || right != 0 {
		ts.Fatalf("expected nil runtime to produce no sample, got [%d %d] ok=%t", left, right, ok)
	}

	runtime := amb.NewTestRuntime(1)
	runtime.SetChannelBuffer(0, []int{123})
	left, right, ok = source.Sample(runtime, 0, 0)
	if ok || left != 0 || right != 0 {
		ts.Fatalf("expected short buffer to produce no sample, got [%d %d] ok=%t", left, right, ok)
	}
}