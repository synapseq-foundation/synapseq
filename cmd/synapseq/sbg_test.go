// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func TestRunSBGConversionUsesDefaultOutputAndWarns(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var status bytes.Buffer
	if err := runSBGConversion([]string{inputPath}, &cli.CLIOptions{}, &status, nil); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !strings.Contains(status.String(), "conversion is approximate") || !strings.Contains(status.String(), "Converted:") || !strings.Contains(status.String(), inputPath[:len(inputPath)-len(".sbg")]+".spsq") {
		t.Fatalf("status = %q", status.String())
	}
	outputPath := filepath.Join(dir, "session.spsq")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(content), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected output:\n%s", content)
	}
}

func TestRunSBGConversionOverwritesOutput(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	outputPath := filepath.Join(dir, "custom.spsq")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("old"), 0o600); err != nil {
		t.Fatalf("write old output: %v", err)
	}
	if err := runSBGConversion([]string{inputPath, outputPath}, &cli.CLIOptions{}, nil, nil); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(content) == "old" {
		t.Fatal("output was not overwritten")
	}
}

func TestRunSBGConversionWritesToStandardOutput(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var output bytes.Buffer
	if err := runSBGConversion([]string{inputPath, "-"}, &cli.CLIOptions{}, nil, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if !strings.Contains(output.String(), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected standard output:\n%s", output.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "session.spsq")); !os.IsNotExist(err) {
		t.Fatalf("default output file was created: %v", err)
	}
}

func TestRunSBGConversionQuietSuppressesStatus(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "session.sbg")
	input := "alpha: 300+10/20\noff: -\nNOW alpha\n+00:01:00 off\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	var status bytes.Buffer
	var output bytes.Buffer
	if err := runSBGConversion([]string{inputPath, "-"}, &cli.CLIOptions{Quiet: true}, &status, &output); err != nil {
		t.Fatalf("runSBGConversion error: %v", err)
	}
	if status.Len() != 0 {
		t.Fatalf("quiet status = %q", status.String())
	}
	if !strings.Contains(output.String(), "00:00:30 alpha steady 0") {
		t.Fatalf("unexpected standard output:\n%s", output.String())
	}
}

func TestSBGConversionWarningRespectsColorSetting(t *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	cli.SetColorEnabled(true)
	var colored bytes.Buffer
	if err := writeSBGConversionWarning(&colored); err != nil {
		t.Fatalf("write colored warning: %v", err)
	}
	if !strings.Contains(colored.String(), "\x1b[") {
		t.Fatalf("expected ANSI warning, got %q", colored.String())
	}

	cli.SetColorEnabled(false)
	var plain bytes.Buffer
	if err := writeSBGConversionWarning(&plain); err != nil {
		t.Fatalf("write plain warning: %v", err)
	}
	if strings.Contains(plain.String(), "\x1b[") {
		t.Fatalf("unexpected ANSI warning: %q", plain.String())
	}
}

func TestSBGConversionSuccessRespectsColorSetting(t *testing.T) {
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	cli.SetColorEnabled(true)
	var colored bytes.Buffer
	if err := writeSBGConversionSuccess(&colored, "session.spsq"); err != nil {
		t.Fatalf("write colored success: %v", err)
	}
	if !strings.Contains(colored.String(), "\x1b[") {
		t.Fatalf("expected ANSI success, got %q", colored.String())
	}

	cli.SetColorEnabled(false)
	var plain bytes.Buffer
	if err := writeSBGConversionSuccess(&plain, "session.spsq"); err != nil {
		t.Fatalf("write plain success: %v", err)
	}
	if strings.Contains(plain.String(), "\x1b[") {
		t.Fatalf("unexpected ANSI success: %q", plain.String())
	}
}
