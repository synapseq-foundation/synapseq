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

	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func transitionSampleCount(period t.Period) int {
	base := 1
	switch period.Transition {
	case t.TransitionEaseOut, t.TransitionEaseIn:
		base = 8
	case t.TransitionSmooth:
		base = 12
	}

	legs := 2*period.Steps + 1
	if legs < 1 {
		legs = 1
	}

	return base * legs
}

func interpolateTrackForPreview(startTrack, endTrack t.Track, alpha float64) t.Track {
	track := startTrack
	track.Amplitude = t.AmplitudeType(float64(startTrack.Amplitude)*(1-alpha) + float64(endTrack.Amplitude)*alpha)
	track.Carrier = startTrack.Carrier*(1-alpha) + endTrack.Carrier*alpha
	track.Resonance = startTrack.Resonance*(1-alpha) + endTrack.Resonance*alpha
	track.NoiseSmooth = startTrack.NoiseSmooth*(1-alpha) + endTrack.NoiseSmooth*alpha
	track.Effect.Value = startTrack.Effect.Value*(1-alpha) + endTrack.Effect.Value*alpha
	track.Effect.Intensity = t.IntensityType(float64(startTrack.Effect.Intensity)*(1-alpha) + float64(endTrack.Effect.Intensity)*alpha)
	return track
}

func applyTransitionAlpha(progress float64, transition t.TransitionType) float64 {
	return tl.ApplyTransitionAlpha(progress, transition)
}

func stepAlphaForPreview(progress float64, period t.Period) float64 {
	return tl.StepAlpha(progress, period.Transition, period.Steps)
}

func interpolateTime(startTime, endTime int, progress float64) int {
	if endTime <= startTime {
		return startTime
	}
	return startTime + int(math.Round(float64(endTime-startTime)*progress))
}

func graphSeriesColors() []string {
	return []string{"#14532d", "#9a3412", "#1d4ed8", "#0f766e", "#7c2d12", "#4338ca", "#b45309", "#be123c"}
}

func carrierBounds(periods []t.Period) (float64, float64, bool) {
	minCarrier := 0.0
	maxCarrier := 0.0
	hasCarrier := false

	for _, period := range periods {
		for ch := range t.NumberOfChannels {
			track := period.TrackStart[ch]
			if track.Carrier <= 0 || !trackHasGraphPoint(track) {
				continue
			}

			if !hasCarrier {
				minCarrier = track.Carrier
				maxCarrier = track.Carrier
				hasCarrier = true
				continue
			}

			if track.Carrier < minCarrier {
				minCarrier = track.Carrier
			}
			if track.Carrier > maxCarrier {
				maxCarrier = track.Carrier
			}
		}
	}

	return minCarrier, maxCarrier, hasCarrier
}

func buildSeries(periods []t.Period, minCarrier, maxCarrier float64, hasCarrier bool, totalDurationMs int) []previewSeriesView {
	if !hasCarrier {
		return []previewSeriesView{}
	}

	colors := graphSeriesColors()
	series := make([]previewSeriesView, 0, t.NumberOfChannels)
	span := maxCarrier - minCarrier
	if span == 0 {
		span = 1
	}

	for ch := range t.NumberOfChannels {
		points := make([]previewGraphPointView, 0, len(periods))
		coordinates := make([]string, 0, len(periods))
		seriesClass := previewClass("off")

		for _, period := range periods {
			track := period.TrackStart[ch]
			if track.Carrier <= 0 || !trackHasGraphPoint(track) {
				continue
			}

			if len(points) == 0 {
				seriesClass = previewClass(trackClassForType(track.Type))
			}

			x := toGraphX(period.Time, totalDurationMs)
			y := int(220 - ((track.Carrier-minCarrier)/span)*180)

			points = append(points, previewGraphPointView{
				X:          x,
				Y:          y,
				TimeLabel:  formatTime(period.Time),
				ValueLabel: buildGraphValueLabel(track),
			})
			coordinates = append(coordinates, fmt.Sprintf("%d,%d", x, y))
		}

		if len(points) == 0 {
			continue
		}

		series = append(series, previewSeriesView{
			ChannelLabel: fmt.Sprintf("CH %02d", ch+1),
			Class:        seriesClass,
			Color:        colors[ch%len(colors)],
			Points:       joinCoordinates(coordinates),
			Markers:      points,
		})
	}

	return series
}