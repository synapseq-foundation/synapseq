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

package spsq

import t "github.com/synapseq-foundation/synapseq/v4/internal/types"

// AddPreset create a new preset with the given name
func (b *Builder) AddPreset(name string) *Builder {
	b.presets[name] = []t.Track{}
	b.lastPreset = name
	return b
}

// AddToneTrack adds a tone track to the last preset
func (b *Builder) AddToneTrack(carrier float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	track := &t.Track{
		Type:    t.TrackPureTone,
		Carrier: carrier,
	}

	b.presets[b.lastPreset] = append(b.presets[b.lastPreset], *track)
	return b
}

// AddNoiseTrack adds a noise track to the last preset
func (b *Builder) AddNoiseTrack() *Builder {
	if b.lastPreset == "" {
		return b
	}

	track := &t.Track{
		Type: t.TrackPinkNoise,
	}

	b.presets[b.lastPreset] = append(b.presets[b.lastPreset], *track)
	return b
}

// AddAmbianceTrack adds an ambiance track to the last preset
func (b *Builder) AddAmbianceTrack(name string) *Builder {
	if b.lastPreset == "" {
		return b
	}

	track := &t.Track{
		Type:         t.TrackAmbiance,
		AmbianceName: name,
	}

	b.presets[b.lastPreset] = append(b.presets[b.lastPreset], *track)
	return b
}

// WithSineWaveform sets the waveform of the last track to sine
func (b *Builder) WithSineWaveform() *Builder {
	if b.lastPreset == "" {
		return b
	}

	lastPreset := b.presets[b.lastPreset]
	if len(lastPreset) == 0 {
		return b
	}

	trackIdx := len(lastPreset) - 1
	b.presets[b.lastPreset][trackIdx].Waveform = t.WaveformSine
	return b
}

// WithSquareWaveform sets the waveform of the last track to square
func (b *Builder) WithSquareWaveform() *Builder {
	if b.lastPreset == "" {
		return b
	}

	lastPreset := b.presets[b.lastPreset]
	if len(lastPreset) == 0 {
		return b
	}

	trackIdx := len(lastPreset) - 1
	b.presets[b.lastPreset][trackIdx].Waveform = t.WaveformSquare
	return b
}

// WithTriangleWaveform sets the waveform of the last track to triangle
func (b *Builder) WithTriangleWaveform() *Builder {
	if b.lastPreset == "" {
		return b
	}

	lastPreset := b.presets[b.lastPreset]
	if len(lastPreset) == 0 {
		return b
	}

	trackIdx := len(lastPreset) - 1
	b.presets[b.lastPreset][trackIdx].Waveform = t.WaveformTriangle
	return b
}

// WithSawtoothWaveform sets the waveform of the last track to sawtooth
func (b *Builder) WithSawtoothWaveform() *Builder {
	if b.lastPreset == "" {
		return b
	}

	lastPreset := b.presets[b.lastPreset]
	if len(lastPreset) == 0 {
		return b
	}

	trackIdx := len(lastPreset) - 1
	b.presets[b.lastPreset][trackIdx].Waveform = t.WaveformSawtooth
	return b
}

// WithBinauralTone adds a binaural tone track to the sequence
func (b *Builder) WithBinauralTone(beat float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Type = t.TrackBinauralBeat
	b.presets[b.lastPreset][trackIdx].Resonance = beat
	return b
}

// WithMonauralTone adds a monaural tone track to the sequence
func (b *Builder) WithMonauralTone(beat float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Type = t.TrackMonauralBeat
	b.presets[b.lastPreset][trackIdx].Resonance = beat
	return b
}

// WithIsochronicTone adds an isochronic tone track to the sequence
func (b *Builder) WithIsochronicTone(beat float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Type = t.TrackIsochronicBeat
	b.presets[b.lastPreset][trackIdx].Resonance = beat
	return b
}

// WithBrownNoise adds a brown noise track to the sequence
func (b *Builder) WithBrownNoise(smooth float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Type = t.TrackBrownNoise
	b.presets[b.lastPreset][trackIdx].NoiseSmooth = smooth
	return b
}

// WithPinkNoise adds a pink noise track to the sequence
func (b *Builder) WithPinkNoise(smooth float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}
	b.presets[b.lastPreset][trackIdx].Type = t.TrackPinkNoise
	b.presets[b.lastPreset][trackIdx].NoiseSmooth = smooth
	return b
}

// WithWhiteNoise adds a white noise track to the sequence
func (b *Builder) WithWhiteNoise(smooth float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Type = t.TrackWhiteNoise
	b.presets[b.lastPreset][trackIdx].NoiseSmooth = smooth
	return b
}

// WithAmplitude sets the amplitude of the last track in the sequence
func (b *Builder) WithAmplitude(amplitude float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Amplitude = t.AmplitudePercentToRaw(amplitude)
	return b
}

// WithPanEffect adds a pan effect to the last track in the sequence
func (b *Builder) WithPanEffect(pan float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Effect.Type = t.EffectPan
	b.presets[b.lastPreset][trackIdx].Effect.Value = pan
	return b
}

// WithModulationEffect adds a modulation effect to the last track in the sequence
func (b *Builder) WithModulationEffect(modulation float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Effect.Type = t.EffectModulation
	b.presets[b.lastPreset][trackIdx].Effect.Value = modulation
	return b
}

// WithDopplerEffect adds a doppler effect to the last track in the sequence
func (b *Builder) WithDopplerEffect(modulation float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Effect.Type = t.EffectDoppler
	b.presets[b.lastPreset][trackIdx].Effect.Value = modulation
	return b
}

// WithIntensity sets the intensity of the effect on the last track in the sequence
func (b *Builder) WithIntensity(value float64) *Builder {
	if b.lastPreset == "" {
		return b
	}

	trackIdx := len(b.presets[b.lastPreset]) - 1
	if trackIdx < 0 {
		return b
	}

	b.presets[b.lastPreset][trackIdx].Effect.Intensity = t.IntensityPercentToRaw(value)
	return b
}
