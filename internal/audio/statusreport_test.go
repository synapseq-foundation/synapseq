/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiRegexp.ReplaceAllString(text, "")
}

func TestStatusReporter_DisplayPeriodChange_PrintsStartAndDash(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endEqual := start
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endEqual

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{}}

	var buf bytes.Buffer
	sr := NewStatusReporter(&buf, false)
	sr.DisplayPeriodChange(r, 0)
	out := buf.String()

	if !strings.Contains(out, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+")") {
		ts.Fatalf("missing start time line: %q", out)
	}
	// We no longer print the end time when start==end
	// if !strings.Contains(out, "  "+p1.TimeString()) {
	// 	ts.Fatalf("missing end time line: %q", out)
	// }
	if !strings.Contains(out, start.String()) {
		ts.Fatalf("missing start track string in output: %q", out)
	}
	// if !strings.Contains(out, "\n       --") {
	// 	ts.Fatalf("expected '--' marker when start==end: %q", out)
	// }
}

func TestStatusReporter_DisplayPeriodChange_ShowsEndTrackWhenChanged(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endChanged := start
	endChanged.Amplitude = t.AmplitudePercentToRaw(20)
	// sanity: ensure IsTrackEqual detects difference
	if s.IsTrackEqual(&start, &endChanged) {
		ts.Fatalf("precondition failed: start and end should not be equal")
	}
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endChanged

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{}}
	var buf bytes.Buffer
	sr := NewStatusReporter(&buf, false)
	sr.DisplayPeriodChange(r, 0)
	out := buf.String()

	if strings.Contains(out, "\n       --") {
		ts.Fatalf("did not expect '--' when start!=end: %q", out)
	}
	if !strings.Contains(out, endChanged.String()) {
		ts.Fatalf("missing end track string when changed: %q", out)
	}
}

func TestStatusReporter_CheckPeriodChange_DetectsTransitions(ts *testing.T) {
	var p0, p1, p2 t.Period
	p0.Time = 0
	p1.Time = 1000
	p2.Time = 2000
	tr := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = tr
	p0.TrackEnd[0] = tr
	p1.TrackStart[0] = tr
	p1.TrackEnd[0] = tr

	r := &AudioRenderer{periods: []t.Period{p0, p1, p2}, AudioRendererOptions: &AudioRendererOptions{}}
	var buf bytes.Buffer
	sr := NewStatusReporter(&buf, false)
	sr.CheckPeriodChange(r, 0)

	out1 := buf.String()
	if !strings.Contains(out1, "- "+p0.TimeString()) {
		ts.Fatalf("expected period 0 output on first check: %q", out1)
	}

	buf.Reset()
	sr.CheckPeriodChange(r, 0)
	out2 := buf.String()
	if out2 != "" {
		ts.Fatalf("expected no output when period index unchanged, got: %q", out2)
	}

	buf.Reset()
	sr.CheckPeriodChange(r, 1)
	out3 := buf.String()
	if !strings.Contains(out3, "- "+p1.TimeString()) {
		ts.Fatalf("expected period 1 output after change: %q", out3)
	}
}

func TestStatusReporter_DisplayPeriodChange_UsesANSIWhenEnabled(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	track := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = track
	p0.TrackEnd[0] = track

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{Colors: true}}

	var buf bytes.Buffer
	sr := NewStatusReporter(&buf, true)
	sr.DisplayPeriodChange(r, 0)
	out := buf.String()

	if !strings.Contains(out, "\x1b[") {
		ts.Fatalf("expected ANSI colors in output, got: %q", out)
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+")") {
		ts.Fatalf("unexpected plain output after stripping ANSI: %q", plain)
	}
}
