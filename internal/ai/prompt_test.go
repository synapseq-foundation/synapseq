// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ai

import (
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestSystemPromptForProvider(ts *testing.T) {
	if got := systemPromptForProvider(t.AIProviderDefault, "sleep"); got != defaultSystemPrompt {
		ts.Fatal("default provider did not select the default prompt")
	}
	if got := systemPromptForProvider(t.AIProviderAppleFoundation, "make a sequence"); got != appleFoundationSystemPrompt {
		ts.Fatal("Apple Foundation provider did not select its base prompt")
	}
	if got := systemPromptForProvider(t.AIProvider("unknown"), "sleep"); got != defaultSystemPrompt {
		ts.Fatal("unknown provider did not fall back to the default prompt")
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
		prompt := systemPromptForProvider(t.AIProviderAppleFoundation, test.request)
		if !strings.Contains(prompt, test.rule) {
			ts.Errorf("prompt for %q does not include %q", test.request, test.rule)
		}
	}
}
