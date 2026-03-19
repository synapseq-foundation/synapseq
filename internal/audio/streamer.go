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

// rendererStreamer streams rendered audio samples
type rendererStreamer struct {
	ch       chan []int
	leftover []int
	done     bool
	err      error
}

func newRendererStreamer(r *AudioRenderer) *rendererStreamer {
	rs := &rendererStreamer{
		ch: make(chan []int, 2),
	}
	go func() {
		defer close(rs.ch)
		// Run the renderer and stream samples
		rs.err = r.Render(func(samples []int) error {
			cpy := make([]int, len(samples))
			copy(cpy, samples)
			rs.ch <- cpy
			return nil
		})
	}()
	return rs
}

// Stream streams audio samples in the range [-1.0, 1.0]
func (rs *rendererStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if rs.done && len(rs.leftover) == 0 {
		return 0, false
	}
	const denom = 32768.0 // 2^15
	for n < len(samples) {
		if len(rs.leftover) < 2 {
			data, more := <-rs.ch
			if !more {
				rs.done = true
				break
			}
			rs.leftover = data
		}
		framesAvail := len(rs.leftover) / 2
		if framesAvail == 0 {
			rs.leftover = nil
			continue
		}
		need := len(samples) - n
		if framesAvail < need {
			need = framesAvail
		}
		for i := 0; i < need; i++ {
			l := rs.leftover[2*i]
			r := rs.leftover[2*i+1]
			samples[n+i][0] = float64(l) / denom
			samples[n+i][1] = float64(r) / denom
		}
		rs.leftover = rs.leftover[need*2:]
		n += need
	}
	ok = !rs.done || len(rs.leftover) > 0 || n > 0
	return
}

// Err implements beep.Streamer.Err and propagates any errors that occurred during rendering
func (rs *rendererStreamer) Err() error {
	return rs.err
}
