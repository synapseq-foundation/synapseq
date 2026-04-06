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
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise, t.TrackAmbiance:
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
	return isToneTrack(track) || track.Type == t.TrackSilence || track.Type == t.TrackAmbiance
}

func usesBeat(track t.Track) bool {
	switch track.Type {
	case t.TrackBinauralBeat, t.TrackMonauralBeat, t.TrackIsochronicBeat:
		return true
	default:
		return false
	}
}
