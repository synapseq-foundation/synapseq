// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package effects

import (
	wt "github.com/synapseq-foundation/synapseq/v4/internal/audio/wavetable"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type WaveformMorph struct {
	Start wt.ID
	End   wt.ID
	Alpha float64
}

func WaveformMorphFromChannel(channel *t.Channel) WaveformMorph {
	if channel.WaveformStart == 0 && channel.WaveformEnd == 0 && channel.WaveformAlpha == 0 {
		waveform, ok := wt.BuiltinID(channel.Track.Waveform)
		if !ok {
			waveform = wt.SineID
		}
		return WaveformMorph{Start: waveform, End: waveform, Alpha: 0}
	}

	return WaveformMorph{Start: wt.ID(channel.WaveformStart), End: wt.ID(channel.WaveformEnd), Alpha: channel.WaveformAlpha}
}

func normalizedWaveformMorph(waveform WaveformMorph) (wt.ID, wt.ID, float64) {
	if waveform.Start == 0 && waveform.End == 0 && waveform.Alpha == 0 {
		return wt.SineID, wt.SineID, 0
	}

	return waveform.Start, waveform.End, waveform.Alpha
}
