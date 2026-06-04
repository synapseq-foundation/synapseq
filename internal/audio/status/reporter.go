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

package status

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/palette"
	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type View struct {
	Periods  []t.Period
	Channels []t.Channel
}

type Reporter struct {
	out             io.Writer
	colors          bool
	lastStatusWidth int
	lastPeriodIdx   int
	updateCounter   int
}

func NewReporter(out io.Writer, colors bool) *Reporter {
	return &Reporter{
		out:           out,
		colors:        colors,
		lastPeriodIdx: -1,
	}
}

func (sr *Reporter) DisplayPeriodChange(view View, periodIdx int) {
	if periodIdx >= len(view.Periods) || sr.out == nil {
		return
	}

	period := view.Periods[periodIdx]
	var nextPeriod *t.Period
	if periodIdx+1 < len(view.Periods) {
		nextPeriod = &view.Periods[periodIdx+1]
	} else {
		nextPeriod = &period
	}

	if sr.lastStatusWidth > 0 {
		fmt.Fprintf(sr.out, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		sr.lastStatusWidth = 0
	}

	line1 := fmt.Sprintf("%s %s %s %s %s",
		sr.statusBullet("-"),
		sr.statusTime(period.TimeString()),
		sr.statusArrow("->"),
		sr.statusTime(nextPeriod.TimeString()),
		sr.formatTransition(period))

	line2 := ""
	duration := tl.CrossfadeDuration(nextPeriod.Time - period.Time)
	for ch := range CountPeriodDisplayChannels(view, period) {
		startTrack := period.TrackStart[ch]
		endTrack := period.TrackEnd[ch]

		if startTrack.Type != t.TrackOff && startTrack.Type != t.TrackSilence {
			line2 += fmt.Sprintf("\n%s %s", strings.Repeat(" ", 6), sr.formatTrackWithFadeLabels(startTrack, sr.crossfadeLabels(period, ch, duration)...))
		}

		if !IsTrackEqual(&startTrack, &endTrack) {
			line2 += fmt.Sprintf("\n   %s  %s", sr.statusArrow("->"), sr.formatTrackWithFadeLabels(endTrack))
		}
	}

	fmt.Fprintf(sr.out, "%s%s\n", line1, line2)
}

func (sr *Reporter) DisplayStatus(view View, currentTimeMs int) {
	if sr.out == nil {
		return
	}

	hh := currentTimeMs / 3600000
	mm := (currentTimeMs % 3600000) / 60000
	ss := (currentTimeMs % 60000) / 1000

	status := sr.statusTime(fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss))
	for ch := range CountActiveChannels(view.Channels) {
		channel := &view.Channels[ch]
		status += sr.statusTrack(channel.Track.ShortString())
	}

	status = "  " + status

	clearStr := ""
	if sr.lastStatusWidth > len(status) {
		clearStr = strings.Repeat(" ", sr.lastStatusWidth-len(status))
	}

	fmt.Fprintf(sr.out, "%s%s\r", status, clearStr)
	sr.lastStatusWidth = len(status)
}

func (sr *Reporter) CheckPeriodChange(view View, periodIdx int) {
	if periodIdx != sr.lastPeriodIdx {
		sr.DisplayPeriodChange(view, periodIdx)
		sr.lastPeriodIdx = periodIdx
	}
}

func (sr *Reporter) ShouldUpdateStatus() bool {
	sr.updateCounter++
	return sr.updateCounter%44 == 0
}

func (sr *Reporter) FinalStatus() {
	if sr.out == nil {
		return
	}
	if sr.lastStatusWidth > 0 {
		fmt.Fprintf(sr.out, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		fmt.Fprintf(sr.out, "\n")
	}
}

func CountActiveChannels(channels []t.Channel) int {
	for i := len(channels) - 1; i >= 0; i-- {
		if channels[i].Track.Type != t.TrackOff {
			return i + 1
		}
	}

	return 1
}

func CountPeriodDisplayChannels(view View, period t.Period) int {
	count := CountActiveChannels(view.Channels)
	for i := t.NumberOfChannels - 1; i >= 0; i-- {
		if period.CrossfadeIn[i].Active || period.CrossfadeOut[i].Active {
			if i+1 > count {
				return i + 1
			}
			return count
		}
	}
	return count
}

func IsTrackEqual(trackA, trackB *t.Track) bool {
	return trackA.Type == trackB.Type &&
		trackA.Amplitude == trackB.Amplitude &&
		trackA.SourceName == trackB.SourceName &&
		trackA.Carrier == trackB.Carrier &&
		trackA.Resonance == trackB.Resonance &&
		trackA.Waveform == trackB.Waveform &&
		trackA.Effect.Intensity == trackB.Effect.Intensity
}

func (sr *Reporter) crossfadeLabels(period t.Period, ch int, durationMs int) []string {
	labels := []string{}
	if period.CrossfadeIn[ch].Active {
		labels = append(labels, sr.formatFadeLabel("fade-in", durationMs))
	}
	if period.CrossfadeOut[ch].Active {
		labels = append(labels, sr.formatFadeLabel("fade-out", durationMs))
	}
	return labels
}

func (sr *Reporter) formatTrackWithFadeLabels(track t.Track, labels ...string) string {
	out := sr.statusTrack(track.String())
	for _, label := range labels {
		if label != "" {
			out += " " + label
		}
	}
	return out
}

func (sr *Reporter) formatFadeLabel(direction string, durationMs int) string {
	return sr.statusFadeLabel(fmt.Sprintf("(%s %s)", direction, formatFadeDuration(durationMs)))
}

func formatFadeDuration(durationMs int) string {
	if durationMs < 1000 {
		return fmt.Sprintf("%dms", durationMs)
	}
	if durationMs%1000 == 0 {
		return fmt.Sprintf("%ds", durationMs/1000)
	}
	return strconv.FormatFloat(float64(durationMs)/1000, 'f', -1, 64) + "s"
}

func (sr *Reporter) statusRGB(text string, token palette.RGBColor, attrs ...color.Attribute) string {
	if !sr.colors {
		return text
	}
	style := color.RGB(token.R(), token.G(), token.B())
	if len(attrs) > 0 {
		style.Add(attrs...)
	}
	style.EnableColor()
	return style.Sprint(text)
}

func (sr *Reporter) statusTime(text string) string {
	return sr.statusRGB(text, palette.Terracotta, color.Bold)
}

func (sr *Reporter) statusArrow(text string) string {
	return sr.statusRGB(text, palette.Ochre, color.Bold)
}

func (sr *Reporter) statusTransition(text string) string {
	return sr.statusRGB(text, palette.MutedWarm)
}

func (sr *Reporter) formatTransition(period t.Period) string {
	transition := sr.statusTransition(period.Transition.String())
	if period.Steps <= 0 {
		return "(" + transition + sr.statusSteps(" - no steps") + ")"
	}

	label := "steps"
	if period.Steps == 1 {
		label = "step"
	}

	return "(" + transition + sr.statusSteps(fmt.Sprintf(" - %d %s", period.Steps, label)) + ")"
}

func (sr *Reporter) statusSteps(text string) string {
	return sr.statusRGB(text, palette.Ochre)
}

func (sr *Reporter) statusTrack(text string) string {
	return sr.statusRGB(text, palette.Green)
}

func (sr *Reporter) statusFadeLabel(text string) string {
	return sr.statusRGB(text, palette.Ochre)
}

func (sr *Reporter) statusBullet(text string) string {
	return sr.statusRGB(text, palette.TerracottaDark, color.Bold)
}
