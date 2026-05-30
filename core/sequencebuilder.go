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

package core

import (
	"fmt"
	"strconv"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// sequenceBuilder is responsible for building a sequence from a string sequence content
type sequenceBuilder struct {
	options    map[string]string
	presets    map[string][]t.Track
	timeline   []string
	lastPreset string
}

// NewSequenceBuilder creates a new sequenceBuilder
func NewSequenceBuilder() *sequenceBuilder {
	defaultOptions := map[string]string{
		t.KeywordOptionSampleRate: "44100",
		t.KeywordOptionVolume:     "100",
	}

	return &sequenceBuilder{
		options:    defaultOptions,
		presets:    make(map[string][]t.Track),
		timeline:   make([]string, 0),
		lastPreset: "",
	}
}

// WithSampleRate sets the sample rate for the sequence
func (sb *sequenceBuilder) WithSampleRate(sampleRate int) *sequenceBuilder {
	sb.options[t.KeywordOptionSampleRate] = strconv.Itoa(sampleRate)
	return sb
}

// WithVolume sets the volume for the sequence
func (sb *sequenceBuilder) WithVolume(volume int) *sequenceBuilder {
	sb.options[t.KeywordOptionVolume] = strconv.Itoa(volume)
	return sb
}

// WithAmbiance sets the ambiance for the sequence
func (sb *sequenceBuilder) WithAmbiance(name, path string) *sequenceBuilder {
	sb.options[t.KeywordOptionAmbiance] = name + " " + path
	return sb
}

// Extends sets the extends for the sequence
func (sb *sequenceBuilder) Extends(extends string) *sequenceBuilder {
	sb.options[t.KeywordOptionExtends] = extends
	return sb
}

// WithPreset sets the last preset for the sequence
func (sb *sequenceBuilder) WithPreset(name string) *sequenceBuilder {
	sb.presets[name] = []t.Track{}
	sb.lastPreset = name
	return sb
}

// WithPureTone adds a pure tone track to the sequence
func (sb *sequenceBuilder) WithPureTone(waveform string, carrier, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:      t.TrackPureTone,
		Waveform:  t.WaveformString(waveform),
		Carrier:   carrier,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithBinauralTone adds a binaural tone track to the sequence
func (sb *sequenceBuilder) WithBinauralTone(waveform string, carrier, beat, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:      t.TrackBinauralBeat,
		Waveform:  t.WaveformString(waveform),
		Carrier:   carrier,
		Resonance: beat,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithMonauralTone adds a monaural tone track to the sequence
func (sb *sequenceBuilder) WithMonauralTone(waveform string, carrier, beat, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:      t.TrackMonauralBeat,
		Waveform:  t.WaveformString(waveform),
		Carrier:   carrier,
		Resonance: beat,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithIsochronicTone adds an isochronic tone track to the sequence
func (sb *sequenceBuilder) WithIsochronicTone(waveform string, carrier, beat, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:      t.TrackIsochronicBeat,
		Waveform:  t.WaveformString(waveform),
		Carrier:   carrier,
		Resonance: beat,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithPinkNoise adds a pink noise track to the sequence
func (sb *sequenceBuilder) WithPinkNoise(smooth, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:        t.TrackPinkNoise,
		Amplitude:   t.AmplitudePercentToRaw(amplitude),
		NoiseSmooth: smooth,
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithBrownNoise adds a brown noise track to the sequence
func (sb *sequenceBuilder) WithBrownNoise(smooth, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:        t.TrackBrownNoise,
		Amplitude:   t.AmplitudePercentToRaw(amplitude),
		NoiseSmooth: smooth,
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// WithWhiteNoise adds a white noise track to the sequence
func (sb *sequenceBuilder) WithWhiteNoise(smooth, amplitude float64) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	track := &t.Track{
		Type:        t.TrackWhiteNoise,
		Amplitude:   t.AmplitudePercentToRaw(amplitude),
		NoiseSmooth: smooth,
	}
	sb.presets[sb.lastPreset] = append(sb.presets[sb.lastPreset], *track)
	return sb
}

// At adds a timeline entry at the given time with the specified transition and steps
func (sb *sequenceBuilder) At(time, transition string, steps int) *sequenceBuilder {
	if sb.lastPreset == "" {
		return sb
	}

	timeline := fmt.Sprintf("%s %s %s %d", time, sb.lastPreset, transition, steps)
	sb.timeline = append(sb.timeline, timeline)
	return sb
}

// String returns the sequence as a string
func (sb *sequenceBuilder) String() string {
	var content strings.Builder

	cmm := t.KeywordComment
	content.WriteString(fmt.Sprintf("%s Generated by Synapseq API\n\n", cmm))

	content.WriteString(fmt.Sprintf("%s Options\n", cmm))
	for key, value := range sb.options {
		optionKey := t.KeywordOption + key
		content.WriteString(fmt.Sprintf("%s %s\n", optionKey, value))
	}

	content.WriteString(fmt.Sprintf("\n%s Presets\n", cmm))
	for preset, tracks := range sb.presets {
		content.WriteString(fmt.Sprintf("%s\n", preset))
		for _, track := range tracks {
			content.WriteString(fmt.Sprintf("  %s\n", track.String()))
		}
	}

	content.WriteString(fmt.Sprintf("\n%s Timeline\n", cmm))
	for _, timeline := range sb.timeline {
		content.WriteString(fmt.Sprintf("%s\n", timeline))
	}

	return content.String()
}
