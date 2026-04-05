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

package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/diag"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
	style "github.com/synapseq-foundation/synapseq/v4/internal/textstyle"
)

// CLIOptions holds command-line options
type CLIOptions struct {
	// Show version information and exit
	ShowVersion bool
	// New starter sequence template type
	New string
	// Preview mode, renders HTML timeline instead of audio
	Preview bool
	// Quiet mode, suppress non-error output
	Quiet bool
	// Test mode, validate syntax without generating output
	Test bool
	// Show help message and exit
	ShowHelp bool
	// Show full manual and exit
	ShowManual bool
	// Disable ANSI colors in CLI output
	NoColor bool
	// Windows file association installation
	InstallFileAssociation bool
	// Clean Windows file association removal
	UninstallFileAssociation bool
	// Play (with ffplay)
	Play bool
	// Hub update index of available sequences
	HubUpdate bool
	// Hub clean up local cache
	HubClean bool
	// Hub list available sequences
	HubList bool
	// Hub search sequences
	HubSearch string
	// Hub download sequences
	HubDownload string
	// Hub info of sequence
	HubInfo string
	// Hub get sequence
	HubGet string
	// Mp3 export with ffmpeg
	Mp3 bool
	// Path to ffplay executable
	FFplayPath string
	// Path to ffmpeg executable
	FFmpegPath string
	// Path to ffprobe executable
	FFprobePath string
}

func init() {
	SetColorEnabled(true)
}

// SetColorEnabled enables or disables ANSI colors for CLI output.
func SetColorEnabled(enabled bool) {
	style.SetColorEnabled(enabled)
}

// Title formats top-level headings.
func Title(text string) string {
	return style.Title(text)
}

// Section formats section headings.
func Section(text string) string {
	return style.Section(text)
}

// Command formats shell commands and paths.
func Command(text string) string {
	return style.Command(text)
}

// Flag formats command-line flags.
func Flag(text string) string {
	return style.Flag(text)
}

func FlagColumn(text string, width int) string {
	return style.FlagColumn(text, width)
}

// Muted formats secondary explanatory text.
func Muted(text string) string {
	return style.Muted(text)
}

// ErrorText formats error text.
func ErrorText(text string) string {
	return style.ErrorText(text)
}

// SuccessText formats success text.
func SuccessText(text string) string {
	return style.SuccessText(text)
}

// Label formats field labels.
func Label(text string) string {
	return style.Label(text)
}

// Accent formats highlighted values.
func Accent(text string) string {
	return style.Accent(text)
}

// Help prints the help message
func Help() {
	fmt.Fprintf(color.Output, "%s\n\n", Title(fmt.Sprintf("SynapSeq %s - Text-Driven Audio Sequencer for Brainwave Entrainment", info.VERSION)))

	fmt.Fprintf(color.Output, "%s\n", Section("Usage:"))
	fmt.Fprintf(color.Output, "  %s\n\n", Command("synapseq [options] <input> [output]"))

	fmt.Fprintf(color.Output, "%s\n", Section("Quick start:"))
	fmt.Fprintf(color.Output, "  %s\n", Label("1. Create a starter file"))
	fmt.Fprintf(color.Output, "     %s\n", Command("synapseq -new meditation starter.spsq"))
	fmt.Fprintf(color.Output, "       %s\n\n", Muted("Create starter.spsq from the meditation template"))
	fmt.Fprintf(color.Output, "  %s\n", Label("2. Render audio"))
	fmt.Fprintf(color.Output, "     %s\n", Command("synapseq starter.spsq"))
	fmt.Fprintf(color.Output, "       %s\n\n", Muted("Generate starter.wav in the current folder"))
	fmt.Fprintf(color.Output, "  %s\n", Label("Available templates"))
	fmt.Fprintf(color.Output, "     %s\n\n", Muted("meditation, focus, sleep, relaxation, example"))

	fmt.Fprintf(color.Output, "%s\n", Section("Next steps:"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -test starter.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Validate syntax and semantics without generating audio"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -preview starter.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Generate starter.html with a visual timeline preview"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -play starter.spsq"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Play the sequence directly with ffplay"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq starter.spsq starter.mp3"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Export to MP3 with ffmpeg"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -manual"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Print the compact syntax reference manual"))

	fmt.Fprintf(color.Output, "%s\n", Section("Input:"))
	fmt.Fprintf(color.Output, "  local file        %s\n", Command("path/to/sequence.spsq"))
	fmt.Fprintf(color.Output, "  URL               %s\n", Command("https://example.com/sequence.spsq"))
	fmt.Fprintf(color.Output, "  standard input    %s\n\n", Command("-"))

	fmt.Fprintf(color.Output, "%s\n", Section("Output:"))
	fmt.Fprintf(color.Output, "  omitted           %s\n", Muted("defaults to <input>.wav"))
	fmt.Fprintf(color.Output, "  WAV file          %s\n", Command("path/to/output.wav"))
	fmt.Fprintf(color.Output, "  MP3 file          %s\n", Command("path/to/output.mp3"))
	fmt.Fprintf(color.Output, "  standard output   %s\n\n", Muted("-   raw PCM (16-bit stereo)"))

	fmt.Fprintf(color.Output, "%s\n", Section("Most common options:"))
	fmt.Fprintf(color.Output, "  %sTemplate type: meditation, focus, sleep, relaxation, example\n", FlagColumn("-new TYPE", 18))
	fmt.Fprintf(color.Output, "  %sCheck syntax only\n", FlagColumn("-test", 18))
	fmt.Fprintf(color.Output, "  %sRender an HTML preview timeline\n", FlagColumn("-preview", 18))
	fmt.Fprintf(color.Output, "  %sPlay audio using ffplay\n", FlagColumn("-play", 18))
	fmt.Fprintf(color.Output, "  %sExport to MP3 with ffmpeg\n", FlagColumn("-mp3", 18))
	fmt.Fprintf(color.Output, "  %sSuppress non-error output\n", FlagColumn("-quiet", 18))
	fmt.Fprintf(color.Output, "  %sDisable ANSI colors in CLI output\n", FlagColumn("-no-color", 18))
	fmt.Fprintf(color.Output, "  %sShow the compact syntax reference manual\n", FlagColumn("-manual", 18))
	fmt.Fprintf(color.Output, "  %sShow version information\n", FlagColumn("-version", 18))
	fmt.Fprintf(color.Output, "  %sShow this help message\n\n", FlagColumn("-help", 18))

	fmt.Fprintf(color.Output, "%s\n", Section("Hub:"))
	fmt.Fprintf(color.Output, "  %s\n\n", Muted("Run -hub-update first to initialize the local Hub index."))
	fmt.Fprintf(color.Output, "  %s Update the local Hub index\n", FlagColumn("-hub-update", 24))
	fmt.Fprintf(color.Output, "  %s List available sequences\n", FlagColumn("-hub-list", 24))
	fmt.Fprintf(color.Output, "  %s Search the Hub\n", FlagColumn("-hub-search WORD", 24))
	fmt.Fprintf(color.Output, "  %s Show information about a sequence\n", FlagColumn("-hub-info NAME", 24))
	fmt.Fprintf(color.Output, "  %s Download a sequence and dependencies\n", FlagColumn("-hub-download NAME [DIR]", 24))
	fmt.Fprintf(color.Output, "  %s Download and generate in one step\n", FlagColumn("-hub-get NAME [OUTPUT]", 24))
	fmt.Fprintf(color.Output, "  %s Clean up local Hub cache\n\n", FlagColumn("-hub-clean", 24))

	fmt.Fprintf(color.Output, "%s\n", Section("Hub quick start:"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-update"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-list"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-search calm-state"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -hub-get calm-state calm-state.wav"))
	fmt.Fprintf(color.Output, "  %s\n\n", Command("synapseq -hub-get calm-state calm-state.mp3"))

	fmt.Fprintf(color.Output, "%s\n", Section("Advanced:"))
	fmt.Fprintf(color.Output, "  %sPath to ffmpeg executable\n", FlagColumn("-ffmpeg-path PATH", 22))
	fmt.Fprintf(color.Output, "  %sPath to ffplay executable\n\n", FlagColumn("-ffplay-path PATH", 22))

	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, "%s\n", Section("Windows-specific options:"))
		fmt.Fprintf(color.Output, "  %sAssociate .spsq files with SynapSeq\n", FlagColumn("-install-file-association", 30))
		fmt.Fprintf(color.Output, "  %sRemove .spsq file association\n\n", FlagColumn("-uninstall-file-association", 30))
	}

	fmt.Fprintf(color.Output, "%s\n", Section("For more information:"))
	fmt.Fprintf(color.Output, "  %s\n", Command("synapseq -manual"))
	fmt.Fprintf(color.Output, "    %s\n\n", Muted("Show the compact syntax reference manual"))
	fmt.Fprintf(color.Output, "  %s\n", Command("https://synapseq.org"))
	fmt.Fprintf(color.Output, "    %s\n", Muted("Visit the website for documentation, examples, and the latest updates"))
}

// ShowVersion prints the version information
func ShowVersion() {
	fmt.Fprintf(
		color.Output,
		"%s %s %s %s %s\n",
		Title("SynapSeq"),
		Accent(info.VERSION),
		Muted(fmt.Sprintf("(%s)", info.GIT_COMMIT)),
		Label("built"),
		Command(fmt.Sprintf("%s for %s/%s", info.BUILD_DATE, runtime.GOOS, runtime.GOARCH)),
	)
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

	// General options
	fs.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	fs.StringVar(&opts.New, "new", "", "Template type: meditation, focus, sleep, relaxation, example")
	fs.BoolVar(&opts.Preview, "preview", false, "Render HTML preview timeline")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.NoColor, "no-color", false, "Disable ANSI colors in CLI output")
	fs.BoolVar(&opts.Test, "test", false, "Validate syntax without generating output")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")
	fs.BoolVar(&opts.ShowManual, "manual", false, "Show the full manual")

	// Hub options
	fs.BoolVar(&opts.HubUpdate, "hub-update", false, "Update index of available sequences")
	fs.BoolVar(&opts.HubClean, "hub-clean", false, "Clean up local cache")
	fs.BoolVar(&opts.HubList, "hub-list", false, "List available sequences")
	fs.StringVar(&opts.HubSearch, "hub-search", "", "Search sequences")
	fs.StringVar(&opts.HubDownload, "hub-download", "", "Download sequence and dependencies")
	fs.StringVar(&opts.HubInfo, "hub-info", "", "Show information about a sequence")
	fs.StringVar(&opts.HubGet, "hub-get", "", "Get sequence")

	// External tool options
	fs.BoolVar(&opts.Play, "play", false, "Play audio using ffplay")
	fs.BoolVar(&opts.Mp3, "mp3", false, "Export to MP3 with ffmpeg")
	fs.StringVar(&opts.FFmpegPath, "ffmpeg-path", "", "Path to ffmpeg executable")
	fs.StringVar(&opts.FFplayPath, "ffplay-path", "", "Path to ffplay executable")

	// Windows-specific options
	fs.BoolVar(&opts.InstallFileAssociation, "install-file-association", false, "Associate .spsq files with SynapSeq (Windows only)")
	fs.BoolVar(&opts.UninstallFileAssociation, "uninstall-file-association", false, "Remove .spsq file association (Windows only)")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, formatFlagParseError(fs, err)
	}

	SetColorEnabled(!opts.NoColor)

	return opts, fs.Args(), err
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
