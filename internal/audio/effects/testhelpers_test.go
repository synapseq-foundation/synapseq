// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"

func newTestProcessor() *Processor {
	return NewProcessor(44100, wt.Init())
}
