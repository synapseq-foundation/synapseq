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
	"bytes"
	"regexp"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiRegexp.ReplaceAllString(text, "")
}

func TestReporterDisplayPeriodChange_PrintsStartAndDash(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.Steps = 3

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endEqual := start
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endEqual

	view := View{Periods: []t.Period{p0, p1}, Channels: make([]t.Channel, t.NumberOfChannels)}

	var buf bytes.Buffer
	sr := NewReporter(&buf, false)
	sr.DisplayPeriodChange(view, 0)
	out := buf.String()

	if !strings.Contains(out, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+" - 3 steps)") {
		ts.Fatalf("missing start time line: %q", out)
	}
	if !strings.Contains(out, start.String()) {
		ts.Fatalf("missing start track string in output: %q", out)
	}
}

func TestReporterDisplayPeriodChange_ShowsEndTrackWhenChanged(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endChanged := start
	endChanged.Amplitude = t.AmplitudePercentToRaw(20)
	if IsTrackEqual(&start, &endChanged) {
		ts.Fatalf("precondition failed: start and end should not be equal")
	}
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endChanged

	view := View{Periods: []t.Period{p0, p1}, Channels: make([]t.Channel, t.NumberOfChannels)}
	var buf bytes.Buffer
	sr := NewReporter(&buf, false)
	sr.DisplayPeriodChange(view, 0)
	out := buf.String()

	if strings.Contains(out, "\n       --") {
		ts.Fatalf("did not expect '--' when start!=end: %q", out)
	}
	if !strings.Contains(out, endChanged.String()) {
		ts.Fatalf("missing end track string when changed: %q", out)
	}
}

func TestReporterCheckPeriodChange_DetectsTransitions(ts *testing.T) {
	var p0, p1, p2 t.Period
	p0.Time = 0
	p1.Time = 1000
	p2.Time = 2000
	tr := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = tr
	p0.TrackEnd[0] = tr
	p1.TrackStart[0] = tr
	p1.TrackEnd[0] = tr

	view := View{Periods: []t.Period{p0, p1, p2}, Channels: make([]t.Channel, t.NumberOfChannels)}
	var buf bytes.Buffer
	sr := NewReporter(&buf, false)
	sr.CheckPeriodChange(view, 0)

	out1 := buf.String()
	if !strings.Contains(out1, "- "+p0.TimeString()) {
		ts.Fatalf("expected period 0 output on first check: %q", out1)
	}

	buf.Reset()
	sr.CheckPeriodChange(view, 0)
	out2 := buf.String()
	if out2 != "" {
		ts.Fatalf("expected no output when period index unchanged, got: %q", out2)
	}

	buf.Reset()
	sr.CheckPeriodChange(view, 1)
	out3 := buf.String()
	if !strings.Contains(out3, "- "+p1.TimeString()) {
		ts.Fatalf("expected period 1 output after change: %q", out3)
	}
}

func TestReporterDisplayPeriodChange_UsesANSIWhenEnabled(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	p0.Steps = 1

	track := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = track
	p0.TrackEnd[0] = track

	view := View{Periods: []t.Period{p0, p1}, Channels: make([]t.Channel, t.NumberOfChannels)}

	var buf bytes.Buffer
	sr := NewReporter(&buf, true)
	sr.DisplayPeriodChange(view, 0)
	out := buf.String()

	if !strings.Contains(out, "\x1b[") {
		ts.Fatalf("expected ANSI colors in output, got: %q", out)
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+" - 1 step)") {
		ts.Fatalf("unexpected plain output after stripping ANSI: %q", plain)
	}
}

func TestReporterDisplayPeriodChange_ShowsNoStepsWhenZero(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	track := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = track
	p0.TrackEnd[0] = track

	view := View{Periods: []t.Period{p0, p1}, Channels: make([]t.Channel, t.NumberOfChannels)}
	var buf bytes.Buffer
	sr := NewReporter(&buf, false)
	sr.DisplayPeriodChange(view, 0)
	out := buf.String()

	if !strings.Contains(out, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+" - no steps)") {
		ts.Fatalf("expected no-steps label in output: %q", out)
	}
}

func TestReporterDisplayStatus_UsesChannelView(ts *testing.T) {
	channels := make([]t.Channel, t.NumberOfChannels)
	channels[0].Track = t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	channels[1].Track = t.Track{Type: t.TrackPinkNoise, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine}
	view := View{Channels: channels}

	var buf bytes.Buffer
	sr := NewReporter(&buf, false)
	sr.DisplayStatus(view, 65_000)
	out := buf.String()

	if !strings.Contains(out, "00:01:05") {
		ts.Fatalf("expected formatted time in output: %q", out)
	}
	if !strings.Contains(out, channels[0].Track.ShortString()) {
		ts.Fatalf("expected first channel track in output: %q", out)
	}
	if !strings.Contains(out, channels[1].Track.ShortString()) {
		ts.Fatalf("expected second channel track in output: %q", out)
	}
}

func TestCountActiveChannels(ts *testing.T) {
	tests := []struct {
		name     string
		channels []t.Channel
		expected int
	}{
		{"empty slice -> at least 1", []t.Channel{}, 1},
		{"all off -> 1", make([]t.Channel, 5), 1},
		{"single active at 0 -> 1", func() []t.Channel { channels := make([]t.Channel, 4); channels[0].Track.Type = t.TrackBinauralBeat; return channels }(), 1},
		{"last active at end -> len", func() []t.Channel { channels := make([]t.Channel, 4); channels[3].Track.Type = t.TrackPinkNoise; return channels }(), 4},
		{"last active in the middle -> index+1", func() []t.Channel { channels := make([]t.Channel, 5); channels[2].Track.Type = t.TrackBrownNoise; return channels }(), 3},
		{"multiple actives -> last index+1", func() []t.Channel { channels := make([]t.Channel, 8); channels[1].Track.Type = t.TrackBinauralBeat; channels[6].Track.Type = t.TrackAmbiance; return channels }(), 7},
		{"all active -> len", func() []t.Channel { channels := make([]t.Channel, 7); for i := range channels { channels[i].Track.Type = t.TrackPinkNoise }; return channels }(), 7},
		{"last off but previous active", func() []t.Channel { channels := make([]t.Channel, 6); channels[4].Track.Type = t.TrackAmbiance; return channels }(), 5},
	}

	for _, test := range tests {
		got := CountActiveChannels(test.channels)
		if got != test.expected {
			ts.Errorf("%s: expected %d, got %d", test.name, test.expected, got)
		}
	}
}

func TestIsTrackEqual(ts *testing.T) {
	base := &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}

	tests := []struct {
		name string
		a    *t.Track
		b    *t.Track
		eq   bool
	}{
		{"identical", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, true},
		{"different amplitude", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(30), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different carrier", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 320, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different resonance", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 12, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different waveform", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformTriangle, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"different intensity", base, &t.Track{Type: t.TrackBinauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(50)}}, false},
		{"different type", base, &t.Track{Type: t.TrackMonauralBeat, Carrier: 300, Resonance: 10, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)}}, false},
		{"ambiance effect type ignored", &t.Track{Type: t.TrackAmbiance, Carrier: 200, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPan, Intensity: t.IntensityPercentToRaw(60)}}, &t.Track{Type: t.TrackAmbiance, Carrier: 200, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(40), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectModulation, Intensity: t.IntensityPercentToRaw(60)}}, true},
	}

	for _, test := range tests {
		got := IsTrackEqual(test.a, test.b)
		if got != test.eq {
			ts.Errorf("%s: expected %v, got %v\nA=%+v\nB=%+v", test.name, test.eq, got, *test.a, *test.b)
		}
	}
}