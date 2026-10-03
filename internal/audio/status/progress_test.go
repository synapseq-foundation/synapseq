// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package status

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestProgressLayout(t *testing.T) {
	p := NewProgress(nil, false, 600000, 100, nil)
	for _, width := range []int{80, 48, 25, 12, 3, 1, 0} {
		line := p.line(40, 240000, width)
		if utf8.RuneCountInString(line) > width {
			t.Fatalf("width %d: %q", width, line)
		}
		if width == 80 && (!strings.Contains(line, "00:04:00 / 00:10:00") || !strings.Contains(line, "█")) {
			t.Fatal(line)
		}
		if width == 25 && strings.Contains(line, "00:") {
			t.Fatal(line)
		}
	}
	p.colors = true
	if !strings.Contains(p.line(40, 240000, 80), "\x1b[") {
		t.Fatal("missing colors")
	}
}

func TestProgressLifecycle(t *testing.T) {
	for _, success := range []bool{true, false} {
		var out bytes.Buffer
		p := NewProgress(&out, true, 1000, 100, nil)
		p.Start()
		p.Update(50)
		p.Finish(100, success)
		got := out.String()
		if strings.Count(got, "\n") != 2 || strings.ContainsAny(got, "\r\x1b") {
			t.Fatal(got)
		}
		if success != strings.Contains(got, "100%") {
			t.Fatal(got)
		}
	}
	var out bytes.Buffer
	p := NewProgress(&out, false, 0, 0, nil)
	p.Start()
	p.Finish(0, true)
	if !strings.Contains(out.String(), "100%") {
		t.Fatal(out.String())
	}
}

func TestProgressThrottleAndResize(t *testing.T) {
	var out bytes.Buffer
	p := NewProgress(&out, false, 1000, 100, nil)
	now := time.Unix(1, 0)
	width := 80
	p.now = func() time.Time { return now }
	p.terminal = true
	p.size = func() (int, int) { return width, 24 }
	p.Start()
	initial := out.Len()
	now = now.Add(99 * time.Millisecond)
	p.Update(20)
	if out.Len() != initial {
		t.Fatal("updated too early")
	}
	width = 10
	now = now.Add(time.Millisecond)
	p.Update(50)
	if !strings.HasSuffix(out.String(), "\r\x1b[2K50%\r\n") {
		t.Fatal(out.String())
	}
	p.Finish(100, true)
	if !strings.HasSuffix(out.String(), "100%\r\n") {
		t.Fatal(out.String())
	}
}

func TestProgressDumbTerminal(t *testing.T) {
	t.Setenv("TERM", "dumb")
	var out bytes.Buffer
	p := NewProgress(&out, true, 1000, 100, nil)
	if p.terminal {
		t.Fatal("dumb terminal animated")
	}
}

func TestProgressContainerGeometry(tst *testing.T) {
	p := NewProgress(nil, false, 600000, 100, []t.Period{{PresetName: "focus-start"}, {Time: 300000, PresetName: "focus-deep"}, {Time: 600000, PresetName: "silence"}})
	for _, tc := range []struct{ width, height, rows int }{
		{80, 24, 3}, {40, 5, 3}, {80, 4, 1}, {39, 24, 2}, {20, 7, 2}, {19, 24, 1}, {80, 6, 3}, {1, 1, 1},
	} {
		lines := p.layout(40, 240000, tc.width, tc.height)
		if len(lines) != tc.rows {
			tst.Fatalf("%+v: %v", tc, lines)
		}
		for _, line := range lines {
			columns := utf8.RuneCountInString(line.text)
			if columns != line.columns || columns > tc.width-1 {
				tst.Fatalf("%+v: %+v", tc, line)
			}
			if tc.rows == 3 && (!strings.HasPrefix(line.text, "  ") || columns != tc.width-3) {
				tst.Fatalf("margins: %+v", line)
			}
		}
		if tc.rows == 3 {
			if !strings.HasPrefix(lines[1].text, "  │ ") || !strings.HasSuffix(lines[1].text, " │") {
				tst.Fatal(lines[1])
			}
			title := lines[0].text
			name, _ := p.presetLine(240000, false, tc.width-11)
			if !strings.Contains(title, " "+name+" ") || strings.Contains(lines[1].text, "00:") {
				tst.Fatal(lines)
			}
			parts := strings.Split(title, " "+name+" ")
			left := utf8.RuneCountInString(strings.TrimPrefix(parts[0], "  ╭"))
			right := utf8.RuneCountInString(strings.TrimSuffix(parts[1], "╮"))
			if right-left < 0 || right-left > 1 {
				tst.Fatal("uncentered title", title)
			}

		}
	}
}

func TestProgressPresetTimeline(tst *testing.T) {
	p := NewProgress(nil, false, 3000, 300, []t.Period{
		{Time: 0, PresetName: "start"}, {Time: 1000, PresetName: "start"}, {Time: 2000, PresetName: "deep"}, {Time: 3000, PresetName: "silence"},
	})
	for _, tc := range []struct {
		ms       int
		complete bool
		want     string
	}{
		{0, false, "start"}, {999, false, "start"}, {1000, false, "start → deep"}, {2000, false, "deep → silence"}, {3000, true, "silence"},
	} {
		got, _ := p.presetLine(tc.ms, tc.complete, 50)
		if got != tc.want {
			tst.Fatalf("%+v: %s", tc, got)
		}
	}
	p.periods[0].PresetName = "abcdefghijklmnopqrst"
	p.periods[1].PresetName = "zyxwvutsrqponmlkjihg"
	for width := 0; width < 45; width++ {
		text, columns := p.presetLine(0, false, width)
		if columns > width || columns != utf8.RuneCountInString(text) {
			tst.Fatalf("width %d: %s", width, text)
		}
		if width < 5 && strings.Contains(text, "→") {
			tst.Fatal(text)
		}
	}
	text, _ := p.presetLine(0, false, 15)
	if text != "abcde… → zyxwv…" {
		tst.Fatal(text)
	}
}

func TestProgressContainerColors(tst *testing.T) {
	p := NewProgress(nil, true, 1000, 100, []t.Period{{PresetName: "a"}, {Time: 1000, PresetName: "b"}})
	colored := p.layout(40, 400, 80, 24)
	p.colors = false
	plain := p.layout(40, 400, 80, 24)
	ansi := regexp.MustCompile("\x1b\\[[0-9;]*m")
	for i, line := range colored {
		if ansi.ReplaceAllString(line.text, "") != plain[i].text || !strings.Contains(line.text, "\x1b[") {
			tst.Fatal(line)
		}
	}
}

func TestProgressClearsPreviousRegion(tst *testing.T) {
	var out bytes.Buffer
	p := NewProgress(&out, false, 3000, 300, []t.Period{{PresetName: "a"}, {Time: 3000, PresetName: "final"}})
	width, height := 80, 24
	p.terminal = true
	p.size = func() (int, int) { return width, height }
	p.Start()
	if !strings.HasPrefix(out.String(), "\n  ╭") {
		tst.Fatal(out.String())
	}
	// Each old 77-column box row reflows into two rows at width 39.
	out.Reset()
	width = 39
	p.draw(100, false)
	if strings.Count(out.String(), "\x1b[2K") != 6 || strings.Count(out.String(), "\x1b[1A") != 6 {
		tst.Fatal(out.String())
	}
	out.Reset()
	width = 19
	p.draw(150, false)
	if strings.Count(out.String(), "\x1b[2K") != 3 {
		tst.Fatal(out.String())
	}
	out.Reset()
	width = 80
	p.Finish(300, true)
	if strings.Count(out.String(), "\x1b[2K") != 1 || !strings.Contains(out.String(), "final") || !strings.HasSuffix(out.String(), "╯\r\n") {
		tst.Fatal(out.String())
	}
	if strings.HasPrefix(out.String(), "\n") {
		tst.Fatal("repeated top margin")
	}
}
