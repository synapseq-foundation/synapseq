// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

var promptDurationPattern = regexp.MustCompile(`(?i)\b(\d+)\s*(hours?|hrs?|h|minutes?|mins?|m)\b`)

func runAI(prompt string, args []string, opts *cli.CLIOptions, statusWriter, outputWriter io.Writer) error {
	if opts == nil {
		return fmt.Errorf("CLI options are nil")
	}
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("AI prompt cannot be empty")
	}
	if len(args) > 1 {
		return fmt.Errorf("invalid AI generation arguments\nUsage: synapseq -ai <prompt> [output.spsq|-]")
	}

	outputPath := aiOutputPath(prompt)
	if len(args) == 1 {
		outputPath = args[0]
	}
	if outputPath != "-" && !strings.EqualFold(filepath.Ext(outputPath), ".spsq") {
		return fmt.Errorf("SPSQ output must use the .spsq extension: %q", outputPath)
	}
	if outputPath != "-" {
		if err := ensureAIOutputAvailable(outputPath); err != nil {
			return err
		}
	}

	progress := startAIProgress(statusWriter, opts.Quiet)

	loaded, err := synapseq.NewAppContext().AI(prompt, &synapseq.AIOptions{
		Model:   opts.AIModel,
		BaseURL: opts.AIBaseURL,
	})
	progress.Stop()
	if err != nil {
		return err
	}

	content := loaded.RawContent()
	if outputPath == "-" {
		if outputWriter == nil {
			return fmt.Errorf("SPSQ output writer is nil")
		}
		_, err := outputWriter.Write(content)
		return err
	}
	if err := writeNewAISequence(outputPath, content); err != nil {
		return err
	}
	if !opts.Quiet && statusWriter != nil {
		_, err := fmt.Fprintln(statusWriter, cli.SuccessText("Generated:"), cli.Command(outputPath))
		return err
	}

	return nil
}

func aiOutputPath(prompt string) string {
	duration := promptDuration(prompt)
	words := strings.FieldsFunc(strings.ToLower(prompt), func(r rune) bool {
		return r < 'a' || r > 'z'
	})

	ignored := map[string]bool{
		"a": true, "an": true, "and": true, "create": true, "generate": true,
		"for": true, "hour": true, "hours": true, "minute": true, "minutes": true,
		"of": true, "sequence": true, "the": true,
	}
	nameParts := []string{}
	for _, word := range words {
		if !ignored[word] {
			nameParts = append(nameParts, word)
		}
	}
	if len(nameParts) == 0 {
		nameParts = append(nameParts, "generated")
	}
	if duration != "" {
		nameParts = append(nameParts, duration)
	}

	name := strings.Join(nameParts, "-")
	if len(name) > 48 {
		name = name[:48]
		name = strings.TrimRight(name, "-")
	}

	return name + ".spsq"
}

func promptDuration(prompt string) string {
	matches := promptDurationPattern.FindStringSubmatch(prompt)
	if len(matches) != 3 {
		return ""
	}

	unit := strings.ToLower(matches[2])
	if strings.HasPrefix(unit, "h") {
		return matches[1] + "h"
	}

	return matches[1] + "m"
}

func writeNewAISequence(outputPath string, content []byte) error {
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("AI output already exists: %q", outputPath)
		}
		return fmt.Errorf("write generated sequence %q: %w", outputPath, err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("write generated sequence %q: %w", outputPath, err)
	}

	return nil
}

func ensureAIOutputAvailable(outputPath string) error {
	_, err := os.Lstat(outputPath)
	if err == nil {
		return fmt.Errorf("AI output already exists: %q", outputPath)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("check AI output %q: %w", outputPath, err)
	}

	return nil
}
