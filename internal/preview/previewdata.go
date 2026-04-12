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

type previewNodeSectionData struct {
	ActiveChannels int
	ToneTracks     int
	TextureTracks  int
	TimelineNodes  []previewNodeMarkerView
	Nodes          []previewNodeView
}

type previewSectionData struct {
	NodeSection     previewNodeSectionData
	Ruler           []previewRulerMarkView
	Segments        []previewSegmentView
	TimelineMetrics []previewGraphMetricView
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

func buildPreviewSectionData(periods []t.Period, totalDurationMs int) previewSectionData {
	nodeSection := buildNodeSectionData(periods, totalDurationMs)
	segments := buildSegmentViews(periods, totalDurationMs)

	return previewSectionData{
		NodeSection:     nodeSection,
		Ruler:           buildRuler(totalDurationMs),
		Segments:        segments,
		TimelineMetrics: buildGraphMetrics(periods, totalDurationMs),
	}
}

func buildPreviewTemplateData(periods []t.Period, totalDurationMs int, sectionData previewSectionData) *previewTemplateData {
	return &previewTemplateData{
		Title:      "SynapSeq Sequence Preview",
		PreviewCSS: previewCSS,
		PreviewJS:  previewJS,
		Stats:      buildPreviewStatsView(periods, totalDurationMs, sectionData.NodeSection),
		Timeline:   buildPreviewTimelineSectionView(periods, totalDurationMs, sectionData),
		Nodes:      buildPreviewNodesSectionView(periods, sectionData.NodeSection),
	}
}

func buildPreviewStatsView(periods []t.Period, totalDurationMs int, nodeSection previewNodeSectionData) previewStatsView {
	return previewStatsView{
		TotalDurationLabel: formatTime(totalDurationMs),
		PeriodCount:        len(periods),
		SegmentCount:       len(periods) - 1,
		ActiveChannels:     nodeSection.ActiveChannels,
		ToneTracks:         nodeSection.ToneTracks,
		TextureTracks:      nodeSection.TextureTracks,
	}
}

func buildPreviewTimelineSectionView(periods []t.Period, totalDurationMs int, sectionData previewSectionData) previewTimelineSectionView {
	return previewTimelineSectionView{
		SegmentCount:       len(periods) - 1,
		TotalDurationLabel: formatTime(totalDurationMs),
		Metrics:            buildPreviewTimelineMetrics(sectionData),
		Segments:           sectionData.Segments,
	}
}

func buildPreviewTimelineMetrics(sectionData previewSectionData) []previewTimelineMetricView {
	metricViews := make([]previewTimelineMetricView, 0, len(sectionData.TimelineMetrics))
	for _, metric := range sectionData.TimelineMetrics {
		metricViews = append(metricViews, previewTimelineMetricView{
			Key:           metric.Key,
			Label:         metric.Label,
			RangeLabel:    metric.RangeLabel,
			MinLabel:      metric.MinLabel,
			MaxLabel:      metric.MaxLabel,
			EmptyLabel:    metric.EmptyLabel,
			Default:       metric.Default,
			HasData:       metric.HasData,
			LegendItems:   metric.LegendItems,
			Series:        metric.Series,
			Segments:      sectionData.Segments,
			Ruler:         sectionData.Ruler,
			TimelineNodes: sectionData.NodeSection.TimelineNodes,
		})
	}

	return metricViews
}

func buildPreviewNodesSectionView(periods []t.Period, nodeSection previewNodeSectionData) previewNodesSectionView {
	return previewNodesSectionView{
		PeriodCount: len(periods),
		Nodes:       nodeSection.Nodes,
	}
}

func buildNodeSectionData(periods []t.Period, totalDurationMs int) previewNodeSectionData {
	sectionData := previewNodeSectionData{
		TimelineNodes: make([]previewNodeMarkerView, 0, len(periods)),
		Nodes:         make([]previewNodeView, 0, len(periods)),
	}

	for idx, period := range periods {
		sectionData.TimelineNodes = append(sectionData.TimelineNodes, previewNodeMarkerView{
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
				sectionData.ToneTracks++
			}
			if isTextureTrack(track) {
				nodeTextureCount++
				sectionData.TextureTracks++
			}
		}

		if len(tracks) > sectionData.ActiveChannels {
			sectionData.ActiveChannels = len(tracks)
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

		sectionData.Nodes = append(sectionData.Nodes, previewNodeView{
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

	return sectionData
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

