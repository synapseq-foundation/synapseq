// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package output

import p "github.com/synapseq-foundation/synapseq/v4/internal/audio/pcm"

type RenderFunc func(consume func(samples []int) error) error

type RendererStreamer struct {
	ch      chan []int
	free    chan []int
	current []int
	offset  int
	done    bool
	err     error
}

func NewRendererStreamer(render RenderFunc) *RendererStreamer {
	rs := &RendererStreamer{
		ch:   make(chan []int, 2),
		free: make(chan []int, 3),
	}
	for range cap(rs.free) {
		rs.free <- nil
	}
	go func() {
		defer close(rs.ch)
		rs.err = render(func(samples []int) error {
			buffer := <-rs.free
			if cap(buffer) < len(samples) {
				buffer = make([]int, 0, len(samples))
			}
			buffer = append(buffer[:0], samples...)
			rs.ch <- buffer
			return nil
		})
	}()
	return rs
}

func (rs *RendererStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if rs.done && rs.current == nil {
		return 0, false
	}

	for n < len(samples) {
		if rs.current == nil {
			data, more := <-rs.ch
			if !more {
				rs.done = true
				break
			}
			rs.current = data
			rs.offset = 0
		}
		framesAvail := (len(rs.current) - rs.offset) / 2
		if framesAvail == 0 {
			rs.releaseCurrent()
			continue
		}
		need := len(samples) - n
		if framesAvail < need {
			need = framesAvail
		}
		for i := 0; i < need; i++ {
			l := rs.current[rs.offset+2*i]
			r := rs.current[rs.offset+2*i+1]
			samples[n+i][0] = p.SampleToUnitFloat64(l)
			samples[n+i][1] = p.SampleToUnitFloat64(r)
		}
		rs.offset += need * 2
		if rs.offset == len(rs.current) {
			rs.releaseCurrent()
		}
		n += need
	}
	ok = !rs.done || rs.current != nil || n > 0
	return
}

func (rs *RendererStreamer) releaseCurrent() {
	select {
	case rs.free <- rs.current[:0]:
	default:
	}
	rs.current = nil
	rs.offset = 0
}

func (rs *RendererStreamer) Err() error {
	return rs.err
}
