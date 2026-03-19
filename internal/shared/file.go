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

package shared

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// readFile reads a file from the given reader up to maxSize bytes
func readFile(r io.Reader, maxSize int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxSize))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// copyFile copies a single file, preserving permissions
func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

// getRemoteFile fetches a remote file and validates its content type and size
func getRemoteFile(url string, maxSize int64) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching remote file: %v", err)
	}
	defer resp.Body.Close()

	data, err := readFile(resp.Body, maxSize)
	if err != nil {
		return nil, fmt.Errorf("error reading remote file: %v", err)
	}

	return data, nil
}

// IsRemoteFile checks if the given file path is a remote URL
func IsRemoteFile(filePath string) bool {
	return strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://")
}

// GetFile retrieves a file from a local path or URL based on the specified type
func GetFile(filePath string, typ t.FileFormat) ([]byte, error) {
	maxSize := int64(0)
	switch typ {
	case t.FormatText:
		maxSize = t.MaxTextFileSize
	case t.FormatWAV:
		maxSize = t.MaxWavFileSize
	}

	if maxSize == 0 {
		return nil, fmt.Errorf("unsupported file type: %s", typ.String())
	}

	switch {
	case filePath == "-":
		data, err := readFile(os.Stdin, maxSize)
		if err != nil {
			return nil, fmt.Errorf("error reading from stdin: %v", err)
		}
		return data, nil

	case IsRemoteFile(filePath):
		return getRemoteFile(filePath, maxSize)

	default:
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer f.Close()

		data, err := readFile(f, maxSize)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
		return data, nil
	}
}

// CopyDir recursively copies a directory from src to dst.
// It preserves file permissions and structure.
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// create target directory if needed
			return os.MkdirAll(targetPath, info.Mode())
		}

		// copy file contents
		if err := copyFile(path, targetPath, info.Mode()); err != nil {
			return err
		}

		return nil
	})
}
