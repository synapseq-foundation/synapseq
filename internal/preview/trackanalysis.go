// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preview

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

func includeTrack(track t.Track) bool {
	return track.Type != t.TrackOff
}

func includeVisibleTrack(track t.Track) bool {
	return includeTrack(track) && track.Type != t.TrackSilence
}

func shouldRenderSegmentItem(startTrack, endTrack t.Track) bool {
	return includeVisibleTrack(startTrack) || includeVisibleTrack(endTrack)
}

func isToneTrack(track t.Track) bool {
	switch track.Type {
	case t.TrackPureTone, t.TrackBinauralBeat, t.TrackMonauralBeat, t.TrackIsochronicBeat:
		return true
	default:
		return false
	}
}

func isTextureTrack(track t.Track) bool {
	switch track.Type {
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance, t.TrackMusic:
		return true
	default:
		return false
	}
}

func isNoiseTrack(track t.Track) bool {
	switch track.Type {
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return true
	default:
		return false
	}
}

func supportsWaveform(track t.Track) bool {
	return isToneTrack(track) || track.Type == t.TrackSilence || track.Type == t.TrackAmbiance || track.Type == t.TrackMusic
}

func usesBeat(track t.Track) bool {
	switch track.Type {
	case t.TrackBinauralBeat, t.TrackMonauralBeat, t.TrackIsochronicBeat:
		return true
	default:
		return false
	}
}
