// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ai

import (
	_ "embed"
	"strings"
)

//go:embed prompts/default.txt
var defaultSystemPromptText string

var defaultSystemPrompt = strings.TrimSuffix(defaultSystemPromptText, "\n")
