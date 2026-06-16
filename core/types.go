// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"io"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// AppContext holds the configuration for the application.
// It provides a safe, immutable context for sequence processing.
// Methods that modify the context return a new instance.
type AppContext struct {
	statusOutput io.Writer
	statusColors bool
}

// LoadedContext holds a loaded sequence and execution settings.
type LoadedContext struct {
	appCtx   *AppContext
	sequence *t.Sequence
}

// Preset holds a preset, including its name and tracks.
type Preset struct {
	Name   string
	Tracks []Track
}

// Track holds a track, including its
// index, type, line, amplitude, carrier, resonance, waveform,
// source name, noise smooth, and effect.
type Track struct {
	Index       int
	Waveform    string
	Type        string
	Carrier     float64
	Resonance   float64
	Amplitude   float64
	SourceName  string
	NoiseSmooth float64
	Effect      Effect
	Line        string
}

// Effect holds an effect, including its type, value, and intensity.
type Effect struct {
	Type      string
	Value     float64
	Intensity float64
}

// TimelineEntry holds a timeline entry for a preset, including its time, preset name, transition, steps, and line.
type TimelineEntry struct {
	Timestamp  string
	PresetName string
	Transition string
	Steps      int
	Line       string
}
