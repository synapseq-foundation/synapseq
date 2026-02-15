//go:build !wasm

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
	"path/filepath"

	"github.com/synapseq-foundation/synapseq/v3/internal/parser"
	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// LoadTextSequence loads a sequence from a text file
func LoadTextSequence(fileName string) (*t.Sequence, error) {
	rawContent, err := s.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	file := NewSequenceFile(rawContent)

	// Get absolute path of input file
	absInputFile, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	presets := make([]t.Preset, 0, t.MaxPresets)

	// Initialize built-in presets
	presets = append(presets, *t.NewBuiltinSilencePreset())

	// Options can only be defined on the top of the file, before any presets
	optionsLocked := false
	// Last loaded preset path from options
	lastLoadedPresetPath := ""
	// Initialize audio options
	options := &t.SequenceOptions{
		SampleRate:     44100,
		Volume:         100,
		BackgroundList: []string{},
		PresetList:     []string{},
		GainLevel:      t.GainLevelOff,
	}

	var (
		periods  []t.Period
		comments []string
	)

	// Parse each line in the file
	for file.NextLine() {
		ln := file.CurrentLine()
		lnn := file.CurrentLineNumber()
		ctx := parser.NewTextParser(ln)

		// Skip empty lines
		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		// Skip comments
		if ctx.HasComment() {
			comment := ctx.ParseComment()
			if comment != "" {
				comments = append(comments, comment)
				// fmt.Fprintf(os.Stderr, "> %s\n", comment)
			}
			continue
		}

		// Option line
		if ctx.HasOption() {
			if optionsLocked {
				return nil, fmt.Errorf("line %d: options must be defined on the top of the file, before any presets or timelines", lnn)
			}

			if err = ctx.ParseOption(options, filepath.Dir(absInputFile)); err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}
			// Validate options
			if err = options.Validate(); err != nil {
				return nil, fmt.Errorf("line %d: %v", lnn, err)
			}

			// Load presets from file if specified in options and not already loaded
			if len(options.PresetList) > 0 {
				lastList := options.PresetList[len(options.PresetList)-1]
				if lastList != lastLoadedPresetPath {
					fpresets, err := loadPresets(lastList)
					if err != nil {
						return nil, err
					}
					presets = append(presets, fpresets...)
					lastLoadedPresetPath = lastList
				}
			}

			continue
		}

		// Preset definition
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

		// Track line
		if ctx.HasTrack() {
			optionsLocked = true

			if len(presets) == 1 { // 1 = silence preset
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

			if track.Type == t.TrackBackground && len(options.BackgroundList) == 0 {
				return nil, fmt.Errorf("line %d: background track defined but no background audio file specified in options", lnn)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		// Track override line
		if ctx.HasTrackOverride() {
			optionsLocked = true

			if len(presets) == 1 { // 1 = silence preset
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

		// Timeline
		if ctx.HasTimeline() {
			optionsLocked = true

			if len(presets) == 1 { // 1 = silence preset
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

		// Check for indentation errors
		tok := ctx.Line.Tokens[0]
		if tok == t.KeywordWaveform ||
			tok == t.KeywordTone ||
			tok == t.KeywordNoise ||
			tok == t.KeywordBackground ||
			tok == t.KeywordTrack {
			return nil, fmt.Errorf("line %d: expected two-space indentation for elements under preset definition\n   %s", lnn, ctx.Line.Raw)
		}

		return nil, fmt.Errorf("line %d: invalid syntax\n    %s", lnn, ctx.Line.Raw)
	}

	// Validate if has one preset (1 = silence preset)
	if len(presets) == 1 {
		return nil, fmt.Errorf("no presets defined")
	}

	// Validate each preset (skip silence preset)
	for i := 1; i < len(presets); i++ {
		p := &presets[i]
		if s.IsPresetEmpty(p) {
			return nil, fmt.Errorf("preset %q is empty", presets[i].String())
		}
	}

	// Validate if has more than two Periods
	if len(periods) < 2 {
		return nil, fmt.Errorf("at least two periods must be defined")
	}

	return &t.Sequence{
		Periods:    periods,
		Options:    options,
		Comments:   comments,
		RawContent: rawContent,
	}, nil
}
