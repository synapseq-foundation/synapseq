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
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
)

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
