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

import (
	"fmt"
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func buildPrimarySummary(items []previewSegmentItemView) string {
	if len(items) == 0 {
		return "No active channels"
	}
	if len(items) == 1 {
		return fmt.Sprintf("%s %s", items[0].ChannelLabel, items[0].Label)
	}
	if len(items) == 2 {
		return fmt.Sprintf("%s %s, %s %s", items[0].ChannelLabel, items[0].Label, items[1].ChannelLabel, items[1].Label)
	}
	return fmt.Sprintf("%s %s, %s %s and %d more", items[0].ChannelLabel, items[0].Label, items[1].ChannelLabel, items[1].Label, len(items)-2)
}

func buildSegmentSummary(startTrack, endTrack t.Track) string {
	parts := make([]string, 0, 4)

	if startTrack.Carrier > 0 || endTrack.Carrier > 0 {
		parts = append(parts, fmt.Sprintf("carrier %s -> %s", formatHz(startTrack.Carrier), formatHz(endTrack.Carrier)))
	}
	if usesBeat(startTrack) || usesBeat(endTrack) || startTrack.Resonance > 0 || endTrack.Resonance > 0 {
		parts = append(parts, fmt.Sprintf("beat %s -> %s", formatHz(startTrack.Resonance), formatHz(endTrack.Resonance)))
	}
	if isNoiseTrack(startTrack) || isNoiseTrack(endTrack) {
		parts = append(parts, fmt.Sprintf("smooth %s -> %s", formatPercent(startTrack.NoiseSmooth), formatPercent(endTrack.NoiseSmooth)))
	}
	parts = append(parts, fmt.Sprintf("amp %s -> %s", formatPercent(startTrack.Amplitude.ToPercent()), formatPercent(endTrack.Amplitude.ToPercent())))

	if len(parts) == 0 {
		return "steady"
	}

	return joinParts(parts)
}

func buildTrackSummary(track t.Track) string {
	switch track.Type {
	case t.TrackOff:
		return "Channel disabled"
	case t.TrackSilence:
		if track.Carrier > 0 {
			return fmt.Sprintf("Fade state with %s carrier", formatHz(track.Carrier))
		}
		return "Silent boundary"
	case t.TrackPureTone:
		return fmt.Sprintf("Pure carrier at %s", formatHz(track.Carrier))
	case t.TrackBinauralBeat, t.TrackMonauralBeat, t.TrackIsochronicBeat:
		return fmt.Sprintf("Carrier %s with beat %s", formatHz(track.Carrier), formatHz(track.Resonance))
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return fmt.Sprintf("%s texture layer with %s smooth", humanTrackType(track), formatPercent(track.NoiseSmooth))
	case t.TrackAmbiance:
		if track.AmbianceName != "" {
			return fmt.Sprintf("Ambiance layer %q", track.AmbianceName)
		}
		return "Ambiance layer"
	default:
		return track.Type.String()
	}
}

func humanTrackType(track t.Track) string {
	switch track.Type {
	case t.TrackOff:
		return "Off"
	case t.TrackSilence:
		return "Silence"
	case t.TrackPureTone:
		return "Pure tone"
	case t.TrackBinauralBeat:
		return "Binaural beat"
	case t.TrackMonauralBeat:
		return "Monaural beat"
	case t.TrackIsochronicBeat:
		return "Isochronic beat"
	case t.TrackWhiteNoise:
		return "White noise"
	case t.TrackPinkNoise:
		return "Pink noise"
	case t.TrackBrownNoise:
		return "Brown noise"
	case t.TrackAmbiance:
		return "Ambiance"
	default:
		return "Unknown"
	}
}

func humanWaveformType(waveform t.WaveformType) string {
	switch waveform {
	case t.WaveformSine:
		return "Sine"
	case t.WaveformSquare:
		return "Square"
	case t.WaveformTriangle:
		return "Triangle"
	case t.WaveformSawtooth:
		return "Sawtooth"
	default:
		return "Unknown"
	}
}

func trackClassForType(trackType t.TrackType) string {
	switch trackType {
	case t.TrackPureTone:
		return "pure"
	case t.TrackBinauralBeat:
		return "binaural"
	case t.TrackMonauralBeat:
		return "monaural"
	case t.TrackIsochronicBeat:
		return "isochronic"
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		return "noise"
	case t.TrackAmbiance:
		return "ambiance"
	case t.TrackSilence:
		return "silence"
	default:
		return "off"
	}
}

func previewClass(name string) string {
	return "track-" + name
}

func buildSeriesLegendLabel(channel int, track t.Track) string {
	return fmt.Sprintf("CH %02d %s", channel+1, humanTrackType(track))
}

func dominantSegmentClass(items []previewSegmentItemView) string {
	if len(items) == 0 {
		return previewClass("off")
	}
	return items[0].Class
}

func formatHz(value float64) string {
	return fmt.Sprintf("%.2f Hz", value)
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatWaveformValue(value float64) string {
	return humanWaveformType(clampWaveformValue(value))
}

func clampWaveformValue(value float64) t.WaveformType {
	waveform := int(math.Round(value))
	if waveform < int(t.WaveformSine) {
		waveform = int(t.WaveformSine)
	}
	if waveform > int(t.WaveformSawtooth) {
		waveform = int(t.WaveformSawtooth)
	}
	return t.WaveformType(waveform)
}

func formatDuration(ms int) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	minutes := (ms % 3600000) / 60000
	seconds := (ms % 60000) / 1000
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func formatTime(ms int) string {
	hours := ms / 3600000
	minutes := (ms % 3600000) / 60000
	seconds := (ms % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func joinParts(items []string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for i := 1; i < len(items); i++ {
		result += " | " + items[i]
	}
	return result
}