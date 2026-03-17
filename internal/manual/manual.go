/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
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
	writeParagraph(&b,
		"synapseq - text-driven audio sequencer for brainwave sessions and ambient sound design.",
	)

	writeSection(&b, "Synopsis")
	writeCodeBlock(&b,
		"synapseq [OPTION]... INPUT [OUTPUT]",
		"synapseq -new TYPE",
		"synapseq -preview INPUT",
		"synapseq -play INPUT",
		"synapseq -test INPUT",
		"synapseq -hub-update",
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
	writeParagraph(&b,
		"For paged reading, use the command that matches your shell environment.",
	)
	writeSubsection(&b, "Linux and macOS")
	writeCodeBlock(&b,
		"synapseq -manual | less",
	)
	writeSubsection(&b, "Windows PowerShell")
	writeCodeBlock(&b,
		"synapseq -manual | more",
	)

	writeSection(&b, "Description")
	writeParagraph(&b,
		"SynapSeq is a text-driven audio sequencer for building repeatable brainwave and ambient listening sessions from plain-text score files. Instead of editing waveforms directly, you describe structure, timing, and sound layers in a small domain-specific language and let the renderer generate the final audio.",
	)
	writeParagraph(&b,
		"A typical sequence combines oscillator-based tones, entrainment beats, filtered noise, ambiance layers, waveform selection, and motion effects. These building blocks are arranged into named presets and then placed on a timeline so the session can evolve predictably over minutes or hours.",
	)
	writeParagraph(&b,
		"The parser is strict by design: options belong at the top, presets come next, and timeline entries close the file. This keeps sessions readable, deterministic, and easier to validate before rendering or playback.",
	)
	writeParagraph(&b,
		"SynapSeq can read local files, remote sequence URLs, or standard input. It can render WAV or MP3 output, stream raw PCM to standard output, preview the session as HTML, play directly through ffplay, and fetch published content from the SynapSeq Hub.",
	)
	writeParagraph(&b,
		"For larger libraries, the language supports modular composition through @extends and reusable template presets, making it practical to share preset families across multiple sessions without duplicating track definitions.",
	)

	writeSection(&b, "Command Line")
	writeParagraph(&b,
		"The command line accepts one primary action at a time plus supporting flags. When INPUT is present, it may be a local .spsq file, a URL, or - for standard input. When OUTPUT is omitted, SynapSeq derives a default output path from the input name and target format.",
	)
	writeParagraph(&b,
		"The following options are available from the CLI.",
	)
	writeSubsection(&b, "General Options")
	writeOption(&b, "-new TYPE", "Create a starter sequence from a built-in template. Supported types are meditation, focus, sleep, relaxation, and example.")
	writeOption(&b, "-test", "Validate sequence syntax and semantics without generating audio output.")
	writeOption(&b, "-preview", "Render an HTML timeline preview instead of audio.")
	writeOption(&b, "-play", "Render and play the result directly with ffplay.")
	writeOption(&b, "-quiet", "Suppress non-error CLI output.")
	writeOption(&b, "-no-color", "Disable ANSI colors in CLI output. Useful for pipes, logs, and pagers configured without raw-control support.")
	writeOption(&b, "-manual", "Print the full manual.")
	writeOption(&b, "-help", "Show the concise command overview.")
	writeOption(&b, "-version", "Print version, build, and platform information.")
	writeParagraph(&b)

	writeSubsection(&b, "Hub Options")
	writeOption(&b, "-hub-update", "Refresh the local index of sequences available in the SynapSeq Hub.")
	writeOption(&b, "-hub-clean", "Remove cached Hub data from the local machine.")
	writeOption(&b, "-hub-list", "List sequences currently available from the local Hub index.")
	writeOption(&b, "-hub-search WORD", "Search the Hub index for matching sequence names or metadata.")
	writeOption(&b, "-hub-info NAME", "Show metadata and descriptive details for a Hub sequence.")
	writeOption(&b, "-hub-download NAME [DIR]", "Download a Hub sequence and its dependencies. If DIR is provided, files are stored there.")
	writeOption(&b, "-hub-get NAME [OUTPUT]", "Download a Hub sequence and render it in one step. If OUTPUT is provided, it is used as the target file name.")
	writeParagraph(&b)

	writeSubsection(&b, "External Tool Options")
	writeOption(&b, "-ffmpeg-path PATH", "Use a specific ffmpeg executable when exporting MP3.")
	writeOption(&b, "-ffplay-path PATH", "Use a specific ffplay executable when playing audio directly.")
	writeParagraph(&b)

	writeSubsection(&b, "Windows Options")
	writeOption(&b, "-install-file-association", "Associate .spsq files with SynapSeq on Windows.")
	writeOption(&b, "-uninstall-file-association", "Remove the SynapSeq .spsq file association on Windows.")

	writeSection(&b, "File Layout")
	writeBullet(&b, "Options", "Start with lines beginning with @. They define global settings such as samplerate, volume, ambiance files, and modular imports.")
	writeBullet(&b, "Presets", "Define named sound setups. A preset line has no indentation. Tracks under that preset must start with exactly two spaces.")
	writeBullet(&b, "Timeline", "Schedule presets over time using HH:MM:SS. Timeline entries are always top-level lines with no indentation.")
	writeParagraph(&b,
		"The first timeline entry must start at 00:00:00, times must be strictly increasing, and a valid sequence needs at least two periods.",
	)
	writeCodeBlock(&b,
		"@volume 90",
		"@samplerate 44100",
		"",
		"focus",
		"  noise white amplitude 10",
		"  tone 240 binaural 16 amplitude 15",
		"",
		"00:00:00 silence",
		"00:00:20 focus",
		"00:19:30 focus",
		"00:20:00 silence",
	)

	writeSection(&b, "Comments")
	writeBullet(&b, "# comment", "Ignored by the parser. Useful for notes and documentation.")
	writeBullet(&b, "## comment", "Stored as a sequence comment and shown by the CLI before rendering when output is not quiet.")
	writeBullet(&b, "Standalone lines only", "Comments must occupy their own lines. Inline comments after options, preset names, track lines, overrides, or timeline entries are invalid syntax.")
	writeCodeBlock(&b,
		"@samplerate 48000 # samplerate",
	)
	writeParagraph(&b,
		"The example above is invalid. Write comments only on otherwise empty lines.",
	)

	writeSection(&b, "Sequence Options")
	writeBullet(&b, "@samplerate NUMBER", "Set the output sample rate. Default is 44100.")
	writeBullet(&b, "@volume NUMBER", "Set the global output volume from 0 to 100. Default is 100.")
	writeBullet(&b, "@ambiance NAME PATH_OR_URL", "Register an ambiance source that can be referenced later by track definitions.")
	writeBullet(&b, "@extends PATH_OR_URL", "Load presets and options from a modular .spsc file before the main sequence is built.")
	writeParagraph(&b,
		"For local @ambiance and @extends paths, use relative paths with forward slashes and omit the file extension. Local ambiance paths resolve to .wav and local extends paths resolve to .spsc.",
		"Absolute paths, parent traversal, Windows drive prefixes, and stdin are rejected for local modular paths.",
	)
	writeCodeBlock(&b,
		"@volume 85",
		"@samplerate 48000",
		"@ambiance rain audio/rain",
		"@ambiance ocean https://example.com/ocean.wav",
		"@extends library/common",
	)

	writeSection(&b, "Preset Names")
	writeParagraph(&b,
		"Named references should be short and predictable. A preset or ambiance name must start with a letter, may contain letters, digits, underscores, or dashes, and may not exceed 20 characters.",
		"Use lowercase names consistently. Presets are normalized to lowercase internally.",
	)

	writeSection(&b, "Presets")
	writeBullet(&b, "Regular preset", "A plain top-level name creates a new preset.")
	writeBullet(&b, "Template preset", "Use as template to define a reusable base that cannot appear directly in the timeline.")
	writeBullet(&b, "Inherited preset", "Use from TEMPLATE_NAME to clone a template preset and then override selected track fields.")
	writeCodeBlock(&b,
		"focus-base as template",
		"  noise white amplitude 15",
		"  tone 240 binaural 14 amplitude 12",
		"",
		"focus-deep from focus-base",
		"  track 1 amplitude 25",
		"  track 2 binaural 18",
	)
	writeParagraph(&b,
		"A preset that inherits from a template cannot define new track lines. It must use track overrides instead.",
	)

	writeSection(&b, "Track Definitions")
	writeParagraph(&b,
		"Track lines belong under a preset and must start with exactly two spaces. SynapSeq allocates tracks automatically in declaration order.",
	)
	writeBullet(&b, "Pure tone", "tone CARRIER amplitude LEVEL")
	writeBullet(&b, "Beat tone", "tone CARRIER binaural|monaural|isochronic BEAT amplitude LEVEL")
	writeBullet(&b, "Noise", "noise white|pink|brown [smooth VALUE] amplitude LEVEL")
	writeBullet(&b, "Ambiance", "ambiance NAME amplitude LEVEL")
	writeParagraph(&b,
		"Amplitude is required on every concrete track definition. Tone carrier and beat values are positive numbers. Noise smooth ranges from 0 to 100.",
	)
	writeCodeBlock(&b,
		"alpha",
		"  tone 300 amplitude 10",
		"  tone 220 binaural 8 amplitude 14",
		"  noise pink smooth 35 amplitude 20",
		"  ambiance rain amplitude 25",
	)

	writeSection(&b, "Sound Concepts")
	writeBullet(&b, "tone", "A steady audible carrier frequency with no beat component. Use it when you want a plain oscillator without binaural, monaural, or isochronic pulsing.")
	writeBullet(&b, "binaural", "A beat created by feeding slightly different frequencies to the left and right channels. In SynapSeq syntax, the number after binaural is the beat frequency, not the carrier.")
	writeBullet(&b, "monaural", "A beat produced by combining close frequencies into a single perceived pulse before it reaches the ears. It still uses a carrier plus a beat value, but behaves differently from binaural playback.")
	writeBullet(&b, "isochronic", "A rhythmic single-tone pulse where the sound is sharply gated on and off. Use it when you want a clearly defined pulse rather than two interacting tones.")
	writeBullet(&b, "noise", "A broadband texture instead of a pitched oscillator. White, pink, and brown noise differ in spectral balance, and smooth controls how gently the noise evolves over time.")
	writeBullet(&b, "ambiance", "A named WAV file or remote WAV resource mixed into the preset. Ambiance is useful for rain, drones, environmental beds, and other non-synth layers.")
	writeParagraph(&b,
		"A practical rule: tone defines the audible pitch, binaural or other beat keywords define how that pitch is modulated or split, amplitude controls loudness, and optional effects add movement.",
	)

	writeSection(&b, "Waveforms")
	writeParagraph(&b,
		"You can prefix tone and ambiance tracks with waveform to change the oscillator shape. Valid values are sine, square, triangle, and sawtooth.",
	)
	writeCodeBlock(&b,
		"shape-demo",
		"  waveform triangle tone 250 monaural 8 amplitude 10",
		"  waveform square ambiance rain amplitude 20",
	)

	writeSection(&b, "Effects")
	writeParagraph(&b,
		"Effects are optional and appear before amplitude. All effects require a value and an intensity percentage.",
	)
	writeBullet(&b, "pan", "Moves the sound across the stereo field over time. The effect value controls the pan rate and intensity controls how far the movement is pushed.")
	writeBullet(&b, "modulation", "Applies cyclic amplitude movement. The effect value acts as the modulation rate and intensity controls how deep the volume swing becomes.")
	writeBullet(&b, "doppler", "Applies motion-based pitch and stereo movement for tone tracks. It is available only on tone-based tracks and is useful for a drifting or orbiting sensation.")
	writeBullet(&b, "Tone effects", "pan, modulation, and doppler are valid for tone-based tracks.")
	writeBullet(&b, "Noise effects", "pan and modulation are valid for noise tracks.")
	writeBullet(&b, "Ambiance effects", "pan and modulation are valid for ambiance tracks.")
	writeCodeBlock(&b,
		"fx-demo",
		"  tone 300 binaural 10 effect modulation 6 intensity 40 amplitude 18",
		"  tone 220 effect doppler 0.9 intensity 70 amplitude 14",
		"  noise white effect pan 0.5 intensity 60 amplitude 12",
		"  ambiance rain effect modulation 4 intensity 35 amplitude 20",
	)

	writeSection(&b, "Track Overrides")
	writeParagraph(&b,
		"Track overrides are only valid inside presets created with from TEMPLATE_NAME. The syntax is 1-based: track 1 targets the first inherited track.",
	)
	writeBullet(&b, "tone VALUE", "Change the carrier frequency of an inherited tone track.")
	writeBullet(&b, "binaural|monaural|isochronic VALUE", "Change the beat value of an inherited matching beat track.")
	writeBullet(&b, "pan|modulation|doppler VALUE", "Change the effect value when that effect already exists on the inherited track.")
	writeBullet(&b, "intensity VALUE", "Change effect intensity.")
	writeBullet(&b, "amplitude VALUE", "Change the track amplitude.")
	writeCodeBlock(&b,
		"focus-strong from focus-base",
		"  track 1 amplitude 30",
		"  track 2 tone 250",
		"  track 2 binaural 18",
		"  track 2 intensity 75",
	)

	writeSection(&b, "Timeline")
	writeBullet(&b, "Format", "HH:MM:SS PRESET_NAME [TRANSITION]")
	writeBullet(&b, "Default transition", "steady")
	writeBullet(&b, "Other transitions", "ease-in, ease-out, smooth")
	writeParagraph(&b,
		"Timeline entries must reference non-template presets. The parser adjusts each period against the next one to build transitions and track interpolation.",
	)
	writeBullet(&b, "steady", "Hold the current preset values without applying an easing curve. This is the neutral transition and the default when you omit the transition keyword.")
	writeBullet(&b, "ease-in", "Start the transition gently and accelerate into the target state. Useful when you want the next section to arrive progressively rather than immediately.")
	writeBullet(&b, "ease-out", "Move quickly at first and then settle more gently as the next state is reached. Useful for soft landings near the end of a section.")
	writeBullet(&b, "smooth", "Apply a balanced easing curve across the whole transition. This is usually the best choice for gradual meditative ramps.")
	writeBullet(&b, "Compatibility rule", "Consecutive timeline entries cannot reuse the same channel with an incompatible track type, effect type, or ambiance source.")
	writeBullet(&b, "Silence bridge", "If you need to switch a channel from one incompatible sound design to another, insert a silence preset between those timeline entries.")
	writeBullet(&b, "Direct on/off changes", "A channel should not jump directly between an active track and off; use silence as the bridge state.")
	writeCodeBlock(&b,
		"00:00:00 silence",
		"00:00:20 focus-light",
		"00:04:00 focus-light smooth",
		"00:07:00 focus",
		"00:19:30 focus",
		"00:20:00 silence",
	)
	writeCodeBlock(&b,
		"# Good: incompatible presets are separated by silence",
		"00:00:00 silence",
		"00:00:15 doppler-preset",
		"00:00:30 silence",
		"00:00:45 pan-preset",
	)

	writeSection(&b, "Extended Files")
	writeParagraph(&b,
		"Files loaded by @extends are modular .spsc files. They may contain options, presets, tracks, and track overrides.",
		"They may not contain timeline entries, and they may not contain another @extends option.",
	)
	writeBullet(&b, "Library file example", "A .spsc file is useful for shared preset libraries. Keep reusable templates or presets there, then import them from your main .spsq session.")
	writeCodeBlock(&b,
		"# library/common.spsc",
		"focus-template as template",
		"  tone 240 binaural 10 amplitude 12",
		"  noise pink smooth 25 amplitude 8",
	)
	writeBullet(&b, "Main sequence example", "Your main .spsq file can import that library with @extends and then build concrete presets or timeline entries from it.")
	writeCodeBlock(&b,
		"@extends library/common",
		"",
		"focus-light from focus-template",
		"  track 1 amplitude 10",
		"  track 2 amplitude 6",
		"",
		"00:00:00 silence",
		"00:00:20 focus-light",
		"00:10:00 silence",
	)

	writeSection(&b, "Examples")
	writeParagraph(&b,
		"Start with a minimal session first. Once that structure is clear, add more layers or modular reuse.",
	)
	writeBullet(&b, "Basic session", "A minimal valid file: one preset, one audible section, and silence at the beginning and end.")
	writeCodeBlock(&b,
		"@volume 90",
		"",
		"focus",
		"  tone 240 binaural 10 amplitude 15",
		"",
		"00:00:00 silence",
		"00:00:20 focus",
		"00:00:40 silence",
	)
	writeBullet(&b, "Layered preset", "Add noise or ambiance when you want more texture without changing the overall preset structure.")
	writeCodeBlock(&b,
		"deep-rest",
		"  noise brown smooth 45 amplitude 12",
		"  tone 180 binaural 6 effect modulation 4 intensity 30 amplitude 16",
		"  ambiance rain effect pan 0.3 intensity 35 amplitude 22",
		"",
		"00:00:00 silence",
		"00:00:20 deep-rest",
		"00:10:00 silence",
	)
	writeBullet(&b, "Reusable templates", "Use @extends and template inheritance when several sessions share the same preset family.")
	writeCodeBlock(&b,
		"@extends library/focus-base",
		"@ambiance rain audio/rain",
		"",
		"focus-light from focus-template",
		"  track 1 amplitude 8",
		"  track 2 binaural 12",
		"",
		"focus-deep from focus-template",
		"  track 1 amplitude 18",
		"  track 2 binaural 18",
		"  track 2 intensity 65",
	)

	writeSection(&b, "Notes")
	writeBullet(&b, "Use silence", "A built-in preset named silence is always available and is ideal for the first or last timeline entry.")
	writeBullet(&b, "Indentation matters", "Preset children must use two leading spaces. A top-level tone/noise/ambiance line is invalid syntax.")
	writeBullet(&b, "Keep options first", "Once presets or timeline entries start, additional options are rejected.")
	writeBullet(&b, "Prefer templates for families", "If several presets share the same structure, declare a template and override only what changes.")
	writeBullet(&b, "Track positions are fixed", "Direct transitions compare tracks by position. Track 1 must stay the same kind of sound across directly connected presets, and the same rule applies to every later track slot.")
	writeBullet(&b, "Structural changes need silence", "Do not switch a track directly between tone, noise, or ambiance types. If the structure must change, insert the built-in silence preset between those timeline entries.")
	writeBullet(&b, "Beat, noise, and effect types must match", "Across directly connected presets, keep binaural versus monaural versus isochronic mode unchanged, keep white versus pink versus brown noise unchanged, and keep pan versus modulation versus doppler unchanged.")
	writeBullet(&b, "Waveforms may change", "Waveform shape is treated as a parameter, not a structural type. It may change between otherwise compatible tone or ambiance tracks.")
	writeParagraph(&b)

	writeSection(&b, "See Also")
	writeParagraph(&b,
		"Use synapseq -help for a concise command overview, or visit the online documentation for installation guides and broader examples.",
	)
	writeCodeBlock(&b,
		"synapseq -help",
	)

	return b.String()
}

func writeTitle(b *strings.Builder, title, subtitle string) {
	b.WriteString(cli.Title(title))
	b.WriteString("\n")
	b.WriteString(cli.Muted(subtitle))
	b.WriteString("\n\n")
}

func writeSection(b *strings.Builder, title string) {
	b.WriteString(cli.Section(strings.ToUpper(title)))
	b.WriteString("\n")
}

func writeSubsection(b *strings.Builder, title string) {
	b.WriteString("    ")
	b.WriteString(cli.Label(title))
	b.WriteString("\n\n")
}

func writeParagraph(b *strings.Builder, lines ...string) {
	for _, line := range lines {
		if line == "" {
			continue
		}
		writeWrappedLine(b, line, "    ", "    ")
	}
	b.WriteString("\n")
}

func writeBullet(b *strings.Builder, label, description string) {
	b.WriteString("    ")
	b.WriteString(cli.Label(label + ":"))
	b.WriteString("\n")
	writeWrappedLine(b, description, "        ", "        ")
	b.WriteString("\n")
}

func writeOption(b *strings.Builder, flag, description string) {
	b.WriteString("    ")
	b.WriteString(cli.Label(flag))
	b.WriteString("\n")
	writeWrappedLine(b, description, "        ", "        ")
	b.WriteString("\n")
}

func writeCodeBlock(b *strings.Builder, lines ...string) {
	for _, line := range lines {
		if line == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString("        ")
		b.WriteString(cli.Command(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func writeWrappedLine(b *strings.Builder, text, firstPrefix, continuationPrefix string) {
	writeWrappedLineWithPrefixes(b, text, firstPrefix, continuationPrefix, len(firstPrefix), len(continuationPrefix))
}

func writeWrappedLineWithPrefixes(b *strings.Builder, text, firstPrefix, continuationPrefix string, firstWidth, continuationWidth int) {
	words := strings.Fields(text)
	if len(words) == 0 {
		b.WriteString("\n")
		return
	}

	prefix := firstPrefix
	available := manualWidth - firstWidth
	lineLen := 0

	b.WriteString(prefix)
	for index, word := range words {
		wordLen := len(word)
		separatorLen := 0
		if lineLen > 0 {
			separatorLen = 1
		}

		if lineLen > 0 && lineLen+separatorLen+wordLen > available {
			b.WriteString("\n")
			prefix = continuationPrefix
			available = manualWidth - continuationWidth
			b.WriteString(prefix)
			b.WriteString(word)
			lineLen = wordLen
			continue
		}

		if index > 0 && lineLen > 0 {
			b.WriteString(" ")
		}
		b.WriteString(word)
		lineLen += separatorLen + wordLen
	}
	if lineLen == 0 {
		b.WriteString(strings.Join(words, " "))
	}
	b.WriteString("\n")
}
