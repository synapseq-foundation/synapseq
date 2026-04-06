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

package sequence

import (
	"fmt"

	p "github.com/synapseq-foundation/synapseq/v4/internal/preset"
	tl "github.com/synapseq-foundation/synapseq/v4/internal/timeline"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func validateSequenceOptionPlacement(sourceFile string, lineNumber int, lineText string, optionsLocked bool) error {
	if !optionsLocked {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, "options must be defined on the top of the file, before any presets or timelines")
}

func mergeSequenceExtends(sourceFile string, lineNumber int, lineText string, parsedOptions *t.ParseOptions, rawOptions *t.ParseOptions, presets *[]t.Preset, loadedExtends map[string]struct{}) error {
	rawOptions.Merge(parsedOptions)

	for _, extFile := range parsedOptions.Extends {
		if _, ok := loadedExtends[extFile]; ok {
			continue
		}

		extendsConfig, err := extends(extFile)
		if err != nil {
			return withSource(err, sourceFile, lineNumber, lineText)
		}

		loadedExtends[extFile] = struct{}{}
		*presets = append(*presets, extendsConfig.Presets...)
		rawOptions.Merge(extendsConfig.Options)
	}

	if _, err := rawOptions.Build(); err != nil {
		return withSource(err, sourceFile, lineNumber, lineText)
	}

	return nil
}

func validateSequencePresetPlacement(sourceFile string, lineNumber int, lineText string, presetCount int, periods []t.Period) error {
	if presetCount >= t.MaxPresets {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "maximum number of presets reached")
	}
	if len(periods) > 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "preset definitions must be before any timeline definitions")
	}

	return nil
}

func validatePresetNotDuplicate(sourceFile string, lineNumber int, lineText string, presetName string, presets []t.Preset) error {
	if p.FindPreset(presetName, presets) == nil {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("duplicate preset definition: %s", presetName))
}

func validateSequenceTrackPlacement(sourceFile string, lineNumber int, lineText string, presetCount int, periods []t.Period, lastPreset *t.Preset) error {
	if presetCount == 1 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("track defined before any preset: %s", lineText))
	}
	if len(periods) > 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "track definitions must be before any timeline definitions")
	}
	if lastPreset.From != nil {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("preset %q inherits from another and cannot define new tracks", lastPreset.String()))
	}

	return nil
}

func validateExtendsTrackPlacement(sourceFile string, lineNumber int, lineText string, presetCount int, lastPreset *t.Preset) error {
	if presetCount == 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("track defined before any preset: %s", lineText))
	}
	if lastPreset.From != nil {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("preset %q inherits from another and cannot define new tracks", lastPreset.String()))
	}

	return nil
}

func validateSequenceTrackOverridePlacement(sourceFile string, lineNumber int, lineText string, presetCount int, periods []t.Period, lastPreset *t.Preset) error {
	if presetCount == 1 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("track override defined before any preset: %s", lineText))
	}
	if len(periods) > 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "track override definitions must be before any timeline definitions")
	}

	return validateTrackOverrideTarget(sourceFile, lineNumber, lineText, lastPreset)
}

func validateExtendsTrackOverridePlacement(sourceFile string, lineNumber int, lineText string, presetCount int, lastPreset *t.Preset) error {
	if presetCount == 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("track override defined before any preset: %s", lineText))
	}

	return validateTrackOverrideTarget(sourceFile, lineNumber, lineText, lastPreset)
}

func validateTrackOverrideTarget(sourceFile string, lineNumber int, lineText string, lastPreset *t.Preset) error {
	if lastPreset.IsTemplate {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("cannot override tracks on template preset %q", lastPreset.String()))
	}
	if lastPreset.From == nil {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("cannot override tracks on preset %q which does not have a 'from' source", lastPreset.String()))
	}

	return nil
}

func validateSequenceTimelinePlacement(sourceFile string, lineNumber int, lineText string, presetCount int) error {
	if presetCount != 1 {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("timeline defined before any preset: %s", lineText))
}

func validateFirstTimelineStartsAtZero(sourceFile string, lineNumber int, lineText string, periods []t.Period, period *t.Period) error {
	if len(periods) != 0 || period.Time == 0 {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, "first timeline must start at 00:00:00")
}

func applyTimelineTransition(sourceFile string, lineNumber int, lineText string, periods []t.Period, period *t.Period) error {
	if len(periods) == 0 {
		return nil
	}

	lastPeriod := &periods[len(periods)-1]
	if lastPeriod.Time >= period.Time {
		return lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("timeline %s overlaps with previous timeline %s", period.TimeString(), lastPeriod.TimeString()))
	}

	if err := tl.AdjustPeriods(lastPeriod, period); err != nil {
		return withSource(err, sourceFile, lineNumber, lineText)
	}

	return nil
}

func unexpectedSequenceLine(sourceFile string, lineNumber int, lineText string, firstToken string) error {
	if firstToken == t.KeywordWaveform ||
		firstToken == t.KeywordTone ||
		firstToken == t.KeywordNoise ||
		firstToken == t.KeywordAmbiance ||
		firstToken == t.KeywordTrack {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "expected two-space indentation for elements under preset definition")
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, "invalid syntax")
}

func finalizeSequence(rawContent []byte, presets []t.Preset, periods []t.Period, comments []string, rawOptions *t.ParseOptions) (*t.Sequence, error) {
	if len(presets) == 1 {
		return nil, fmt.Errorf("no presets defined")
	}

	for i := range presets {
		if p.IsPresetEmpty(&presets[i]) {
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

func validateExtendsOptionPlacement(sourceFile string, lineNumber int, lineText string, optionsLocked bool) error {
	if !optionsLocked {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, "options must be defined before any presets")
}

func validateExtendsOptionContent(sourceFile string, lineNumber int, lineText string, parsedOptions *t.ParseOptions, rawOptions *t.ParseOptions) error {
	if len(parsedOptions.Extends) > 0 {
		return lineDiagnostic(sourceFile, lineNumber, lineText, "extends option is not supported in extended files")
	}

	rawOptions.Merge(parsedOptions)
	if _, err := rawOptions.Build(); err != nil {
		return withSource(err, sourceFile, lineNumber, lineText)
	}

	return nil
}

func validateExtendsPresetPlacement(sourceFile string, lineNumber int, lineText string, presetCount int) error {
	if presetCount < t.MaxPresets {
		return nil
	}

	return lineDiagnostic(sourceFile, lineNumber, lineText, "maximum number of presets reached")
}

func finalizeExtends(rawOptions *t.ParseOptions, presets []t.Preset) (*t.Extends, error) {
	for i := range presets {
		if p.IsPresetEmpty(&presets[i]) {
			return nil, fmt.Errorf("spsc file: preset %q is empty", presets[i].String())
		}
	}

	return &t.Extends{
		Options: rawOptions,
		Presets: presets,
	}, nil
}
