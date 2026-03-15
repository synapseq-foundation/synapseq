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

package sequence

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// parseSequenceContent parses the raw content of a sequence file and returns a Sequence struct
func parseSequenceContent(rawContent []byte, sourceFile string, baseRef string) (*t.Sequence, error) {
	file := NewSequenceFile(rawContent)

	presets := make([]t.Preset, 0, t.MaxPresets)
	presets = append(presets, *t.NewBuiltinSilencePreset())

	rawOptions := t.NewParseOptions()
	loadedExtends := make(map[string]struct{})
	optionsLocked := false

	var (
		periods  []t.Period
		comments []string
	)

	for file.NextLine() {
		ln := file.CurrentLine()
		lnn := file.CurrentLineNumber()
		ctx := parser.NewTextParser(ln)

		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		if ctx.HasComment() {
			comment := ctx.ParseComment()
			if comment != "" {
				comments = append(comments, comment)
			}
			continue
		}

		if ctx.HasOption() {
			if optionsLocked {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "options must be defined on the top of the file, before any presets or timelines")
			}

			parsedOptions, err := ctx.ParseOption(baseRef)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			rawOptions.Merge(parsedOptions)

			for _, extFile := range parsedOptions.Extends {
				if _, ok := loadedExtends[extFile]; ok {
					continue
				}

				extendsConfig, err := extends(extFile)
				if err != nil {
					return nil, withSource(err, sourceFile, lnn, ln)
				}

				loadedExtends[extFile] = struct{}{}
				presets = append(presets, extendsConfig.Presets...)
				rawOptions.Merge(extendsConfig.Options)
			}

			if _, err := rawOptions.Build(); err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			continue
		}

		if ctx.HasPreset() {
			optionsLocked = true

			if len(presets) >= t.MaxPresets {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "maximum number of presets reached")
			}

			if len(periods) > 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "preset definitions must be before any timeline definitions")
			}

			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			pName := preset.String()
			p := s.FindPreset(pName, presets)
			if p != nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("duplicate preset definition: %s", pName))
			}

			presets = append(presets, *preset)
			continue
		}

		if ctx.HasTrack() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("track defined before any preset: %s", ctx.Line.Raw))
			}

			if len(periods) > 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "track definitions must be before any timeline definitions")
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("preset %q inherits from another and cannot define new tracks", lastPreset.String()))
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		if ctx.HasTrackOverride() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("track override defined before any preset: %s", ctx.Line.Raw))
			}

			if len(periods) > 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "track override definitions must be before any timeline definitions")
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("cannot override tracks on template preset %q", lastPreset.String()))
			}
			if lastPreset.From == nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("cannot override tracks on preset %q which does not have a 'from' source", lastPreset.String()))
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			continue
		}

		if ctx.HasTimeline() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("timeline defined before any preset: %s", ctx.Line.Raw))
			}

			period, err := ctx.ParseTimeline(&presets)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			if len(periods) == 0 && period.Time != 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "first timeline must start at 00:00:00")
			}

			if len(periods) > 0 {
				lastPeriod := &periods[len(periods)-1]

				if lastPeriod.Time >= period.Time {
					return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("timeline %s overlaps with previous timeline %s", period.TimeString(), lastPeriod.TimeString()))
				}

				if err := s.AdjustPeriods(lastPeriod, period); err != nil {
					return nil, withSource(err, sourceFile, lnn, ln)
				}
			}

			periods = append(periods, *period)
			continue
		}

		tok := ctx.Line.Tokens[0]
		if tok == t.KeywordWaveform ||
			tok == t.KeywordTone ||
			tok == t.KeywordNoise ||
			tok == t.KeywordAmbiance ||
			tok == t.KeywordTrack {
			return nil, lineDiagnostic(sourceFile, lnn, ln, "expected two-space indentation for elements under preset definition")
		}

		return nil, lineDiagnostic(sourceFile, lnn, ln, "invalid syntax")
	}

	if len(presets) == 1 {
		return nil, fmt.Errorf("no presets defined")
	}

	for i := range presets {
		if s.IsPresetEmpty(&presets[i]) {
			return nil, fmt.Errorf("preset %q is empty", presets[i].String())
		}
	}

	if len(periods) < 2 {
		return nil, fmt.Errorf("at least two periods must be defined")
	}

	options, err := rawOptions.Build()
	if err != nil {
		return nil, err
	}

	return &t.Sequence{
		Periods:    periods,
		Options:    options,
		Comments:   comments,
		RawContent: rawContent,
	}, nil
}

// parseExtendsContent parses the raw content of an extended sequence file and returns an Extends struct
func parseExtendsContent(rawContent []byte, sourceFile string, baseRef string) (*t.Extends, error) {
	file := NewSequenceFile(rawContent)
	presets := make([]t.Preset, 0, t.MaxPresets)
	rawOptions := t.NewParseOptions()
	optionsLocked := false

	for file.NextLine() {
		ln := file.CurrentLine()
		lnn := file.CurrentLineNumber()
		ctx := parser.NewTextParser(ln)

		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		if ctx.HasComment() {
			continue
		}

		if ctx.HasOption() {
			if optionsLocked {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "options must be defined before any presets")
			}

			parsedOptions, err := ctx.ParseOption(baseRef)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			if len(parsedOptions.Extends) > 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "extends option is not supported in extended files")
			}

			rawOptions.Merge(parsedOptions)

			if _, err := rawOptions.Build(); err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			continue
		}

		if ctx.HasPreset() {
			optionsLocked = true

			if len(presets) >= t.MaxPresets {
				return nil, lineDiagnostic(sourceFile, lnn, ln, "maximum number of presets reached")
			}

			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			pName := preset.String()
			p := s.FindPreset(pName, presets)
			if p != nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("duplicate preset definition: %s", pName))
			}

			presets = append(presets, *preset)
			continue
		}

		if ctx.HasTrack() {
			optionsLocked = true

			if len(presets) == 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("track defined before any preset: %s", ctx.Line.Raw))
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("preset %q inherits from another and cannot define new tracks", lastPreset.String()))
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		if ctx.HasTrackOverride() {
			optionsLocked = true

			if len(presets) == 0 {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("track override defined before any preset: %s", ctx.Line.Raw))
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("cannot override tracks on template preset %q", lastPreset.String()))
			}
			if lastPreset.From == nil {
				return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("cannot override tracks on preset %q which does not have a 'from' source", lastPreset.String()))
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, withSource(err, sourceFile, lnn, ln)
			}

			continue
		}

		return nil, lineDiagnostic(sourceFile, lnn, ln, fmt.Sprintf("unexpected content: %s", ln))
	}

	for i := range presets {
		if s.IsPresetEmpty(&presets[i]) {
			return nil, fmt.Errorf("spsc file: preset %q is empty", presets[i].String())
		}
	}

	return &t.Extends{
		Options: rawOptions,
		Presets: presets,
	}, nil
}

func lineDiagnostic(sourceFile string, lineNumber int, lineText string, message string) error {
	return diag.Parse(message).WithSpan(diag.Span{
		File:      sourceFile,
		Line:      lineNumber,
		Column:    1,
		EndColumn: 2,
		LineText:  lineText,
	})
}

func withSource(err error, sourceFile string, lineNumber int, lineText string) error {
	if diagnostic, ok := diag.As(err); ok {
		if diagnostic.Span.File == "" {
			diagnostic.Span.File = sourceFile
		}
		if diagnostic.Span.Line == 0 {
			diagnostic.Span.Line = lineNumber
		}
		if diagnostic.Span.LineText == "" {
			diagnostic.Span.LineText = lineText
		}
		if diagnostic.Span.Column < 1 {
			diagnostic.Span.Column = 1
		}
		if diagnostic.Span.EndColumn < diagnostic.Span.Column+1 {
			diagnostic.Span.EndColumn = diagnostic.Span.Column + 1
		}
		return diagnostic
	}

	return diag.Wrap(diag.KindParse, err.Error(), err).WithSpan(diag.Span{
		File:      sourceFile,
		Line:      lineNumber,
		Column:    1,
		EndColumn: 2,
		LineText:  lineText,
	})
}
