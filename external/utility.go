//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"runtime"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
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
		return "", fmt.Errorf("utility path cannot be empty")
	}

	filePath, err := exec.LookPath(utilPath)
	if err == nil {
		return filePath, nil
	}

	fileInfo, err := os.Stat(utilPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("executable not found at custom path: %s", utilPath)
		}
		return "", fmt.Errorf("error checking path: %s, error: %v", utilPath, err)
	}

	if fileInfo.Mode().IsRegular() {
		if runtime.GOOS == "windows" || fileInfo.Mode()&0111 != 0 {
			return utilPath, nil
		}
	}

	return "", fmt.Errorf("file at path is not executable: %s", utilPath)
}

// startPipeCmd starts the given command and pipes loadedCtx streaming audio to its stdin.
func startPipeCmd(cmd *exec.Cmd, loadedCtx *synapseq.LoadedContext) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return err
	}

	streamErr := loadedCtx.Stream(stdin)

	stdin.Close()

	waitErr := cmd.Wait()

	if streamErr != nil {
		return streamErr
	}

	if waitErr != nil {
		return waitErr
	}

	return nil
}
