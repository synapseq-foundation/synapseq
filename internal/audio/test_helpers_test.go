// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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