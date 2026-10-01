// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import (
	"fmt"
	"path/filepath"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func build(ctx *synapseq.AppContext, parsedSequence *sequence) (*spsq.Builder, error) {
	builder, err := spsq.New(ctx)
	if err != nil {
		return nil, err
	}
	builder.SampleRate(parsedSequence.sampleRate)
	var musicName string
	if parsedSequence.musicPath != "" {
		musicPath, err := currentRelativeMusicPath(parsedSequence.musicPath)
		if err != nil {
			return nil, lineError(parsedSequence.source, parsedSequence.musicLine, err.Error())
		}
		musicName = filepath.Base(musicPath)
		builder.Music(musicName, musicPath)
	}

	presets, silentDefinitions, err := convertDefinitions(builder, parsedSequence, musicName)
	if err != nil {
		return nil, err
	}
	if err := appendTimeline(builder, parsedSequence, presets, silentDefinitions); err != nil {
		return nil, err
	}

	return builder, nil
}

func convertDefinitions(
	builder *spsq.Builder,
	parsedSequence *sequence,
	musicName string,
) (map[string]*spsq.Preset, map[string]struct{}, error) {
	presets := make(map[string]*spsq.Preset, len(parsedSequence.definitions))
	names := convertedNames(parsedSequence.definitions)
	silentDefinitions := make(map[string]struct{})
	for _, definition := range parsedSequence.definitions {
		hasTrack := false
		hasOff := false
		for _, voice := range definition.voices {
			if voice.kind == voiceMix && parsedSequence.musicPath == "" {
				return nil, nil, lineError(parsedSequence.source, definition.line, "mix voice requires a -m music source")
			}
			if voice.kind == voiceOff {
				hasOff = true
				continue
			}
			hasTrack = true
		}
		if !hasTrack {
			if hasOff {
				silentDefinitions[definition.name] = struct{}{}
				continue
			}
			message := fmt.Sprintf("NameDef %q has no convertible voices", definition.name)
			return nil, nil, lineError(parsedSequence.source, definition.line, message)
		}

		preset := builder.NewPreset(names[definition.name])
		for _, voice := range definition.voices {
			switch voice.kind {
			case voiceOff:
				continue
			case voicePink:
				preset.Pink(0).Amplitude(voice.amplitude)
			case voiceTone:
				preset.Tone(voice.carrier).Amplitude(voice.amplitude)
			case voiceBinaural:
				preset.Tone(voice.carrier).Binaural(voice.beat).Amplitude(voice.amplitude)
			case voiceSpin:
				preset.Pink(0).Pan(voice.beat).Intensity(100).Amplitude(voice.amplitude)
			case voiceMix:
				preset.Music(musicName).Amplitude(voice.amplitude)
			}
		}
		presets[definition.name] = preset
	}
	return presets, silentDefinitions, nil
}
