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

	writeTitle(&b, fmt.Sprintf("SynapSeq Manual %s", info.VERSION))
	writeParagraph(&b,
		"SynapSeq is a text-driven audio sequencer. A sequence file usually has three major sections: options, presets, and timeline.",
		"The parser is strict by design: options belong at the top, presets come next, and timeline entries close the file.",
	)

	writeSection(&b, "Command Line")
	writeParagraph(&b,
		"Use the manual directly from the terminal whenever you need a full language reference.",
	)
	writeCodeBlock(&b,
		"synapseq -manual",
		"synapseq -manual -no-color",
	)
	writeParagraph(&b,
		"The manual is printed in English and follows the same color rules as the rest of the CLI.",
	)

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

	writeSection(&b, "Options")
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
	writeCodeBlock(&b,
		"00:00:00 silence",
		"00:00:20 focus-light",
		"00:04:00 focus-light smooth",
		"00:07:00 focus",
		"00:19:30 focus",
		"00:20:00 silence",
	)

	writeSection(&b, "Extended Files")
	writeParagraph(&b,
		"Files loaded by @extends are modular .spsc files. They may contain options, presets, tracks, and track overrides.",
		"They may not contain timeline entries, and they may not contain another @extends option.",
	)

	writeSection(&b, "Advanced Examples")
	writeParagraph(&b,
		"The following examples show how the main building blocks combine in real-world usage.",
	)
	writeBullet(&b, "Reusable template library", "Keep common preset structures in a .spsc file and import them with @extends so multiple sessions can share the same sound design vocabulary.")
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
	writeBullet(&b, "Long-form session", "Use several timeline points with smooth transitions to guide the listener through stages without abrupt jumps.")
	writeCodeBlock(&b,
		"00:00:00 silence",
		"00:00:30 focus-light",
		"00:05:00 focus-light smooth",
		"00:08:00 focus-deep",
		"00:18:00 focus-deep smooth",
		"00:20:00 silence",
	)
	writeBullet(&b, "Layered preset", "Blend a beat track with noise and ambiance when you want pitch, texture, and environmental depth at the same time.")
	writeCodeBlock(&b,
		"deep-rest",
		"  noise brown smooth 45 amplitude 12",
		"  tone 180 binaural 6 effect modulation 4 intensity 30 amplitude 16",
		"  ambiance rain effect pan 0.3 intensity 35 amplitude 22",
	)

	writeSection(&b, "Practical Notes")
	writeBullet(&b, "Use silence", "A built-in preset named silence is always available and is ideal for the first or last timeline entry.")
	writeBullet(&b, "Indentation matters", "Preset children must use two leading spaces. A top-level tone/noise/ambiance line is invalid syntax.")
	writeBullet(&b, "Keep options first", "Once presets or timeline entries start, additional options are rejected.")
	writeBullet(&b, "Prefer templates for families", "If several presets share the same structure, declare a template and override only what changes.")
	writeParagraph(&b)

	writeSection(&b, "More Help")
	writeParagraph(&b,
		"Use synapseq -help for a concise command overview, or visit the online documentation for installation guides and broader examples.",
	)
	writeCodeBlock(&b,
		"synapseq -help",
		info.DOC_URL,
	)

	return b.String()
}

func writeTitle(b *strings.Builder, title string) {
	b.WriteString(cli.Title(title))
	b.WriteString("\n\n")
}

func writeSection(b *strings.Builder, title string) {
	b.WriteString(cli.Section(title))
	b.WriteString("\n")
}

func writeParagraph(b *strings.Builder, lines ...string) {
	for _, line := range lines {
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func writeBullet(b *strings.Builder, label, description string) {
	b.WriteString("- ")
	b.WriteString(cli.Label(label + ":"))
	b.WriteString(" ")
	b.WriteString(description)
	b.WriteString("\n")
}

func writeCodeBlock(b *strings.Builder, lines ...string) {
	for _, line := range lines {
		if line == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString("  ")
		b.WriteString(cli.Command(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}
