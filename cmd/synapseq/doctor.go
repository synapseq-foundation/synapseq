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
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

type ToolCheck struct {
	Name      string
	Category  string
	Installed bool
	Error     string
}

func DoctorCategories() map[string][]string {
	return map[string][]string{
		"Export":   {"ffmpeg"},
		"Playback": {"ffplay"},
	}
}

func Doctor() []ToolCheck {
	categories := DoctorCategories()
	var checks []ToolCheck

	for category, tools := range categories {
		for _, tool := range tools {
			checks = append(checks, CheckToolWithCategory(tool, category))
		}
	}

	return checks
}

func CheckTool(name string) ToolCheck {
	return CheckToolWithCategory(name, "")
}

func CheckToolWithCategory(name, category string) ToolCheck {
	_, err := exec.LookPath(name)
	if err == nil {
		return ToolCheck{Name: name, Category: category, Installed: true}
	}
	return ToolCheck{Name: name, Category: category, Installed: false, Error: err.Error()}
}

func FormatDoctorOutput(checks []ToolCheck) {
	FormatDoctorOutputTo(checks, os.Stdout)
}

func FormatDoctorOutputTo(checks []ToolCheck, w io.Writer) {
	if w == nil {
		w = io.Discard
	}

	categorized := groupByCategory(checks)
	categories := []string{"Export", "Playback"}

	for _, category := range categories {
		tools := categorized[category]
		if len(tools) == 0 {
			continue
		}

		fmt.Fprintln(w, cli.Label(category))

		for _, check := range tools {
			status := cli.SuccessText("✓")
			if !check.Installed {
				status = cli.ErrorText("✖")
			}
			fmt.Fprintf(w, " %s %s\n", status, check.Name)
		}
		fmt.Fprintln(w)
	}

	hasMissing := false
	for _, check := range checks {
		if !check.Installed {
			hasMissing = true
			break
		}
	}

	if hasMissing {
		fmt.Fprintln(w, cli.ErrorText("Some tools are missing. Install them to use all features."))
	} else {
		fmt.Fprintln(w, cli.SuccessText("All tools installed. You're ready to go!"))
	}
}

func groupByCategory(checks []ToolCheck) map[string][]ToolCheck {
	result := make(map[string][]ToolCheck)
	for _, check := range checks {
		result[check.Category] = append(result[check.Category], check)
	}
	return result
}

func runDoctor() error {
	checks := Doctor()
	FormatDoctorOutputTo(checks, os.Stdout)
	return nil
}
