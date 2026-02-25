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

	"github.com/synapseq-foundation/synapseq/v3/internal/parser"
	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// loadPresets loads presets from a given file path
func loadPresets(filename string) ([]t.Preset, error) {
	rawContent, err := s.GetFile(filename, t.FormatText)
	if err != nil {
		return nil, err
	}

	f := NewSequenceFile(rawContent)

	presets := make([]t.Preset, 0, t.MaxPresets)
	for f.NextLine() {
		ln := f.CurrentLine()
		lnn := f.CurrentLineNumber()
		ctx := parser.NewTextParser(ln)

		// Skip empty lines
		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		// Skip comments
		if ctx.HasComment() {
			continue
		}

		// Parse preset lines
		if ctx.HasPreset() {
			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", lnn, err)
			}
			presets = append(presets, *preset)
			continue
		}

		// Track line
		if ctx.HasTrack() {
			if len(presets) == 0 {
				return nil, fmt.Errorf("preset file, line %d: track defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, fmt.Errorf("preset file, line %d: preset %q inherits from another and cannot define new tracks", lnn, lastPreset.String())
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", lnn, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", lnn, err)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		// Track override line
		if ctx.HasTrackOverride() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, fmt.Errorf("preset file, line %d: track override defined before any preset: %s", lnn, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, fmt.Errorf("preset file, line %d: cannot override tracks on template preset %q", lnn, lastPreset.String())
			}
			if lastPreset.From == nil {
				return nil, fmt.Errorf("preset file, line %d: cannot override tracks on preset %q which does not have a 'from' source", lnn, lastPreset.String())
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", lnn, err)
			}

			continue
		}

		return nil, fmt.Errorf("preset file, line %d: unexpected content: %s", lnn, ln)
	}

	// Validate if has one preset
	if len(presets) == 0 {
		return nil, fmt.Errorf("preset file: no presets defined")
	}

	// Validate each preset (skip silence preset)
	for _, p := range presets {
		if s.IsPresetEmpty(&p) {
			return nil, fmt.Errorf("preset file: preset %q is empty", p.String())
		}
	}

	return presets, nil
}
