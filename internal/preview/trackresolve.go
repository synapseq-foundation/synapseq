// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preview

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

func resolveNodeTrack(periods []t.Period, periodIndex int, channel int) (t.Track, bool) {
	track, ok := resolveGraphTrack(periods, periodIndex, channel)
	if ok {
		return track, true
	}

	track = periods[periodIndex].TrackStart[channel]
	if includeVisibleTrack(track) {
		return track, true
	}

	return t.Track{}, false
}

func resolveGraphTrack(periods []t.Period, periodIndex int, channel int) (t.Track, bool) {
	track := periods[periodIndex].TrackStart[channel]
	if includeVisibleTrack(track) {
		return track, true
	}

	if track.Type != t.TrackSilence {
		return t.Track{}, false
	}

	for idx := periodIndex - 1; idx >= 0; idx-- {
		previous := periods[idx].TrackStart[channel]
		if !includeVisibleTrack(previous) {
			continue
		}

		resolved := previous
		resolved.Carrier = track.Carrier
		resolved.Resonance = track.Resonance
		resolved.NoiseSmooth = track.NoiseSmooth
		resolved.Amplitude = track.Amplitude
		resolved.Effect = track.Effect
		return resolved, true
	}

	return t.Track{}, false
}

func preferredTrack(startTrack, endTrack t.Track) t.Track {
	if includeVisibleTrack(startTrack) {
		return startTrack
	}
	if includeVisibleTrack(endTrack) {
		return endTrack
	}
	return endTrack
}

func primaryTrackType(startTrack, endTrack t.Track) t.TrackType {
	return preferredTrack(startTrack, endTrack).Type
}