// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	nr "github.com/synapseq-foundation/synapseq/v4/internal/nameref"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HasOption checks if the first element is an option.
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and returns raw parsed option values.
func (ctx *TextParser) ParseOption(_ string) (*t.ParseOptions, error) {
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "option")
	}
	span, _ := ctx.Line.LastTokenSpan()

	if string(tok[0]) != t.KeywordOption {
		return nil, diag.Parse("expected option").WithSpan(span).WithFound(tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return nil, diag.Parse("expected option name").WithSpan(span).WithFound(tok)
	}

	parsed := t.NewParseOptions()
	validOptions := []string{
		t.KeywordOptionSampleRate,
		t.KeywordOptionVolume,
		t.KeywordOptionAmbiance,
		t.KeywordOptionMusic,
		t.KeywordOptionExtends,
		t.KeywordOptionWaveform,
		t.KeywordOptionTransition,
	}

	switch option {
	case t.KeywordOptionSampleRate:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "samplerate value")
		}

		parsed.Values[t.KeywordOptionSampleRate] = value
	case t.KeywordOptionVolume:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "volume value")
		}

		parsed.Values[t.KeywordOptionVolume] = value
	case t.KeywordOptionAmbiance:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "ambiance name")
		}
		nameSpan, _ := ctx.Line.LastTokenSpan()

		if err := nr.IsValid(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			content = name // allow shorthand ambiance name as path
		}

		parsed.Ambiance[name] = content
	case t.KeywordOptionMusic:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "music name")
		}
		nameSpan, _ := ctx.Line.LastTokenSpan()

		if err := nr.IsValid(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			content = name
		}

		parsed.Music[name] = content
	case t.KeywordOptionExtends:
		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "extends path")
		}

		parsed.Extends = append(parsed.Extends, content)
	case t.KeywordOptionWaveform:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "waveform name")
		}
		nameSpan, _ := ctx.Line.LastTokenSpan()
		if err := nr.IsValid(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}
		if t.IsBuiltinWaveformName(name) {
			return nil, diag.Validation(fmt.Sprintf("waveform name %q is reserved for a built-in waveform", name)).WithSpan(nameSpan).WithFound(name)
		}

		points := make([]float64, 0)
		for {
			_, ok := ctx.Line.Peek()
			if !ok {
				break
			}
			value, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			valueSpan, _ := ctx.Line.LastTokenSpan()
			if value < 0 || value > 100 {
				return nil, diag.Validation("waveform points must be between 0 and 100").WithSpan(valueSpan).WithFound(fmt.Sprintf("%g", value))
			}
			if len(points) == t.MaxWaveformPoints {
				return nil, diag.Validation(fmt.Sprintf("waveform cannot contain more than %d points", t.MaxWaveformPoints)).WithSpan(valueSpan)
			}
			points = append(points, t.NormalizeWaveformPoint(value))
		}
		if len(points) < t.MinWaveformPoints {
			return nil, diag.Validation(fmt.Sprintf("waveform must contain at least %d points", t.MinWaveformPoints)).WithSpan(nameSpan).WithFound(name)
		}
		parsed.Waveforms = append(parsed.Waveforms, t.WaveformDefinition{
			Name:   t.WaveformName(name),
			Points: points,
		})
	case t.KeywordOptionTransition:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "transition name")
		}
		nameSpan, _ := ctx.Line.LastTokenSpan()
		if err := nr.IsValid(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}
		if t.IsBuiltinTransitionName(name) {
			return nil, diag.Validation(fmt.Sprintf("transition name %q is reserved for a built-in transition", name)).WithSpan(nameSpan).WithFound(name)
		}

		points := make([]float64, 0)
		for {
			_, ok := ctx.Line.Peek()
			if !ok {
				break
			}
			value, err := ctx.Line.NextFloat64Strict()
			if err != nil {
				return nil, err
			}
			valueSpan, _ := ctx.Line.LastTokenSpan()
			if value < 0 || value > 100 {
				return nil, diag.Validation("transition points must be between 0 and 100").WithSpan(valueSpan).WithFound(fmt.Sprintf("%g", value))
			}
			if len(points) == t.MaxTransitionPoints {
				return nil, diag.Validation(fmt.Sprintf("transition cannot contain more than %d points", t.MaxTransitionPoints)).WithSpan(valueSpan)
			}
			if len(points) > 0 && value < points[len(points)-1]*100 {
				return nil, diag.Validation("transition points must be monotonic").WithSpan(valueSpan).WithFound(fmt.Sprintf("%g", value))
			}
			points = append(points, value/100)
		}
		if len(points) < t.MinTransitionPoints {
			return nil, diag.Validation(fmt.Sprintf("transition must contain at least %d points", t.MinTransitionPoints)).WithSpan(nameSpan).WithFound(name)
		}
		if points[0] != 0 || points[len(points)-1] != 1 {
			return nil, diag.Validation("transition must start at 0 and end at 100").WithSpan(nameSpan).WithFound(name)
		}
		parsed.Transitions = append(parsed.Transitions, t.TransitionDefinition{Name: name, Points: points})
	default:
		diagnostic := diag.Parse("invalid option").WithSpan(span).WithFound(option).WithExpected(validOptions...)
		if suggestion, ok := diag.ClosestMatch(option, validOptions, diag.DefaultSuggestionDistance(option)); ok {
			diagnostic.WithSuggestion(fmt.Sprintf("did you mean %q?", suggestion))
		}
		return nil, diagnostic
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		unknownSpan, _ := ctx.Line.PeekSpan()
		return nil, diag.Parse("unexpected token after option definition").WithSpan(unknownSpan).WithFound(unknown)
	}

	return parsed, nil
}
