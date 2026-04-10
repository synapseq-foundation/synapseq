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

package effects

import "testing"

func TestAdvancePhaseWrapsAtMask(ts *testing.T) {
	got := advancePhase(phaseMask-5, 10)
	want := 4
	if got != want {
		ts.Fatalf("unexpected wrapped phase: got %d, want %d", got, want)
	}
}
