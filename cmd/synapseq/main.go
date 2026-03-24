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
	"os"
	"path/filepath"
	"strings"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/manual"
)

// main is the entry point of the SynapSeq application
func main() {
	opts, args, err := cli.ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, formatCLIError(err))
		os.Exit(1)
	}

	if err := run(opts, args); err != nil {
		fmt.Fprintln(os.Stderr, formatCLIError(err))
		os.Exit(1)
	}
}

// run executes the main application logic based on CLI options and arguments
func run(opts *cli.CLIOptions, args []string) error {
	// --version
	if opts.ShowVersion {
		cli.ShowVersion()
		return nil
	}

	if opts.ShowManual {
		manual.Show()
		return nil
	}

	// --hub-update
	if opts.HubUpdate {
		return hubRunUpdate(opts.Quiet)
	}

	// --hub-clean
	if opts.HubClean {
		return hubRunClean(opts.Quiet)
	}

	// --hub-get
	if opts.HubGet != "" {
		var outputFile string
		if len(args) == 1 {
			outputFile = args[0]
		}
		return hubRunGet(opts.HubGet, outputFile, opts)
	}

	// --hub-list
	if opts.HubList {
		return hubRunList()
	}

	// --hub-search
	if opts.HubSearch != "" {
		return hubRunSearch(opts.HubSearch)
	}

	// --hub-download
	if opts.HubDownload != "" {
		targetDir := ""
		if len(args) == 1 {
			targetDir = args[0]
		}
		return hubRunDownload(opts.HubDownload, targetDir, opts.Quiet)
	}

	// --hub-info
	if opts.HubInfo != "" {
		return hubRunInfo(opts.HubInfo)
	}

	// --install-file-association (Windows only)
	if opts.InstallFileAssociation {
		return installWindowsFileAssociation(opts.Quiet)
	}

	// --uninstall-file-association (Windows only)
	if opts.UninstallFileAssociation {
		return uninstallWindowsFileAssociation(opts.Quiet)
	}

	// --new template generation
	if opts.New != "" {
		outputFile := opts.New + ".spsq"
		if len(args) == 1 {
			outputFile = args[0]
		}
		return generateTemplate(opts.New, outputFile)
	}

	// --help or missing args
	if opts.ShowHelp || len(args) == 0 {
		cli.Help()
		return nil
	}

	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid number of flags\nUse -help for usage information")
	}

	// Default: process input file and generate output
	outputFormat := "wav"
	if opts.Preview {
		outputFormat = "html"
	}
	if opts.Mp3 {
		outputFormat = "mp3"
	}

	inputFile := args[0]
	outputFile := getDefaultOutputFile(inputFile, outputFormat)
	if len(args) == 2 {
		outputFile = args[1]
		outputFormat = strings.ToLower(filepath.Ext(outputFile))
	}

	appCtx := synapseq.NewAppContext()

	if !opts.Quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stderr, !opts.NoColor)
	}

	// Load sequence file
	loadedCtx, err := appCtx.Load(inputFile)
	if err != nil {
		return err
	}

	// --- Handle Test mode (no output required)
	if opts.Test {
		if !opts.Quiet {
			fmt.Println(cli.SuccessText("Sequence is valid."))
		}
		return nil
	}

	// --- Process output using centralized handler
	outputOpts := &outputOptions{
		OutputFile: outputFile,
		Quiet:      opts.Quiet,
		Preview:    opts.Preview,
		Play:       opts.Play,
		Mp3:        outputFormat == ".mp3" || opts.Mp3,
		FFplayPath: opts.FFplayPath,
		FFmpegPath: opts.FFmpegPath,
	}

	return processSequenceOutput(loadedCtx, outputOpts)
}

// getDefaultOutputFile generates a default output filename based on the input filename
func getDefaultOutputFile(inputFile string, extension string) string {
	base := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	return base + "." + extension
}
