//go:build !js && !wasm

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
	"runtime"
	"strings"
	"testing"

	clistyle "github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func TestCheckToolInstalled(t *testing.T) {
	check := CheckTool("ffmpeg")
	if !check.Installed {
		t.Fatalf("expected ffmpeg to be installed on this system")
	}
	if check.Name != "ffmpeg" {
		t.Fatalf("expected name ffmpeg, got %q", check.Name)
	}
}

func TestCheckToolNotInstalled(t *testing.T) {
	check := CheckTool("nonexistent-tool-xyz")
	if check.Installed {
		t.Fatalf("expected nonexistent-tool-xyz to not be installed")
	}
	if check.Name != "nonexistent-tool-xyz" {
		t.Fatalf("expected name nonexistent-tool-xyz, got %q", check.Name)
	}
	if check.Error == "" {
		t.Fatalf("expected error message for missing tool")
	}
}

func TestDoctorReturnsAllTools(t *testing.T) {
	checks := Doctor()
	expectedTools := []string{"ffmpeg", "ffplay", "git", "gh"}

	if len(checks) != len(expectedTools) {
		t.Fatalf("expected %d tools, got %d", len(expectedTools), len(checks))
	}

	names := make(map[string]bool)
	for _, check := range checks {
		names[check.Name] = true
	}

	for _, tool := range expectedTools {
		if !names[tool] {
			t.Fatalf("expected tool %q in results", tool)
		}
	}
}

func TestGetInstallCommandBrew(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("skipping brew test on windows")
	}

	cmd := getInstallCommand()
	if cmd != "brew" {
		t.Fatalf("expected brew on darwin/linux, got %q", cmd)
	}
}

func TestGetInstallCommandWinget(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping winget test on non-windows")
	}

	cmd := getInstallCommand()
	if cmd != "winget" {
		t.Fatalf("expected winget on windows, got %q", cmd)
	}
}

func TestGetInstallSuggestion(t *testing.T) {
	suggestion := getInstallSuggestion("ffmpeg")
	if !strings.Contains(suggestion, "ffmpeg") {
		t.Fatalf("expected suggestion to contain ffmpeg, got %q", suggestion)
	}
	if runtime.GOOS == "windows" {
		if !strings.Contains(suggestion, "winget") {
			t.Fatalf("expected suggestion to contain winget on windows, got %q", suggestion)
		}
	} else {
		if !strings.Contains(suggestion, "brew") {
			t.Fatalf("expected suggestion to contain brew on non-windows, got %q", suggestion)
		}
	}
}

func TestFormatDoctorOutputInstalled(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Installed: true},
		{Name: "ffplay", Installed: true},
	}

	var buf strings.Builder
	FormatDoctorOutputTo(checks, &buf)
	output := buf.String()

	if !strings.Contains(output, "ffmpeg installed") {
		t.Fatalf("expected ffmpeg installed in output, got %s", output)
	}
	if !strings.Contains(output, "ffplay installed") {
		t.Fatalf("expected ffplay installed in output, got %s", output)
	}
	if strings.Contains(output, "Suggested fixes") {
		t.Fatalf("expected no suggested fixes when all tools installed, got %s", output)
	}
}

func TestFormatDoctorOutputMissing(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Installed: true},
		{Name: "gh", Installed: false},
	}

	var buf strings.Builder
	FormatDoctorOutputTo(checks, &buf)
	output := buf.String()

	if !strings.Contains(output, "ffmpeg installed") {
		t.Fatalf("expected ffmpeg installed in output, got %s", output)
	}
	if !strings.Contains(output, "gh not installed") {
		t.Fatalf("expected gh not installed in output, got %s", output)
	}
	if !strings.Contains(output, "Suggested fixes") {
		t.Fatalf("expected suggested fixes section, got %s", output)
	}
	if !strings.Contains(output, "brew") && !strings.Contains(output, "winget") {
		t.Fatalf("expected install suggestion in output, got %s", output)
	}
}

func TestResolveSpecialCommandDoctor(t *testing.T) {
	command := clistyle.ResolveSpecialCommand(&clistyle.CLIOptions{ShowDoctor: true}, nil)
	if command.Kind != clistyle.SpecialCommandDoctor {
		t.Fatalf("expected doctor command, got %q", command.Kind)
	}
}

func TestFormatDoctorOutputNilWriter(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Installed: true},
	}

	FormatDoctorOutputTo(checks, nil)
}

func TestNewDoctorOutputNilWriter(t *testing.T) {
	d := NewDoctorOutput(nil)
	if d.writer == nil {
		t.Fatalf("expected nil writer to default to discard, got nil")
	}
}