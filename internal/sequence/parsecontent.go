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
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	"github.com/synapseq-foundation/synapseq/v4/internal/parser"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// parseSequenceContent parses the raw content of a sequence file and returns a Sequence struct
func parseSequenceContent(rawContent []byte, sourceFile string, baseRef string) (*t.Sequence, error) {
	builder := newSequenceBuilder(rawContent, sourceFile)

	if err := parseIntoBuilder(NewSequenceFile(rawContent), builder, baseRef); err != nil {
		return nil, err
	}

	return builder.buildSequence()
}

// parseExtendsContent parses the raw content of an extended sequence file and returns an Extends struct
func parseExtendsContent(rawContent []byte, sourceFile string, baseRef string) (*t.Extends, error) {
	builder := newExtendsBuilder(sourceFile)

	if err := parseIntoBuilder(NewSequenceFile(rawContent), builder, baseRef); err != nil {
		return nil, err
	}

	return builder.buildExtends()
}

func parseIntoBuilder(file *SequenceFile, builder *sequenceBuilder, baseRef string) error {
	for file.NextLine() {
		lineText := file.CurrentLine()
		lineNumber := file.CurrentLineNumber()
		ctx := parser.NewTextParser(lineText)

		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		if ctx.HasComment() {
			builder.handleComment(ctx.ParseComment())
			continue
		}

		if ctx.HasOption() {
			parsedOptions, err := ctx.ParseOption(baseRef)
			if err != nil {
				return withSource(err, builder.sourceFile, lineNumber, lineText)
			}
			if err := resolveParsedOptions(baseRef, parsedOptions); err != nil {
				return withSource(err, builder.sourceFile, lineNumber, lineText)
			}
			if err := builder.handleOption(lineNumber, lineText, parsedOptions); err != nil {
				return err
			}
			continue
		}

		if ctx.HasPreset() {
			if err := builder.handlePreset(lineNumber, lineText, ctx); err != nil {
				return err
			}
			continue
		}

		if ctx.HasTrack() {
			if err := builder.handleTrack(lineNumber, lineText, ctx); err != nil {
				return err
			}
			continue
		}

		if ctx.HasTrackOverride() {
			if err := builder.handleTrackOverride(lineNumber, lineText, ctx); err != nil {
				return err
			}
			continue
		}

		if ctx.HasTimeline() && !builder.extendsMode {
			if err := builder.handleTimeline(lineNumber, lineText, ctx); err != nil {
				return err
			}
			continue
		}

		return builder.handleUnexpectedLine(lineNumber, lineText, ctx.Line.Tokens[0])
	}

	return nil
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
