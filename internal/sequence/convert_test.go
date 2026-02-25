//go:build ignore

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

package sequence

import (
	"strings"
	"testing"

	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

func TestConvertToText_BasicSequence(ts *testing.T) {
	var periods []t.Period

	// Period 0
	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period0)

	// Period 1
	period1 := t.Period{
		Time:       15000,
		Transition: t.TransitionSmooth,
	}
	period1.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period1)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
		Comments: []string{"Test sequence", "Basic conversion"},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify header
	if !strings.Contains(result, "# GENERATED FROM SYNAPSEQ STRUCTURED SEQUENCE FILE") {
		ts.Errorf("expected header not found in output")
	}

	// Verify comments
	if !strings.Contains(result, "## Test sequence") {
		ts.Errorf("expected first comment not found in output")
	}
	if !strings.Contains(result, "## Basic conversion") {
		ts.Errorf("expected second comment not found in output")
	}

	// Verify options section
	if !strings.Contains(result, "# Options") {
		ts.Errorf("expected options section not found")
	}
	if !strings.Contains(result, "@samplerate 44100") {
		ts.Errorf("expected samplerate option not found")
	}
	if !strings.Contains(result, "@volume 100") {
		ts.Errorf("expected volume option not found")
	}

	// Verify presets section
	if !strings.Contains(result, "# Presets") {
		ts.Errorf("expected presets section not found")
	}
	if !strings.Contains(result, "tone-set-001") {
		ts.Errorf("expected first preset not found")
	}
	if !strings.Contains(result, "tone-set-002") {
		ts.Errorf("expected second preset not found")
	}

	// Verify timeline section
	if !strings.Contains(result, "# Timeline") {
		ts.Errorf("expected timeline section not found")
	}
	if !strings.Contains(result, "0:00 tone-set-001 steady") {
		ts.Errorf("expected first timeline entry not found")
	}
	if !strings.Contains(result, "0:15 tone-set-002 smooth") {
		ts.Errorf("expected second timeline entry not found")
	}

	// Verify track content
	if !strings.Contains(result, "binaural") {
		ts.Errorf("expected binaural track type not found")
	}
}

func TestConvertToText_WithBackgroundOptions(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackWhiteNoise,
		Amplitude: t.AmplitudePercentToRaw(25),
	}
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate:     48000,
			Volume:         90,
			BackgroundPath: "sounds/pink-noise.wav",
			GainLevel:      t.GainLevelHigh,
		},
		Comments: []string{"Background test"},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify background options
	if !strings.Contains(result, "@background sounds/pink-noise.wav") {
		ts.Errorf("expected background path not found")
	}
	if !strings.Contains(result, "@gainlevel high") {
		ts.Errorf("expected gainlevel not found")
	}
	if !strings.Contains(result, "@samplerate 48000") {
		ts.Errorf("expected custom samplerate not found")
	}
	if !strings.Contains(result, "@volume 90") {
		ts.Errorf("expected custom volume not found")
	}
}

func TestConvertToText_MultipleTracksPerPeriod(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   200,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	period0.TrackStart[1] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   300,
		Resonance: 7,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSquare,
	}
	period0.TrackStart[2] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(20),
	}
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify all track types are present
	if !strings.Contains(result, "binaural") {
		ts.Errorf("expected binaural track not found")
	}
	if !strings.Contains(result, "isochronic") {
		ts.Errorf("expected isochronic track not found")
	}
	if !strings.Contains(result, "noise pink") {
		ts.Errorf("expected pink noise track not found")
	}
}

func TestConvertToText_EmptyTracks(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	// TrackStart[1] and [2] are TrackOff (default)
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Count occurrences of binaural (should only be one track)
	binauralCount := strings.Count(result, "binaural")
	if binauralCount != 1 {
		ts.Errorf("expected 1 binaural track, found %d", binauralCount)
	}

	// Verify preset has content
	if !strings.Contains(result, "tone-set-001") {
		ts.Errorf("expected preset not found")
	}
}

func TestConvertToText_DifferentWaveforms(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSquare,
	}
	periods = append(periods, period0)

	period1 := t.Period{
		Time:       10000,
		Transition: t.TransitionSmooth,
	}
	period1.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformTriangle,
	}
	periods = append(periods, period1)

	period2 := t.Period{
		Time:       20000,
		Transition: t.TransitionEaseIn,
	}
	period2.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   350,
		Resonance: 12,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformSawtooth,
	}
	periods = append(periods, period2)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify different waveforms
	if !strings.Contains(result, "square") {
		ts.Errorf("expected square waveform not found")
	}
	if !strings.Contains(result, "triangle") {
		ts.Errorf("expected triangle waveform not found")
	}
	if !strings.Contains(result, "sawtooth") {
		ts.Errorf("expected sawtooth waveform not found")
	}

	// Verify different transitions
	if !strings.Contains(result, "ease-in") {
		ts.Errorf("expected ease-in transition not found")
	}
}

func TestConvertToText_DifferentTransitions(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period0)

	period1 := t.Period{
		Time:       10000,
		Transition: t.TransitionSmooth,
	}
	period1.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period1)

	period2 := t.Period{
		Time:       20000,
		Transition: t.TransitionEaseOut,
	}
	period2.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   350,
		Resonance: 12,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period2)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify different transitions in timeline
	if !strings.Contains(result, "steady") {
		ts.Errorf("expected steady transition not found")
	}
	if !strings.Contains(result, "smooth") {
		ts.Errorf("expected smooth transition not found")
	}
	if !strings.Contains(result, "ease-out") {
		ts.Errorf("expected ease-out transition not found")
	}
}

func TestConvertToText_NoComments(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
		Comments: []string{}, // No comments
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Should still have header but no comment lines with ##
	if !strings.Contains(result, "# GENERATED FROM SYNAPSEQ STRUCTURED SEQUENCE FILE") {
		ts.Errorf("expected header not found")
	}

	// Count comment markers (should only be section headers with single #)
	doubleHashCount := strings.Count(result, "##")
	if doubleHashCount != 0 {
		ts.Errorf("expected no ## comment markers, found %d", doubleHashCount)
	}
}

func TestConvertToText_NilOptions(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: nil, // Nil options
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Should not have options section
	if strings.Contains(result, "# Options") {
		ts.Errorf("unexpected options section found for nil options")
	}

	// Should still have other sections
	if !strings.Contains(result, "# Presets") {
		ts.Errorf("expected presets section not found")
	}
	if !strings.Contains(result, "# Timeline") {
		ts.Errorf("expected timeline section not found")
	}
}

func TestConvertToText_LongSequence(ts *testing.T) {
	var periods []t.Period

	// Create 10 periods
	for i := 0; i < 10; i++ {
		period := t.Period{
			Time:       i * 60000, // Every minute
			Transition: t.TransitionSmooth,
		}
		period.TrackStart[0] = t.Track{
			Type:      t.TrackBinauralBeat,
			Carrier:   float64(200 + i*10),
			Resonance: float64(5 + i),
			Amplitude: t.AmplitudePercentToRaw(float64(10 + i)),
			Waveform:  t.WaveformSine,
		}
		periods = append(periods, period)
	}

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify all presets were generated
	for i := 1; i <= 10; i++ {
		presetID := strings.Contains(result, "tone-set-")
		if !presetID {
			ts.Errorf("expected preset tone-set-%03d not found", i)
			break
		}
	}

	// Verify timeline entries
	if !strings.Contains(result, "0:00:00") {
		ts.Errorf("expected first timeline entry not found")
	}
	if !strings.Contains(result, "0:09:00") {
		ts.Errorf("expected last timeline entry not found")
	}
}

func TestConvertToText_AllNoiseTypes(ts *testing.T) {
	var periods []t.Period

	period0 := t.Period{
		Time:       0,
		Transition: t.TransitionSteady,
	}
	period0.TrackStart[0] = t.Track{
		Type:      t.TrackWhiteNoise,
		Amplitude: t.AmplitudePercentToRaw(10),
	}
	period0.TrackStart[1] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(15),
	}
	period0.TrackStart[2] = t.Track{
		Type:      t.TrackBrownNoise,
		Amplitude: t.AmplitudePercentToRaw(20),
	}
	periods = append(periods, period0)

	seq := &t.Sequence{
		Periods: periods,
		Options: &t.SequenceOptions{
			SampleRate: 44100,
			Volume:     100,
		},
	}

	result, err := ConvertToText(seq)
	if err != nil {
		ts.Fatalf("ConvertToText() error: %v", err)
	}

	// Verify all noise types
	if !strings.Contains(result, "noise white") {
		ts.Errorf("expected white noise track not found")
	}
	if !strings.Contains(result, "noise pink") {
		ts.Errorf("expected pink noise track not found")
	}
	if !strings.Contains(result, "noise brown") {
		ts.Errorf("expected brown noise track not found")
	}
}
