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
		"SynapSeq parses a line-oriented score language and renders deterministic audio. Ordering is strict. Sequence options must appear first. Presets must appear next. Timeline entries must appear last.",
	)
	writeParagraph(&b,
		"INPUT may be a local .spsq file, a sequence URL, or -. Output is WAV by default. -mp3 selects MP3. -preview selects HTML. -play selects direct playback. -test validates only.",
	)

	writeSection(&b, "Sound Concepts")
	writeBullet(&b, "binaural", "A pair of nearby tones, one per ear, that creates a perceived beat inside the listener. Headphones are recommended because each ear must receive a different carrier.")
	writeBullet(&b, "monaural", "A beat created by mixing the tones before playback. The pulse is present in the audio signal itself, so it can be heard on speakers or headphones.")
	writeBullet(&b, "isochronic", "A single carrier that is gated on and off at the beat rate. The result is a clearly pulsed rhythm with sharp entrainment cues.")
	writeBullet(&b, "doppler", "An effect for tone tracks that shifts motion and pitch perspective so the sound feels like it is moving toward and away from the listener.")
	writeBullet(&b, "pan", "An effect that moves energy between the left and right channels. Use it to widen static sounds or create gentle spatial motion.")
	writeBullet(&b, "modulation", "An effect that varies amplitude over time. Use it to add movement, shimmer, or slow breathing to a track.")

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
	writeExample(&b, "synapseq -new meditation")
	writeBullet(&b, "TYPE", "meditation\nfocus\nsleep\nrelaxation\nexample")
	writeBullet(&b, "-test", "validate syntax and semantics only")
	writeExample(&b, "synapseq -test session.spsq")
	writeBullet(&b, "-preview", "render HTML preview")
	writeExample(&b, "synapseq -preview session.spsq")
	writeBullet(&b, "-play", "render and play through ffplay")
	writeExample(&b, "synapseq -play session.spsq")
	writeBullet(&b, "-mp3", "render MP3 output")
	writeExample(&b, "synapseq -mp3 session.spsq")
	writeBullet(&b, "-quiet", "suppress non-error CLI output")
	writeExample(&b, "synapseq -quiet session.spsq")
	writeBullet(&b, "-no-color", "disable ANSI color output")
	writeExample(&b, "synapseq -manual -no-color")
	writeBullet(&b, "-manual", "print full manual")
	writeExample(&b, "synapseq -manual")
	writeBullet(&b, "-help", "print concise command overview")
	writeExample(&b, "synapseq -help")
	writeBullet(&b, "-version", "print version build and platform information")
	writeExample(&b, "synapseq -version")

	writeSubsection(&b, "Hub")
	writeBullet(&b, "-hub-update", "refresh local Hub index")
	writeExample(&b, "synapseq -hub-update")
	writeBullet(&b, "-hub-clean", "remove cached Hub data")
	writeExample(&b, "synapseq -hub-clean")
	writeBullet(&b, "-hub-list", "list sequences from local Hub index")
	writeExample(&b, "synapseq -hub-list")
	writeBullet(&b, "-hub-search WORD", "search local Hub index")
	writeExample(&b, "synapseq -hub-search calm")
	writeBullet(&b, "-hub-info NAME", "print Hub sequence metadata")
	writeExample(&b, "synapseq -hub-info calm-state")
	writeBullet(&b, "-hub-download NAME [DIR]", "download sequence and dependencies\nDIR optional")
	writeExample(&b, "synapseq -hub-download calm-state downloads")
	writeBullet(&b, "-hub-get NAME [OUTPUT]", "download and render sequence\nOUTPUT optional")
	writeExample(&b, "synapseq -hub-get calm-state calm-state.wav")

	writeSubsection(&b, "External tools")
	writeBullet(&b, "-ffmpeg-path PATH", "use specific ffmpeg executable")
	writeExample(&b, "synapseq -ffmpeg-path /usr/local/bin/ffmpeg -mp3 session.spsq")
	writeBullet(&b, "-ffplay-path PATH", "use specific ffplay executable")
	writeExample(&b, "synapseq -ffplay-path /usr/local/bin/ffplay -play session.spsq")

	writeSubsection(&b, "Windows")
	writeBullet(&b, "-install-file-association", "associate .spsq files with SynapSeq")
	writeExample(&b, "synapseq -install-file-association")
	writeBullet(&b, "-uninstall-file-association", "remove .spsq file association")
	writeExample(&b, "synapseq -uninstall-file-association")

	writeSection(&b, "Sequence File")
	writeSubsection(&b, "File type")
	writeLineBlock(&b,
		"plain-text .spsq",
		"parsed top to bottom",
	)

	writeSubsection(&b, "Order")
	writeLineBlock(&b,
		"sequence options",
		"presets",
		"timeline",
	)

	writeSubsection(&b, "Top level")
	writeLineBlock(&b,
		"no indentation",
		"sequence options must start in column 1",
		"preset declarations must start in column 1",
		"timeline entries must start in column 1",
	)

	writeSubsection(&b, "Preset body")
	writeLineBlock(&b,
		"exactly two leading spaces required",
		"track definitions use this indentation",
		"track overrides use this indentation",
	)

	writeSubsection(&b, "Comments")
	writeBullet(&b, "# comment", "ignored")
	writeBullet(&b, "## comment", "stored and printed unless -quiet is set")
	writeBullet(&b, "Rules", "full line only\ninline comment not permitted")

	writeSubsection(&b, "Sequence Options")
	writeParagraph(&b,
		"Sequence options configure the whole file and must appear before the first preset or timeline entry.",
	)
	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"valid only before first preset",
		"valid only before first timeline entry",
	)
	writeNestedSubsection(&b, "Commands")
	writeNestedBullet(&b, "@samplerate NUMBER", "set output sample rate\ndefault 44100")
	writeNestedExample(&b, "@samplerate 48000")
	writeNestedBullet(&b, "@volume NUMBER", "set master volume\nrange 0 to 100\ndefault 100")
	writeNestedExample(&b, "@volume 80")
	writeNestedBullet(&b, "@ambiance NAME PATH_OR_URL", "register a named ambiance source for later use in presets")
	writeNestedExample(&b, "@ambiance rain audio/rain")
	writeNestedBullet(&b, "@extends PATH_OR_URL", "load modular .spsc content before the main file is parsed")
	writeNestedExample(&b, "@extends library/common")
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
		"allowed: options presets tracks track overrides",
		"not permitted: timeline entries nested @extends",
	)

	writeSubsection(&b, "Presets")
	writeParagraph(&b,
		"Presets define reusable track groups. Declare them at the top level and place track lines inside the body with exactly two leading spaces.",
	)
	writeNestedSubsection(&b, "Identifiers")
	writeNestedLineBlock(&b,
		"first character must be a letter",
		"remaining characters may be letters digits underscores dashes",
		"maximum length 20",
		"preset references are case-insensitive",
		"preset names normalize to lowercase",
	)
	writeNestedSubsection(&b, "Preset declarations")
	writeNestedCodeBlock(&b,
		"NAME",
		"NAME as template",
		"NAME from TEMPLATE_NAME",
	)
	writeNestedSubsection(&b, "Preset example")
	writeNestedCodeBlock(&b,
		"focus",
		"  noise pink amplitude 30",
		"  tone 180 binaural 10 amplitude 18",
	)
	writeNestedSubsection(&b, "Template rules")
	writeNestedLineBlock(&b,
		"template preset cannot appear on timeline",
	)
	writeNestedSubsection(&b, "Inheritance rules")
	writeNestedLineBlock(&b,
		"inherited preset cannot define new track lines",
		"inherited preset may contain track overrides only",
	)
	writeNestedSubsection(&b, "Inherited preset example")
	writeNestedCodeBlock(&b,
		"base as template",
		"  noise pink smooth 20 amplitude 30",
		"  tone 220 binaural 10 amplitude 18",
		"",
		"focus from base",
		"  track 1 smooth 35",
		"  track 1 amplitude +5",
		"  track 2 tone 180",
		"  track 2 binaural -2",
		"  track 2 amplitude +4",
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
		"waveform prefix allowed on tone track",
		"waveform prefix allowed on ambiance track",
		"effect TYPE VALUE intensity PERCENT appears before amplitude",
	)
	writeNestedSubsection(&b, "Waveforms")
	writeNestedCodeBlock(&b,
		"waveform VALUE",
	)
	writeIndentedSubsection(&b, 12, "VALUE")
	writeIndentedLineBlock(&b, 16,
		"sine square triangle sawtooth",
	)
	writeNestedSubsection(&b, "Effects")
	writeNestedCodeBlock(&b,
		"effect TYPE VALUE intensity PERCENT",
	)
	writeIndentedSubsection(&b, 12, "TYPE")
	writeIndentedLineBlock(&b, 16,
		"pan modulation doppler",
	)
	writeIndentedSubsection(&b, 12, "Allowed")
	writeIndentedBullet(&b, 16, 20, "tone", "pan modulation doppler")
	writeIndentedBullet(&b, 16, 20, "noise", "pan modulation")
	writeIndentedBullet(&b, 16, 20, "ambiance", "pan modulation")
	writeNestedSubsection(&b, "Track overrides")
	writeIndentedSubsection(&b, 12, "Scope")
	writeIndentedLineBlock(&b, 16,
		"inherited presets only",
	)
	writeIndentedSubsection(&b, 12, "Index")
	writeIndentedLineBlock(&b, 16,
		"1-based",
	)
	writeIndentedSubsection(&b, 12, "Syntax")
	writeIndentedCodeBlock(&b, 16,
		"track N tone VALUE",
		"track N binaural|monaural|isochronic VALUE",
		"track N smooth VALUE",
		"track N pan|modulation|doppler VALUE",
		"track N intensity VALUE",
		"track N amplitude VALUE",
	)
	writeIndentedSubsection(&b, 12, "Rules")
	writeIndentedLineBlock(&b, 16,
		"VALUE may be absolute or signed relative delta",
		"+VALUE adds to inherited value",
		"-VALUE subtracts from inherited value",
		"matching existing effect type only",
	)

	writeSubsection(&b, "Timeline")
	writeParagraph(&b,
		"Timeline entries schedule presets over time. Each entry starts with a timestamp and may include a transition keyword.",
	)
	writeNestedSubsection(&b, "Syntax")
	writeNestedCodeBlock(&b,
		"HH:MM:SS PRESET_NAME [TRANSITION]",
	)
	writeNestedSubsection(&b, "Time")
	writeNestedLineBlock(&b,
		"HH:MM:SS required",
	)
	writeNestedSubsection(&b, "Rules")
	writeNestedLineBlock(&b,
		"top-level lines only",
		"first entry zero time",
		"strictly increasing",
		"at least two periods",
		"non-template preset only",
	)
	writeNestedSubsection(&b, "TRANSITION")
	writeNestedLineBlock(&b,
		"steady ease-in ease-out smooth",
	)
	writeNestedSubsection(&b, "Timeline examples")
	writeNestedCodeBlock(&b,
		"00:00:00 silence",
		"00:00:30 focus ease-in",
		"00:10:00 focus-deep smooth",
		"00:20:00 silence ease-out",
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
