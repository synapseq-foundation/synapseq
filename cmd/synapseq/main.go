//go:build !js && !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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

	"github.com/fatih/color"
	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	"github.com/synapseq-foundation/synapseq/v4/internal/manual"
)

// main is the entry point of the SynapSeq application
func main() {
	opts, args, err := cli.ParseFlags()
	if err != nil {
		os.Exit(1)
	}

	if err := run(opts, args); err != nil {
		fmt.Fprintln(os.Stderr, formatCLIError(err))
		os.Exit(1)
	}
}

func formatCLIError(err error) string {
	if diagnostic, ok := diag.As(err); ok {
		var lines []string
		location := formatDiagnosticLocation(diagnostic)
		if location != "" {
			lines = append(lines, cli.ErrorText("synapseq:")+" "+cli.Accent(location)+": "+cli.ErrorText(diagnostic.Message))
		} else {
			lines = append(lines, cli.ErrorText("synapseq:")+" "+cli.ErrorText(diagnostic.Message))
		}

		if diagnostic.Span.LineText != "" {
			lines = append(lines, color.New(color.FgWhite).Sprint(diagnostic.Span.LineText))
			lines = append(lines, cli.ErrorText(strings.Repeat(" ", diagnostic.Span.Column-1)+strings.Repeat("^", max(1, diagnostic.Span.EndColumn-diagnostic.Span.Column))))
		}
		if diagnostic.Found != "" {
			lines = append(lines, cli.Label("found:")+" "+cli.Accent(fmt.Sprintf("%q", diagnostic.Found)))
		}
		if len(diagnostic.Expected) > 0 {
			lines = append(lines, cli.Label("expected:")+" "+formatExpectedList(diagnostic.Expected))
		}
		if diagnostic.Suggestion != "" {
			lines = append(lines, cli.SuccessText(diagnostic.Suggestion))
		}
		if diagnostic.Hint != "" {
			lines = append(lines, cli.Muted(diagnostic.Hint))
		}

		return strings.Join(lines, "\n")
	}
	return cli.ErrorText("synapseq:") + " " + err.Error()
}

func formatDiagnosticLocation(diagnostic *diag.Diagnostic) string {
	if diagnostic == nil {
		return ""
	}
	span := diagnostic.Span
	switch {
	case span.File != "" && span.Line > 0 && span.Column > 0:
		return fmt.Sprintf("%s:%d:%d", span.File, span.Line, span.Column)
	case span.Line > 0 && span.Column > 0:
		return fmt.Sprintf("line %d:%d", span.Line, span.Column)
	case span.Line > 0:
		return fmt.Sprintf("line %d", span.Line)
	case span.File != "":
		return span.File
	default:
		return ""
	}
}

func formatExpectedList(values []string) string {
	colored := make([]string, len(values))
	for i, value := range values {
		colored[i] = cli.Accent(fmt.Sprintf("%q", value))
	}
	if len(colored) == 1 {
		return colored[0]
	}
	if len(colored) == 2 {
		return colored[0] + " or " + colored[1]
	}
	return strings.Join(colored[:len(colored)-1], ", ") + ", or " + colored[len(colored)-1]
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
		Mp3:        outputFormat == ".mp3",
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
