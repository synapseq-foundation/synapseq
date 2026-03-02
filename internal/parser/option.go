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
	"os"
	"path/filepath"
	"strings"

	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// getFullPath resolves the full path of a given file path
func getFullPath(path, basePath string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		expanded := filepath.Join(homeDir, strings.TrimPrefix(path, "~"))
		return filepath.Clean(expanded), nil
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	fullPath := filepath.Join(basePath, path)
	return filepath.Clean(fullPath), nil
}

// HasOption checks if the first element is an option
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and applies the option from the elements
func (ctx *TextParser) ParseOption(options *t.SequenceOptions, filePath string) error {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return fmt.Errorf("expected option, got EOF: %s", ln)
	}

	if string(tok[0]) != t.KeywordOption {
		return fmt.Errorf("expected option. Received: %s", tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return fmt.Errorf("expected option name: %s", ln)
	}

	switch option {
	case t.KeywordOptionSampleRate:
		sampleRate, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("samplerate: %v", err)
		}
		options.SampleRate = sampleRate
	case t.KeywordOptionVolume:
		volume, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("volume: %v", err)
		}
		options.Volume = volume
	case t.KeywordOptionAmbiance:
		name, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected name for ambiance audio file: %s", ln)
		}

		if err := s.IsValidNamedRef(name); err != nil {
			return err
		}

		content := strings.Join(ctx.Line.Tokens[2:], " ")

		if content == "-" {
			return fmt.Errorf("stdin (-) is not supported for ambiance list")
		}

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = getFullPath(content, filePath)
			if err != nil {
				return fmt.Errorf("path: %v", err)
			}
		}

		options.AmbianceList[name] = fullPath
	case t.KeywordOptionPresetList:
		_, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected path: %s", ln)
		}

		content := strings.Join(ctx.Line.Tokens[1:], " ")

		if content == "-" {
			return fmt.Errorf("stdin (-) is not supported for preset list")
		}

		fullPath := content
		if !s.IsRemoteFile(content) {
			var err error
			fullPath, err = getFullPath(content, filePath)
			if err != nil {
				return fmt.Errorf("path: %v", err)
			}
		}

		options.PresetList = append(options.PresetList, fullPath)
	default:
		return fmt.Errorf("invalid option: %q", option)
	}

	// If the option is not ambiance, ensure no extra tokens are present
	if option != t.KeywordOptionAmbiance {
		unknown, ok := ctx.Line.Peek()
		if ok {
			return fmt.Errorf("unexpected token after option definition: %q", unknown)
		}
	}

	return nil
}
