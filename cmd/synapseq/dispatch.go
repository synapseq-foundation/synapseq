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
	case cli.SpecialCommandDoctor:
		return true, runDoctor()
	default:
		return false, nil
	}
}
