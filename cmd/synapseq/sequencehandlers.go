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
)

type sequenceCommand struct {
	inputFile    string
	outputFile   string
	outputFormat string
}

func handleSequenceCommand(args []string, opts *cli.CLIOptions) error {
	command, err := prepareSequenceCommand(args, opts)
	if err != nil {
		return err
	}

	return runSequenceInput(command.inputFile, command.outputFile, command.outputFormat, os.Stderr, opts)
}

func prepareSequenceCommand(args []string, opts *cli.CLIOptions) (*sequenceCommand, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("invalid number of flags\nUse -help for usage information")
	}

	inputFile := args[0]
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	outputFile, outputFormat := resolveOutputTarget(baseName, "", opts)
	if len(args) == 2 {
		outputFile, outputFormat = resolveOutputTarget(baseName, args[1], opts)
	}

	return &sequenceCommand{
		inputFile:    inputFile,
		outputFile:   outputFile,
		outputFormat: outputFormat,
	}, nil
}

func loadSequenceContext(inputFile, outputFile string, verboseWriter io.Writer, opts *cli.CLIOptions) (*synapseq.LoadedContext, error) {
	appCtx := synapseq.NewAppContext()
	if !opts.Quiet && outputFile != "-" && verboseWriter != nil {
		appCtx = appCtx.WithVerbose(verboseWriter, !opts.NoColor)
	}

	return appCtx.LoadFile(inputFile)
}

func runLoadedSequence(loadedCtx *synapseq.LoadedContext, outputFile, outputFormat string, opts *cli.CLIOptions) error {
	if opts.Test {
		if !opts.Quiet {
			fmt.Println(cli.SuccessText("Sequence is valid."))
		}
		return nil
	}

	return processSequenceOutput(loadedCtx, buildOutputOptions(outputFile, outputFormat, opts))
}

func runSequenceInput(inputFile, outputFile, outputFormat string, verboseWriter io.Writer, opts *cli.CLIOptions) error {
	loadedCtx, err := loadSequenceContext(inputFile, outputFile, verboseWriter, opts)
	if err != nil {
		return err
	}

	return runLoadedSequence(loadedCtx, outputFile, outputFormat, opts)
}
