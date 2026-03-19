//go:build !windows && !js && !wasm

/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

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
