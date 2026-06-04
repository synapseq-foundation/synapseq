// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sequence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var windowsDrivePathPattern = regexp.MustCompile(`^[a-zA-Z]:`)

func resolveParsedOptions(baseRef string, parsedOptions *t.ParseOptions) error {
	for name, path := range parsedOptions.Ambiance {
		resolved, err := resolveAmbianceOptionFile(baseRef, path)
		if err != nil {
			return err
		}
		parsedOptions.Ambiance[name] = resolved
	}

	for name, path := range parsedOptions.Music {
		resolved, err := resolveMusicOptionFile(baseRef, path)
		if err != nil {
			return err
		}
		parsedOptions.Music[name] = resolved
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

func resolveAmbianceOptionFile(baseRef, content string) (string, error) {
	if r.IsRemoteFile(content) {
		return content, nil
	}

	basePath, err := resolveOptionFileBase(baseRef, content, "ambiance")
	if err != nil {
		return "", err
	}

	wavPath := basePath + ".wav"
	mp3Path := basePath + ".mp3"

	if exists, err := fileExists(wavPath); err != nil {
		return "", err
	} else if exists {
		return wavPath, nil
	}

	if exists, err := fileExists(mp3Path); err != nil {
		return "", err
	} else if exists {
		return mp3Path, nil
	}

	return "", fmt.Errorf("ambiance file not found; tried %q and %q", wavPath, mp3Path)
}

func resolveMusicOptionFile(baseRef, content string) (string, error) {
	if r.IsRemoteFile(content) {
		return content, nil
	}

	basePath, err := resolveOptionFileBase(baseRef, content, "music")
	if err != nil {
		return "", err
	}

	mp3Path := basePath + ".mp3"
	wavPath := basePath + ".wav"

	if exists, err := fileExists(mp3Path); err != nil {
		return "", err
	} else if exists {
		return mp3Path, nil
	}

	if exists, err := fileExists(wavPath); err != nil {
		return "", err
	} else if exists {
		return wavPath, nil
	}

	return "", fmt.Errorf("music file not found; tried %q and %q", mp3Path, wavPath)
}

func resolveOptionFile(baseRef, content, ext, optionName string) (string, error) {
	if r.IsRemoteFile(content) {
		return content, nil
	}
	basePath, err := resolveOptionFileBase(baseRef, content, optionName)
	if err != nil {
		return "", err
	}

	return basePath + ext, nil
}

func resolveOptionFileBase(baseRef, content, optionName string) (string, error) {
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

	return filepath.Join(baseRef, cleanPath), nil
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("failed to inspect audio file %q: %w", path, err)
}
