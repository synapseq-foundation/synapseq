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

package manual

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/synapseq-foundation/synapseq/v4/internal/info"
)

const manualWidth = 78

// Show prints the SynapSeq manual to the standard color-aware output.
func Show() {
	Print(color.Output)
}

// Print writes the SynapSeq manual to the provided writer.
func Print(w io.Writer) {
	fmt.Fprint(w, Render())
}

// Render returns the SynapSeq manual as a formatted string.
func Render() string {
	var b strings.Builder

	writeTitle(&b, "SYNAPSEQ(1)", fmt.Sprintf("SynapSeq Manual %s", info.VERSION))

	writeSection(&b, "Name")
	writeSubsection(&b, "synapseq")
	writeParagraph(&b,
		"text-driven audio sequencer for brainwave sessions and ambient sound design",
	)

	writeSection(&b, "Synopsis")
	writeCodeBlock(&b,
		"synapseq [OPTION]... INPUT [OUTPUT]",
		"synapseq -new TYPE",
		"synapseq -preview INPUT",
		"synapseq -play INPUT",
		"synapseq -mp3 INPUT",
		"synapseq -test INPUT",
		"synapseq -hub-update",
		"synapseq -hub-clean",
		"synapseq -hub-list",
		"synapseq -hub-search WORD",
		"synapseq -hub-info NAME",
		"synapseq -hub-download NAME [DIR]",
		"synapseq -hub-get NAME [OUTPUT]",
		"synapseq -manual",
		"synapseq -help",
		"synapseq -version",
		"synapseq -manual -no-color",
	)
	writeSubsection(&b, "Paged reading")
	writeBullet(&b, "Linux and macOS", "synapseq -manual | less")
	writeBullet(&b, "Windows PowerShell", "synapseq -manual | more")

	writeSection(&b, "Description")
	writeParagraph(&b,
		"SynapSeq parses a line-oriented score language and renders deterministic audio. Ordering is strict. Options must appear first. Presets must appear next. Timeline entries must appear last.",
	)
	writeParagraph(&b,
		"INPUT may be a local .spsq file, a sequence URL, or -. Output is WAV by default. -mp3 selects MP3. -preview selects HTML. -play selects direct playback. -test validates only.",
	)

	writeSection(&b, "Options")
	writeSubsection(&b, "Rules")
	writeLineBlock(&b,
		"one primary action per command",
		"INPUT required except when not applicable",
		"OUTPUT optional",
		"derived output path when OUTPUT is omitted",
	)

	writeSubsection(&b, "General")
	writeBullet(&b, "-new TYPE", "create a starter sequence")
	writeBullet(&b, "TYPE", "meditation\nfocus\nsleep\nrelaxation\nexample")
	writeBullet(&b, "-test", "validate syntax and semantics only")
	writeBullet(&b, "-preview", "render HTML preview")
	writeBullet(&b, "-play", "render and play through ffplay")
	writeBullet(&b, "-mp3", "render MP3 output")
	writeBullet(&b, "-quiet", "suppress non-error CLI output")
	writeBullet(&b, "-no-color", "disable ANSI color output")
	writeBullet(&b, "-manual", "print full manual")
	writeBullet(&b, "-help", "print concise command overview")
	writeBullet(&b, "-version", "print version build and platform information")

	writeSubsection(&b, "Hub")
	writeBullet(&b, "-hub-update", "refresh local Hub index")
	writeBullet(&b, "-hub-clean", "remove cached Hub data")
	writeBullet(&b, "-hub-list", "list sequences from local Hub index")
	writeBullet(&b, "-hub-search WORD", "search local Hub index")
	writeBullet(&b, "-hub-info NAME", "print Hub sequence metadata")
	writeBullet(&b, "-hub-download NAME [DIR]", "download sequence and dependencies\nDIR optional")
	writeBullet(&b, "-hub-get NAME [OUTPUT]", "download and render sequence\nOUTPUT optional")

	writeSubsection(&b, "External tools")
	writeBullet(&b, "-ffmpeg-path PATH", "use specific ffmpeg executable")
	writeBullet(&b, "-ffplay-path PATH", "use specific ffplay executable")

	writeSubsection(&b, "Windows")
	writeBullet(&b, "-install-file-association", "associate .spsq files with SynapSeq")
	writeBullet(&b, "-uninstall-file-association", "remove .spsq file association")

	writeSection(&b, "Sequence File")
	writeSubsection(&b, "File type")
	writeLineBlock(&b,
		"plain-text .spsq",
		"parsed top to bottom",
	)

	writeSubsection(&b, "Order")
	writeLineBlock(&b,
		"options",
		"presets",
		"timeline",
	)

	writeSubsection(&b, "Top level")
	writeLineBlock(&b,
		"no indentation",
		"preset declarations must start in column 1",
		"timeline entries must start in column 1",
	)

	writeSubsection(&b, "Preset body")
	writeLineBlock(&b,
		"exactly two leading spaces required",
		"track definitions use this indentation",
		"track overrides use this indentation",
	)

	writeSubsection(&b, "Options")
	writeLineBlock(&b,
		"valid only before first preset",
		"valid only before first timeline entry",
	)

	writeSubsection(&b, "Timeline")
	writeNestedBullet(&b, "Syntax", "HH:MM:SS PRESET_NAME [TRANSITION]")
	writeNestedBullet(&b, "Time", "HH:MM:SS required")
	writeNestedBullet(&b, "Rules", "top-level lines only\nfirst entry zero time\nstrictly increasing\nat least two periods\nnon-template preset only")
	writeNestedSubsection(&b, "TRANSITION")
	writeNestedLineBlock(&b,
		"steady ease-in ease-out smooth",
	)

	writeSubsection(&b, "Comments")
	writeNestedBullet(&b, "# comment", "ignored")
	writeNestedBullet(&b, "## comment", "stored and printed unless -quiet is set")
	writeNestedBullet(&b, "Rules", "full line only\ninline comment not permitted")

	writeSubsection(&b, "Global options")
	writeNestedBullet(&b, "@samplerate NUMBER", "default 44100")
	writeNestedBullet(&b, "@volume NUMBER", "range 0 to 100\ndefault 100")
	writeNestedBullet(&b, "@ambiance NAME PATH_OR_URL", "register named ambiance source")
	writeNestedBullet(&b, "@extends PATH_OR_URL", "load modular .spsc content before main file")

	writeSubsection(&b, "Local paths")
	writeLineBlock(&b,
		"relative only",
		"forward slashes only",
		"no extension",
		"ambiance resolves to .wav",
		"extends resolves to .spsc",
		"absolute path not permitted",
		"parent traversal not permitted",
		"Windows drive prefix not permitted",
		"stdin not permitted",
	)

	writeSubsection(&b, "Identifiers")
	writeLineBlock(&b,
		"first character must be a letter",
		"remaining characters may be letters digits underscores dashes",
		"maximum length 20",
		"preset references are case-insensitive",
		"preset names normalize to lowercase",
	)

	writeSubsection(&b, "Preset declarations")
	writeCodeBlock(&b,
		"NAME",
		"NAME as template",
		"NAME from TEMPLATE_NAME",
	)

	writeSubsection(&b, "Template rules")
	writeLineBlock(&b,
		"template preset cannot appear on timeline",
	)

	writeSubsection(&b, "Inheritance rules")
	writeLineBlock(&b,
		"inherited preset cannot define new track lines",
		"inherited preset may contain track overrides only",
	)

	writeSubsection(&b, "Track syntax")
	writeCodeBlock(&b,
		"tone CARRIER amplitude LEVEL",
		"tone CARRIER binaural|monaural|isochronic BEAT amplitude LEVEL",
		"noise white|pink|brown [smooth VALUE] amplitude LEVEL",
		"ambiance NAME amplitude LEVEL",
	)

	writeSubsection(&b, "Track rules")
	writeLineBlock(&b,
		"amplitude required",
		"carrier positive only",
		"beat positive only",
		"smooth 0 to 100 only",
		"waveform prefix allowed on tone track",
		"waveform prefix allowed on ambiance track",
		"effect TYPE VALUE intensity PERCENT appears before amplitude",
	)

	writeSubsection(&b, "Waveforms")
	writeCodeBlock(&b,
		"waveform VALUE",
	)
	writeNestedSubsection(&b, "VALUE")
	writeNestedLineBlock(&b,
		"sine square triangle sawtooth",
	)

	writeSubsection(&b, "Effects")
	writeCodeBlock(&b,
		"effect TYPE VALUE intensity PERCENT",
	)

	writeNestedSubsection(&b, "TYPE")
	writeNestedLineBlock(&b,
		"pan modulation doppler",
	)

	writeNestedSubsection(&b, "Allowed")
	writeDeepBullet(&b, "tone", "pan modulation doppler")
	writeDeepBullet(&b, "noise", "pan modulation")
	writeDeepBullet(&b, "ambiance", "pan modulation")

	writeSubsection(&b, "Track overrides")
	writeNestedSubsection(&b, "Scope")
	writeNestedLineBlock(&b,
		"inherited presets only",
	)

	writeNestedSubsection(&b, "Index")
	writeNestedLineBlock(&b,
		"1-based",
	)

	writeNestedSubsection(&b, "Syntax")
	writeNestedCodeBlock(&b,
		"track N tone VALUE",
		"track N binaural|monaural|isochronic VALUE",
		"track N smooth VALUE",
		"track N pan|modulation|doppler VALUE",
		"track N intensity VALUE",
		"track N amplitude VALUE",
	)

	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"VALUE may be absolute or signed relative delta",
		"+VALUE adds to inherited value",
		"-VALUE subtracts from inherited value",
		"matching existing effect type only",
	)

	writeSubsection(&b, "Extended files")
	writeNestedSubsection(&b, "type")
	writeNestedLineBlock(&b,
		".spsc",
	)

	writeNestedSubsection(&b, "Allowed")
	writeNestedLineBlock(&b,
		"options",
		"presets",
		"tracks",
		"track overrides",
	)

	writeNestedSubsection(&b, "Not permitted")
	writeNestedLineBlock(&b,
		"timeline entries",
		"nested @extends",
	)

	writeSection(&b, "Compatibility")
	writeSubsection(&b, "Scope")
	writeLineBlock(&b,
		"checked between consecutive timeline entries",
		"checked per channel",
		"channels match by track declaration order",
	)

	writeSubsection(&b, "Required matches")
	writeBullet(&b, "track kind", "tone only with tone\nnoise only with noise\nambiance only with ambiance")
	writeBullet(&b, "beat mode", "binaural only with binaural\nmonaural only with monaural\nisochronic only with isochronic")
	writeBullet(&b, "noise color", "white only with white\npink only with pink\nbrown only with brown")
	writeBullet(&b, "effect type", "pan only with pan\nmodulation only with modulation\ndoppler only with doppler")
	writeBullet(&b, "ambiance source", "same source only")

	writeSubsection(&b, "Not permitted")
	writeLineBlock(&b,
		"active track to off",
		"off to active track",
		"direct switch between incompatible track kinds",
		"direct switch between incompatible effect types",
		"direct switch between different ambiance sources",
	)

	writeSubsection(&b, "Allowed")
	writeLineBlock(&b,
		"waveform change between otherwise compatible tone tracks",
		"waveform change between otherwise compatible ambiance tracks",
	)

	writeSubsection(&b, "Bridge")
	writeLineBlock(&b,
		"use built-in silence preset between incompatible states",
	)

	writeSection(&b, "Common Errors")
	writeBullet(&b, "Option after preset or timeline", "move all @options to file start")
	writeBullet(&b, "Indented top-level line", "remove indentation from preset declarations and timeline entries")
	writeBullet(&b, "Wrong preset indentation", "use exactly two leading spaces for preset body lines")
	writeBullet(&b, "Inline comment", "move comment to its own line")
	writeBullet(&b, "Template on timeline", "replace template with concrete preset")
	writeBullet(&b, "First timeline entry not zero time", "set first entry to zero time")
	writeBullet(&b, "Non-increasing timeline", "reorder entries into strictly increasing time order")
	writeBullet(&b, "Single period only", "add a second timeline entry")
	writeBullet(&b, "New track in inherited preset", "remove track definition\nuse track override only")
	writeBullet(&b, "Incompatible direct transition", "insert silence between presets")
	writeBullet(&b, "Invalid local path", "use relative path\nuse forward slashes\nomit extension")

	writeSection(&b, "Notes")
	writeBullet(&b, "silence", "built-in preset\nalways available")
	writeBullet(&b, "Channels", "assigned by track declaration order")
	writeBullet(&b, "Reuse", "@extends loads reusable .spsc content")

	writeSection(&b, "See Also")
	writeCodeBlock(&b,
		"synapseq -help",
	)

	return b.String()
}
