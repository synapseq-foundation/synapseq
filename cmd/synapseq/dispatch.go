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
	case cli.SpecialCommandRemoteSync:
		return true, remoteRunSync(opts.Quiet)
	case cli.SpecialCommandRemoteClean:
		return true, remoteRunClean(opts.Quiet)
	case cli.SpecialCommandRemoteGet:
		return true, remoteRunGet(opts.RemoteGet, command.OptionalArg, opts)
	case cli.SpecialCommandRemoteList:
		return true, remoteRunList()
	case cli.SpecialCommandRemoteSearch:
		return true, remoteRunSearch(opts.RemoteSearch)
	case cli.SpecialCommandRemoteDownload:
		return true, remoteRunDownload(opts.RemoteDownload, command.OptionalArg, opts.Quiet)
	case cli.SpecialCommandRemoteInfo:
		return true, remoteRunInfo(opts.RemoteInfo)
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
