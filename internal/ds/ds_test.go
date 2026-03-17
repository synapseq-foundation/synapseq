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

package ds

import "testing"

func TestRGBColorChannels(ts *testing.T) {
	color := RGBColor(0x123456)

	if color.R() != 0x12 {
		ts.Fatalf("expected red channel 0x12, got %#x", color.R())
	}
	if color.G() != 0x34 {
		ts.Fatalf("expected green channel 0x34, got %#x", color.G())
	}
	if color.B() != 0x56 {
		ts.Fatalf("expected blue channel 0x56, got %#x", color.B())
	}
}

func TestSharedTerminalTokens(ts *testing.T) {
	tests := []struct {
		name  string
		color RGBColor
		r     int
		g     int
		b     int
	}{
		{name: "Terracotta", color: Terracotta, r: 177, g: 77, b: 42},
		{name: "TerracottaDark", color: TerracottaDark, r: 127, g: 45, b: 24},
		{name: "Ochre", color: Ochre, r: 160, g: 113, b: 38},
		{name: "Green", color: Green, r: 47, g: 107, b: 69},
		{name: "MutedWarm", color: MutedWarm, r: 107, g: 98, b: 89},
		{name: "DangerRed", color: DangerRed, r: 139, g: 46, b: 46},
	}

	for _, test := range tests {
		if test.color.R() != test.r || test.color.G() != test.g || test.color.B() != test.b {
			ts.Fatalf(
				"%s: expected RGB (%d, %d, %d), got (%d, %d, %d)",
				test.name,
				test.r,
				test.g,
				test.b,
				test.color.R(),
				test.color.G(),
				test.color.B(),
			)
		}
	}
}
