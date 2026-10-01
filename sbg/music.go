// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func currentRelativeMusicPath(value string) (string, error) {
	if strings.Contains(value, "://") {
		return "", fmt.Errorf("music source must be a local path")
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("music source must use '/' path separators")
	}
	if hasParentTraversal(value) {
		return "", fmt.Errorf("music source must not contain parent directory traversal")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}
	absolutePath, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("resolve music source %q: %w", value, err)
	}
	relativePath, err := filepath.Rel(cwd, absolutePath)
	if err != nil {
		return "", fmt.Errorf("relativize music source %q: %w", value, err)
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("music source must be inside the current directory")
	}

	basePath := strings.TrimSuffix(relativePath, filepath.Ext(relativePath))
	if basePath == "." || basePath == "" {
		return "", fmt.Errorf("music source must name a file")
	}
	return filepath.ToSlash(basePath), nil
}

func hasParentTraversal(path string) bool {
	for _, part := range strings.Split(path, "/") {
		if part == ".." {
			return true
		}
	}
	return false
}
