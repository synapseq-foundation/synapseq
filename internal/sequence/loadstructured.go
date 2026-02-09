//go:build !wasm

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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"

	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// initializeTracks initializes an array of off tracks
func initializeTracks() [t.NumberOfChannels]t.Track {
	var tracks [t.NumberOfChannels]t.Track
	for ch := range t.NumberOfChannels {
		tracks[ch].Type = t.TrackOff
		tracks[ch].Carrier = 0.0
		tracks[ch].Resonance = 0.0
		tracks[ch].Amplitude = 0.0
		tracks[ch].Waveform = t.WaveformSine
		tracks[ch].Effect = t.Effect{Type: t.EffectOff, Intensity: 0.0}
	}
	return tracks
}

// resolveBackgroundPath resolves the background audio path
func resolveBackgroundPath(path, basePath string) (string, error) {
	if path == "-" {
		return "", fmt.Errorf("stdin (-) is not supported for background audio")
	}

	if s.IsRemoteFile(path) {
		return path, nil
	}

	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		expanded := filepath.Join(homeDir, strings.TrimPrefix(path, "~"))
		return filepath.Clean(expanded), nil
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	fullPath := filepath.Join(basePath, path)
	return filepath.Clean(fullPath), nil
}

// LoadStructuredSequence loads and parses a json/xml/yaml sequence file
func LoadStructuredSequence(filename string, format t.FileFormat) (*t.Sequence, error) {
	data, err := s.GetFile(filename, format)
	if err != nil {
		return nil, fmt.Errorf("failed to load structured sequence file: %v", err)
	}

	var input t.SynapSeqInput
	switch format {
	case t.FormatJSON:
		if err := json.Unmarshal(data, &input); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}
	case t.FormatXML:
		if err := xml.Unmarshal(data, &input); err != nil {
			return nil, fmt.Errorf("error unmarshalling XML: %v", err)
		}
	case t.FormatYAML:
		if err := yaml.Unmarshal(data, &input); err != nil {
			return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s (use json | xml | yaml)", format.String())
	}

	if len(input.Sequence) < 2 {
		return nil, fmt.Errorf("not enough sequence data found in input file")
	}

	var backgroundPath string
	if input.Options.Background != "" {
		var err error

		path := input.Options.Background
		basePath := filepath.Dir(filename)

		backgroundPath, err = resolveBackgroundPath(path, basePath)
		if err != nil {
			return nil, err
		}
	}

	gainLevel := t.GainLevelOff
	if input.Options.GainLevel != "" {
		switch strings.ToLower(strings.TrimSpace(input.Options.GainLevel)) {
		case t.KeywordOptionGainLevelHigh:
			gainLevel = t.GainLevelHigh
		case t.KeywordOptionGainLevelMedium:
			gainLevel = t.GainLevelMedium
		case t.KeywordOptionGainLevelLow:
			gainLevel = t.GainLevelLow
		default:
			return nil, fmt.Errorf("invalid gain level: %s", input.Options.GainLevel)
		}
	}

	// Initialize audio options
	options := &t.SequenceOptions{
		SampleRate:     input.Options.Samplerate,
		Volume:         input.Options.Volume,
		BackgroundPath: backgroundPath,
		GainLevel:      gainLevel,
	}

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// Initialize off tracks
	offTracks := initializeTracks()

	var periods []t.Period
	for idx, seq := range input.Sequence {
		if idx == 0 && seq.Time != 0 {
			return nil, fmt.Errorf("first timeline must start at 0ms (00:00:00)")
		}
		if idx >= 1 && seq.Time <= input.Sequence[idx-1].Time {
			return nil, fmt.Errorf("timeline %d time must be greater than previous timeline time", idx+1)
		}
		if len(seq.Track.Tones)+len(seq.Track.Noises) > t.NumberOfChannels {
			return nil, fmt.Errorf("too many elements defined (max %d)", t.NumberOfChannels)
		}

		tracks := offTracks
		trackIdx := 0

		for _, tone := range seq.Track.Tones {
			var mode t.TrackType
			// Get mode
			switch tone.Mode {
			case t.KeywordBinaural:
				mode = t.TrackBinauralBeat
			case t.KeywordMonaural:
				mode = t.TrackMonauralBeat
			case t.KeywordIsochronic:
				mode = t.TrackIsochronicBeat
			case t.KeywordPure:
				mode = t.TrackPureTone
			default:
				return nil, fmt.Errorf("invalid tone mode: %s", tone.Mode)
			}

			var waveForm t.WaveformType
			switch tone.Waveform {
			case t.KeywordSine:
				waveForm = t.WaveformSine
			case t.KeywordSquare:
				waveForm = t.WaveformSquare
			case t.KeywordTriangle:
				waveForm = t.WaveformTriangle
			case t.KeywordSawtooth:
				waveForm = t.WaveformSawtooth
			default:
				return nil, fmt.Errorf("invalid waveform type: %s", tone.Waveform)
			}

			tr := t.Track{
				Type:      mode,
				Carrier:   tone.Carrier,
				Resonance: tone.Resonance,
				Amplitude: t.AmplitudePercentToRaw(tone.Amplitude),
				Waveform:  waveForm,
			}

			if err := tr.Validate(); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			tracks[trackIdx] = tr
			trackIdx++
		}
		for _, noise := range seq.Track.Noises {
			var mode t.TrackType
			// Get mode
			switch noise.Mode {
			case t.KeywordWhite:
				mode = t.TrackWhiteNoise
			case t.KeywordPink:
				mode = t.TrackPinkNoise
			case t.KeywordBrown:
				mode = t.TrackBrownNoise
			default:
				return nil, fmt.Errorf("invalid noise mode: %s", noise.Mode)
			}

			tr := t.Track{
				Type:      mode,
				Amplitude: t.AmplitudePercentToRaw(noise.Amplitude),
			}

			if err := tr.Validate(); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			tracks[trackIdx] = tr
			trackIdx++
		}

		if backgroundPath != "" {
			if seq.Track.Background == nil {
				return nil, fmt.Errorf("background audio defined but no background settings found in timeline %d", idx+1)
			}

			var waveForm t.WaveformType
			switch seq.Track.Background.Waveform {
			case t.KeywordSine:
				waveForm = t.WaveformSine
			case t.KeywordSquare:
				waveForm = t.WaveformSquare
			case t.KeywordTriangle:
				waveForm = t.WaveformTriangle
			case t.KeywordSawtooth:
				waveForm = t.WaveformSawtooth
			default:
				return nil, fmt.Errorf("invalid background waveform type: %s", seq.Track.Background.Waveform)
			}

			effect := t.Effect{Type: t.EffectOff, Intensity: 0.0}
			if seq.Track.Background.Effect != nil {
				effect.Intensity = t.IntensityPercentToRaw(seq.Track.Background.Effect.Intensity)

				if seq.Track.Background.Effect.Spin != nil {
					effect.Type = t.EffectPan
				} else if seq.Track.Background.Effect.Pulse != nil {
					effect.Type = t.EffectModulation
				} else {
					return nil, fmt.Errorf("invalid background effect type")
				}
			}

			var bgTrack t.Track
			switch effect.Type {
			case t.EffectPan:
				bgTrack = t.Track{
					Type:      t.TrackBackground,
					Carrier:   seq.Track.Background.Effect.Spin.Width,
					Resonance: seq.Track.Background.Effect.Spin.Rate,
					Amplitude: t.AmplitudePercentToRaw(seq.Track.Background.Amplitude),
					Waveform:  waveForm,
					Effect:    effect,
				}
			case t.EffectModulation:
				bgTrack = t.Track{
					Type:      t.TrackBackground,
					Resonance: seq.Track.Background.Effect.Pulse.Resonance,
					Amplitude: t.AmplitudePercentToRaw(seq.Track.Background.Amplitude),
					Waveform:  waveForm,
					Effect:    effect,
				}
			default:
				bgTrack = t.Track{
					Type:      t.TrackBackground,
					Amplitude: t.AmplitudePercentToRaw(seq.Track.Background.Amplitude),
					Waveform:  waveForm,
					Effect:    effect,
				}
			}

			if err := bgTrack.Validate(); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			tracks[trackIdx] = bgTrack
			trackIdx++
		}

		var transition t.TransitionType
		switch seq.Transition {
		case t.KeywordTransitionSteady:
			transition = t.TransitionSteady
		case t.KeywordTransitionEaseIn:
			transition = t.TransitionEaseIn
		case t.KeywordTransitionEaseOut:
			transition = t.TransitionEaseOut
		case t.KeywordTransitionSmooth:
			transition = t.TransitionSmooth
		default:
			return nil, fmt.Errorf("invalid transition type: %s", seq.Transition)
		}

		// Process Period
		period := t.Period{
			Time:       seq.Time,
			TrackStart: tracks,
			TrackEnd:   tracks,
			Transition: transition,
		}
		// Adjust previous period end if needed
		var lastPeriod *t.Period
		if len(periods) > 0 {
			lastPeriod = &periods[len(periods)-1]
		}
		if lastPeriod != nil {
			if err := s.AdjustPeriods(lastPeriod, &period); err != nil {
				return nil, fmt.Errorf("%v", err)
			}
		}
		periods = append(periods, period)
	}

	return &t.Sequence{
		Periods:  periods,
		Options:  options,
		Comments: input.Description,
	}, nil
}
