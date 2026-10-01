// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	"fmt"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	p "github.com/synapseq-foundation/synapseq/v4/internal/preset"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func buildPresetFromDeclaration(sourceFile string, lineNumber int, lineText string, decl *parser.ParsedPresetDeclaration, presets []t.Preset) (*t.Preset, error) {
	if err := validatePresetNotDuplicate(sourceFile, lineNumber, lineText, strings.ToLower(decl.Name), presets); err != nil {
		return nil, err
	}

	var fromPreset *t.Preset
	if decl.FromName != "" {
		fromPreset = p.FindPreset(decl.FromName, presets)
		if fromPreset == nil {
			return nil, lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("unknown preset to inherit from: %q", decl.FromName))
		}
		if !fromPreset.IsTemplate {
			return nil, lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("can only inherit from a template preset, but %q is not a template", decl.FromName))
		}
	}

	preset, err := t.NewPreset(decl.Name, decl.IsTemplate, fromPreset)
	if err != nil {
		return nil, lineDiagnostic(sourceFile, lineNumber, lineText, err.Error())
	}

	return preset, nil
}

func buildPeriodFromDeclaration(sourceFile string, lineNumber int, lineText string, decl *parser.ParsedTimelineDeclaration, presets []t.Preset, transitions []t.TransitionDefinition) (*t.Period, error) {
	selectedPreset := p.FindPreset(strings.ToLower(decl.PresetName), presets)
	if selectedPreset == nil {
		return nil, lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("preset %q not found", decl.PresetName))
	}
	if selectedPreset.IsTemplate {
		return nil, lineDiagnostic(sourceFile, lineNumber, lineText, fmt.Sprintf("cannot use template preset %q in timeline", selectedPreset.String()))
	}
	if err := validateTransitionReference(decl.TransitionName, transitions); err != nil {
		return nil, lineDiagnostic(sourceFile, lineNumber, lineText, err.Error())
	}

	return &t.Period{
		Time:           decl.Time,
		PresetName:     selectedPreset.String(),
		TrackStart:     selectedPreset.Track,
		TrackEnd:       selectedPreset.Track,
		Transition:     decl.Transition,
		TransitionName: decl.TransitionName,
		Steps:          decl.Steps,
	}, nil
}
