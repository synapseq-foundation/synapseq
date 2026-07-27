// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/sbg"
)

const sbgConversionWarning = "SBaGen conversion is approximate and may contain errors; review the generated SPSQ before use."

func runSBGConversion(args []string, quiet bool, statusWriter, outputWriter io.Writer) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid SBaGen conversion arguments\nUsage: synapseq -sbg <file.sbg> [output.spsq]")
	}
	inputPath := args[0]
	if !strings.EqualFold(filepath.Ext(inputPath), ".sbg") {
		return fmt.Errorf("SBaGen input must use the .sbg extension: %q", inputPath)
	}
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".spsq"
	if len(args) == 2 {
		outputPath = args[1]
	}
	if outputPath != "-" && !strings.EqualFold(filepath.Ext(outputPath), ".spsq") {
		return fmt.Errorf("SPSQ output must use the .spsq extension: %q", outputPath)
	}

	if !quiet && statusWriter != nil {
		if err := writeSBGConversionWarning(statusWriter); err != nil {
			return err
		}
	}
	converter, err := sbg.New(synapseq.NewAppContext())
	if err != nil {
		return err
	}
	loaded, err := converter.LoadFile(inputPath)
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
	if err := os.WriteFile(outputPath, content, 0o644); err != nil {
		return fmt.Errorf("write converted sequence %q: %w", outputPath, err)
	}
	if !quiet && statusWriter != nil {
		return writeSBGConversionSuccess(statusWriter, outputPath)
	}
	return nil
}

func writeSBGConversionWarning(writer io.Writer) error {
	_, err := fmt.Fprintln(writer, cli.ErrorText("Warning:"), cli.Muted(sbgConversionWarning))
	return err
}

func writeSBGConversionSuccess(writer io.Writer, outputPath string) error {
	_, err := fmt.Fprintln(writer, cli.SuccessText("Converted:"), cli.Command(outputPath))
	return err
}
