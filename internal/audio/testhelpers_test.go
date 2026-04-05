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

package audio

type constStreamer struct {
	framesLeft int
	val        float64
}

func (cs *constStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n = len(samples)
	if n > cs.framesLeft {
		n = cs.framesLeft
	}
	for i := 0; i < n; i++ {
		samples[i][0] = cs.val
		samples[i][1] = cs.val
	}
	cs.framesLeft -= n
	ok = cs.framesLeft > 0
	return
}

func (cs *constStreamer) Err() error { return nil }