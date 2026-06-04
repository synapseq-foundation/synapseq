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
	newAudio     func(paths []string, sampleRate int) (SampleAudio, error)
	paths        []string
	sampleRate   int
	sourceKind   string
	trackType    t.TrackType
	perChannel   bool
	samplesByIdx [][]int
	samplesByCh  [t.NumberOfChannels][]int
	activeIdx    []int
	activeMask   []bool
	activeCh     []int
	activeChMask [t.NumberOfChannels]bool
	channelIdx   [t.NumberOfChannels]int
	channelAudio [t.NumberOfChannels]SampleAudio
	periodStart  [][]int
}

func NewRuntime(periods []t.Period, ambiance map[string]string, sampleRate int, newAudio func(paths []string, sampleRate int) (SampleAudio, error)) (*Runtime, error) {
	return newRuntime(periods, ambiance, sampleRate, t.TrackAmbiance, "ambiance", false, newAudio)
}

func NewMusicRuntime(periods []t.Period, music map[string]string, sampleRate int, newAudio func(paths []string, sampleRate int) (SampleAudio, error)) (*Runtime, error) {
	return newRuntime(periods, music, sampleRate, t.TrackMusic, "music", true, newAudio)
}

func newRuntime(periods []t.Period, sources map[string]string, sampleRate int, trackType t.TrackType, sourceKind string, perChannel bool, newAudio func(paths []string, sampleRate int) (SampleAudio, error)) (*Runtime, error) {
	paths, nameToIndex, err := BuildIndex(sources, sourceKind)
	if err != nil {
		return nil, err
	}

	periodStart, err := PrecomputePeriodStart(periods, nameToIndex, trackType, sourceKind)
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
		newAudio:     newAudio,
		paths:        paths,
		sampleRate:   sampleRate,
		sourceKind:   sourceKind,
		trackType:    trackType,
		perChannel:   perChannel,
		samplesByIdx: make([][]int, len(paths)),
		activeIdx:    make([]int, 0, t.NumberOfChannels),
		activeMask:   make([]bool, len(paths)),
		activeCh:     make([]int, 0, t.NumberOfChannels),
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

func BuildIndex(sources map[string]string, sourceKind string) ([]string, map[string]int, error) {
	if len(sources) == 0 {
		return nil, map[string]int{}, nil
	}

	names := make([]string, 0, len(sources))
	for name := range sources {
		names = append(names, name)
	}
	sort.Strings(names)

	paths := make([]string, len(names))
	nameToIndex := make(map[string]int, len(names))

	for i, name := range names {
		path := sources[name]
		if path == "" {
			return nil, nil, fmt.Errorf("%s %q has empty path", sourceKind, name)
		}
		paths[i] = path
		nameToIndex[name] = i
	}

	return paths, nameToIndex, nil
}

func PrecomputePeriodStart(periods []t.Period, nameToIndex map[string]int, trackType t.TrackType, sourceKind string) ([][]int, error) {
	out := make([][]int, len(periods))
	for pIdx := range periods {
		row := make([]int, t.NumberOfChannels)
		for ch := range t.NumberOfChannels {
			row[ch] = -1
			if ch >= len(periods[pIdx].TrackStart) {
				continue
			}
			tr := periods[pIdx].TrackStart[ch]
			if tr.Type != trackType {
				continue
			}
			idx, ok := nameToIndex[tr.SourceName]
			if !ok {
				return nil, fmt.Errorf("unknown %s name %q (period %d, channel %d)", sourceKind, tr.SourceName, pIdx, ch)
			}
			row[ch] = idx
		}
		out[pIdx] = row
	}
	return out, nil
}

func (ar *Runtime) Close() error {
	if ar == nil {
		return nil
	}

	var firstErr error
	if ar.audio != nil {
		if err := ar.audio.Close(); err != nil {
			firstErr = err
		}
		ar.audio = nil
	}
	for ch := range ar.channelAudio {
		if ar.channelAudio[ch] != nil {
			if err := ar.channelAudio[ch].Close(); err != nil && firstErr == nil {
				firstErr = err
			}
			ar.channelAudio[ch] = nil
		}
	}
	return firstErr
}

func (ar *Runtime) UpdateChannelIndex(ch int, periodIdx int, trackType t.TrackType) {
	if ar == nil || ch < 0 || ch >= len(ar.channelIdx) {
		return
	}

	nextIdx := -1
	if trackType == ar.trackType {
		nextIdx = ar.periodStart[periodIdx][ch]
	}

	if ar.perChannel && nextIdx != ar.channelIdx[ch] {
		ar.closeChannelAudio(ch)
	}
	ar.channelIdx[ch] = nextIdx
}

func (ar *Runtime) CollectActiveIndices(channels []t.Channel) {
	if ar == nil {
		return
	}

	if ar.perChannel {
		for _, ch := range ar.activeCh {
			ar.activeChMask[ch] = false
		}
		ar.activeCh = ar.activeCh[:0]

		for ch := range channels {
			if ch >= len(ar.channelIdx) || channels[ch].Track.Type != ar.trackType {
				continue
			}
			idx := ar.channelIdx[ch]
			if idx < 0 || idx >= len(ar.paths) {
				continue
			}
			if !ar.activeChMask[ch] {
				ar.activeChMask[ch] = true
				ar.activeCh = append(ar.activeCh, ch)
			}
		}
		return
	}

	for _, idx := range ar.activeIdx {
		ar.activeMask[idx] = false
	}
	ar.activeIdx = ar.activeIdx[:0]

	for ch := range channels {
		if channels[ch].Track.Type != ar.trackType {
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

	if ar.perChannel {
		ar.prepareChannelBuffers(bufferSize)
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

func (ar *Runtime) prepareChannelBuffers(bufferSize int) {
	need := bufferSize * stereoChannels
	for _, ch := range ar.activeCh {
		idx := ar.channelIdx[ch]
		if idx < 0 || idx >= len(ar.paths) {
			continue
		}

		buf := ar.samplesByCh[ch]
		if len(buf) != need {
			buf = make([]int, need)
			ar.samplesByCh[ch] = buf
		}

		audio, err := ar.audioForChannel(ch)
		if err != nil || audio == nil {
			zeroSamples(buf)
			continue
		}

		if _, err := audio.ReadSamplesAt(idx, buf, need); err != nil {
			zeroSamples(buf)
		}
	}
}

func (ar *Runtime) audioForChannel(ch int) (SampleAudio, error) {
	if ch < 0 || ch >= len(ar.channelAudio) {
		return nil, fmt.Errorf("invalid channel index: %d", ch)
	}
	if ar.channelAudio[ch] != nil {
		return ar.channelAudio[ch], nil
	}
	if ar.newAudio == nil {
		return nil, nil
	}

	audio, err := ar.newAudio(ar.paths, ar.sampleRate)
	if err != nil {
		return nil, err
	}
	ar.channelAudio[ch] = audio
	return audio, nil
}

func (ar *Runtime) closeChannelAudio(ch int) {
	if ch < 0 || ch >= len(ar.channelAudio) || ar.channelAudio[ch] == nil {
		return
	}
	_ = ar.channelAudio[ch].Close()
	ar.channelAudio[ch] = nil
	ar.samplesByCh[ch] = nil
}

func (ar *Runtime) ChannelBuffer(ch int) []int {
	if ar == nil {
		return nil
	}

	if ar.perChannel {
		if ch < 0 || ch >= len(ar.samplesByCh) {
			return nil
		}
		return ar.samplesByCh[ch]
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
