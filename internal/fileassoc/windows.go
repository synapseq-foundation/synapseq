//go:build windows

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

package fileassoc

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// Helper that recursively deletes a key tree
func deleteRegistryTree(base registry.Key, path string) error {
	k, err := registry.OpenKey(base, path, registry.ALL_ACCESS)
	if err != nil {
		return nil // If not exist, fine
	}
	defer k.Close()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err == nil {
		for _, sub := range subKeys {
			_ = deleteRegistryTree(base, path+`\`+sub)
		}
	}

	err = registry.DeleteKey(base, path)
	if err != nil {
		return fmt.Errorf("failed to delete registry key %s: %w", path, err)
	}

	return nil
}

// CleanSynapSeqWindowsRegistry removes all SynapSeq-related registry keys.
// Safe to run even if some keys don't exist.
func CleanSynapSeqWindowsRegistry() error {
	extKeyPath := `Software\Classes\.spsq`

	extKey, err := registry.OpenKey(registry.CURRENT_USER, extKeyPath, registry.READ)
	if err == nil {
		defer extKey.Close()

		val, _, err := extKey.GetStringValue("")
		if err == nil && val == "SynapSeq.File" {
			registry.DeleteKey(registry.CURRENT_USER, extKeyPath)
		}
	}

	_ = deleteRegistryTree(registry.CURRENT_USER, `Software\Classes\SynapSeq.File`)
	_ = deleteRegistryTree(registry.CURRENT_USER,
		`Software\Classes\SystemFileAssociations\.wav\shell\SynapSeqExtract`)
	_ = deleteRegistryTree(registry.CURRENT_USER,
		`Software\Classes\SystemFileAssociations\.mp3\shell\SynapSeqExtract`)

	return nil
}

// InstallWindowsFileAssociation sets up the Windows registry to associate .spsq files with SynapSeq
func InstallWindowsFileAssociation() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Clean(exe)

	extKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\.spsq`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer extKey.Close()

	if err := extKey.SetStringValue("", "SynapSeq.File"); err != nil {
		return err
	}

	progIDKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer progIDKey.Close()

	progIDKey.SetStringValue("", "SynapSeq Sequence File")

	iconKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File\DefaultIcon`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer iconKey.Close()

	iconKey.SetStringValue("", exePath+",0")

	cmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File\shell\open\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer cmdKey.Close()

	openCmd := `cmd.exe /C synapseq -play "%1" & echo. & pause`
	cmdKey.SetStringValue("", openCmd)

	return nil
}

// InstallWindowsContextMenu adds SynapSeq options to the Windows context menu for .spsq files
func InstallWindowsContextMenu() error {
	base := `Software\Classes\SynapSeq.File\shell`

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get exe path: %w", err)
	}
	exePath := filepath.Clean(exe)

	// ===============================
	// Edit Sequence
	// ===============================
	editKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\EditSequence`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create EditSequence menu: %w", err)
	}
	defer editKey.Close()

	editKey.SetStringValue("", "SynapSeq: Edit sequence")
	editKey.SetStringValue("Icon", exePath+",0")

	editCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\EditSequence\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create EditSequence command: %w", err)
	}
	defer editCmdKey.Close()

	editCmd := `notepad.exe "%1"`
	editCmdKey.SetStringValue("", editCmd)

	// ===============================
	// Test Sequence
	// ===============================
	testKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\TestSequence`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create TestSequence menu: %w", err)
	}
	defer testKey.Close()

	testKey.SetStringValue("", "SynapSeq: Test sequence")
	testKey.SetStringValue("Icon", exePath+",0")

	testCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\TestSequence\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create TestSequence command: %w", err)
	}
	defer testCmdKey.Close()

	testCmd := `cmd.exe /C synapseq -test "%1" & echo. & pause`
	testCmdKey.SetStringValue("", testCmd)

	// ===============================
	// Convert to WAV
	// ===============================
	wavKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\ConvertToWAV`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create ConvertToWAV menu: %w", err)
	}
	defer wavKey.Close()

	wavKey.SetStringValue("", "SynapSeq: Convert to WAV")
	wavKey.SetStringValue("Icon", exePath+",0")

	wavCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\ConvertToWAV\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create ConvertToWAV command: %w", err)
	}
	defer wavCmdKey.Close()

	wavCmd := `cmd.exe /C synapseq "%1" & echo. & pause`
	wavCmdKey.SetStringValue("", wavCmd)

	// ===============================
	// Convert to MP3
	// ===============================
	mp3Key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\ConvertToMP3`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create ConvertToMP3 menu: %w", err)
	}
	defer mp3Key.Close()

	mp3Key.SetStringValue("", "SynapSeq: Convert to MP3")
	mp3Key.SetStringValue("Icon", exePath+",0")

	mp3CmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\ConvertToMP3\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create ConvertToMP3 command: %w", err)
	}
	defer mp3CmdKey.Close()

	mp3Cmd := `cmd.exe /C synapseq -mp3 "%1" & echo. & pause`
	mp3CmdKey.SetStringValue("", mp3Cmd)

	return nil
}

// InstallWindowsExtractMenu adds an "Extract sequence" option to the Windows context menu for .wav and .mp3 files
func InstallWindowsExtractMenu() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Clean(exe)

	// ===============================
	// Extract from WAV
	// ===============================
	wavBase := `Software\Classes\SystemFileAssociations\.wav\shell\SynapSeqExtract`

	wavKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		wavBase,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create wav extract menu: %w", err)
	}
	defer wavKey.Close()

	wavKey.SetStringValue("", "SynapSeq: Extract sequence")
	wavKey.SetStringValue("Icon", exePath+",0")

	wavCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		wavBase+`\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create wav extract command: %w", err)
	}
	defer wavCmdKey.Close()

	wavExtractCmd := `cmd.exe /C synapseq -extract "%1" & echo. & pause`
	wavCmdKey.SetStringValue("", wavExtractCmd)

	// ===============================
	// Extract from MP3
	// ===============================
	mp3Base := `Software\Classes\SystemFileAssociations\.mp3\shell\SynapSeqExtract`

	mp3Key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		mp3Base,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create mp3 extract menu: %w", err)
	}
	defer mp3Key.Close()

	mp3Key.SetStringValue("", "SynapSeq: Extract sequence")
	mp3Key.SetStringValue("Icon", exePath+",0")

	mp3CmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		mp3Base+`\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create mp3 extract command: %w", err)
	}
	defer mp3CmdKey.Close()

	mp3ExtractCmd := `cmd.exe /C synapseq -extract -mp3 "%1" & echo. & pause`
	mp3CmdKey.SetStringValue("", mp3ExtractCmd)

	return nil
}
