// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build !darwin || !arm64

package core

import "fmt"

func appleFoundationAIModel() string {
	return ""
}

func validateAppleFoundationProvider() error {
	return fmt.Errorf("AI provider %q is supported only on macOS with Apple Silicon (darwin/arm64)", AIProviderAppleFoundation)
}
