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

// RGBColor stores a terminal-friendly RGB token derived from the SynapSeq design system.
// The value is packed as 0xRRGGBB so it can be declared as a constant.
type RGBColor uint32

func (c RGBColor) R() int {
	return int((c >> 16) & 0xff)
}

func (c RGBColor) G() int {
	return int((c >> 8) & 0xff)
}

func (c RGBColor) B() int {
	return int(c & 0xff)
}

const (
	// CLI and terminal tokens approximate the warm SynapSeq design system palette.
	Terracotta     RGBColor = 0xb14d2a
	TerracottaDark RGBColor = 0x7f2d18
	Ochre          RGBColor = 0xa07126
	Green          RGBColor = 0x2f6b45
	MutedWarm      RGBColor = 0x6b6259
	DangerRed      RGBColor = 0x8b2e2e
)
