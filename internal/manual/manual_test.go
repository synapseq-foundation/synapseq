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
		"synapseq [OPTION]... INPUT [OUTPUT]",
		"synapseq -hub-download NAME [DIR]",
		"synapseq -hub-get NAME [OUTPUT]",
		"FILE LAYOUT",
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
		"-hub-clean",
		"-hub-list",
		"-hub-search WORD",
		"-hub-info NAME",
		"-ffmpeg-path PATH",
		"-ffplay-path PATH",
		"-install-file-association",
		"-uninstall-file-association",
		"repeatable",
		"ambient listening sessions",
		"Options",
		"SEQUENCE OPTIONS",
		"Presets",
		"TRACK DEFINITIONS",
		"SOUND CONCEPTS",
		"binaural",
		"isochronic",
		"pan",
		"doppler",
		"TRACK OVERRIDES",
		"TIMELINE",
		"steady",
		"ease-in",
		"ease-out",
		"smooth",
		"Compatibility rule",
		"Silence bridge",
		"incompatible track type",
		"effect type",
		"ambiance source",
		"EXTENDED FILES",
		"EXAMPLES",
		"SEE ALSO",
		"synapseq -manual",
		"synapseq -manual | less",
		"synapseq -manual | more",
		"Basic session",
		"00:00:40 silence",
		"Reusable templates",
		"Track positions are fixed",
		"Structural changes need silence",
		"Beat, noise, and effect types must match",
		"Waveforms may change",
		"@ambiance rain audio/rain",
		"00:00:00 silence",
		"# library/common.spsc",
		"focus-template as template",
		"@extends library/common",
		"@extends library/focus-base",
		"Standalone lines only",
		"@samplerate 48000 # samplerate",
		"The example above is invalid.",
	}

	for _, check := range checks {
		if !strings.Contains(manual, check) {
			ts.Fatalf("expected manual to contain %q, got:\n%s", check, manual)
		}
	}
	if strings.Contains(manual, "\x1b[") {
		ts.Fatalf("expected manual without ANSI colors when color is disabled, got: %q", manual)
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
	if !strings.Contains(plain, "COMMAND LINE") {
		ts.Fatalf("expected stripped manual to keep text content, got: %q", plain)
	}
}
