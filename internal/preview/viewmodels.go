// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preview

import "html/template"

type previewTemplateData struct {
	Title       string
	PreviewCSS  template.CSS
	PreviewJS   template.JS
	LucideJS    template.JS
	ChartJS     template.JS
	LogoDataURL template.URL
	Stats       previewStatsView
	Timeline    previewTimelineSectionView
	Nodes       previewNodesSectionView
}

type previewStatsView struct {
	TotalDurationLabel string
	PeriodCount        int
	SegmentCount       int
	ActiveChannels     int
	ToneTracks         int
	TextureTracks      int
}

type previewTimelineSectionView struct {
	SegmentCount       int
	TotalDurationLabel string
	Metrics            []previewTimelineMetricView
	Segments           []previewSegmentView
}

type previewTimelineMetricView struct {
	Key           string
	Label         string
	RangeLabel    string
	MinLabel      string
	MaxLabel      string
	EmptyLabel    string
	Default       bool
	HasData       bool
	LegendItems   []previewGraphLegendItemView
	Series        []previewSeriesView
	Segments      []previewSegmentView
	Ruler         []previewRulerMarkView
	TimelineNodes []previewNodeMarkerView
}

type previewNodesSectionView struct {
	PeriodCount int
	Nodes       []previewNodeView
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
	ChartData    template.JS
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
	Crossfade    bool
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
