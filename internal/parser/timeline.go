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

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type ParsedTimelineDeclaration struct {
	Time       int
	PresetName string
	Transition t.TransitionType
	Steps      int
}

// parseTime parses a time string in HH:MM:SS format to milliseconds
func parseTime(s string) (int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format: expected HH:MM:SS")
	}

	for _, p := range parts {
		if len(p) != 2 {
			return 0, fmt.Errorf("invalid time format: each field must have 2 digits")
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

func (ctx *TextParser) ParseTimelineDeclaration() (*ParsedTimelineDeclaration, error) {
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "time")
	}
	timeSpan, _ := ctx.Line.LastTokenSpan()

	timeMs, err := parseTime(tok)
	if err != nil {
		return nil, diag.Parse(err.Error()).WithSpan(timeSpan).WithFound(tok).WithExpected("HH:MM:SS")
	}

	presetName, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "preset name")
	}

	decl := &ParsedTimelineDeclaration{
		Time:       timeMs,
		PresetName: presetName,
		Transition: t.TransitionSteady,
	}

	transition, ok := ctx.Line.NextToken()
	if ok {
		transitionSpan, _ := ctx.Line.LastTokenSpan()
		switch transition {
		case t.KeywordTransitionSteady:
			decl.Transition = t.TransitionSteady
		case t.KeywordTransitionEaseOut:
			decl.Transition = t.TransitionEaseOut
		case t.KeywordTransitionEaseIn:
			decl.Transition = t.TransitionEaseIn
		case t.KeywordTransitionSmooth:
			decl.Transition = t.TransitionSmooth
		default:
			return nil, diag.UnexpectedToken(transitionSpan, transition, t.KeywordTransitionSteady, t.KeywordTransitionEaseOut, t.KeywordTransitionEaseIn, t.KeywordTransitionSmooth)
		}

		if stepToken, ok := ctx.Line.Peek(); ok {
			if _, convErr := strconv.Atoi(stepToken); convErr == nil {
				decl.Steps, err = ctx.Line.NextIntStrict()
				if err != nil {
					return nil, err
				}
				if decl.Steps < 0 {
					stepSpan, _ := ctx.Line.LastTokenSpan()
					return nil, diag.Validation("steps must be non-negative").WithSpan(stepSpan).WithFound(stepToken)
				}
			}
		}
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token on timeline").WithSpan(unknownSpan).WithFound(unknown)
	}

	return decl, nil
}
