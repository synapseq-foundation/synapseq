// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import "testing"

func TestAdvancePhaseWrapsAtMask(ts *testing.T) {
	got := advancePhase(phaseMask-5, 10)
	want := 4
	if got != want {
		ts.Fatalf("unexpected wrapped phase: got %d, want %d", got, want)
	}
}
