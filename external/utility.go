//go:build !wasm

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

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

// newUtility creates a new baseUtility instance after validating the utility path
func newUtility(utilPath string) (*baseUtility, error) {
	path, err := utilityPath(utilPath)
	if err != nil {
		return nil, err
	}

	return &baseUtility{path: path}, nil
}

// utilityPath checks and returns the absolute path of the given utility executable
func utilityPath(utilPath string) (string, error) {
	if utilPath == "" {
		return "", diag.Validation("external utility path cannot be empty").WithHint("pass a utility name from PATH or a custom executable path")
	}

	filePath, err := exec.LookPath(utilPath)
	if err == nil {
		return filePath, nil
	}

	fileInfo, err := os.Stat(utilPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", diag.Validation(fmt.Sprintf("external utility not found: %s", utilPath)).WithHint("install the utility or pass a valid executable path")
		}
		return "", diag.Wrap(diag.KindIO, fmt.Sprintf("failed to inspect external utility path: %s", utilPath), err).WithHint("check file permissions and that the path is accessible")
	}

	if fileInfo.Mode().IsRegular() {
		if runtime.GOOS == "windows" || fileInfo.Mode()&0111 != 0 {
			return utilPath, nil
		}
	}

	return "", diag.Validation(fmt.Sprintf("external utility is not executable: %s", utilPath)).WithHint("mark the file as executable or point to the correct binary")
}

// startPipeCmd starts the given command and pipes loadedCtx streaming audio to its stdin.
func startPipeCmd(cmd *exec.Cmd, loadedCtx *synapseq.LoadedContext) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return diag.Wrap(diag.KindIO, fmt.Sprintf("failed to prepare stdin pipe for %s", utilityDisplayName(cmd)), err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return diag.Wrap(diag.KindIO, fmt.Sprintf("failed to start %s", utilityDisplayName(cmd)), err).WithHint("confirm that the external utility exists and can be executed in this environment")
	}

	streamErr := loadedCtx.Stream(stdin)

	stdin.Close()

	waitErr := cmd.Wait()

	if streamErr != nil {
		return diag.Wrap(diag.KindIO, fmt.Sprintf("failed while streaming audio to %s", utilityDisplayName(cmd)), streamErr)
	}

	if waitErr != nil {
		return diag.Wrap(diag.KindIO, fmt.Sprintf("%s exited with an error", utilityDisplayName(cmd)), waitErr).WithHint("see the external utility output above for details")
	}

	return nil
}

func utilityDisplayName(cmd *exec.Cmd) string {
	if cmd == nil || cmd.Path == "" {
		return "external utility"
	}
	return filepath.Base(cmd.Path)
}
