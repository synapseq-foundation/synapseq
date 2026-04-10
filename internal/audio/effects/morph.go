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

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

type WaveformMorph struct {
	Start t.WaveformType
	End   t.WaveformType
	Alpha float64
}

func WaveformMorphFromChannel(channel *t.Channel) WaveformMorph {
	if channel.WaveformStart == 0 && channel.WaveformEnd == 0 && channel.WaveformAlpha == 0 {
		return WaveformMorph{Start: channel.Track.Waveform, End: channel.Track.Waveform, Alpha: 0}
	}

	return WaveformMorph{Start: channel.WaveformStart, End: channel.WaveformEnd, Alpha: channel.WaveformAlpha}
}

func normalizedWaveformMorph(waveform WaveformMorph) (t.WaveformType, t.WaveformType, float64) {
	if waveform.Start == 0 && waveform.End == 0 && waveform.Alpha == 0 {
		return t.WaveformSine, t.WaveformSine, 0
	}

	return waveform.Start, waveform.End, waveform.Alpha
}