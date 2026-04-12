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

package ambiance

import (
	"fmt"
	"sort"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

const stereoChannels = 2

type SampleAudio interface {
	ReadSamplesAt(index int, samples []int, numSamples int) (int, error)
	Close() error
}

type Runtime struct {
	audio        SampleAudio
	samplesByIdx [][]int
	activeIdx    []int
	activeMask   []bool
	channelIdx   [t.NumberOfChannels]int
	periodStart  [][]int
}

func NewRuntime(periods []t.Period, ambiance map[string]string, sampleRate int, newAudio func(paths []string, sampleRate int) (SampleAudio, error)) (*Runtime, error) {
	paths, nameToIndex, err := BuildIndex(ambiance)
	if err != nil {
		return nil, err
	}

	periodStart, err := PrecomputePeriodStart(periods, nameToIndex)
	if err != nil {
		return nil, err
	}

	var audio SampleAudio
	if newAudio != nil {
		audio, err = newAudio(paths, sampleRate)
		if err != nil {
			return nil, err
		}
	}

	runtime := &Runtime{
		audio:        audio,
		samplesByIdx: make([][]int, len(paths)),
		activeIdx:    make([]int, 0, t.NumberOfChannels),
		activeMask:   make([]bool, len(paths)),
		periodStart:  periodStart,
	}

	for i := range runtime.channelIdx {
		runtime.channelIdx[i] = -1
	}

	return runtime, nil
}

func NewTestRuntime(sampleCount int) *Runtime {
	runtime := &Runtime{
		samplesByIdx: make([][]int, sampleCount),
		activeIdx:    make([]int, 0, t.NumberOfChannels),
		activeMask:   make([]bool, sampleCount),
	}

	for i := range runtime.channelIdx {
		runtime.channelIdx[i] = -1
	}

	return runtime
}

func BuildIndex(ambiance map[string]string) ([]string, map[string]int, error) {
	if len(ambiance) == 0 {
		return nil, map[string]int{}, nil
	}

	names := make([]string, 0, len(ambiance))
	for name := range ambiance {
		names = append(names, name)
	}
	sort.Strings(names)

	paths := make([]string, len(names))
	nameToIndex := make(map[string]int, len(names))

	for i, name := range names {
		path := ambiance[name]
		if path == "" {
			return nil, nil, fmt.Errorf("ambiance %q has empty path", name)
		}
		paths[i] = path
		nameToIndex[name] = i
	}

	return paths, nameToIndex, nil
}

func PrecomputePeriodStart(periods []t.Period, nameToIndex map[string]int) ([][]int, error) {
	out := make([][]int, len(periods))
	for pIdx := range periods {
		row := make([]int, t.NumberOfChannels)
		for ch := range t.NumberOfChannels {
			row[ch] = -1
			if ch >= len(periods[pIdx].TrackStart) {
				continue
			}
			tr := periods[pIdx].TrackStart[ch]
			if tr.Type != t.TrackAmbiance {
				continue
			}
			idx, ok := nameToIndex[tr.AmbianceName]
			if !ok {
				return nil, fmt.Errorf("unknown ambiance name %q (period %d, channel %d)", tr.AmbianceName, pIdx, ch)
			}
			row[ch] = idx
		}
		out[pIdx] = row
	}
	return out, nil
}

func (ar *Runtime) Close() error {
	if ar == nil || ar.audio == nil {
		return nil
	}

	return ar.audio.Close()
}

func (ar *Runtime) UpdateChannelIndex(ch int, periodIdx int, trackType t.TrackType) {
	if ar == nil {
		return
	}

	if trackType == t.TrackAmbiance {
		ar.channelIdx[ch] = ar.periodStart[periodIdx][ch]
		return
	}

	ar.channelIdx[ch] = -1
}

func (ar *Runtime) CollectActiveIndices(channels []t.Channel) {
	if ar == nil {
		return
	}

	for _, idx := range ar.activeIdx {
		ar.activeMask[idx] = false
	}
	ar.activeIdx = ar.activeIdx[:0]

	for ch := range channels {
		if channels[ch].Track.Type != t.TrackAmbiance {
			continue
		}

		idx := ar.channelIdx[ch]
		if idx < 0 || idx >= len(ar.samplesByIdx) {
			continue
		}

		if !ar.activeMask[idx] {
			ar.activeMask[idx] = true
			ar.activeIdx = append(ar.activeIdx, idx)
		}
	}
}

func (ar *Runtime) PrepareBuffers(bufferSize int) {
	if ar == nil {
		return
	}

	need := bufferSize * stereoChannels
	for _, idx := range ar.activeIdx {
		buf := ar.samplesByIdx[idx]
		if len(buf) != need {
			buf = make([]int, need)
			ar.samplesByIdx[idx] = buf
		}

		if ar.audio == nil {
			zeroSamples(buf)
			continue
		}

		if _, err := ar.audio.ReadSamplesAt(idx, buf, need); err != nil {
			zeroSamples(buf)
		}
	}
}

func (ar *Runtime) ChannelBuffer(ch int) []int {
	if ar == nil {
		return nil
	}

	idx := ar.channelIdx[ch]
	if idx < 0 || idx >= len(ar.samplesByIdx) {
		return nil
	}

	return ar.samplesByIdx[idx]
}

func (ar *Runtime) SetChannelBuffer(idx int, samples []int) {
	if ar == nil || idx < 0 || idx >= len(ar.samplesByIdx) {
		return
	}

	ar.samplesByIdx[idx] = append([]int(nil), samples...)
}

func (ar *Runtime) SetChannelIndex(ch int, idx int) {
	if ar == nil || ch < 0 || ch >= len(ar.channelIdx) {
		return
	}

	ar.channelIdx[ch] = idx
}

func zeroSamples(samples []int) {
	for i := range samples {
		samples[i] = 0
	}
}