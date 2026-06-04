//go:build !windows

// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
)

// installWindowsFileAssociation is disabled for non-Windows builds
func installWindowsFileAssociation(quiet bool) error {
	return fmt.Errorf("this build does not support Windows file association installation")
}

// uninstallWindowsFileAssociation is disabled for non-Windows builds
func uninstallWindowsFileAssociation(quiet bool) error {
	return fmt.Errorf("this build does not support Windows file association removal")
}
