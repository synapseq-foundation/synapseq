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

	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// parseSequenceContent parses the raw content of a sequence file and returns a Sequence struct
func parseSequenceContent(rawContent []byte, baseRef string) (*t.Sequence, error) {
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
				return nil, fmt.Errorf("line %d: options must be defined on the top of the file, before any presets or timelines", lnn)
			}

			parsedOptions, err := ctx.ParseOption(baseRef)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			rawOptions.Merge(parsedOptions)

			for _, extFile := range parsedOptions.Extends {
				if _, ok := loadedExtends[extFile]; ok {
					continue
				}

				extendsConfig, err := extends(extFile)
				if err != nil {
					return nil, fmt.Errorf("line %d: %v", lnn, err)
				}

				loadedExtends[extFile] = struct{}{}
				presets = append(presets, extendsConfig.Presets...)
				rawOptions.Merge(extendsConfig.Options)
			}

			if _, err := rawOptions.Build(); err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			continue
		}

		if ctx.HasPreset() {
			optionsLocked = true

			if len(presets) >= t.MaxPresets {
				return nil, fmt.Errorf("line %d: maximum number of presets reached", lnn)
			}

			if len(periods) > 0 {
				return nil, fmt.Errorf("line %d: preset definitions must be before any timeline definitions", lnn)
			}

			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			pName := preset.String()
			p := s.FindPreset(pName, presets)
			if p != nil {
				return nil, fmt.Errorf("line %d: duplicate preset definition: %s", lnn, pName)
			}

			presets = append(presets, *preset)
			continue
		}

		if ctx.HasTrack() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, fmt.Errorf("line %d: track defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			if len(periods) > 0 {
				return nil, fmt.Errorf("line %d: track definitions must be before any timeline definitions", lnn)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, fmt.Errorf("line %d: preset %q inherits from another and cannot define new tracks", lnn, lastPreset.String())
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		if ctx.HasTrackOverride() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, fmt.Errorf("line %d: track override defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			if len(periods) > 0 {
				return nil, fmt.Errorf("line %d: track override definitions must be before any timeline definitions", lnn)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, fmt.Errorf("line %d: cannot override tracks on template preset %q", lnn, lastPreset.String())
			}
			if lastPreset.From == nil {
				return nil, fmt.Errorf("line %d: cannot override tracks on preset %q which does not have a 'from' source", lnn, lastPreset.String())
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			continue
		}

		if ctx.HasTimeline() {
			optionsLocked = true

			if len(presets) == 1 {
				return nil, fmt.Errorf("line %d: timeline defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			period, err := ctx.ParseTimeline(&presets)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			if len(periods) == 0 && period.Time != 0 {
				return nil, fmt.Errorf("line %d: first timeline must start at 00:00:00", lnn)
			}

			if len(periods) > 0 {
				lastPeriod := &periods[len(periods)-1]

				if lastPeriod.Time >= period.Time {
					return nil, fmt.Errorf("line %d: timeline %s overlaps with previous timeline %s", lnn, period.TimeString(), lastPeriod.TimeString())
				}

				if err := s.AdjustPeriods(lastPeriod, period); err != nil {
					return nil, fmt.Errorf("line %d: %v", lnn, err)
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
			return nil, fmt.Errorf("line %d: expected two-space indentation for elements under preset definition\n   %s", lnn, ctx.Line.Raw)
		}

		return nil, fmt.Errorf("line %d: invalid syntax\n    %s", lnn, ctx.Line.Raw)
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
func parseExtendsContent(rawContent []byte, baseRef string) (*t.Extends, error) {
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
				return nil, fmt.Errorf("spsc file, line %d: options must be defined before any presets", lnn)
			}

			parsedOptions, err := ctx.ParseOption(baseRef)
			if err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			if len(parsedOptions.Extends) > 0 {
				return nil, fmt.Errorf("spsc file, line %d: extends option is not supported in extended files", lnn)
			}

			rawOptions.Merge(parsedOptions)

			if _, err := rawOptions.Build(); err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			continue
		}

		if ctx.HasPreset() {
			optionsLocked = true

			if len(presets) >= t.MaxPresets {
				return nil, fmt.Errorf("spsc file, line %d: maximum number of presets reached", lnn)
			}

			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			pName := preset.String()
			p := s.FindPreset(pName, presets)
			if p != nil {
				return nil, fmt.Errorf("spsc file, line %d: duplicate preset definition: %s", lnn, pName)
			}

			presets = append(presets, *preset)
			continue
		}

		if ctx.HasTrack() {
			optionsLocked = true

			if len(presets) == 0 {
				return nil, fmt.Errorf("spsc file, line %d: track defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, fmt.Errorf("spsc file, line %d: preset %q inherits from another and cannot define new tracks", lnn, lastPreset.String())
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		if ctx.HasTrackOverride() {
			optionsLocked = true

			if len(presets) == 0 {
				return nil, fmt.Errorf("spsc file, line %d: track override defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, fmt.Errorf("spsc file, line %d: cannot override tracks on template preset %q", lnn, lastPreset.String())
			}
			if lastPreset.From == nil {
				return nil, fmt.Errorf("spsc file, line %d: cannot override tracks on preset %q which does not have a 'from' source", lnn, lastPreset.String())
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, fmt.Errorf("spsc file, line %d: %v", lnn, err)
			}

			continue
		}

		return nil, fmt.Errorf("spsc file, line %d: unexpected content: %s", lnn, ln)
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
