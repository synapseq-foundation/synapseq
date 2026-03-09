//go:build wasm

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

	s "github.com/synapseq-foundation/synapseq/v4/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HasOption checks if the first element is an option
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and returns raw parsed option values.
func (ctx *TextParser) ParseOption(_ string) (*t.ParseOptions, error) {
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
			return nil, fmt.Errorf("invalid ambiance name: %v", err)
		}

		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected path for ambiance audio file: %s", ln)
		}

		if !s.IsRemoteFile(content) {
			return nil, fmt.Errorf("file paths are not supported in WASM for ambiance audio: %s", content)
		}

		parsed.Ambiance[name] = content
	case t.KeywordOptionExtends:
		content, ok := ctx.Line.NextToken()
		if !ok {
			return nil, fmt.Errorf("expected path: %s", ln)
		}

		if !s.IsRemoteFile(content) {
			return nil, fmt.Errorf("file paths are not supported in WASM for extends: %s", content)
		}

		parsed.Extends = append(parsed.Extends, content)
	default:
		return nil, fmt.Errorf("invalid option: %q", option)
	}

	if unknown, ok := ctx.Line.Peek(); ok {
		return nil, fmt.Errorf("unexpected token after option definition: %q", unknown)
	}

	return parsed, nil
}
