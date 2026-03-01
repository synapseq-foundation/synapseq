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

package parser

import (
	"fmt"
	"strconv"
	"strings"

	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// parseTime parses a time string in HH:MM:SS format to milliseconds
func parseTime(s string) (int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format (must be HH:MM:SS): %s", s)
	}

	for _, p := range parts {
		if len(p) != 2 {
			return 0, fmt.Errorf("each field must have 2 digits: %s", s)
		}
	}

	hh, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", parts[0])
	}
	mm, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", parts[1])
	}
	ss, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid second: %s", parts[2])
	}

	if ss < 0 || mm < 0 || hh < 0 || ss > 59 || mm > 59 || hh > 23 {
		return 0, fmt.Errorf("invalid time value: %s", s)
	}

	return (hh*3600 + mm*60 + ss) * 1000, nil
}

// HasTimeline checks if the current line is a timeline entry
func (ctx *TextParser) HasTimeline() bool {
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	if ctx.Line.Raw[0] == ' ' {
		return false
	}

	if _, err := parseTime(tok); err != nil {
		return false
	}

	return true
}

// ParseTimeline parses a timeline line and returns a Period
func (ctx *TextParser) ParseTimeline(presets *[]t.Preset) (*t.Period, error) {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected time, got EOF: %s", ln)
	}

	timeMs, err := parseTime(tok)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	tok, ok = ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected preset name, got EOF: %s", ln)
	}

	// default transition type
	transitionType := t.TransitionSteady
	transition, ok := ctx.Line.NextToken()
	if ok {
		switch transition {
		case t.KeywordTransitionSteady:
			transitionType = t.TransitionSteady
		case t.KeywordTransitionEaseOut:
			transitionType = t.TransitionEaseOut
		case t.KeywordTransitionEaseIn:
			transitionType = t.TransitionEaseIn
		case t.KeywordTransitionSmooth:
			transitionType = t.TransitionSmooth
		default:
			return nil, fmt.Errorf("unknown transition mode %q: %s", transition, ln)
		}
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token on timeline %q: %s", unknown, ln)
	}

	p := s.FindPreset(strings.ToLower(tok), *presets)
	if p == nil {
		return nil, fmt.Errorf("preset %q not found: %s", tok, ln)
	}

	if p.IsTemplate {
		return nil, fmt.Errorf("cannot use template preset %q in timeline: %s", p.String(), ln)
	}

	period := &t.Period{
		Time:       timeMs,
		TrackStart: p.Track,
		TrackEnd:   p.Track,
		Transition: transitionType,
	}

	return period, nil
}
