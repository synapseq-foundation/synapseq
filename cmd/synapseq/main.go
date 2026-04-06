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

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
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
	handled, err := dispatchSpecialCommand(opts, args)
	if err != nil {
		return err
	}
	if handled {
		return nil
	}

	// --help or missing args
	if opts.ShowHelp || len(args) == 0 {
		cli.Help()
		return nil
	}

	return handleSequenceCommand(args, opts)
}
