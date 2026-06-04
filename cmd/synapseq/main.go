// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
