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

func TestCheckToolWithCategory(t *testing.T) {
	check := CheckToolWithCategory("ffmpeg", "Export")
	if check.Name != "ffmpeg" {
		t.Fatalf("expected name ffmpeg, got %q", check.Name)
	}
	if check.Category != "Export" {
		t.Fatalf("expected category Export, got %q", check.Category)
	}
}

func TestDoctorReturnsAllTools(t *testing.T) {
	checks := Doctor()
	expectedTools := map[string]string{
		"ffmpeg": "Export",
		"ffplay": "Playback",
		"git":    "Hub",
		"gh":     "Hub",
	}

	if len(checks) != len(expectedTools) {
		t.Fatalf("expected %d tools, got %d", len(expectedTools), len(checks))
	}

	for _, check := range checks {
		expectedCategory, ok := expectedTools[check.Name]
		if !ok {
			t.Fatalf("unexpected tool %q", check.Name)
		}
		if check.Category != expectedCategory {
			t.Fatalf("expected category %q for %s, got %q", expectedCategory, check.Name, check.Category)
		}
	}
}

func TestDoctorCategories(t *testing.T) {
	categories := DoctorCategories()
	expected := map[string][]string{
		"Export":   {"ffmpeg"},
		"Playback": {"ffplay"},
		"Hub":      {"git", "gh"},
	}

	for cat, tools := range expected {
		gotTools, ok := categories[cat]
		if !ok {
			t.Fatalf("expected category %q", cat)
		}
		if len(gotTools) != len(tools) {
			t.Fatalf("expected %d tools in %s, got %d", len(tools), cat, len(gotTools))
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

func TestGetInstallSuggestionForCategory(t *testing.T) {
	suggestion := getInstallSuggestionForCategory("Export")
	if !strings.Contains(suggestion, "ffmpeg") {
		t.Fatalf("expected suggestion to contain ffmpeg, got %q", suggestion)
	}

	suggestion = getInstallSuggestionForCategory("Hub")
	if !strings.Contains(suggestion, "git") || !strings.Contains(suggestion, "gh") {
		t.Fatalf("expected suggestion to contain git and gh, got %q", suggestion)
	}
}

func TestFormatDoctorOutputAllInstalled(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Category: "Export", Installed: true},
		{Name: "ffplay", Category: "Playback", Installed: true},
		{Name: "git", Category: "Hub", Installed: true},
		{Name: "gh", Category: "Hub", Installed: true},
	}

	var buf strings.Builder
	FormatDoctorOutputTo(checks, &buf)
	output := buf.String()

	if !strings.Contains(output, "Export") {
		t.Fatalf("expected Export category in output, got %s", output)
	}
	if !strings.Contains(output, "Playback") {
		t.Fatalf("expected Playback category in output, got %s", output)
	}
	if !strings.Contains(output, "Hub") {
		t.Fatalf("expected Hub category in output, got %s", output)
	}
	if !strings.Contains(output, "ffmpeg") {
		t.Fatalf("expected ffmpeg in output, got %s", output)
	}
	if !strings.Contains(output, "✓ ffmpeg") {
		t.Fatalf("expected ✓ ffmpeg in output, got %s", output)
	}
	if !strings.Contains(output, "All tools installed") {
		t.Fatalf("expected all tools installed message in output, got %s", output)
	}
	if strings.Contains(output, "Some tools are missing") {
		t.Fatalf("expected no missing tools message when all installed, got %s", output)
	}
}

func TestFormatDoctorOutputSomeMissing(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Category: "Export", Installed: true},
		{Name: "ffplay", Category: "Playback", Installed: false},
		{Name: "git", Category: "Hub", Installed: true},
		{Name: "gh", Category: "Hub", Installed: false},
	}

	var buf strings.Builder
	FormatDoctorOutputTo(checks, &buf)
	output := buf.String()

	if !strings.Contains(output, "Some tools are missing") {
		t.Fatalf("expected missing tools message in output, got %s", output)
	}
	if !strings.Contains(output, "✖ ffplay") {
		t.Fatalf("expected ✖ ffplay in output, got %s", output)
	}
	if !strings.Contains(output, "✖ gh") {
		t.Fatalf("expected ✖ gh in output, got %s", output)
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
		{Name: "ffmpeg", Category: "Export", Installed: true},
	}

	FormatDoctorOutputTo(checks, nil)
}

func TestGroupByCategory(t *testing.T) {
	checks := []ToolCheck{
		{Name: "ffmpeg", Category: "Export", Installed: true},
		{Name: "ffplay", Category: "Playback", Installed: true},
		{Name: "git", Category: "Hub", Installed: false},
	}

	grouped := groupByCategory(checks)

	if len(grouped["Export"]) != 1 {
		t.Fatalf("expected 1 tool in Export, got %d", len(grouped["Export"]))
	}
	if len(grouped["Playback"]) != 1 {
		t.Fatalf("expected 1 tool in Playback, got %d", len(grouped["Playback"]))
	}
	if len(grouped["Hub"]) != 1 {
		t.Fatalf("expected 1 tool in Hub, got %d", len(grouped["Hub"]))
	}
}

func TestFormatDoctorOutputAllMissing(t *testing.T) {
	clistyle.SetColorEnabled(false)
	defer clistyle.SetColorEnabled(true)

	checks := []ToolCheck{
		{Name: "ffmpeg", Category: "Export", Installed: false},
		{Name: "ffplay", Category: "Playback", Installed: false},
		{Name: "git", Category: "Hub", Installed: false},
		{Name: "gh", Category: "Hub", Installed: false},
	}

	var buf strings.Builder
	FormatDoctorOutputTo(checks, &buf)
	output := buf.String()

	if !strings.Contains(output, "Some tools are missing") {
		t.Fatalf("expected missing tools message in output, got %s", output)
	}
	if !strings.Contains(output, "✖ ffmpeg") {
		t.Fatalf("expected ✖ ffmpeg in output, got %s", output)
	}
	if !strings.Contains(output, "✖ ffplay") {
		t.Fatalf("expected ✖ ffplay in output, got %s", output)
	}
	if !strings.Contains(output, "✖ git") {
		t.Fatalf("expected ✖ git in output, got %s", output)
	}
	if !strings.Contains(output, "✖ gh") {
		t.Fatalf("expected ✖ gh in output, got %s", output)
	}
}