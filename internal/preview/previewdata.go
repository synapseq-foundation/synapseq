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

type previewNodeData struct {
	ActiveChannels int
	ToneTracks     int
	TextureTracks  int
	TimelineNodes  []previewNodeMarkerView
	Nodes          []previewNodeView
}

type previewRenderData struct {
	NodeData          previewNodeData
	Series            []previewSeriesView
	Segments          []previewSegmentView
	GraphMetrics      []previewGraphMetricView
	CarrierRangeLabel string
	GraphMinLabel     string
	GraphMaxLabel     string
}

func validatePreviewPeriods(periods []t.Period) (int, error) {
	if len(periods) < 2 {
		return 0, fmt.Errorf("preview requires at least two periods")
	}

	totalDurationMs := periods[len(periods)-1].Time
	if totalDurationMs <= 0 {
		return 0, fmt.Errorf("preview requires a positive total duration")
	}

	return totalDurationMs, nil
}

func buildPreviewRenderData(periods []t.Period, totalDurationMs int) previewRenderData {
	minCarrier, maxCarrier, hasCarrier := carrierBounds(periods)
	nodeData := buildNodeData(periods, totalDurationMs)
	segments := buildSegmentViews(periods, totalDurationMs)
	series := buildSeries(periods, minCarrier, maxCarrier, hasCarrier, totalDurationMs)
	graphMetrics := buildGraphMetrics(periods, totalDurationMs)
	carrierRangeLabel, graphMinLabel, graphMaxLabel := buildCarrierRangeLabels(minCarrier, maxCarrier, hasCarrier)

	return previewRenderData{
		NodeData:          nodeData,
		Series:            series,
		Segments:          segments,
		GraphMetrics:      graphMetrics,
		CarrierRangeLabel: carrierRangeLabel,
		GraphMinLabel:     graphMinLabel,
		GraphMaxLabel:     graphMaxLabel,
	}
}

func buildPreviewTemplateData(periods []t.Period, totalDurationMs int, renderData previewRenderData) *previewTemplateData {
	return &previewTemplateData{
		Title:              "SynapSeq Sequence Preview",
		PreviewCSS:         previewCSS,
		PreviewJS:          previewJS,
		TotalDurationLabel: formatTime(totalDurationMs),
		TotalDurationMs:    totalDurationMs,
		PeriodCount:        len(periods),
		SegmentCount:       len(periods) - 1,
		ActiveChannels:     renderData.NodeData.ActiveChannels,
		ToneTracks:         renderData.NodeData.ToneTracks,
		TextureTracks:      renderData.NodeData.TextureTracks,
		CarrierRangeLabel:  renderData.CarrierRangeLabel,
		GraphMinLabel:      renderData.GraphMinLabel,
		GraphMaxLabel:      renderData.GraphMaxLabel,
		Ruler:              buildRuler(totalDurationMs),
		TimelineNodes:      renderData.NodeData.TimelineNodes,
		GraphMetrics:       renderData.GraphMetrics,
		Series:             renderData.Series,
		Segments:           renderData.Segments,
		Nodes:              renderData.NodeData.Nodes,
	}
}

func buildNodeData(periods []t.Period, totalDurationMs int) previewNodeData {
	data := previewNodeData{
		TimelineNodes: make([]previewNodeMarkerView, 0, len(periods)),
		Nodes:         make([]previewNodeView, 0, len(periods)),
	}

	for idx, period := range periods {
		data.TimelineNodes = append(data.TimelineNodes, previewNodeMarkerView{
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
				data.ToneTracks++
			}
			if isTextureTrack(track) {
				nodeTextureCount++
				data.TextureTracks++
			}
		}

		if len(tracks) > data.ActiveChannels {
			data.ActiveChannels = len(tracks)
		}

		if len(tracks) == 0 {
			continue
		}

		transition := period.Transition.String()
		if idx == len(periods)-1 {
			transition = "end"
		} else {
			transition = formatPreviewTransition(period)
		}

		data.Nodes = append(data.Nodes, previewNodeView{
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

	return data
}

func buildSegmentViews(periods []t.Period, totalDurationMs int) []previewSegmentView {
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
			Transition:     formatPreviewTransition(current),
			LeftPct:        toGraphPercent(current.Time, totalDurationMs),
			WidthPct:       toGraphWidth(next.Time-current.Time, totalDurationMs),
			Class:          dominantSegmentClass(items),
			ChannelCount:   len(items),
			PrimarySummary: buildPrimarySummary(items),
			Items:          items,
		})
	}

	return segments
}

func buildCarrierRangeLabels(minCarrier, maxCarrier float64, hasCarrier bool) (string, string, string) {
	if !hasCarrier {
		return "No tonal carriers", "0 Hz", "0 Hz"
	}

	return fmt.Sprintf("%s - %s", formatHz(minCarrier), formatHz(maxCarrier)), formatHz(minCarrier), formatHz(maxCarrier)
}
