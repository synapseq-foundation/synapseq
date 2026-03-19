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

package external

import "os/exec"

// baseUtility represents a base external utility
type baseUtility struct{ path string }

// Path returns the path of the external utility
func (bu *baseUtility) Path() string {
	return bu.path
}

// Command creates an exec.Cmd for the utility with given arguments
func (bu *baseUtility) Command(args ...string) *exec.Cmd {
	return exec.Command(bu.path, args...)
}
