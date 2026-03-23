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

package manual

import (
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
	clistyle "github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}

func TestRenderIncludesCoreSections(ts *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	clistyle.SetColorEnabled(false)
	manual := Render()

	checks := []string{
		"SYNAPSEQ(1)",
		"NAME",
		"SYNOPSIS",
		"DESCRIPTION",
		"OPTIONS",
		"SEQUENCE FILE",
		"COMPATIBILITY",
		"COMMON ERRORS",
		"synapseq [OPTION]... INPUT [OUTPUT]",
		"synapseq -hub-clean",
		"synapseq -hub-download NAME [DIR]",
		"synapseq -hub-get NAME [OUTPUT]",
		"Paged reading",
		"-new TYPE",
		"-test",
		"-preview",
		"-play",
		"-quiet",
		"-no-color",
		"-manual",
		"-help",
		"-version",
		"-hub-update",
		"-hub-list",
		"-hub-search WORD",
		"-hub-info NAME",
		"-ffmpeg-path PATH",
		"-ffplay-path PATH",
		"-install-file-association",
		"-uninstall-file-association",
		"line-oriented score language",
		"Timeline entries must appear last.",
		"INPUT may be a local .spsq file",
		"Rules",
		"Global options",
		"Identifiers",
		"Preset declarations",
		"Track syntax",
		"binaural",
		"isochronic",
		"pan",
		"doppler",
		"Track overrides",
		"smooth VALUE",
		"Timeline",
		"steady",
		"ease-in",
		"ease-out",
		"beat mode",
		"track kind",
		"effect type",
		"ambiance source",
		"Extended files",
		"SEE ALSO",
		"synapseq -manual | less",
		"synapseq -manual | more",
		"plain-text .spsq",
		"VALUE",
		"Allowed",
		"Bridge",
		"Rules",
		"Inline comment",
		"Option after preset or timeline",
		"New track in inherited preset",
		"Invalid local path",
		"HH MM SS required",
		"tone only with tone",
		"same source only",
	}

	for _, check := range checks {
		if !strings.Contains(manual, check) {
			ts.Fatalf("expected manual to contain %q, got:\n%s", check, manual)
		}
	}
	if strings.Contains(manual, "\x1b[") {
		ts.Fatalf("expected manual without ANSI colors when color is disabled, got: %q", manual)
	}

	removed := []string{
		"EXAMPLES",
		"Basic session",
		"Reusable templates",
		"00:00:40 silence",
		"@extends library/focus-base",
		"@ambiance rain audio/rain",
		"# library/common.spsc",
		"TRACK KEYWORDS",
		"focus-base as template",
		"track 3 smooth 40",
		"@extends library/common",
		"@samplerate 48000 # samplerate",
		"The example above is invalid.",
		"COMMAND LINE",
		"SEQUENCE OPTIONS",
		"PRESET NAMES",
		"PRESETS",
		"Compatibility rule",
		"Silence bridge",
		"Standalone lines only",
		"Track Definitions",
		"->",
		":",
		"00:00:00",
		"HH:MM:SS",
	}

	for _, check := range removed {
		if strings.Contains(manual, check) {
			ts.Fatalf("expected manual to omit %q, got:\n%s", check, manual)
		}
	}
}

func TestRenderEmitsANSIWhenEnabled(ts *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	clistyle.SetColorEnabled(true)
	manual := Render()

	if !strings.Contains(manual, "\x1b[") {
		ts.Fatalf("expected ANSI sequences when color is enabled, got: %q", manual)
	}

	plain := stripANSI(manual)
	if !strings.Contains(plain, "OPTIONS") {
		ts.Fatalf("expected stripped manual to keep text content, got: %q", plain)
	}
}
