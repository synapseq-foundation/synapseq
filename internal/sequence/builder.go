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
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	p "github.com/synapseq-foundation/synapseq/v4/internal/preset"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type sequenceBuilder struct {
	sourceFile    string
	rawContent    []byte
	rawOptions    *t.ParseOptions
	presets       []t.Preset
	periods       []t.Period
	comments      []string
	loadedExtends map[string]struct{}
	optionsLocked bool
	extendsMode   bool
}

func newSequenceBuilder(rawContent []byte, sourceFile string) *sequenceBuilder {
	builder := &sequenceBuilder{
		sourceFile:    sourceFile,
		rawContent:    rawContent,
		rawOptions:    t.NewParseOptions(),
		loadedExtends: make(map[string]struct{}),
	}
	builder.presets = append(builder.presets, *t.NewBuiltinSilencePreset())
	return builder
}

func newExtendsBuilder(sourceFile string) *sequenceBuilder {
	return &sequenceBuilder{
		sourceFile:  sourceFile,
		rawOptions:  t.NewParseOptions(),
		extendsMode: true,
	}
}

func (b *sequenceBuilder) handleComment(comment string) {
	if b.extendsMode || comment == "" {
		return
	}
	b.comments = append(b.comments, comment)
}

func (b *sequenceBuilder) handleOption(lineNumber int, lineText string, parsedOptions *t.ParseOptions) error {
	if b.extendsMode {
		if err := validateExtendsOptionPlacement(b.sourceFile, lineNumber, lineText, b.optionsLocked); err != nil {
			return err
		}
		return validateExtendsOptionContent(b.sourceFile, lineNumber, lineText, parsedOptions, b.rawOptions)
	}

	if err := validateSequenceOptionPlacement(b.sourceFile, lineNumber, lineText, b.optionsLocked); err != nil {
		return err
	}

	return mergeSequenceExtends(b.sourceFile, lineNumber, lineText, parsedOptions, b.rawOptions, &b.presets, b.loadedExtends)
}

func (b *sequenceBuilder) handlePreset(lineNumber int, lineText string, ctx *parser.TextParser) error {
	b.optionsLocked = true

	if b.extendsMode {
		if err := validateExtendsPresetPlacement(b.sourceFile, lineNumber, lineText, len(b.presets)); err != nil {
			return err
		}
	} else {
		if err := validateSequencePresetPlacement(b.sourceFile, lineNumber, lineText, len(b.presets), b.periods); err != nil {
			return err
		}
	}

	decl, err := ctx.ParsePresetDeclaration()
	if err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	preset, err := buildPresetFromDeclaration(b.sourceFile, lineNumber, lineText, decl, b.presets)
	if err != nil {
		return err
	}

	b.presets = append(b.presets, *preset)
	return nil
}

func (b *sequenceBuilder) handleTrack(lineNumber int, lineText string, ctx *parser.TextParser) error {
	b.optionsLocked = true

	lastPreset := b.lastPreset()
	var err error
	if b.extendsMode {
		err = validateExtendsTrackPlacement(b.sourceFile, lineNumber, lineText, len(b.presets), lastPreset)
	} else {
		err = validateSequenceTrackPlacement(b.sourceFile, lineNumber, lineText, len(b.presets), b.periods, lastPreset)
	}
	if err != nil {
		return err
	}

	trackIndex, err := p.AllocateTrack(lastPreset)
	if err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	decl, err := ctx.ParseTrackDeclaration()
	if err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	track, err := buildTrackFromDeclaration(b.sourceFile, lineNumber, lineText, decl)
	if err != nil {
		return err
	}

	lastPreset.Track[trackIndex] = *track
	return nil
}

func (b *sequenceBuilder) handleTrackOverride(lineNumber int, lineText string, ctx *parser.TextParser) error {
	b.optionsLocked = true

	lastPreset := b.lastPreset()
	var err error
	if b.extendsMode {
		err = validateExtendsTrackOverridePlacement(b.sourceFile, lineNumber, lineText, len(b.presets), lastPreset)
	} else {
		err = validateSequenceTrackOverridePlacement(b.sourceFile, lineNumber, lineText, len(b.presets), b.periods, lastPreset)
	}
	if err != nil {
		return err
	}

	decl, err := ctx.ParseTrackOverrideDeclaration()
	if err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	if err := p.ApplyTrackOverride(lastPreset, decl); err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	return nil
}

func (b *sequenceBuilder) handleTimeline(lineNumber int, lineText string, ctx *parser.TextParser) error {
	b.optionsLocked = true

	if err := validateSequenceTimelinePlacement(b.sourceFile, lineNumber, lineText, len(b.presets)); err != nil {
		return err
	}

	decl, err := ctx.ParseTimelineDeclaration()
	if err != nil {
		return withSource(err, b.sourceFile, lineNumber, lineText)
	}

	period, err := buildPeriodFromDeclaration(b.sourceFile, lineNumber, lineText, decl, b.presets)
	if err != nil {
		return err
	}

	if err := validateFirstTimelineStartsAtZero(b.sourceFile, lineNumber, lineText, b.periods, period); err != nil {
		return err
	}
	if err := applyTimelineTransition(b.sourceFile, lineNumber, lineText, b.periods, period); err != nil {
		return err
	}

	b.periods = append(b.periods, *period)
	return nil
}

func (b *sequenceBuilder) buildSequence() (*t.Sequence, error) {
	return finalizeSequence(b.rawContent, b.presets, b.periods, b.comments, b.rawOptions)
}

func (b *sequenceBuilder) buildExtends() (*t.Extends, error) {
	return finalizeExtends(b.rawOptions, b.presets)
}

func (b *sequenceBuilder) handleUnexpectedLine(lineNumber int, lineText string, firstToken string) error {
	if b.extendsMode {
		return lineDiagnostic(b.sourceFile, lineNumber, lineText, "unexpected content: "+lineText)
	}

	return unexpectedSequenceLine(b.sourceFile, lineNumber, lineText, firstToken)
}

func (b *sequenceBuilder) lastPreset() *t.Preset {
	if len(b.presets) == 0 {
		return nil
	}

	return &b.presets[len(b.presets)-1]
}
