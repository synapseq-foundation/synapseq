/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	"fmt"
	"sort"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// Render generates the audio and passes buffers to the consume function
func buildBackgroundIndex(backgrounds map[string]string) ([]string, map[string]int, error) {
	if len(backgrounds) == 0 {
		return nil, map[string]int{}, nil
	}

	names := make([]string, 0, len(backgrounds))
	for name := range backgrounds {
		names = append(names, name)
	}
	sort.Strings(names)

	paths := make([]string, len(names))
	nameToIndex := make(map[string]int, len(names))

	for i, name := range names {
		path := backgrounds[name]
		if path == "" {
			return nil, nil, fmt.Errorf("background %q has empty path", name)
		}
		paths[i] = path
		nameToIndex[name] = i
	}

	return paths, nameToIndex, nil
}

// precomputePeriodBGStart precomputes the background audio indices for each period and channel to optimize lookups during rendering
func precomputePeriodBGStart(periods []t.Period, nameToIndex map[string]int) ([][]int, error) {
	out := make([][]int, len(periods))
	for pIdx := range periods {
		row := make([]int, t.NumberOfChannels)
		for ch := range t.NumberOfChannels {
			row[ch] = -1
			if ch >= len(periods[pIdx].TrackStart) {
				continue
			}
			tr := periods[pIdx].TrackStart[ch]
			if tr.Type != t.TrackBackground {
				continue
			}
			idx, ok := nameToIndex[tr.BackgroundName]
			if !ok {
				return nil, fmt.Errorf("unknown background name %q (period %d, channel %d)", tr.BackgroundName, pIdx, ch)
			}
			row[ch] = idx
		}
		out[pIdx] = row
	}
	return out, nil
}

// collectActiveBackgroundIndices identifies which background audio tracks are currently active based on the channels' configurations
func (r *AudioRenderer) collectActiveBackgroundIndices() {
	for _, idx := range r.activeBGIndices {
		r.activeBGMask[idx] = false
	}
	r.activeBGIndices = r.activeBGIndices[:0]

	for ch := range t.NumberOfChannels {
		c := &r.channels[ch]
		if c.Track.Type != t.TrackBackground {
			continue
		}
		idx := r.channelBGIndex[ch]
		if idx < 0 || idx >= len(r.backgroundSamplesByIndex) {
			continue
		}
		if !r.activeBGMask[idx] {
			r.activeBGMask[idx] = true
			r.activeBGIndices = append(r.activeBGIndices, idx)
		}
	}
}

// prepareBackgroundBuffers ensures that the background audio buffers are ready for mixing, loading data from the background audio source as needed
func (r *AudioRenderer) prepareBackgroundBuffers() {
	need := t.BufferSize * audioChannels

	for _, idx := range r.activeBGIndices {
		buf := r.backgroundSamplesByIndex[idx]
		if len(buf) != need {
			buf = make([]int, need)
			r.backgroundSamplesByIndex[idx] = buf
		}

		if r.backgroundAudio == nil || !r.backgroundAudio.IsEnabled() {
			for i := range buf {
				buf[i] = 0
			}
			continue
		}

		if _, err := r.backgroundAudio.ReadSamplesAt(idx, buf, need); err != nil {
			for i := range buf {
				buf[i] = 0
			}
		}
	}
}
