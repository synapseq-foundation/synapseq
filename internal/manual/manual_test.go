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
	style "github.com/synapseq-foundation/synapseq/v4/internal/textstyle"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}

func TestRenderIncludesCoreSections(ts *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	style.SetColorEnabled(false)
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
		"synapseq -new TYPE [OUTPUT]",
		"synapseq -preview INPUT [OUTPUT]",
		"synapseq -mp3 INPUT [OUTPUT]",
		"synapseq -hub-clean",
		"synapseq -hub-download NAME [DIR]",
		"synapseq -hub-get NAME [OUTPUT]",
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
		"synapseq -manual -no-color",
		"line-oriented .spsq files",
		"INPUT may be a local .spsq file",
		"Rules",
		"Render modes",
		"Creation",
		"Tools and system",
		"Information",
		"Order",
		"Sequence Options",
		"Syntax",
		"Identifiers",
		"Preset declarations",
		"Track syntax",
		"binaural",
		"monaural",
		"isochronic",
		"pan",
		"doppler",
		"Track overrides",
		"VALUE may be absolute or signed relative delta",
		"+VALUE adds to inherited value",
		"-VALUE subtracts from inherited value",
		"track N amplitude VALUE",
		"track N waveform sine|square|triangle|sawtooth",
		"@samplerate NUMBER",
		"@volume NUMBER",
		"@ambiance NAME PATH_OR_URL",
		"@extends PATH_OR_URL",
		"tone CARRIER amplitude LEVEL",
		"tone CARRIER binaural|monaural|isochronic BEAT amplitude LEVEL",
		"noise white|pink|brown [smooth VALUE] amplitude LEVEL",
		"ambiance NAME amplitude LEVEL",
		"smooth VALUE",
		"effect pan VALUE intensity PERCENT",
		"effect modulation VALUE intensity PERCENT",
		"effect doppler VALUE intensity PERCENT",
		"Remaining characters may be:",
		"letters",
		"digits",
		"underscores",
		"dashes",
		"track overrides",
		"timeline entries",
		"nested @extends",
		"waveform VALUE",
		"VALUE",
		"tone tracks",
		"ambiance tracks",
		"Timeline",
		"HH:MM:SS PRESET_NAME [TRANSITION [STEPS]]",
		"steady",
		"ease-in",
		"ease-out",
		"steps require explicit transition",
		"steps 0 keeps the normal transition",
		"limited by 5 seconds per leg, with a hard cap of 12",
		"beat mode",
		"track kind",
		"effect type",
		"ambiance source",
		"Extended files",
		"plain-text .spsq",
		"all top-level lines start in column 1",
		"# comment: ignored",
		"## comment: printed unless -quiet is set",
		"waveform prefix allowed on:",
		"effect appears before amplitude",
		"top-level lines only; HH:MM:SS required",
		"Allowed",
		"Bridge",
		"Inline comment",
		"Option after preset or timeline",
		"New track in inherited preset",
		"Invalid local path",
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
		"SOUND CONCEPTS",
		"Paged reading",
		"SEE ALSO",
		"Preset example",
		"Inherited preset example",
		"A pair of nearby tones, one per ear",
		"A beat created by mixing the tones before playback",
		"A single carrier that is gated on and off at the beat rate",
		"synapseq -manual | less",
		"synapseq -manual | more",
		"Timeline examples",
		"00:00:00 silence",
		"00:00:30 focus ease-in",
		"00:10:00 focus-deep smooth 3",
		"00:20:00 silence ease-out",
		"parsed top to bottom",
		"sequence options must start in column 1",
		"preset declarations must start in column 1",
		"timeline entries must start in column 1",
		"# comments are ignored",
		"## comments are printed unless -quiet is set",
		"remaining characters may be letters digits underscores dashes",
		"allowed: options presets tracks track overrides",
		"not permitted: timeline entries nested @extends",
		"tone: pan modulation doppler",
		"effect TYPE VALUE intensity PERCENT appears before amplitude",
		"EXAMPLES",
		"Basic session",
		"Reusable templates",
		"TRACK KEYWORDS",
		"track 3 smooth 40",
		"@samplerate 48000 # samplerate",
		"The example above is invalid.",
		"COMMAND LINE",
		"PRESET NAMES",
		"Compatibility rule",
		"Silence bridge",
		"Standalone lines only",
		"Track Definitions",
		"->",
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

	style.SetColorEnabled(true)
	manual := Render()

	if !strings.Contains(manual, "\x1b[") {
		ts.Fatalf("expected ANSI sequences when color is enabled, got: %q", manual)
	}

	plain := stripANSI(manual)
	if !strings.Contains(plain, "OPTIONS") {
		ts.Fatalf("expected stripped manual to keep text content, got: %q", plain)
	}
}
