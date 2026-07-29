// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"strings"
	"testing"
	"time"
)

func TestParseSequence(t *testing.T) {
	input := strings.Join([]string{
		"-SE -m audio/background.mp3",
		"alpha: pink/40 300+10/20 200/5 spin:300+4.2/10 mix/25",
		"off: -",
		"NOW+00:00:15 <> alpha ->",
		"+00:01:00 == off",
	}, "\n")

	parsed, err := parse("test.sbg", strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed.musicPath != "audio/background.mp3" || parsed.musicLine != 1 {
		t.Fatalf("music source = %q at line %d", parsed.musicPath, parsed.musicLine)
	}
	if len(parsed.definitions) != 2 || len(parsed.definitions[0].voices) != 5 {
		t.Fatalf("definitions = %#v", parsed.definitions)
	}
	if got := parsed.definitions[0].voices[3]; got.kind != voiceSpin || got.carrier != 300 || got.beat != 4.2 || got.amplitude != 10 {
		t.Fatalf("spin voice = %#v", got)
	}
	if got := parsed.timeline[0]; !got.initial || got.at != 15*time.Second || got.name != "alpha" {
		t.Fatalf("initial timeline = %#v", got)
	}
}

func TestParseRejectsUndefinedTimelineName(t *testing.T) {
	_, err := parse("test.sbg", strings.NewReader("alpha: 300+10/20\nNOW missing\n"))
	if err == nil || !strings.Contains(err.Error(), "undefined NameDef") {
		t.Fatalf("error = %v", err)
	}
}

func TestParseAbsoluteTimelineTimes(t *testing.T) {
	input := "off: -\nts1: 300+10/20\nalpha: 300+10/20\n00:00:00 off ->\n00:00:15 ts1\n00:20:00 alpha\n"
	parsed, err := parse("test.sbg", strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if got := parsed.timeline[2].at; got != 20*time.Minute {
		t.Fatalf("absolute timeline time = %v", got)
	}
}

func TestParseRejectsNegativeTimelineFields(t *testing.T) {
	input := "alpha: 300+10/20\n-1:00 alpha\n"
	_, err := parse("test.sbg", strings.NewReader(input))
	if err == nil || !strings.Contains(err.Error(), "invalid timeline time") {
		t.Fatalf("error = %v", err)
	}
}
