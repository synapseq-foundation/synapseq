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
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func dispatchSpecialCommand(opts *cli.CLIOptions, args []string) (bool, error) {
	if opts.CompletionBash {
		cli.PrintBashCompletion()
		return true, nil
	}
	if opts.CompletionZsh {
		cli.PrintZshCompletion()
		return true, nil
	}
	if opts.CompletionArgs {
		cli.PrintCompletionArgs()
		return true, nil
	}

	command := cli.ResolveSpecialCommand(opts, args)

	switch command.Kind {
	case cli.SpecialCommandShowVersion:
		cli.ShowVersion()
		return true, nil
	case cli.SpecialCommandShowManual:
		cli.ShowManual()
		return true, nil
	case cli.SpecialCommandHubUpdate:
		return true, hubRunUpdate(opts.Quiet)
	case cli.SpecialCommandHubClean:
		return true, hubRunClean(opts.Quiet)
	case cli.SpecialCommandHubGet:
		return true, hubRunGet(opts.HubGet, command.OptionalArg, opts)
	case cli.SpecialCommandHubList:
		return true, hubRunList()
	case cli.SpecialCommandHubSearch:
		return true, hubRunSearch(opts.HubSearch)
	case cli.SpecialCommandHubDownload:
		return true, hubRunDownload(opts.HubDownload, command.OptionalArg, opts.Quiet)
	case cli.SpecialCommandHubInfo:
		return true, hubRunInfo(opts.HubInfo)
	case cli.SpecialCommandInstallFileAssociation:
		return true, installWindowsFileAssociation(opts.Quiet)
	case cli.SpecialCommandUninstallFileAssociation:
		return true, uninstallWindowsFileAssociation(opts.Quiet)
	case cli.SpecialCommandGenerateTemplate:
		outputFile := opts.New + ".spsq"
		if command.OptionalArg != "" {
			outputFile = command.OptionalArg
		}
		return true, generateTemplate(opts.New, outputFile)
	case cli.SpecialCommandDoctor:
		return true, runDoctor()
	default:
		return false, nil
	}
}
