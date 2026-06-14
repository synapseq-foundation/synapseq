// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
)

type flagValueKind int

const (
	flagValueBool flagValueKind = iota
	flagValueString
)

type flagBinding struct {
	Name           string
	Usage          string
	ValueKind      flagValueKind
	BindBool       func(*CLIOptions) *bool
	BindString     func(*CLIOptions) *string
	SpecialCommand SpecialCommandKind
}

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if hasNoColorArg(os.Args[1:]) {
		SetColorEnabled(false)
	}

	fs.Usage = func() {}
	bindFlags(fs, opts)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, formatFlagParseError(fs, err)
	}

	SetColorEnabled(!opts.NoColor)

	return opts, fs.Args(), err
}

func bindFlags(fs *flag.FlagSet, opts *CLIOptions) {
	for _, binding := range flagBindings() {
		switch binding.ValueKind {
		case flagValueBool:
			fs.BoolVar(binding.BindBool(opts), binding.Name, false, binding.Usage)
		case flagValueString:
			fs.StringVar(binding.BindString(opts), binding.Name, "", binding.Usage)
		}
	}
}

func ResolveSpecialCommand(opts *CLIOptions, args []string) SpecialCommand {
	optionalArg := firstOptionalArg(args)

	for _, binding := range flagBindings() {
		if binding.SpecialCommand == SpecialCommandNone {
			continue
		}

		switch binding.ValueKind {
		case flagValueBool:
			if *binding.BindBool(opts) {
				return SpecialCommand{Kind: binding.SpecialCommand, OptionalArg: optionalArg}
			}
		case flagValueString:
			if *binding.BindString(opts) != "" {
				return SpecialCommand{Kind: binding.SpecialCommand, OptionalArg: optionalArg}
			}
		}
	}

	return SpecialCommand{Kind: SpecialCommandNone}
}

func hasNoColorArg(args []string) bool {
	for _, arg := range args {
		if arg == "-no-color" {
			return true
		}
	}
	return false
}

func formatFlagParseError(fs *flag.FlagSet, err error) error {
	if err == nil {
		return nil
	}

	message := err.Error()
	knownFlags := flagNames(fs)

	switch {
	case strings.HasPrefix(message, "flag provided but not defined: "):
		found := strings.TrimSpace(strings.TrimPrefix(message, "flag provided but not defined: "))
		diagnostic := diag.Validation("unknown command-line flag").WithFound(found).WithHint("use -help to see valid command-line options")
		if suggestion, ok := diag.ClosestMatch(found, knownFlags, diag.DefaultSuggestionDistance(found)); ok {
			diagnostic.WithSuggestion(fmt.Sprintf("did you mean %q?", suggestion))
		}
		return diagnostic
	case strings.HasPrefix(message, "flag needs an argument: "):
		found := strings.TrimSpace(strings.TrimPrefix(message, "flag needs an argument: "))
		return diag.Validation("missing value for command-line flag").WithFound(found).WithHint("pass a value for this flag or use -help to review its syntax")
	case strings.HasPrefix(message, "invalid boolean value "):
		return diag.Validation("invalid value for command-line flag").WithHint(message + "; use -help to review accepted flag values")
	default:
		return diag.Validation("invalid command-line arguments").WithHint(message + "; use -help for usage information")
	}
}

func flagNames(fs *flag.FlagSet) []string {
	flags := make([]string, 0, 16)
	fs.VisitAll(func(f *flag.Flag) {
		flags = append(flags, "-"+f.Name)
	})
	return flags
}

func flagBindings() []flagBinding {
	return []flagBinding{
		{Name: "version", Usage: "Show version information", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.ShowVersion }, SpecialCommand: SpecialCommandShowVersion},
		{Name: "preview", Usage: "Render HTML preview timeline", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Preview }},
		{Name: "dump", Usage: "Render JSON sequence data", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Dump }},
		{Name: "quiet", Usage: "Enable quiet mode", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Quiet }},
		{Name: "no-color", Usage: "Disable ANSI colors in CLI output", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.NoColor }},
		{Name: "test", Usage: "Validate syntax without generating output", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Test }},
		{Name: "help", Usage: "Show help", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.ShowHelp }},
		{Name: "remote-sync", Usage: "Sync index of available sequences", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.RemoteSync }, SpecialCommand: SpecialCommandRemoteSync},
		{Name: "remote-clean", Usage: "Clean up local cache", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.RemoteClean }, SpecialCommand: SpecialCommandRemoteClean},
		{Name: "remote-get", Usage: "Get remote sequence", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.RemoteGet }, SpecialCommand: SpecialCommandRemoteGet},
		{Name: "remote-list", Usage: "List remote sequences", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.RemoteList }, SpecialCommand: SpecialCommandRemoteList},
		{Name: "remote-search", Usage: "Search remote sequences", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.RemoteSearch }, SpecialCommand: SpecialCommandRemoteSearch},
		{Name: "remote-download", Usage: "Download remote sequence", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.RemoteDownload }, SpecialCommand: SpecialCommandRemoteDownload},
		{Name: "remote-info", Usage: "Show information about a remote sequence", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.RemoteInfo }, SpecialCommand: SpecialCommandRemoteInfo},
		{Name: "play", Usage: "Play audio using ffplay", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Play }},
		{Name: "mp3", Usage: "Export to MP3 with ffmpeg", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.Mp3 }},
		{Name: "install-file-association", Usage: "Associate .spsq files with SynapSeq (Windows only)", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.InstallFileAssociation }, SpecialCommand: SpecialCommandInstallFileAssociation},
		{Name: "uninstall-file-association", Usage: "Remove .spsq file association (Windows only)", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.UninstallFileAssociation }, SpecialCommand: SpecialCommandUninstallFileAssociation},
		{Name: "new", Usage: "Template type: meditation, focus, sleep, relaxation, example", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.New }, SpecialCommand: SpecialCommandGenerateTemplate},
		{Name: "ffmpeg-path", Usage: "Path to ffmpeg executable", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.FFmpegPath }},
		{Name: "ffplay-path", Usage: "Path to ffplay executable", ValueKind: flagValueString, BindString: func(opts *CLIOptions) *string { return &opts.FFplayPath }},
		{Name: "doctor", Usage: "Check environment for required tools", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.ShowDoctor }, SpecialCommand: SpecialCommandDoctor},
		{Name: "completion-bash", Usage: "Print bash completion script", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.CompletionBash }},
		{Name: "completion-zsh", Usage: "Print zsh completion script", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.CompletionZsh }},
		{Name: "completion-args", Usage: "Print flags with descriptions", ValueKind: flagValueBool, BindBool: func(opts *CLIOptions) *bool { return &opts.CompletionArgs }},
	}
}

func firstOptionalArg(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}
