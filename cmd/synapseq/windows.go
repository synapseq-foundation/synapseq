//go:build windows

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

package main

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/fileassoc"
)

// installWindowsFileAssociation sets up the file association for .spsq files on Windows
func installWindowsFileAssociation(quiet bool) error {
	_ = fileassoc.CleanSynapSeqWindowsRegistry()

	if err := fileassoc.InstallWindowsFileAssociation(); err != nil {
		return err
	}
	if err := fileassoc.InstallWindowsContextMenu(); err != nil {
		return err
	}
	if err := fileassoc.InstallWindowsExtractMenu(); err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Successfully installed .spsq file association with SynapSeq.")
	}
	return nil
}

// uninstallWindowsFileAssociation removes the file association for .spsq files on Windows
func uninstallWindowsFileAssociation(quiet bool) error {
	if err := fileassoc.CleanSynapSeqWindowsRegistry(); err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Successfully removed .spsq file association with SynapSeq.")
	}
	return nil
}
