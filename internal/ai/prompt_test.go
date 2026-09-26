// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ai

import "testing"

func TestDefaultSystemPrompt(ts *testing.T) {
	if defaultSystemPrompt == "" {
		ts.Fatal("default system prompt is empty")
	}
}
