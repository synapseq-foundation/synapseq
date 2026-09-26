// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build darwin && arm64

package ai

import (
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestAppleFoundationDoesNotRequireAPIKey(ts *testing.T) {
	_, err := New(Config{Model: "system", Provider: t.AIProviderAppleFoundation})
	if err != nil && strings.Contains(err.Error(), "SYNAPSEQ_AI_API_KEY") {
		ts.Fatalf("Apple Foundation unexpectedly required an API key: %v", err)
	}
}

func TestAppleFoundationPromptIncludesProfileRules(ts *testing.T) {
	tests := []struct {
		request string
		rule    string
	}{
		{request: "Create relaxation", rule: "beta first and alpha second"},
		{request: "Create sleep", rule: "theta first and delta second"},
		{request: "Create focus", rule: "alpha first and beta second"},
		{request: "Create alert session", rule: "beta first and gamma second"},
		{request: "Create alertness session", rule: "beta first and gamma second"},
	}

	for _, test := range tests {
		prompt := appleFoundationSystemPromptForRequest(test.request)
		if !strings.Contains(prompt, test.rule) {
			ts.Errorf("prompt for %q does not include %q", test.request, test.rule)
		}
	}
}
