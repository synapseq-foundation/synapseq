// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
	case cli.SpecialCommandSync:
		return true, remoteRunSync(opts.Quiet)
	case cli.SpecialCommandClean:
		return true, remoteRunClean(opts.Quiet)
	case cli.SpecialCommandGet:
		return true, remoteRunGet(opts.RemoteGet, command.OptionalArg, opts)
	case cli.SpecialCommandList:
		return true, remoteRunList()
	case cli.SpecialCommandSearch:
		return true, remoteRunSearch(opts.RemoteSearch)
	case cli.SpecialCommandDownload:
		return true, remoteRunDownload(opts.RemoteDownload, command.OptionalArg, opts.Quiet)
	case cli.SpecialCommandInfo:
		return true, remoteRunInfo(opts.RemoteInfo)
	case cli.SpecialCommandInstallFileAssociation:
		return true, installWindowsFileAssociation(opts.Quiet)
	case cli.SpecialCommandUninstallFileAssociation:
		return true, uninstallWindowsFileAssociation(opts.Quiet)
	case cli.SpecialCommandDoctor:
		return true, runDoctor()
	default:
		return false, nil
	}
}
