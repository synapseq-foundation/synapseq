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

type graphMetricDefinition struct {
	Key                   string
	Label                 string
	RangePrefix           string
	EmptyLabel            string
	Default               bool
	SelectValue           func(track t.Track) (float64, bool)
	SelectTransitionValue func(startTrack, endTrack t.Track, alpha float64) (float64, bool)
	BuildLegendItems      func(series []previewSeriesView, periods []t.Period) []previewGraphLegendItemView
	BuildPointLabel       func(channel int, track t.Track) string
	FormatValue           func(value float64) string
}

type rawGraphSeriesPoint struct {
	Time  int
	Value float64
	Label string
}

type rawGraphSeries struct {
	Channel int
	Legend  string
	Class   string
	Curve   []rawGraphSeriesPoint
	Markers []rawGraphSeriesPoint
}

const (
	graphInsetPct   = 4.0
	graphViewWidth  = 1000.0
	graphInnerPct   = 100.0 - (graphInsetPct * 2)
	graphInnerWidth = graphViewWidth * (graphInnerPct / 100.0)
	graphMinX       = graphViewWidth * (graphInsetPct / 100.0)
)

func buildGraphMetrics(periods []t.Period, totalDurationMs int) []previewGraphMetricView {
	definitions := []graphMetricDefinition{
		{
			Key:         "resonance",
			Label:       "Resonance",
			RangePrefix: "Beat range",
			EmptyLabel:  "No beat or resonance data available for this sequence.",
			Default:     true,
			SelectValue: func(track t.Track) (float64, bool) {
				if !usesBeat(track) {
					return 0, false
				}
				return track.Resonance, true
			},
			FormatValue: formatHz,
		},
		{
			Key:         "carrier",
			Label:       "Carrier",
			RangePrefix: "Carrier range",
			EmptyLabel:  "No carrier frequency data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !isToneTrack(track) || track.Carrier <= 0 {
					return 0, false
				}
				return track.Carrier, true
			},
			FormatValue: formatHz,
		},
		{
			Key:         "waveform",
			Label:       "Waveform",
			RangePrefix: "Waveform range",
			EmptyLabel:  "No waveform data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !supportsWaveform(track) {
					return 0, false
				}
				return float64(track.Waveform), true
			},
			SelectTransitionValue: func(startTrack, endTrack t.Track, alpha float64) (float64, bool) {
				if !supportsWaveform(startTrack) {
					return 0, false
				}

				startValue := float64(startTrack.Waveform)
				endValue := startValue
				if supportsWaveform(endTrack) {
					endValue = float64(endTrack.Waveform)
				}

				return startValue*(1-alpha) + endValue*alpha, true
			},
			BuildLegendItems: buildWaveformLegendItems,
			BuildPointLabel:  buildWaveformPointLabel,
			FormatValue:      formatWaveformValue,
		},
		{
			Key:         "amplitude",
			Label:       "Amplitude",
			RangePrefix: "Amplitude range",
			EmptyLabel:  "No amplitude data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !includeVisibleTrack(track) {
					return 0, false
				}
				return track.Amplitude.ToPercent(), true
			},
			FormatValue: formatPercent,
		},
		{
			Key:         "smooth",
			Label:       "Smooth",
			RangePrefix: "Smooth range",
			EmptyLabel:  "No noise smooth data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !isNoiseTrack(track) {
					return 0, false
				}
				return track.NoiseSmooth, true
			},
			FormatValue: formatPercent,
		},
		{
			Key:         "effect",
			Label:       "Effect",
			RangePrefix: "Effect range",
			EmptyLabel:  "No effect value data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !includeVisibleTrack(track) || track.Effect.Type == t.EffectOff {
					return 0, false
				}
				return track.Effect.Value, true
			},
			FormatValue: formatFloat,
		},
		{
			Key:         "intensity",
			Label:       "Intensity",
			RangePrefix: "Effect intensity range",
			EmptyLabel:  "No effect intensity data available for this sequence.",
			SelectValue: func(track t.Track) (float64, bool) {
				if !includeVisibleTrack(track) || track.Effect.Type == t.EffectOff {
					return 0, false
				}
				return track.Effect.Intensity.ToPercent(), true
			},
			FormatValue: formatPercent,
		},
	}

	metrics := make([]previewGraphMetricView, 0, len(definitions))
	for _, definition := range definitions {
		metric := buildGraphMetricView(definition, periods, totalDurationMs)
		metrics = append(metrics, metric)
	}

	return metrics
}

func buildGraphMetricView(definition graphMetricDefinition, periods []t.Period, totalDurationMs int) previewGraphMetricView {
	rawSeries, minValue, maxValue, hasData := collectGraphMetricSeries(definition, periods)
	series := make([]previewSeriesView, 0, len(rawSeries))

	if hasData {
		series = buildPreviewSeriesFromRaw(rawSeries, minValue, maxValue, totalDurationMs, definition.FormatValue)
	}

	legendItems := buildMetricLegendItems(definition, series, periods)

	rangeLabel := definition.EmptyLabel
	minLabel := "0"
	maxLabel := "0"
	if hasData {
		minLabel = definition.FormatValue(minValue)
		maxLabel = definition.FormatValue(maxValue)
		rangeLabel = fmt.Sprintf("%s: %s - %s", definition.RangePrefix, minLabel, maxLabel)
	}

	return previewGraphMetricView{
		Key:         definition.Key,
		Label:       definition.Label,
		RangeLabel:  rangeLabel,
		MinLabel:    minLabel,
		MaxLabel:    maxLabel,
		EmptyLabel:  definition.EmptyLabel,
		Default:     definition.Default,
		HasData:     hasData,
		LegendItems: legendItems,
		Series:      series,
	}
}

func collectGraphMetricSeries(definition graphMetricDefinition, periods []t.Period) ([]rawGraphSeries, float64, float64, bool) {
	rawSeries := make([]rawGraphSeries, 0, t.NumberOfChannels)
	hasData := false
	minValue := 0.0
	maxValue := 0.0

	for ch := range t.NumberOfChannels {
		curvePoints := make([]rawGraphSeriesPoint, 0, len(periods)*8)
		markerPoints := make([]rawGraphSeriesPoint, 0, len(periods))
		seriesClass := previewClass("off")
		legendLabel := fmt.Sprintf("CH %02d", ch+1)
		seriesSet := false

		for idx := 0; idx < len(periods)-1; idx++ {
			period := periods[idx]
			next := periods[idx+1]
			startTrack := period.TrackStart[ch]
			endTrack := period.TrackEnd[ch]

			if !includeVisibleTrack(startTrack) {
				continue
			}

			if !seriesSet {
				seriesClass = previewClass(trackClassForType(startTrack.Type))
				legendLabel = buildSeriesLegendLabel(ch, startTrack)
				seriesSet = true
			}

			samples := transitionSampleCount(period)
			for step := 0; step <= samples; step++ {
				if idx > 0 && step == 0 {
					continue
				}

				progress := float64(step) / float64(samples)
				alpha := stepAlphaForPreview(progress, period)
				var value float64
				var ok bool
				if definition.SelectTransitionValue != nil {
					value, ok = definition.SelectTransitionValue(startTrack, endTrack, alpha)
				} else {
					track := interpolateTrackForPreview(startTrack, endTrack, alpha)
					value, ok = definition.SelectValue(track)
				}
				if !ok {
					continue
				}

				if !hasData {
					minValue = value
					maxValue = value
					hasData = true
				} else {
					if value < minValue {
						minValue = value
					}
					if value > maxValue {
						maxValue = value
					}
				}

				time := interpolateTime(period.Time, next.Time, progress)
				curvePoints = append(curvePoints, rawGraphSeriesPoint{
					Time:  time,
					Value: value,
					Label: buildGraphPointLabel(definition, ch, trackForMetricLabel(definition, startTrack, endTrack, alpha)),
				})
			}
		}

		for idx := range periods {
			period := periods[idx]
			track, ok := resolveGraphTrack(periods, idx, ch)
			if !ok {
				continue
			}

			value, ok := definition.SelectValue(track)
			if !ok {
				continue
			}

			markerPoints = append(markerPoints, rawGraphSeriesPoint{
				Time:  period.Time,
				Value: value,
				Label: buildGraphPointLabel(definition, ch, track),
			})
		}

		if len(curvePoints) == 0 {
			continue
		}

		rawSeries = append(rawSeries, rawGraphSeries{
			Channel: ch,
			Legend:  legendLabel,
			Class:   seriesClass,
			Curve:   curvePoints,
			Markers: markerPoints,
		})
	}

	return rawSeries, minValue, maxValue, hasData
}

func buildPreviewSeriesFromRaw(rawSeries []rawGraphSeries, minValue, maxValue float64, totalDurationMs int, formatValue func(value float64) string) []previewSeriesView {
	colors := graphSeriesColors()
	series := make([]previewSeriesView, 0, len(rawSeries))
	span := maxValue - minValue
	if span == 0 {
		span = 1
	}

	for _, raw := range rawSeries {
		markers := make([]previewGraphPointView, 0, len(raw.Markers))
		coordinates := make([]string, 0, len(raw.Curve))
		color := colors[raw.Channel%len(colors)]

		for _, point := range raw.Curve {
			x := toGraphX(point.Time, totalDurationMs)
			y := int(220 - ((point.Value-minValue)/span)*180)
			coordinates = append(coordinates, fmt.Sprintf("%d,%d", x, y))
		}

		for _, point := range raw.Markers {
			x := toGraphX(point.Time, totalDurationMs)
			y := int(220 - ((point.Value-minValue)/span)*180)

			markers = append(markers, previewGraphPointView{
				X:          x,
				Y:          y,
				TimeLabel:  formatTime(point.Time),
				ValueLabel: point.Label,
			})
		}

		series = append(series, previewSeriesView{
			ChannelLabel: fmt.Sprintf("CH %02d", raw.Channel+1),
			LegendLabel:  raw.Legend,
			Class:        raw.Class,
			Color:        color,
			Points:       joinCoordinates(coordinates),
			Markers:      markers,
		})
	}

	_ = formatValue
	return series
}

func buildMetricLegendItems(definition graphMetricDefinition, series []previewSeriesView, periods []t.Period) []previewGraphLegendItemView {
	if definition.BuildLegendItems != nil {
		return definition.BuildLegendItems(series, periods)
	}

	items := make([]previewGraphLegendItemView, 0, len(series))
	for _, current := range series {
		items = append(items, previewGraphLegendItemView{
			Label: current.LegendLabel,
			Color: current.Color,
		})
	}

	return items
}

func buildWaveformLegendItems(series []previewSeriesView, _ []t.Period) []previewGraphLegendItemView {
	items := make([]previewGraphLegendItemView, 0, len(series))
	for _, current := range series {
		items = append(items, previewGraphLegendItemView{
			Label: current.LegendLabel,
			Color: current.Color,
		})
	}

	return items
}

func buildGraphPointLabel(definition graphMetricDefinition, channel int, track t.Track) string {
	if definition.BuildPointLabel != nil {
		return definition.BuildPointLabel(channel, track)
	}

	value, ok := definition.SelectValue(track)
	if !ok {
		return ""
	}

	return definition.FormatValue(value)
}

func trackForMetricLabel(definition graphMetricDefinition, startTrack, endTrack t.Track, alpha float64) t.Track {
	if definition.SelectTransitionValue == nil {
		return startTrack
	}

	track := interpolateTrackForPreview(startTrack, endTrack, alpha)
	if definition.Key == "waveform" {
		waveformValue, ok := definition.SelectTransitionValue(startTrack, endTrack, alpha)
		if ok {
			track.Waveform = clampWaveformValue(waveformValue)
		}
	}

	return track
}

func buildWaveformPointLabel(channel int, track t.Track) string {
	return fmt.Sprintf("%s • %s", buildSeriesLegendLabel(channel, track), humanWaveformType(track.Waveform))
}

func buildRuler(totalDurationMs int) []previewRulerMarkView {
	marks := make([]previewRulerMarkView, 0, 6)
	for i := 0; i <= 5; i++ {
		ms := (totalDurationMs * i) / 5
		marks = append(marks, previewRulerMarkView{
			Label:   formatTime(ms),
			LeftPct: toGraphPercent(ms, totalDurationMs),
		})
	}
	return marks
}

func toPercent(part, total int) float64 {
	if total <= 0 {
		return 0
	}
	return (float64(part) / float64(total)) * 100
}

func toGraphPercent(part, total int) float64 {
	return graphInsetPct + ((toPercent(part, total) / 100.0) * graphInnerPct)
}

func toGraphWidth(part, total int) float64 {
	return (toPercent(part, total) / 100.0) * graphInnerPct
}

func toGraphX(part, total int) int {
	if total <= 0 {
		return int(graphMinX)
	}

	return int(math.Round(graphMinX + (float64(part)/float64(total))*graphInnerWidth))
}

func joinCoordinates(items []string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for i := 1; i < len(items); i++ {
		result += " " + items[i]
	}
	return result
}
