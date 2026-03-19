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
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// writeTempFile creates a temporary file with the given content
func writeTempFile(ts *testing.T, name string, content []byte) string {
	ts.Helper()

	tmpDir := ts.TempDir()
	path := filepath.Join(tmpDir, name)

	if err := os.WriteFile(path, content, 0644); err != nil {
		ts.Fatalf("failed to write temp file: %v", err)
	}

	return path
}

// makeBigContent generates a slice of lines totaling at least minBytes
func makeBigContent(lineLen, minBytes int) []byte {
	var buf bytes.Buffer
	line := strings.Repeat("x", lineLen) + "\n"

	for buf.Len() < minBytes {
		buf.WriteString(line)
	}

	return buf.Bytes()
}

func TestGetFile_LocalFile_Text(ts *testing.T) {
	content := []byte("line 1\nline 2\nline 3\n")
	path := writeTempFile(ts, "test.spsq", content)

	got, err := GetFile(path, t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if !bytes.Equal(got, content) {
		ts.Errorf("expected %q, got %q", content, got)
	}
}

func TestGetFile_LocalFile_NotFound(ts *testing.T) {
	path := filepath.Join(ts.TempDir(), "missing.spsq")

	_, err := GetFile(path, t.FormatText)
	if err == nil {
		ts.Fatal("expected error for missing file, got nil")
	}

	if !strings.Contains(err.Error(), "error opening file") {
		ts.Errorf("unexpected error message: %v", err)
	}
}

func TestGetFile_LocalFile_Truncate_Text(ts *testing.T) {
	const maxSize = t.MaxTextFileSize
	bigContent := makeBigContent(100, int(maxSize+8192))
	path := writeTempFile(ts, "big.spsq", bigContent)

	got, err := GetFile(path, t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if len(got) != int(maxSize) {
		ts.Errorf("expected %d bytes (truncated), got %d", maxSize, len(got))
	}

	if !bytes.Equal(got, bigContent[:maxSize]) {
		ts.Error("truncated content does not match expected prefix")
	}
}

func TestGetFile_LocalFile_Truncate_WAV(ts *testing.T) {
	const maxSize = t.MaxWavFileSize
	bigContent := makeBigContent(100, int(maxSize+4096))
	path := writeTempFile(ts, "big.wav", bigContent)

	got, err := GetFile(path, t.FormatWAV)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if len(got) != int(maxSize) {
		ts.Errorf("expected %d bytes (truncated), got %d", maxSize, len(got))
	}
}

func TestGetFile_Stdin(ts *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		ts.Fatalf("failed to create pipe: %v", err)
	}
	defer r.Close()

	// Replace stdin
	os.Stdin = r

	// Write content to pipe
	content := []byte("stdin content\n")
	go func() {
		defer w.Close()
		_, _ = w.Write(content)
	}()

	// Test GetFile with stdin
	got, err := GetFile("-", t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if !bytes.Equal(got, content) {
		ts.Errorf("expected %q, got %q", content, got)
	}
}

func TestGetFile_Stdin_Truncate(ts *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		ts.Fatalf("failed to create pipe: %v", err)
	}
	defer r.Close()

	os.Stdin = r

	const maxSize = t.MaxTextFileSize
	bigContent := makeBigContent(100, int(maxSize+8192))

	go func() {
		defer w.Close()
		_, _ = w.Write(bigContent)
	}()

	got, err := GetFile("-", t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if len(got) != int(maxSize) {
		ts.Errorf("expected %d bytes (truncated), got %d", maxSize, len(got))
	}
}

func TestGetFile_HTTP_Text(ts *testing.T) {
	content := []byte("remote content\n")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	got, err := GetFile(srv.URL+"/test.spsq", t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if !bytes.Equal(got, content) {
		ts.Errorf("expected %q, got %q", content, got)
	}
}

func TestGetFile_HTTP_WAV(ts *testing.T) {
	// Minimal WAV header (44 bytes)
	content := []byte("RIFF" + string([]byte{0, 0, 0, 0}) + "WAVEfmt " +
		string(make([]byte, 32)))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	got, err := GetFile(srv.URL+"/test.wav", t.FormatWAV)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if !bytes.Equal(got, content) {
		ts.Errorf("content mismatch")
	}
}

func TestGetFile_HTTP_Truncate(ts *testing.T) {
	const maxSize = t.MaxTextFileSize
	bigContent := makeBigContent(100, int(maxSize+8192))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(bigContent)
	}))
	defer srv.Close()

	got, err := GetFile(srv.URL+"/big.spsq", t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if len(got) != int(maxSize) {
		ts.Errorf("expected %d bytes (truncated), got %d", maxSize, len(got))
	}
}

func TestGetFile_HTTP_NetworkError(ts *testing.T) {
	// Use invalid URL
	_, err := GetFile("http://localhost:1/nonexistent", t.FormatText)
	if err == nil {
		ts.Fatal("expected network error, got nil")
	}

	if !strings.Contains(err.Error(), "error fetching remote file") {
		ts.Errorf("unexpected error: %v", err)
	}
}

func TestGetFile_UnsupportedFormat(ts *testing.T) {
	path := writeTempFile(ts, "test.txt", []byte("content"))

	// Use an invalid format value
	_, err := GetFile(path, t.FileFormat(999))
	if err == nil {
		ts.Fatal("expected error for unsupported format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported file type") {
		ts.Errorf("unexpected error: %v", err)
	}
}

func TestGetFile_EmptyFile(ts *testing.T) {
	path := writeTempFile(ts, "empty.spsq", []byte{})

	got, err := GetFile(path, t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	if len(got) != 0 {
		ts.Errorf("expected empty content, got %d bytes", len(got))
	}
}

func TestIsRemoteFile(ts *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"http URL", "http://example.com/file.spsq", true},
		{"https URL", "https://example.com/file.spsq", true},
		{"local path", "/path/to/file.spsq", false},
		{"relative path", "file.spsq", false},
		{"stdin", "-", false},
		{"ftp URL", "ftp://example.com/file.spsq", false},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func(ts *testing.T) {
			got := IsRemoteFile(tt.filePath)
			if got != tt.want {
				ts.Errorf("IsRemoteFile(%q) = %v, want %v", tt.filePath, got, tt.want)
			}
		})
	}
}
