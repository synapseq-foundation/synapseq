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
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

//go:embed *.html
var files embed.FS

var previewTemplate = template.Must(template.ParseFS(files, "preview.html"))

type previewTemplateData struct {
	Title              string
	TotalDurationLabel string
	TotalDurationMs    int
	PeriodCount        int
	SegmentCount       int
	ActiveChannels     int
	ToneTracks         int
	TextureTracks      int
	CarrierRangeLabel  string
	GraphMinLabel      string
	GraphMaxLabel      string
	Ruler              []previewRulerMarkView
	TimelineNodes      []previewNodeMarkerView
	GraphMetrics       []previewGraphMetricView
	Series             []previewSeriesView
	Segments           []previewSegmentView
	Nodes              []previewNodeView
}

type previewRulerMarkView struct {
	Label   string
	LeftPct float64
}

type previewNodeMarkerView struct {
	TimeLabel   string
	PositionPct float64
}

type previewGraphMetricView struct {
	Key         string
	Label       string
	RangeLabel  string
	MinLabel    string
	MaxLabel    string
	EmptyLabel  string
	Default     bool
	HasData     bool
	LegendItems []previewGraphLegendItemView
	Series      []previewSeriesView
}

type previewGraphLegendItemView struct {
	Label string
	Color string
}

type previewSeriesView struct {
	ChannelLabel string
	LegendLabel  string
	Class        string
	Color        string
	Points       string
	Markers      []previewGraphPointView
}

type previewGraphPointView struct {
	X          int
	Y          int
	TimeLabel  string
	ValueLabel string
}

type previewSegmentView struct {
	Label          string
	StartLabel     string
	EndLabel       string
	DurationLabel  string
	Transition     string
	LeftPct        float64
	WidthPct       float64
	Class          string
	ChannelCount   int
	PrimarySummary string
	Items          []previewSegmentItemView
}

type previewSegmentItemView struct {
	ChannelLabel string
	Class        string
	Label        string
	Summary      string
}

type previewNodeView struct {
	ID           string
	TimeLabel    string
	Transition   string
	PositionPct  float64
	TrackCount   int
	ToneCount    int
	TextureCount int
	Tracks       []previewTrackView
}

type previewTrackView struct {
	ChannelLabel string
	Class        string
	TypeLabel    string
	Summary      string
	Meta         []previewMetaView
}

type previewMetaView struct {
	Label string
	Value string
}

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

func GetPreviewContent(periods []t.Period) ([]byte, error) {
	data, err := buildPreviewData(periods)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	if err := previewTemplate.Execute(&output, data); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func buildPreviewData(periods []t.Period) (*previewTemplateData, error) {
	if len(periods) < 2 {
		return nil, fmt.Errorf("preview requires at least two periods")
	}

	totalDurationMs := periods[len(periods)-1].Time
	if totalDurationMs <= 0 {
		return nil, fmt.Errorf("preview requires a positive total duration")
	}

	activeChannels := 0
	toneTracks := 0
	textureTracks := 0
	minCarrier, maxCarrier, hasCarrier := carrierBounds(periods)

	timelineNodes := make([]previewNodeMarkerView, 0, len(periods))
	nodes := make([]previewNodeView, 0, len(periods))
	for idx, period := range periods {
		timelineNodes = append(timelineNodes, previewNodeMarkerView{
			TimeLabel:   formatTime(period.Time),
			PositionPct: toGraphPercent(period.Time, totalDurationMs),
		})

		tracks := make([]previewTrackView, 0, t.NumberOfChannels)
		nodeToneCount := 0
		nodeTextureCount := 0

		for ch := range t.NumberOfChannels {
			track, ok := resolveNodeTrack(periods, idx, ch)
			if !ok {
				continue
			}

			tracks = append(tracks, buildTrackView(ch, track))
			if isToneTrack(track) {
				nodeToneCount++
				toneTracks++
			}
			if isTextureTrack(track) {
				nodeTextureCount++
				textureTracks++
			}
		}

		if len(tracks) > activeChannels {
			activeChannels = len(tracks)
		}

		if len(tracks) == 0 {
			continue
		}

		transition := period.Transition.String()
		if idx == len(periods)-1 {
			transition = "end"
		}

		nodes = append(nodes, previewNodeView{
			ID:           fmt.Sprintf("node-%d", idx),
			TimeLabel:    formatTime(period.Time),
			Transition:   transition,
			PositionPct:  toGraphPercent(period.Time, totalDurationMs),
			TrackCount:   len(tracks),
			ToneCount:    nodeToneCount,
			TextureCount: nodeTextureCount,
			Tracks:       tracks,
		})
	}

	segments := make([]previewSegmentView, 0, len(periods)-1)
	for idx := 0; idx < len(periods)-1; idx++ {
		current := periods[idx]
		next := periods[idx+1]
		items := make([]previewSegmentItemView, 0, t.NumberOfChannels)

		for ch := range t.NumberOfChannels {
			startTrack := current.TrackStart[ch]
			endTrack := current.TrackEnd[ch]
			if !shouldRenderSegmentItem(startTrack, endTrack) {
				continue
			}

			items = append(items, buildSegmentItemView(ch, startTrack, endTrack))
		}

		segments = append(segments, previewSegmentView{
			Label:          fmt.Sprintf("%s -> %s", formatTime(current.Time), formatTime(next.Time)),
			StartLabel:     formatTime(current.Time),
			EndLabel:       formatTime(next.Time),
			DurationLabel:  formatDuration(next.Time - current.Time),
			Transition:     current.Transition.String(),
			LeftPct:        toGraphPercent(current.Time, totalDurationMs),
			WidthPct:       toGraphWidth(next.Time-current.Time, totalDurationMs),
			Class:          dominantSegmentClass(items),
			ChannelCount:   len(items),
			PrimarySummary: buildPrimarySummary(items),
			Items:          items,
		})
	}

	series := buildSeries(periods, minCarrier, maxCarrier, hasCarrier, totalDurationMs)
	graphMetrics := buildGraphMetrics(periods, totalDurationMs)

	carrierRangeLabel := "No tonal carriers"
	graphMinLabel := "0 Hz"
	graphMaxLabel := "0 Hz"
	if hasCarrier {
		carrierRangeLabel = fmt.Sprintf("%s - %s", formatHz(minCarrier), formatHz(maxCarrier))
		graphMinLabel = formatHz(minCarrier)
		graphMaxLabel = formatHz(maxCarrier)
	}

	return &previewTemplateData{
		Title:              "SynapSeq Sequence Preview",
		TotalDurationLabel: formatTime(totalDurationMs),
		TotalDurationMs:    totalDurationMs,
		PeriodCount:        len(periods),
		SegmentCount:       len(periods) - 1,
		ActiveChannels:     activeChannels,
		ToneTracks:         toneTracks,
		TextureTracks:      textureTracks,
		CarrierRangeLabel:  carrierRangeLabel,
		GraphMinLabel:      graphMinLabel,
		GraphMaxLabel:      graphMaxLabel,
		Ruler:              buildRuler(totalDurationMs),
		TimelineNodes:      timelineNodes,
		GraphMetrics:       graphMetrics,
		Series:             series,
		Segments:           segments,
		Nodes:              nodes,
	}, nil
}

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

			samples := transitionSampleCount(period.Transition)
			for step := 0; step <= samples; step++ {
				if idx > 0 && step == 0 {
					continue
				}

				progress := float64(step) / float64(samples)
				alpha := applyTransitionAlpha(progress, period.Transition)
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

func transitionSampleCount(transition t.TransitionType) int {
	switch transition {
	case t.TransitionEaseOut, t.TransitionEaseIn:
		return 8
	case t.TransitionSmooth:
		return 12
	default:
		return 1
	}
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
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	alpha := progress
	switch transition {
	case t.TransitionEaseOut:
		alpha = math.Log1p(math.Expm1(t.TransitionCurveK)*progress) / t.TransitionCurveK
	case t.TransitionEaseIn:
		alpha = math.Expm1(t.TransitionCurveK*progress) / math.Expm1(t.TransitionCurveK)
	case t.TransitionSmooth:
		raw := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*(progress-0.5)))
		min := 1.0 / (1.0 + math.Exp(t.TransitionCurveK*0.5))
		max := 1.0 / (1.0 + math.Exp(-t.TransitionCurveK*0.5))
		alpha = (raw - min) / (max - min)
	}

	return alpha
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

func buildGraphValueLabel(track t.Track) string {
	if usesBeat(track) {
		return fmt.Sprintf("%s / %s", formatHz(track.Carrier), formatHz(track.Resonance))
	}
	return formatHz(track.Carrier)
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

func includeTrack(track t.Track) bool {
	return track.Type != t.TrackOff
}

func includeVisibleTrack(track t.Track) bool {
	return includeTrack(track) && track.Type != t.TrackSilence
}

func shouldRenderSegmentItem(startTrack, endTrack t.Track) bool {
	return includeVisibleTrack(startTrack) || includeVisibleTrack(endTrack)
}

func trackHasGraphPoint(track t.Track) bool {
	return isToneTrack(track) || track.Type == t.TrackSilence
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
