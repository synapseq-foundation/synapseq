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
func buildAmbianceIndex(ambiance map[string]string) ([]string, map[string]int, error) {
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

// precomputePeriodAmbianceStart precomputes the ambiance audio indices for each period and channel to optimize lookups during rendering
func precomputePeriodAmbianceStart(periods []t.Period, nameToIndex map[string]int) ([][]int, error) {
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

// collectActiveAmbianceIndices identifies which ambiance audio tracks are currently active based on the channels' configurations
func (r *AudioRenderer) collectActiveAmbianceIndices() {
	for _, idx := range r.activeAmbianceIndices {
		r.activeAmbianceMask[idx] = false
	}
	r.activeAmbianceIndices = r.activeAmbianceIndices[:0]

	for ch := range t.NumberOfChannels {
		c := &r.channels[ch]
		if c.Track.Type != t.TrackAmbiance {
			continue
		}
		idx := r.channelAmbianceIndex[ch]
		if idx < 0 || idx >= len(r.ambianceSamplesByIndex) {
			continue
		}
		if !r.activeAmbianceMask[idx] {
			r.activeAmbianceMask[idx] = true
			r.activeAmbianceIndices = append(r.activeAmbianceIndices, idx)
		}
	}
}

// prepareAmbianceBuffers ensures that the ambiance audio buffers are ready for mixing, loading data from the ambiance audio source as needed
func (r *AudioRenderer) prepareAmbianceBuffers() {
	need := t.BufferSize * audioChannels

	for _, idx := range r.activeAmbianceIndices {
		buf := r.ambianceSamplesByIndex[idx]
		if len(buf) != need {
			buf = make([]int, need)
			r.ambianceSamplesByIndex[idx] = buf
		}

		if r.ambianceAudio == nil || !r.ambianceAudio.IsEnabled() {
			for i := range buf {
				buf[i] = 0
			}
			continue
		}

		if _, err := r.ambianceAudio.ReadSamplesAt(idx, buf, need); err != nil {
			for i := range buf {
				buf[i] = 0
			}
		}
	}
}
