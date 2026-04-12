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

package sequence

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var windowsDrivePathPattern = regexp.MustCompile(`^[a-zA-Z]:`)

func resolveParsedOptions(baseRef string, parsedOptions *t.ParseOptions) error {
	for name, path := range parsedOptions.Ambiance {
		resolved, err := resolveOptionFile(baseRef, path, ".wav", "ambiance")
		if err != nil {
			return err
		}
		parsedOptions.Ambiance[name] = resolved
	}

	for i := range parsedOptions.Extends {
		resolved, err := resolveOptionFile(baseRef, parsedOptions.Extends[i], ".spsc", "extends")
		if err != nil {
			return err
		}
		parsedOptions.Extends[i] = resolved
	}

	return nil
}

func resolveOptionFile(baseRef, content, ext, optionName string) (string, error) {
	if r.IsRemoteFile(content) {
		return content, nil
	}
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

	return filepath.Join(baseRef, cleanPath) + ext, nil
}
