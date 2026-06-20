// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
