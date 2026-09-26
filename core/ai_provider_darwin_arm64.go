// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build darwin && arm64

package core

func appleFoundationAIModel() string {
	return "system"
}

func validateAppleFoundationProvider() error {
	return nil
}
