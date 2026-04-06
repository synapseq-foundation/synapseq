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

import "html/template"

type previewTemplateData struct {
	Title              string
	PreviewCSS         template.CSS
	PreviewJS          template.JS
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
