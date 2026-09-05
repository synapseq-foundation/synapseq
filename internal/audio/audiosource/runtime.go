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

const dopplerReadFrames = 256

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

type bufferedCursorSampleAudio interface {
	RewindBufferedFrames(index, frames int) error
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
	dopplerAudio [t.NumberOfChannels]SampleAudio
	dopplerBuf   [t.NumberOfChannels][]int
	dopplerHead  [t.NumberOfChannels]int
	dopplerCount [t.NumberOfChannels]int
	dopplerPos   [t.NumberOfChannels]float64
	dopplerMask  [t.NumberOfChannels]bool
	periodStart  [][]int
	err          error
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
	for ch := range ar.dopplerAudio {
		ar.closeDopplerAudio(ch)
	}
	return firstErr
}

// Err returns the first external audio failure encountered during rendering.
func (ar *Runtime) Err() error {
	if ar == nil {
		return nil
	}

	return ar.err
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
	if changedSource {
		ar.ResetDoppler(ch)
	}
	if ar.scope == BufferScopeChannel && changedSource {
		if nextIdx >= 0 {
			if audio, ok := ar.channelAudio[ch].(selectableSampleAudio); ok {
				if err := audio.SelectAudio(nextIdx); err == nil {
					ar.channelIdx[ch] = nextIdx
					return
				} else {
					ar.recordError(err)
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
			ar.dopplerMask[ch] = channels[ch].Track.Effect.Type == t.EffectDoppler
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
		ar.dopplerMask[ch] = channels[ch].Track.Effect.Type == t.EffectDoppler
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
		if !ar.sourceNeedsFixedBuffer(idx) {
			continue
		}
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
			ar.recordError(err)
			zeroSamples(buf)
		}
	}
}

func (ar *Runtime) sourceNeedsFixedBuffer(idx int) bool {
	for ch, channelIdx := range ar.channelIdx {
		if channelIdx == idx && !ar.dopplerMask[ch] {
			return true
		}
	}
	return false
}

func (ar *Runtime) prepareChannelBuffers(bufferSize int) {
	need := bufferSize * stereoChannels
	for _, ch := range ar.activeCh {
		if ar.dopplerMask[ch] {
			continue
		}
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
			ar.recordError(err)
			zeroSamples(buf)
			continue
		}

		if _, err := audio.ReadSamplesAt(idx, buf, need); err != nil {
			ar.recordError(err)
			zeroSamples(buf)
		}
	}
}

// SampleDoppler reads one stereo frame at a fractional playback rate. Its
// cursor and buffered source data are isolated per sequencer channel.
func (ar *Runtime) SampleDoppler(ch int, rate float64) (int, int, bool) {
	if ar == nil || ch < 0 || ch >= len(ar.channelIdx) || ar.channelIdx[ch] < 0 {
		return 0, 0, false
	}
	if rate <= 0 {
		rate = 1
	}

	if !ar.ensureDopplerFrames(ch, int(ar.dopplerPos[ch])+2) {
		return 0, 0, false
	}

	position := ar.dopplerPos[ch]
	frame := int(position)
	fraction := position - float64(frame)
	left := interpolateSample(
		ar.dopplerFrameSample(ch, frame, 0),
		ar.dopplerFrameSample(ch, frame+1, 0),
		fraction,
	)
	right := interpolateSample(
		ar.dopplerFrameSample(ch, frame, 1),
		ar.dopplerFrameSample(ch, frame+1, 1),
		fraction,
	)

	next := position + rate
	consumed := int(next)
	ar.dopplerPos[ch] = next - float64(consumed)
	if consumed > 0 {
		ar.consumeDopplerFrames(ch, consumed)
	}
	return left, right, true
}

// ResetDoppler discards the independent PCM cursor for one sequencer channel.
func (ar *Runtime) ResetDoppler(ch int) {
	if ar == nil || ch < 0 || ch >= len(ar.dopplerBuf) {
		return
	}
	if audio, ok := ar.dopplerPlaybackAudio(ch).(bufferedCursorSampleAudio); ok {
		ar.recordError(audio.RewindBufferedFrames(ar.channelIdx[ch], ar.dopplerCount[ch]))
	}
	ar.closeDopplerAudio(ch)
	ar.dopplerBuf[ch] = nil
	ar.dopplerHead[ch] = 0
	ar.dopplerCount[ch] = 0
	ar.dopplerPos[ch] = 0
}

func (ar *Runtime) ensureDopplerFrames(ch, needed int) bool {
	if needed <= ar.dopplerCount[ch] {
		return true
	}
	if ar.dopplerBuf[ch] == nil {
		ar.dopplerBuf[ch] = make([]int, dopplerReadFrames*2*stereoChannels)
	}
	audio, err := ar.dopplerAudioForChannel(ch)
	if err != nil || audio == nil {
		ar.recordError(err)
		return false
	}

	for ar.dopplerCount[ch] < needed {
		capacity := ar.dopplerFrameCapacity(ch) - ar.dopplerCount[ch]
		if capacity == 0 {
			return false
		}
		frames := min(dopplerReadFrames, capacity)
		if !ar.readDopplerFrames(ch, audio, frames) {
			return false
		}
	}
	return true
}

func (ar *Runtime) consumeDopplerFrames(ch, count int) {
	if count >= ar.dopplerCount[ch] {
		ar.dopplerHead[ch] = 0
		ar.dopplerCount[ch] = 0
		return
	}
	ar.dopplerHead[ch] = (ar.dopplerHead[ch] + count) % ar.dopplerFrameCapacity(ch)
	ar.dopplerCount[ch] -= count
}

func (ar *Runtime) readDopplerFrames(ch int, audio SampleAudio, frames int) bool {
	capacity := ar.dopplerFrameCapacity(ch)
	tail := (ar.dopplerHead[ch] + ar.dopplerCount[ch]) % capacity
	firstFrames := min(frames, capacity-tail)
	if !ar.readDopplerFrameRange(ch, audio, tail, firstFrames) {
		return false
	}
	remaining := frames - firstFrames
	if remaining > 0 && !ar.readDopplerFrameRange(ch, audio, 0, remaining) {
		return false
	}
	ar.dopplerCount[ch] += frames
	return true
}

func (ar *Runtime) readDopplerFrameRange(ch int, audio SampleAudio, startFrame, frames int) bool {
	start := startFrame * stereoChannels
	end := start + frames*stereoChannels
	_, err := audio.ReadSamplesAt(ar.channelIdx[ch], ar.dopplerBuf[ch][start:end], end-start)
	if err != nil {
		ar.recordError(err)
		return false
	}
	return true
}

func (ar *Runtime) dopplerFrameCapacity(ch int) int {
	return len(ar.dopplerBuf[ch]) / stereoChannels
}

func (ar *Runtime) dopplerFrameSample(ch, offset, channel int) int {
	frame := (ar.dopplerHead[ch] + offset) % ar.dopplerFrameCapacity(ch)
	return ar.dopplerBuf[ch][frame*stereoChannels+channel]
}

func (ar *Runtime) dopplerAudioForChannel(ch int) (SampleAudio, error) {
	if ar.scope == BufferScopeChannel {
		return ar.audioForChannel(ch)
	}
	if ar.dopplerAudio[ch] != nil {
		return ar.dopplerAudio[ch], nil
	}
	if template, ok := ar.audio.(cloneableSampleAudio); ok {
		audio, err := template.CloneAudio()
		if err != nil {
			return nil, err
		}
		if selectable, ok := audio.(selectableSampleAudio); ok {
			if err := selectable.SelectAudio(ar.channelIdx[ch]); err != nil {
				_ = audio.Close()
				return nil, err
			}
		}
		ar.dopplerAudio[ch] = audio
		return audio, nil
	}
	if ar.newAudio == nil {
		return nil, nil
	}
	audio, err := ar.newAudio(ar.paths, ar.sampleRate)
	if err != nil {
		return nil, err
	}
	if selectable, ok := audio.(selectableSampleAudio); ok {
		if err := selectable.SelectAudio(ar.channelIdx[ch]); err != nil {
			_ = audio.Close()
			return nil, err
		}
	}
	ar.dopplerAudio[ch] = audio
	return audio, nil
}

func (ar *Runtime) dopplerPlaybackAudio(ch int) SampleAudio {
	if ar.scope == BufferScopeChannel {
		return ar.channelAudio[ch]
	}
	return ar.dopplerAudio[ch]
}

func (ar *Runtime) closeDopplerAudio(ch int) {
	if ar.dopplerAudio[ch] == nil {
		return
	}
	_ = ar.dopplerAudio[ch].Close()
	ar.dopplerAudio[ch] = nil
}

func interpolateSample(start, end int, fraction float64) int {
	return start + int(float64(end-start)*fraction)
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

func (ar *Runtime) recordError(err error) {
	if ar.err == nil && err != nil {
		ar.err = err
	}
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
