//go:build !wasm

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

package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var windowsDrivePathPattern = regexp.MustCompile(`^[a-zA-Z]:`)

// HasOption checks if the first element is an option.
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// resolveLocalOptionFile resolves a local modular file path relative to dirPath.
func resolveLocalOptionFile(dirPath, content, ext, optionName string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("expected path for %s option", optionName)
	}

	if content == "-" {
		return "", fmt.Errorf("stdin (-) is not supported for %s option", optionName)
	}

	if strings.Contains(content, "\\") {
		return "", fmt.Errorf("invalid path separator '\\'. Use '/' in %s option paths", optionName)
	}

	if strings.HasPrefix(content, "/") {
		return "", fmt.Errorf("absolute paths are not allowed in %s option paths", optionName)
	}

	if windowsDrivePathPattern.MatchString(content) {
		return "", fmt.Errorf("drive paths are not allowed in %s option paths", optionName)
	}

	cleanPath := filepath.Clean(content)
	if cleanPath == "." || cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("parent directory traversal is not allowed in %s option paths", optionName)
	}

	if filepath.Ext(cleanPath) != "" {
		return "", fmt.Errorf("%s local path must not include file extension", optionName)
	}

	return filepath.Join(dirPath, cleanPath) + ext, nil
}

// ParseOption extracts and returns raw parsed option values.
func (ctx *TextParser) ParseOption(dirPath string) (*t.ParseOptions, error) {
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
		t.KeywordOptionExtends,
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

		if err := s.IsValidNamedRef(name); err != nil {
			return nil, diag.Validation(err.Error()).WithSpan(nameSpan).WithFound(name).WithCause(err)
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			content = name // allow shorthand ambiance name as path
			// return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "ambiance path")
		}
		contentSpan, _ := ctx.Line.LastTokenSpan()

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = resolveLocalOptionFile(dirPath, content, ".wav", "ambiance")
			if err != nil {
				return nil, diag.Validation(err.Error()).WithSpan(contentSpan).WithFound(content).WithCause(err)
			}
		}

		parsed.Ambiance[name] = fullPath
	case t.KeywordOptionExtends:
		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, diag.UnexpectedEOF(ctx.Line.EOFSpan(), "extends path")
		}
		contentSpan, _ := ctx.Line.LastTokenSpan()

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = resolveLocalOptionFile(dirPath, content, ".spsc", "extends")
			if err != nil {
				return nil, diag.Validation(err.Error()).WithSpan(contentSpan).WithFound(content).WithCause(err)
			}
		}

		parsed.Extends = append(parsed.Extends, fullPath)
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
