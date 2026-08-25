// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package audiosource

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

type cloneableSampleAudio interface {
	CloneAudio() (SampleAudio, error)
}

type selectableSampleAudio interface {
	SelectAudio(index int) error
}

type NewAudioFunc func(paths []string, sampleRate int) (SampleAudio, error)

type BufferScope int

const (
	BufferScopeSource BufferScope = iota
	BufferScopeChannel
)

type RuntimeOptions struct {
	TrackType  t.TrackType
	SourceKind string
	Scope      BufferScope
	NewAudio   NewAudioFunc
}

type Runtime struct {
	audio        SampleAudio
	newAudio     NewAudioFunc
	paths        []string
	sampleRate   int
	sourceKind   string
	trackType    t.TrackType
	scope        BufferScope
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

func NewRuntime(periods []t.Period, sources map[string]string, sampleRate int, opts RuntimeOptions) (*Runtime, error) {
	sourceKind := opts.SourceKind
	if sourceKind == "" {
		sourceKind = "external audio"
	}

	reachable := ReachableSources(periods, sources, opts.TrackType)
	paths, nameToIndex, err := buildIndex(reachable, sourceKind, opts.Scope == BufferScopeChannel)
	if err != nil {
		return nil, err
	}

	periodStart, err := PrecomputePeriodStart(periods, nameToIndex, opts.TrackType, sourceKind)
	if err != nil {
		return nil, err
	}

	var audio SampleAudio
	if opts.NewAudio != nil {
		audio, err = opts.NewAudio(paths, sampleRate)
		if err != nil {
			return nil, err
		}
	}

	runtime := &Runtime{
		audio:        audio,
		newAudio:     opts.NewAudio,
		paths:        paths,
		sampleRate:   sampleRate,
		sourceKind:   sourceKind,
		trackType:    opts.TrackType,
		scope:        opts.Scope,
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
	return buildIndex(sources, sourceKind, false)
}

func buildIndex(sources map[string]string, sourceKind string, deduplicatePaths bool) ([]string, map[string]int, error) {
	if len(sources) == 0 {
		return nil, map[string]int{}, nil
	}

	names := make([]string, 0, len(sources))
	for name := range sources {
		names = append(names, name)
	}
	sort.Strings(names)

	paths := make([]string, 0, len(names))
	nameToIndex := make(map[string]int, len(names))
	pathToIndex := make(map[string]int, len(names))

	for _, name := range names {
		path := sources[name]
		if path == "" {
			return nil, nil, fmt.Errorf("%s %q has empty path", sourceKind, name)
		}
		if deduplicatePaths {
			if index, ok := pathToIndex[path]; ok {
				nameToIndex[name] = index
				continue
			}
		}
		nameToIndex[name] = len(paths)
		pathToIndex[path] = len(paths)
		paths = append(paths, path)
	}

	return paths, nameToIndex, nil
}

// ReachableSources returns declared sources referenced by periods or boundary crossfades.
func ReachableSources(periods []t.Period, sources map[string]string, trackType t.TrackType) map[string]string {
	reachable := make(map[string]string)
	for _, period := range periods {
		for channel := range t.NumberOfChannels {
			tracks := []t.Track{
				period.TrackStart[channel],
				period.TrackEnd[channel],
				period.CrossfadeOut[channel].Track,
				period.CrossfadeIn[channel].Track,
			}
			for _, track := range tracks {
				if track.Type != trackType {
					continue
				}
				if path, ok := sources[track.SourceName]; ok {
					reachable[track.SourceName] = path
				}
			}
		}
	}

	return reachable
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

	changedSource := nextIdx != ar.channelIdx[ch]
	if ar.scope == BufferScopeChannel && changedSource {
		if nextIdx >= 0 {
			if audio, ok := ar.channelAudio[ch].(selectableSampleAudio); ok {
				if err := audio.SelectAudio(nextIdx); err == nil {
					ar.channelIdx[ch] = nextIdx
					return
				}
			}
		}
		ar.closeChannelAudio(ch)
	}
	ar.channelIdx[ch] = nextIdx
}

func (ar *Runtime) CollectActiveIndices(channels []t.Channel) {
	if ar == nil {
		return
	}

	if ar.scope == BufferScopeChannel {
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
	ar.releaseInactiveSourceBuffers()
}

func (ar *Runtime) releaseInactiveSourceBuffers() {
	for index := range ar.samplesByIdx {
		if !ar.activeMask[index] {
			ar.samplesByIdx[index] = nil
		}
	}
}

func (ar *Runtime) PrepareBuffers(bufferSize int) {
	if ar == nil {
		return
	}

	if ar.scope == BufferScopeChannel {
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
	if template, ok := ar.audio.(cloneableSampleAudio); ok {
		audio, err := template.CloneAudio()
		if err != nil {
			return nil, err
		}
		ar.channelAudio[ch] = audio
		return audio, nil
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

	if ar.scope == BufferScopeChannel {
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
