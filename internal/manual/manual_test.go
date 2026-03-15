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
		"SynapSeq Manual",
		"File Layout",
		"Options",
		"Presets",
		"Track Definitions",
		"Sound Concepts",
		"binaural",
		"isochronic",
		"pan",
		"doppler",
		"Track Overrides",
		"Timeline",
		"steady",
		"ease-in",
		"ease-out",
		"smooth",
		"Compatibility rule",
		"Silence bridge",
		"incompatible track type, effect type, or ambiance source",
		"Extended Files",
		"Advanced Examples",
		"synapseq -manual",
		"@ambiance rain audio/rain",
		"00:00:00 silence",
		"@extends library/focus-base",
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
	if !strings.Contains(plain, "Command Line") {
		ts.Fatalf("expected stripped manual to keep text content, got: %q", plain)
	}
}
