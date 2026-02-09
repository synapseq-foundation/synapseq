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

package shared

import (
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// Equal checks if two track configurations are identical
func IsTrackEqual(tr1, tr2 *t.Track) bool {
	return tr1.Type == tr2.Type &&
		tr1.Amplitude == tr2.Amplitude &&
		tr1.Carrier == tr2.Carrier &&
		tr1.Resonance == tr2.Resonance &&
		tr1.Waveform == tr2.Waveform &&
		tr1.Effect.Intensity == tr2.Effect.Intensity
}
