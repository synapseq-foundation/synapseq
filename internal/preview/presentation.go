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

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func buildTrackView(channel int, track t.Track) previewTrackView {
	meta := make([]previewMetaView, 0, 8)

	if track.Carrier > 0 {
		meta = append(meta, previewMetaView{Label: "Carrier", Value: formatHz(track.Carrier)})
	}
	if usesBeat(track) || track.Resonance > 0 {
		meta = append(meta, previewMetaView{Label: "Beat", Value: formatHz(track.Resonance)})
	}
	if waveform := track.Waveform.String(); waveform != "" && supportsWaveform(track) {
		meta = append(meta, previewMetaView{Label: "Waveform", Value: waveform})
	}
	if track.AmbianceName != "" {
		meta = append(meta, previewMetaView{Label: "Ambiance", Value: track.AmbianceName})
	}
	if isNoiseTrack(track) {
		meta = append(meta, previewMetaView{Label: "Smooth", Value: formatPercent(track.NoiseSmooth)})
	}
	meta = append(meta, previewMetaView{Label: "Amplitude", Value: formatPercent(track.Amplitude.ToPercent())})
	if track.Effect.Type != t.EffectOff {
		meta = append(meta,
			previewMetaView{Label: "Effect", Value: track.Effect.Type.String()},
			previewMetaView{Label: "Effect value", Value: formatFloat(track.Effect.Value)},
			previewMetaView{Label: "Intensity", Value: formatPercent(track.Effect.Intensity.ToPercent())},
		)
	}

	return previewTrackView{
		ChannelLabel: fmt.Sprintf("CH %02d", channel+1),
		Class:        previewClass(trackClassForType(track.Type)),
		TypeLabel:    humanTrackType(track),
		Summary:      buildTrackSummary(track),
		Meta:         meta,
	}
}

func buildSegmentItemView(channel int, startTrack, endTrack t.Track) previewSegmentItemView {
	class := trackClassForType(primaryTrackType(startTrack, endTrack))

	return previewSegmentItemView{
		ChannelLabel: fmt.Sprintf("CH %02d", channel+1),
		Class:        previewClass(class),
		Label:        humanTrackType(preferredTrack(startTrack, endTrack)),
		Summary:      buildSegmentSummary(startTrack, endTrack),
	}
}

func formatPreviewTransition(period t.Period) string {
	if period.Steps <= 0 {
		return period.Transition.String() + " - no steps"
	}

	label := "steps"
	if period.Steps == 1 {
		label = "step"
	}

	return fmt.Sprintf("%s - %d %s", period.Transition.String(), period.Steps, label)
}