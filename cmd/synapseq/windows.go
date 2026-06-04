//go:build windows

// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
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

	if !quiet {
		fmt.Println(cli.SuccessText("Successfully installed .spsq file association with SynapSeq."))
	}
	return nil
}

// uninstallWindowsFileAssociation removes the file association for .spsq files on Windows
func uninstallWindowsFileAssociation(quiet bool) error {
	if err := fileassoc.CleanSynapSeqWindowsRegistry(); err != nil {
		return err
	}

	if !quiet {
		fmt.Println(cli.SuccessText("Successfully removed .spsq file association with SynapSeq."))
	}
	return nil
}
