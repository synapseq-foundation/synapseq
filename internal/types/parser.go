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

package types

const (
	// Represents an off state
	KeywordOff = "off"
	// Represents silence
	KeywordSilence = "silence"
	// Represents a comment
	KeywordComment = "#"
	// Represents an option
	KeywordOption = "@"
	// Represents a sample rate option
	KeywordOptionSampleRate = "samplerate"
	// Represents a volume option
	KeywordOptionVolume = "volume"
	// Represents an ambiance option
	KeywordOptionAmbiance = "ambiance"
	// Represents a presetlist option
	KeywordOptionPresetList = "presetlist"
	// Represents a waveform option
	KeywordWaveform = "waveform"
	// Represents a sine wave
	KeywordSine = "sine"
	// Represents a square wave
	KeywordSquare = "square"
	// Represents a triangle wave
	KeywordTriangle = "triangle"
	// Represents a sawtooth wave
	KeywordSawtooth = "sawtooth"
	// Represents a tone
	KeywordTone = "tone"
	// Represents a binaural tone
	KeywordBinaural = "binaural"
	// Represents a monaural tone
	KeywordMonaural = "monaural"
	// Represents an isochronic tone
	KeywordIsochronic = "isochronic"
	// Represents an amplitude
	KeywordAmplitude = "amplitude"
	// Represents a noise
	KeywordNoise = "noise"
	// Represents a white noise
	KeywordWhite = "white"
	// Represents a pink noise
	KeywordPink = "pink"
	// Represents a brown noise
	KeywordBrown = "brown"
	// Represents a pan effect
	KeywordPan = "pan"
	// Represents an effect
	KeywordEffect = "effect"
	// Represents an ambiance sound
	KeywordAmbiance = "ambiance"
	// Represents a modulation effect
	KeywordModulation = "modulation"
	// Represents an intensity parameter
	KeywordIntensity = "intensity"
	// Represents an pure tone
	KeywordPure = "pure"
	// Represents a steady transition
	KeywordTransitionSteady = "steady"
	// Represents a ease-out transition
	KeywordTransitionEaseOut = "ease-out"
	// Represents an ease-in transition
	KeywordTransitionEaseIn = "ease-in"
	// Represents a smooth transition
	KeywordTransitionSmooth = "smooth"
	// Represents a from to copy preset
	KeywordFrom = "from"
	// Represents a track parameter
	KeywordTrack = "track"
	// Represents an "as" keyword
	KeywordAs = "as"
	// Represents a template preset
	KeywordTemplate = "template"
	// Represents a doppler effect
	KeywordDoppler = "doppler"
)

// Parser defines the interface for parsing different content types
type Parser interface {
	// HasComment checks if the content is a comment
	HasComment() bool
	// HasOption checks if the content is an option
	HasOption() bool
	// HasPreset checks if the content is a preset
	HasPreset() bool
	// HasTrack checks if the content is a track
	HasTrack() bool
	// HasTrackOverride checks if the content is a track override
	HasTrackOverride() bool
	// HasTimeline checks if the content is a timeline
	HasTimeline() bool

	// ParseComment parses a comment content
	ParseComment() string
	// ParseOption parses an option content
	ParseOption(*SequenceOptions, string) error
	// ParsePreset parses a preset content
	ParsePreset(*[]Preset) (*Preset, error)
	// ParseTrack parses a track content
	ParseTrack() (*Track, error)
	// ParseTrackOverride parses a track override content
	ParseTrackOverride(*Preset) error
	// ParseTimeline parses a timeline content
	ParseTimeline(*[]Preset) (*Period, error)
}
