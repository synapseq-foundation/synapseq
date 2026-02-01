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

import "fmt"

// TrackType represents the type of track/sound
type TrackType int

const (
	// Track is off
	TrackOff TrackType = iota
	// Track is silence
	TrackSilence
	// Track is a pure tone (no beat)
	TrackPureTone
	// Track is a binaural beat
	TrackBinauralBeat
	// Track is a monaural beat
	TrackMonauralBeat
	// Track is an isochronic beat
	TrackIsochronicBeat
	// Track is white noise
	TrackWhiteNoise
	// Track is pink noise
	TrackPinkNoise
	// Track is brown noise
	TrackBrownNoise
	// Track is a background noise
	TrackBackground
)

// String returns the string representation of the TrackType
func (tr TrackType) String() string {
	switch tr {
	case TrackOff:
		return KeywordOff
	case TrackSilence:
		return KeywordSilence
	case TrackPureTone:
		return KeywordTone
	case TrackBinauralBeat:
		return KeywordBinaural
	case TrackMonauralBeat:
		return KeywordMonaural
	case TrackIsochronicBeat:
		return KeywordIsochronic
	case TrackWhiteNoise:
		return KeywordWhite
	case TrackPinkNoise:
		return KeywordPink
	case TrackBrownNoise:
		return KeywordBrown
	case TrackBackground:
		return KeywordBackground
	default:
		return "unknown"
	}
}

// Track represents a track configuration
type Track struct {
	// Track type
	Type TrackType
	// Amplitude level (0-4096 for 0-100%)
	Amplitude AmplitudeType
	// Carrier frequency
	Carrier float64
	// Resonance frequency
	Resonance float64
	// Waveform shape
	Waveform WaveformType
	// Effect configuration
	Effect
}

// Validate checks if the track configuration is valid
func (tr *Track) Validate() error {
	if tr.Amplitude < 0 || tr.Amplitude > 4096 {
		return fmt.Errorf("amplitude must be between 0 and 100. Received: %.2f", tr.Amplitude.ToPercent())
	}
	if tr.Carrier < 0 {
		return fmt.Errorf("carrier frequency must be positive. Received: %.2f", tr.Carrier)
	}
	if tr.Resonance < 0 {
		return fmt.Errorf("resonance frequency must be positive. Received: %.2f", tr.Resonance)
	}
	if tr.Intensity < 0 || tr.Intensity > 1.0 {
		return fmt.Errorf("intensity must be between 0 and 100. Received: %.2f", tr.Intensity.ToPercent())
	}
	return nil
}

// String returns the string representation of the Track configuration
func (tr *Track) String() string {
	switch tr.Type {
	case TrackOff, TrackSilence:
		return "--"
	case TrackPureTone:
		return fmt.Sprintf("%s %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		return fmt.Sprintf("%s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		return fmt.Sprintf("%s %s %s %.2f", KeywordNoise, tr.Type.String(), KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackBackground:
		// Special handling for background effects
		switch tr.Effect.Type {
		case EffectSpin:
			ec := tr.Effect.Configuration.(EffectSpinConfiguration)
			return fmt.Sprintf("%s %s %s %s %.2f %s %.2f %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordBackground, KeywordSpin, ec.Width, KeywordRate, ec.Rate, KeywordIntensity, tr.Intensity.ToPercent(), KeywordAmplitude, tr.Amplitude.ToPercent())
		case EffectPulse:
			ec := tr.Effect.Configuration.(EffectPulseConfiguration)
			return fmt.Sprintf("%s %s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordBackground, KeywordPulse, ec.Pulse, KeywordIntensity, tr.Intensity.ToPercent(), KeywordAmplitude, tr.Amplitude.ToPercent())
		default:
			return fmt.Sprintf("%s %s %.2f", KeywordBackground, KeywordAmplitude, tr.Amplitude.ToPercent())
		}
	default:
		return " ???"
	}
}

// ShortString returns a compact string representation of the track configuration
func (tr *Track) ShortString() string {
	switch tr.Type {
	case TrackOff, TrackSilence:
		return " -"
	case TrackPureTone:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f)", KeywordTone, tr.Carrier, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		return fmt.Sprintf(" (%s:%.2f)", KeywordNoise, tr.Amplitude.ToPercent())
	case TrackBackground:
		// Special handling for background effects
		switch tr.Effect.Type {
		case EffectSpin:
			cfg := tr.Effect.Configuration.(EffectSpinConfiguration)
			return fmt.Sprintf(" (%s:%s %s:%.2f %s:%.2f %s:%.2f %s:%.2f)",
				KeywordEffect, tr.Effect.Type.String(), KeywordWidth, cfg.Width, KeywordRate, cfg.Rate, KeywordIntensity, tr.Intensity.ToPercent(), KeywordAmplitude, tr.Amplitude.ToPercent())
		case EffectPulse:
			cfg := tr.Effect.Configuration.(EffectPulseConfiguration)
			return fmt.Sprintf(" (%s:%s %s:%.2f %s:%.2f %s:%.2f)",
				KeywordEffect, tr.Effect.Type.String(), KeywordPulse, cfg.Pulse, KeywordIntensity, tr.Intensity.ToPercent(), KeywordAmplitude, tr.Amplitude.ToPercent())
		default:
			return fmt.Sprintf(" (%s:%.2f)", KeywordAmplitude, tr.Amplitude.ToPercent())
		}
	default:
		return " ???"
	}
}
