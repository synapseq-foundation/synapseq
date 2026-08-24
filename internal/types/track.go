// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
	// Track is a ambiance
	TrackAmbiance
	// Track is music
	TrackMusic
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
	case TrackAmbiance:
		return KeywordAmbiance
	case TrackMusic:
		return KeywordMusic
	default:
		return "unknown"
	}
}

// Track represents a track configuration
type Track struct {
	// Track type
	Type TrackType
	// Amplitude levels for the left and right channels (0-4096 for 0-100%).
	Amplitude [2]AmplitudeType
	// Carrier frequency
	Carrier float64
	// Resonance frequency
	Resonance float64
	// Waveform shape
	Waveform WaveformName
	// Named audio source
	SourceName string
	// Noise smooth (0-100, for 0-100%)
	NoiseSmooth float64
	// Effect configuration
	Effect Effect
}

// Validate checks if the track configuration is valid
func (tr *Track) Validate() error {
	effect := &tr.Effect

	for channel, amplitude := range tr.Amplitude {
		if amplitude < 0 || amplitude > 4096 {
			return fmt.Errorf("%s amplitude must be between 0 and 100. Received: %.2f", channelName(channel), amplitude.ToPercent())
		}
	}
	if tr.Carrier < 0 {
		return fmt.Errorf("carrier frequency must be a positive number. Received: %.2f", tr.Carrier)
	}
	if tr.Resonance < 0 {
		return fmt.Errorf("resonance frequency must be a positive number. Received: %.2f", tr.Resonance)
	}
	if effect.Value < 0 {
		return fmt.Errorf("effect value must be greater than or equal to 0. Received: %.2f", effect.Value)
	}
	if effect.Intensity < 0 || effect.Intensity > 1.0 {
		return fmt.Errorf("intensity must be between 0 and 100. Received: %.2f", effect.Intensity.ToPercent())
	}

	// Track-type specific validation
	switch tr.Type {
	case TrackPureTone:
		if tr.Resonance != 0 {
			return fmt.Errorf("tone does not use beat/resonance (use binaural/monaural/isochronic). Received: %.2f", tr.Resonance)
		}
	case TrackBinauralBeat, TrackMonauralBeat:
		if tr.Resonance >= 2*tr.Carrier {
			return fmt.Errorf("binaural beat and monaural beat must be < 2*carrier (carrier - beat/2 must be > 0). Received beat: %.2f, carrier: %.2f", tr.Resonance, tr.Carrier)
		}
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		if tr.NoiseSmooth < 0 || tr.NoiseSmooth > 100 {
			return fmt.Errorf("noise smooth must be between 0 and 100. Received: %.2f", tr.NoiseSmooth)
		}
	}

	return nil
}

func channelName(channel int) string {
	if channel == 0 {
		return "left"
	}
	return "right"
}

// String returns the string representation of the Track configuration
func (tr *Track) String() string {
	switch tr.Type {
	case TrackOff, TrackSilence:
		return "--"
	case TrackPureTone:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf("%s %s %s %.2f %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf("%s %s %s %.2f %s %s %.2f %s %.2f %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, KeywordEffect, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf("%s %s %s %.2f %s %.2f %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf("%s %s %s %.2f %s %.2f %s %s %.2f %s %.2f %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordEffect, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf("%s %s %s %.2f %s %s %.2f %s %.2f", KeywordNoise, tr.Type.String(), KeywordSmooth, tr.NoiseSmooth, KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf("%s %s %s %.2f %s %s %.2f %s %.2f %s %s %.2f %s %.2f", KeywordNoise, tr.Type.String(), KeywordSmooth, tr.NoiseSmooth, KeywordEffect, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
	case TrackAmbiance, TrackMusic:
		keyword := tr.Type.String()
		if tr.Effect.Type == EffectOff {
			if tr.Waveform.Effective() == WaveformSine {
				return fmt.Sprintf("%s %s %s %s %.2f %s %.2f", keyword, tr.SourceName, KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
			}
			return fmt.Sprintf("%s %s %s %s %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), keyword, tr.SourceName, KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
		if tr.Waveform.Effective() == WaveformSine {
			return fmt.Sprintf("%s %s %s %s %.2f %s %.2f %s %s %.2f %s %.2f", keyword, tr.SourceName, KeywordEffect, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
		return fmt.Sprintf("%s %s %s %s %s %s %.2f %s %.2f %s %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), keyword, tr.SourceName, KeywordEffect, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordAmplitude, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
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
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)", KeywordTone, tr.Carrier, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f %s:%.2f %s:%.2f)", KeywordTone, tr.Carrier, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f %s:%.2f)", KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f %s:%.2f %s:%.2f %s:%.2f)", KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent(), KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		}
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf(" (%s %s:%.2f %s:%.2f %s:%.2f)", tr.Type.String(), KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent(), KeywordSmooth, tr.NoiseSmooth)
		} else {
			return fmt.Sprintf(" (%s %s:%.2f %s:%.2f %s:%.2f %s:%.2f %s:%.2f)", tr.Type.String(), KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent(), KeywordSmooth, tr.NoiseSmooth, tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent())
		}
	case TrackAmbiance, TrackMusic:
		if tr.Effect.Type == EffectOff {
			return fmt.Sprintf(" (%s %s:%.2f %s:%.2f)", tr.SourceName, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent())
		} else {
			return fmt.Sprintf(" (%s %s:%.2f %s:%.2f %s:%.2f %s:%.2f)", tr.SourceName, KeywordLeft, tr.Amplitude[0].ToPercent(), KeywordRight, tr.Amplitude[1].ToPercent(), tr.Effect.Type.String(), tr.Effect.Value, KeywordIntensity, tr.Effect.Intensity.ToPercent())
		}
	default:
		return " ???"
	}
}
