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

package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

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
	ln := ctx.Line.Raw

	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected option, got EOF: %s", ln)
	}

	if string(tok[0]) != t.KeywordOption {
		return nil, fmt.Errorf("expected option. Received: %s", tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return nil, fmt.Errorf("expected option name: %s", ln)
	}

	parsed := t.NewParseOptions()

	switch option {
	case t.KeywordOptionSampleRate:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected samplerate value: %s", ln)
		}

		parsed.Values[t.KeywordOptionSampleRate] = value
	case t.KeywordOptionVolume:
		value, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected volume value: %s", ln)
		}

		parsed.Values[t.KeywordOptionVolume] = value
	case t.KeywordOptionAmbiance:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected name for ambiance audio file: %s", ln)
		}

		if err := s.IsValidNamedRef(name); err != nil {
			return nil, err
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected path for ambiance audio file: %s", ln)
		}

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = resolveLocalOptionFile(dirPath, content, ".wav", "ambiance")
			if err != nil {
				return nil, err
			}
		}

		parsed.Ambiance[name] = fullPath
	case t.KeywordOptionExtends:
		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected path: %s", ln)
		}

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = resolveLocalOptionFile(dirPath, content, ".spsc", "extends")
			if err != nil {
				return nil, err
			}
		}

		parsed.Extends = append(parsed.Extends, fullPath)
	default:
		return nil, fmt.Errorf("invalid option: %q", option)
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		return nil, fmt.Errorf("unexpected token after option definition: %q", unknown)
	}

	return parsed, nil
}
