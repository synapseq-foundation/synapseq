// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ai

import (
	_ "embed"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

var systemPrompts = map[t.AIProvider]string{
	t.AIProviderDefault:         defaultSystemPrompt,
	t.AIProviderAppleFoundation: appleFoundationSystemPrompt,
}

//go:embed prompts/apple-foundation.txt
var appleFoundationSystemPromptText string

var appleFoundationSystemPrompt = strings.TrimSuffix(appleFoundationSystemPromptText, "\n")

//go:embed prompts/apple-relaxation.txt
var appleRelaxationProfilePromptText string

//go:embed prompts/apple-sleep.txt
var appleSleepProfilePromptText string

//go:embed prompts/apple-focus.txt
var appleFocusProfilePromptText string

//go:embed prompts/apple-alert.txt
var appleAlertProfilePromptText string

var appleFoundationProfilePrompts = map[string]string{
	"relaxation": strings.TrimSuffix(appleRelaxationProfilePromptText, "\n"),
	"sleep":      strings.TrimSuffix(appleSleepProfilePromptText, "\n"),
	"focus":      strings.TrimSuffix(appleFocusProfilePromptText, "\n"),
	"alert":      strings.TrimSuffix(appleAlertProfilePromptText, "\n"),
}

func systemPromptForProvider(provider t.AIProvider, request string) string {
	prompt, ok := systemPrompts[provider]
	if !ok {
		return systemPrompts[t.AIProviderDefault]
	}
	if provider != t.AIProviderAppleFoundation {
		return prompt
	}

	profile := appleFoundationProfile(request)
	if profile == "" {
		return prompt
	}

	return prompt + "\n\n" + appleFoundationProfilePrompts[profile]
}

func appleFoundationProfile(request string) string {
	words := map[string]bool{}
	for _, word := range strings.FieldsFunc(strings.ToLower(request), func(r rune) bool {
		return r < 'a' || r > 'z'
	}) {
		words[word] = true
	}

	for _, profile := range []string{"sleep", "focus", "relaxation", "alert"} {
		if words[profile] {
			return profile
		}
	}
	if words["alertness"] {
		return "alert"
	}

	return ""
}

//go:embed prompts/default.txt
var defaultSystemPromptText string

var defaultSystemPrompt = strings.TrimSuffix(defaultSystemPromptText, "\n")
