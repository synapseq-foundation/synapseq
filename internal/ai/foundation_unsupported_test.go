// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build !darwin || !arm64

package ai

import (
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestAppleFoundationIsUnavailableOutsideDarwinARM64(ts *testing.T) {
	_, err := New(Config{Model: "system", Provider: t.AIProviderAppleFoundation})
	if err == nil || !strings.Contains(err.Error(), "supported only on macOS with Apple Silicon") {
		ts.Fatalf("unexpected Apple Foundation error: %v", err)
	}
}
