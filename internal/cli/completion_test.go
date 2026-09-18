// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cli

import "testing"

func TestCompletionFlagsIncludeAIOptions(ts *testing.T) {
	for _, flag := range []string{
		"ai",
		"ai-model",
		"ai-base-url",
		"ai-provider",
		"ai-temperature",
		"ai-timeout",
	} {
		if _, ok := completionFlags[flag]; !ok {
			ts.Errorf("completion flags do not include %q", flag)
		}
	}
}
