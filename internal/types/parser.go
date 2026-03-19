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

package types

import (
	"fmt"
	"maps"
	"strconv"
)

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
	// Represents an extends option
	KeywordOptionExtends = "extends"
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
	// Represents a smooth
	KeywordSmooth = "smooth"
)

// ParseOptions stores raw option values parsed from text input
type ParseOptions struct {
	Values   map[string]string
	Ambiance map[string]string
	Extends  []string
}

// NewParseOptions creates an empty ParseOptions instance
func NewParseOptions() *ParseOptions {
	return &ParseOptions{
		Values:   make(map[string]string),
		Ambiance: make(map[string]string),
		Extends:  []string{},
	}
}

// Merge merges parsed option values into the current instance
func (po *ParseOptions) Merge(other *ParseOptions) {
	if po == nil || other == nil {
		return
	}

	if po.Values == nil {
		po.Values = make(map[string]string)
	}
	if po.Ambiance == nil {
		po.Ambiance = make(map[string]string)
	}
	if po.Extends == nil {
		po.Extends = []string{}
	}

	maps.Copy(po.Values, other.Values)
	maps.Copy(po.Ambiance, other.Ambiance)
	po.Extends = append(po.Extends, other.Extends...)
}

// Build converts parsed raw options into validated SequenceOptions
func (po *ParseOptions) Build() (*SequenceOptions, error) {
	options := &SequenceOptions{
		SampleRate: 44100,
		Volume:     100,
		Ambiance:   make(map[string]string),
		Extends:    []string{},
	}

	if po == nil {
		return options, nil
	}

	if value, ok := po.Values[KeywordOptionSampleRate]; ok {
		sampleRate, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid samplerate value %q", value)
		}
		options.SampleRate = sampleRate
	}

	if value, ok := po.Values[KeywordOptionVolume]; ok {
		volume, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid volume value %q", value)
		}
		options.Volume = volume
	}

	maps.Copy(options.Ambiance, po.Ambiance)
	options.Extends = append(options.Extends, po.Extends...)

	if err := options.Validate(); err != nil {
		return nil, err
	}

	return options, nil
}

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
	ParseOption(string) (*ParseOptions, error)
	// ParsePreset parses a preset content
	ParsePreset(*[]Preset) (*Preset, error)
	// ParseTrack parses a track content
	ParseTrack() (*Track, error)
	// ParseTrackOverride parses a track override content
	ParseTrackOverride(*Preset) error
	// ParseTimeline parses a timeline content
	ParseTimeline(*[]Preset) (*Period, error)
}
