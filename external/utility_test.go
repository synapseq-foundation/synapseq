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

package external

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

func TestUtilityPath_MissingCustomPathReturnsDiagnostic(ts *testing.T) {
	missing := filepath.Join(ts.TempDir(), "missing-utility")

	_, err := utilityPath(missing)
	if err == nil {
		ts.Fatalf("expected error for missing utility path")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}

	if !strings.Contains(diagnostic.Message, "external utility not found") {
		ts.Fatalf("unexpected diagnostic message: %q", diagnostic.Message)
	}
	if !strings.Contains(diagnostic.Hint, "install the utility") {
		ts.Fatalf("unexpected diagnostic hint: %q", diagnostic.Hint)
	}
}

func TestUtilityPath_NotExecutableReturnsDiagnostic(ts *testing.T) {
	path := filepath.Join(ts.TempDir(), "fake-tool")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho hi\n"), 0644); err != nil {
		ts.Fatalf("failed to create fake utility: %v", err)
	}

	_, err := utilityPath(path)
	if err == nil {
		ts.Fatalf("expected error for non-executable utility")
	}

	diagnostic, ok := diag.As(err)
	if !ok {
		ts.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}

	if !strings.Contains(diagnostic.Message, "not executable") {
		ts.Fatalf("unexpected diagnostic message: %q", diagnostic.Message)
	}
	if !strings.Contains(diagnostic.Hint, "mark the file as executable") {
		ts.Fatalf("unexpected diagnostic hint: %q", diagnostic.Hint)
	}
}
