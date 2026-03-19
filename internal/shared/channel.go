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

package shared

import (
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// CountActiveChannels counts the number of active channels
func CountActiveChannels(chs []t.Channel) int {
	for i := len(chs) - 1; i >= 0; i-- {
		if chs[i].Track.Type != t.TrackOff {
			return i + 1
		}
	}
	return 1 // At least 1 channel always
}
