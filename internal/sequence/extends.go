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

	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// extends loads preset and option definitions from a .spsc file.
func extends(fileName string) (*t.Extends, error) {
	rawContent, err := s.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, err
	}

	file := NewSequenceFile(rawContent)

	absInputFile, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	baseDir := filepath.Dir(absInputFile)

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

			parsedOptions, err := ctx.ParseOption(baseDir)
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
