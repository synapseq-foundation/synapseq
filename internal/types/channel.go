// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package types

const NumberOfChannels = 16

// Channel represents a channel state
type Channel struct {
	Track         Track
	WaveformStart WaveformType
	WaveformEnd   WaveformType
	WaveformAlpha float64
	Type          TrackType
	Amplitude     [2]int
	Increment     [2]int
	Offset        [2]int
	Effect        EffectState
}
