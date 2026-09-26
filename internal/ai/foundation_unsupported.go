// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build !darwin || !arm64

package ai

import (
	"context"
	"fmt"
)

type foundationClient struct{}

func newFoundationClient() (foundationClient, error) {
	return foundationClient{}, fmt.Errorf("AI provider %q is supported only on macOS with Apple Silicon (darwin/arm64)", "apple-foundation")
}

func (foundationClient) generate(context.Context, Config, string) (string, error) {
	return "", fmt.Errorf("AI provider %q is supported only on macOS with Apple Silicon (darwin/arm64)", "apple-foundation")
}
