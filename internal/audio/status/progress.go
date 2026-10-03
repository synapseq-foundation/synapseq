// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package status

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/palette"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
	"golang.org/x/term"
)

// Progress displays consumed audio frames, independently of render speed.
// All calls are made synchronously by the render loop.
type Progress struct {
	out         io.Writer
	colors      bool
	totalMs     int
	totalFrames int64
	terminal    bool
	size        func() (int, int)
	periods     []t.Period
	lastLines   []int
	now         func() time.Time
	lastUpdate  time.Time
}

func NewProgress(out io.Writer, colors bool, totalMs int, totalFrames int64, periods []t.Period) *Progress {
	p := &Progress{out: out, colors: colors, totalMs: totalMs, totalFrames: totalFrames, now: time.Now, periods: periods}
	if file, ok := out.(*os.File); ok && os.Getenv("TERM") != "dumb" && term.IsTerminal(int(file.Fd())) {
		p.terminal = true
		p.size = func() (int, int) {
			width, height, err := term.GetSize(int(file.Fd()))
			if err != nil || width <= 0 || height <= 0 {
				return 80, 24
			}
			return width, height
		}
	}
	return p
}

func (p *Progress) Start() {
	if p.out == nil {
		return
	}
	p.lastUpdate = p.now()
	if p.terminal {
		fmt.Fprintln(p.out)
		p.draw(0, false)
	} else {
		fmt.Fprintf(p.out, "Rendering audio: 0%%  %s / %s\n", progressTime(0), progressTime(p.totalMs))
	}
}

func (p *Progress) Update(frames int64) {
	if p.out == nil || !p.terminal {
		return
	}
	now := p.now()
	if now.Sub(p.lastUpdate) < 100*time.Millisecond {
		return
	}
	p.lastUpdate = now
	p.draw(frames, false)
}

func (p *Progress) Finish(frames int64, success bool) {
	if p.out == nil {
		return
	}
	complete := success && frames >= p.totalFrames
	if p.terminal {
		p.draw(frames, complete)
		return
	}
	percent, elapsed := p.values(frames, complete)
	label := "Audio rendering failed"
	if complete {
		label = "Audio rendered"
	}
	fmt.Fprintf(p.out, "%s: %d%%  %s / %s\n", label, percent, progressTime(elapsed), progressTime(p.totalMs))
}

func (p *Progress) values(frames int64, complete bool) (int, int) {
	if complete {
		return 100, p.totalMs
	}
	if p.totalFrames <= 0 {
		return 0, 0
	}
	ratio := min(1.0, max(0.0, float64(frames)/float64(p.totalFrames)))
	// Only Finish may announce completion.
	return min(99, int(ratio*100)), int(ratio * float64(p.totalMs))
}

// draw leaves the cursor on a blank row below the component. This also lets
// terminals reflow every component row on resize (the cursor row is often
// excluded from reflow). Only previously drawn rows are erased.
func (p *Progress) draw(frames int64, complete bool) {
	width, height := p.size()
	width = max(1, width)
	percent, elapsed := p.values(frames, complete)
	lines := p.layout(percent, elapsed, width, height)
	rows := 0
	for _, columns := range p.lastLines {
		rows += max(1, (columns+width-1)/width)
	}
	// Rows scrolled out by a height reduction are no longer addressable.
	rows = min(rows, max(0, height-1))
	var output strings.Builder
	if rows > 0 {
		for row := 0; row < rows; row++ {
			output.WriteString("\x1b[1A\r\x1b[2K")
		}
	}
	p.lastLines = p.lastLines[:0]
	for _, line := range lines {
		output.WriteString(line.text)
		output.WriteString("\r\n")
		p.lastLines = append(p.lastLines, line.columns)
	}
	fmt.Fprint(p.out, output.String())
}

type progressLine struct {
	text    string
	columns int
}

func (p *Progress) layout(percent, elapsed, width, height int) []progressLine {
	available := max(0, width-1)
	if width < 20 || height < 5 {
		plain := *p
		plain.colors = false
		return []progressLine{{p.line(percent, elapsed, available), utf8.RuneCountInString(plain.line(percent, elapsed, available))}}
	}
	if width < 40 {
		inner := max(0, available-2)
		name, columns := p.presetLine(elapsed, percent == 100, inner)
		bar := p.bar(percent, inner)
		return []progressLine{{" " + name, 1 + columns}, {" " + bar, 1 + inner}}
	}
	// Two outer margin columns per side, borders, and one padding column per side.
	boxWidth := available - 4
	inner := boxWidth - 4
	border := func(left, right string) progressLine {
		return progressLine{"  " + p.paint(left+strings.Repeat("─", boxWidth-2)+right, palette.MutedWarm), 2 + boxWidth}
	}
	row := func(text string, columns int) progressLine {
		return progressLine{"  " + p.paint("│", palette.MutedWarm) + " " + text + strings.Repeat(" ", max(0, inner-columns)) + " " + p.paint("│", palette.MutedWarm), 2 + boxWidth}
	}
	name, columns := p.presetLine(elapsed, percent == 100, inner-2)
	title := border("╭", "╮")
	if columns > 0 {
		remaining := boxWidth - 4 - columns
		left := remaining / 2
		right := remaining - left
		title.text = "  " + p.paint("╭"+strings.Repeat("─", left)+" ", palette.MutedWarm) + name + p.paint(" "+strings.Repeat("─", right)+"╮", palette.MutedWarm)
	}
	return []progressLine{title, row(p.bar(percent, inner), inner), border("╰", "╯")}

}

func (p *Progress) presetLine(elapsed int, complete bool, width int) (string, int) {
	current, next := p.presetNames(elapsed, complete)
	if current == next || next == "" || width < 5 {
		current = truncateProgress(current, width)
		return p.paint(current, palette.Terracotta), utf8.RuneCountInString(current)
	}
	left := (width - 3 + 1) / 2
	right := width - 3 - left
	// Give unused space from a short name to the other name.
	if len(current) < left {
		right += left - len(current)
		left = len(current)
	}
	if len(next) < right {
		left += right - len(next)
		right = len(next)
	}
	current = truncateProgress(current, left)
	next = truncateProgress(next, right)
	return p.paint(current, palette.Terracotta) + p.paint(" → ", palette.Ochre) + p.paint(next, palette.Green), utf8.RuneCountInString(current) + 3 + utf8.RuneCountInString(next)
}

func (p *Progress) presetNames(elapsed int, complete bool) (string, string) {
	if len(p.periods) == 0 {
		return "", ""
	}
	index := 0
	if complete {
		index = len(p.periods) - 1
	} else {
		for index+1 < len(p.periods) && elapsed >= p.periods[index+1].Time {
			index++
		}
	}
	current := p.periods[index].PresetName
	if index+1 < len(p.periods) {
		return current, p.periods[index+1].PresetName
	}
	return current, ""
}

func (p *Progress) bar(percent, width int) string {
	label := fmt.Sprintf("%d%%", percent)
	cells := max(0, width-len(label)-2)
	filled := cells * percent / 100
	return p.paint(strings.Repeat("█", filled), palette.Green) + p.paint(strings.Repeat("░", cells-filled), palette.MutedWarm) + "  " + p.paint(label, palette.Terracotta)
}

func truncateProgress(text string, width int) string {
	runes := []rune(text)
	if width <= 0 {
		return ""
	}
	if len(runes) <= width {
		return text
	}
	return string(runes[:width-1]) + "…"
}

func (p *Progress) line(percent, elapsed, width int) string {
	label := fmt.Sprintf("%d%%", percent)
	times := "  " + progressTime(elapsed) + " / " + progressTime(p.totalMs)
	barWidth := width - 3 - len(label) - len(times)
	if barWidth < 8 {
		times = ""
		barWidth = width - 3 - len(label)
	}
	if barWidth < 8 {
		if len(label) > width {
			label = label[:width]
		}
		return p.paint(label, palette.Terracotta)
	}
	filled := barWidth * percent / 100
	return "[" + p.paint(strings.Repeat("█", filled), palette.Green) +
		p.paint(strings.Repeat("░", barWidth-filled), palette.MutedWarm) + "] " +
		p.paint(label, palette.Terracotta) + p.paint(times, palette.MutedWarm)
}

func (p *Progress) paint(text string, token palette.RGBColor) string {
	if !p.colors {
		return text
	}
	style := color.RGB(token.R(), token.G(), token.B())
	style.EnableColor()
	return style.Sprint(text)
}

func progressTime(ms int) string {
	seconds := max(0, ms) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", seconds/3600, seconds/60%60, seconds%60)
}
