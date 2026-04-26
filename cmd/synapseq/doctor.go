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
	"runtime"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

type ToolCheck struct {
	Name      string
	Installed bool
	Error     string
}

type DoctorOutput struct {
	writer io.Writer
}

func NewDoctorOutput(w io.Writer) *DoctorOutput {
	if w == nil {
		w = io.Discard
	}
	return &DoctorOutput{writer: w}
}

const checkMark = "\u2714"
const crossMark = "\u2716"

func Doctor() []ToolCheck {
	tools := []string{"ffmpeg", "ffplay", "git", "gh"}
	checks := make([]ToolCheck, len(tools))

	for i, tool := range tools {
		checks[i] = CheckTool(tool)
	}

	return checks
}

func CheckTool(name string) ToolCheck {
	_, err := exec.LookPath(name)
	if err == nil {
		return ToolCheck{Name: name, Installed: true}
	}
	return ToolCheck{Name: name, Installed: false, Error: err.Error()}
}

func FormatDoctorOutput(checks []ToolCheck) {
	FormatDoctorOutputTo(checks, os.Stdout)
}

func FormatDoctorOutputTo(checks []ToolCheck, w io.Writer) {
	if w == nil {
		w = io.Discard
	}

	for _, check := range checks {
		if check.Installed {
			fmt.Fprintf(w, "%s %s installed\n", cli.SuccessText(checkMark), check.Name)
		} else {
			fmt.Fprintf(w, "%s %s not installed\n", cli.ErrorText(crossMark), check.Name)
		}
	}

	hasMissing := false
	for _, check := range checks {
		if !check.Installed {
			hasMissing = true
			break
		}
	}

	if hasMissing {
		fmt.Fprintln(w)
		fmt.Fprintln(w, cli.Section("Suggested fixes:"))
		for _, check := range checks {
			if !check.Installed {
				suggestion := getInstallSuggestion(check.Name)
				fmt.Fprintf(w, "  %s\n", cli.Command(suggestion))
			}
		}
	}
}

func getInstallSuggestion(tool string) string {
	installCmd := getInstallCommand()
	return fmt.Sprintf("%s install %s", installCmd, tool)
}

func getInstallCommand() string {
	switch runtime.GOOS {
	case "windows":
		return "winget"
	default:
		return "brew"
	}
}

func runDoctor() error {
	checks := Doctor()
	FormatDoctorOutputTo(checks, os.Stdout)
	return nil
}