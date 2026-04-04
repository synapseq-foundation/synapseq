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
		"synapseq -new TYPE [OUTPUT]",
		"synapseq -preview INPUT [OUTPUT]",
		"synapseq -play INPUT",
		"synapseq -mp3 INPUT [OUTPUT]",
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

	writeSection(&b, "Description")
	writeParagraph(&b,
		"SynapSeq parses line-oriented .spsq files and renders deterministic audio.",
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

	writeSubsection(&b, "Render modes")
	writeBullet(&b, "-test", "validate syntax and semantics only")
	writeBullet(&b, "-preview", "render HTML preview")
	writeBullet(&b, "-play", "render and play through ffplay")
	writeBullet(&b, "-mp3", "render MP3 output")
	writeBullet(&b, "-quiet", "suppress non-error CLI output")

	writeSubsection(&b, "Creation")
	writeBullet(&b, "-new TYPE", "create a starter sequence")
	writeBullet(&b, "TYPE", "meditation\nfocus\nsleep\nrelaxation\nexample")

	writeSubsection(&b, "Hub")
	writeBullet(&b, "-hub-update", "refresh local Hub index")
	writeBullet(&b, "-hub-clean", "remove cached Hub data")
	writeBullet(&b, "-hub-list", "list sequences from local Hub index")
	writeBullet(&b, "-hub-search WORD", "search local Hub index")
	writeBullet(&b, "-hub-info NAME", "print Hub sequence metadata")
	writeBullet(&b, "-hub-download NAME [DIR]", "download sequence and dependencies")
	writeBullet(&b, "-hub-get NAME [OUTPUT]", "download and render sequence")

	writeSubsection(&b, "Tools and system")
	writeBullet(&b, "-ffmpeg-path PATH", "use specific ffmpeg executable")
	writeBullet(&b, "-ffplay-path PATH", "use specific ffplay executable")
	writeBullet(&b, "-install-file-association", "associate .spsq files with SynapSeq on Windows")
	writeBullet(&b, "-uninstall-file-association", "remove .spsq file association on Windows")

	writeSubsection(&b, "Information")
	writeBullet(&b, "-manual", "print compact syntax reference")
	writeBullet(&b, "-help", "print concise command overview")
	writeBullet(&b, "-version", "print version build and platform information")
	writeBullet(&b, "-no-color", "disable ANSI color output")

	writeSection(&b, "Sequence File")
	writeSubsection(&b, "Order")
	writeLineBlock(&b,
		"plain-text .spsq",
		"sequence options",
		"presets",
		"timeline",
	)

	writeSubsection(&b, "Top level")
	writeLineBlock(&b,
		"all top-level lines start in column 1",
		"preset body lines require exactly two leading spaces",
	)

	writeSubsection(&b, "Comments")
	writeLineBlock(&b,
		"# comment: ignored",
		"## comment: printed unless -quiet is set",
		"full line only",
		"inline comment not permitted",
	)

	writeSubsection(&b, "Sequence Options")
	writeNestedSubsection(&b, "Syntax")
	writeNestedCodeBlock(&b,
		"@samplerate NUMBER",
		"@volume NUMBER",
		"@ambiance NAME PATH_OR_URL",
		"@extends PATH_OR_URL",
	)
	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"valid only before first preset",
		"valid only before first timeline entry",
		"@samplerate default 44100",
		"@volume range 0 to 100; default 100",
	)
	writeNestedSubsection(&b, "Local paths")
	writeNestedLineBlock(&b,
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
	writeNestedSubsection(&b, "Extended files")
	writeNestedLineBlock(&b,
		"type .spsc",
	)
	writeIndentedSubsection(&b, 12, "Allowed")
	writeIndentedLineBlock(&b, 16,
		"options",
		"presets",
		"tracks",
		"track overrides",
	)
	writeIndentedSubsection(&b, 12, "Not permitted")
	writeIndentedLineBlock(&b, 16,
		"timeline entries",
		"nested @extends",
	)

	writeSubsection(&b, "Presets")
	writeNestedSubsection(&b, "Identifiers")
	writeNestedLineBlock(&b,
		"first character must be a letter",
		"maximum length 20",
		"preset references are case-insensitive",
		"preset names normalize to lowercase",
	)
	writeIndentedSubsection(&b, 12, "Remaining characters may be:")
	writeIndentedLineBlock(&b, 16,
		"letters",
		"digits",
		"underscores",
		"dashes",
	)
	writeNestedSubsection(&b, "Preset declarations")
	writeNestedCodeBlock(&b,
		"NAME",
		"NAME as template",
		"NAME from TEMPLATE_NAME",
	)
	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"template preset cannot appear on timeline",
		"inherited preset cannot define new track lines",
		"inherited preset may contain track overrides only",
	)
	writeNestedSubsection(&b, "Track syntax")
	writeNestedCodeBlock(&b,
		"tone CARRIER amplitude LEVEL",
		"tone CARRIER binaural|monaural|isochronic BEAT amplitude LEVEL",
		"noise white|pink|brown [smooth VALUE] amplitude LEVEL",
		"ambiance NAME amplitude LEVEL",
	)
	writeNestedSubsection(&b, "Track rules")
	writeNestedLineBlock(&b,
		"amplitude required",
		"carrier positive only",
		"beat positive only",
		"smooth 0 to 100 only",
		"waveform prefix allowed on:",
		"tone tracks",
		"ambiance tracks",
		"effect appears before amplitude",
	)
	writeNestedSubsection(&b, "Waveforms")
	writeNestedCodeBlock(&b, "waveform VALUE")
	writeIndentedSubsection(&b, 12, "VALUE")
	writeIndentedLineBlock(&b, 16,
		"sine",
		"square",
		"triangle",
		"sawtooth",
	)
	writeNestedSubsection(&b, "Effects")
	writeNestedCodeBlock(&b,
		"effect pan VALUE intensity PERCENT",
		"effect modulation VALUE intensity PERCENT",
		"effect doppler VALUE intensity PERCENT",
	)
	writeNestedSubsection(&b, "Allowed")
	writeIndentedSubsection(&b, 12, "tone")
	writeIndentedLineBlock(&b, 16,
		"pan",
		"modulation",
		"doppler",
	)
	writeIndentedSubsection(&b, 12, "noise")
	writeIndentedLineBlock(&b, 16,
		"pan",
		"modulation",
	)
	writeIndentedSubsection(&b, 12, "ambiance")
	writeIndentedLineBlock(&b, 16,
		"pan",
		"modulation",
	)
	writeNestedSubsection(&b, "Track overrides")
	writeNestedCodeBlock(&b,
		"track N tone VALUE",
		"track N binaural|monaural|isochronic VALUE",
		"track N smooth VALUE",
		"track N pan|modulation|doppler VALUE",
		"track N intensity VALUE",
		"track N amplitude VALUE",
	)
	writeNestedLineBlock(&b,
		"inherited presets only",
		"1-based index",
		"VALUE may be absolute or signed relative delta",
		"+VALUE adds to inherited value",
		"-VALUE subtracts from inherited value",
		"matching existing effect type only",
	)

	writeSubsection(&b, "Timeline")
	writeNestedSubsection(&b, "Syntax")
	writeNestedCodeBlock(&b,
		"HH:MM:SS PRESET_NAME [TRANSITION [STEPS]]",
	)
	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"top-level lines only; HH:MM:SS required",
		"first entry zero time",
		"strictly increasing",
		"at least two periods",
		"non-template preset only",
		"steps require explicit transition",
		"steps 0 keeps the normal transition",
		"maximum steps depend on the time until the next entry",
	)
	writeNestedSubsection(&b, "TRANSITION")
	writeNestedLineBlock(&b,
		"steady",
		"ease-in",
		"ease-out",
		"smooth",
	)
	writeNestedSubsection(&b, "STEPS")
	writeNestedLineBlock(&b,
		"integer 0 or greater",
		"0 means no step alternation",
		"uses the selected transition curve on each leg",
		"limited by 5 seconds per leg, with a hard cap of 12",
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

	return b.String()
}
