// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dump

import (
	"encoding/json"
	"fmt"
	"maps"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type document struct {
	Comments []string        `json:"comments"`
	Options  options         `json:"options"`
	Presets  []preset        `json:"presets"`
	Timeline []timelineEntry `json:"timeline"`
}

type options struct {
	SampleRate int               `json:"sampleRate"`
	Volume     int               `json:"volume"`
	Ambiance   map[string]string `json:"ambiance"`
	Music      map[string]string `json:"music"`
	Extends    []string          `json:"extends"`
}

type preset struct {
	Name   string  `json:"name"`
	Tracks []track `json:"tracks"`
}

type track struct {
	Index       int     `json:"index"`
	Waveform    string  `json:"waveform"`
	Type        string  `json:"type"`
	Carrier     float64 `json:"carrier"`
	Resonance   float64 `json:"resonance"`
	Amplitude   float64 `json:"amplitude"`
	SourceName  string  `json:"sourceName"`
	NoiseSmooth float64 `json:"noiseSmooth"`
	Effect      effect  `json:"effect"`
	Line        string  `json:"line"`
}

type effect struct {
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	Intensity float64 `json:"intensity"`
}

type timelineEntry struct {
	Timestamp  string `json:"timestamp"`
	PresetName string `json:"presetName"`
	Transition string `json:"transition"`
	Steps      int    `json:"steps"`
	Line       string `json:"line"`
}

// JSON returns an indented JSON representation of a loaded sequence.
func JSON(sequence *t.Sequence) ([]byte, error) {
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}
	if sequence.Options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	doc := document{
		Comments: sequenceComments(sequence),
		Options:  sequenceOptions(sequence.Options),
		Presets:  sequencePresets(sequence.Presets),
		Timeline: sequenceTimeline(sequence.Periods),
	}

	return json.MarshalIndent(doc, "", "  ")
}

func sequenceComments(sequence *t.Sequence) []string {
	if len(sequence.Comments) == 0 {
		return []string{}
	}

	comments := make([]string, len(sequence.Comments))
	copy(comments, sequence.Comments)
	return comments
}

func sequenceOptions(opts *t.SequenceOptions) options {
	extends := []string{}
	if len(opts.Extends) > 0 {
		extends = make([]string, len(opts.Extends))
		copy(extends, opts.Extends)
	}

	return options{
		SampleRate: opts.SampleRate,
		Volume:     opts.Volume,
		Ambiance:   stringMap(opts.Ambiance),
		Music:      stringMap(opts.Music),
		Extends:    extends,
	}
}

func stringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}

	result := make(map[string]string, len(values))
	maps.Copy(result, values)

	return result
}

func sequencePresets(presets []t.Preset) []preset {
	result := make([]preset, 0, len(presets))

	for i := range presets {
		src := &presets[i]
		if src.String() == "silence" {
			continue
		}

		tracks := presetTracks(src)
		if len(tracks) == 0 {
			continue
		}

		result = append(result, preset{
			Name:   src.String(),
			Tracks: tracks,
		})
	}

	return result
}

func presetTracks(preset *t.Preset) []track {
	tracks := make([]track, 0, len(preset.Track))

	for i := range preset.Track {
		if preset.Track[i].Type == t.TrackOff || preset.Track[i].Type == t.TrackSilence {
			continue
		}
		tracks = append(tracks, convertTrack(i+1, &preset.Track[i]))
	}

	return tracks
}

func convertTrack(index int, src *t.Track) track {
	return track{
		Index:       index,
		Waveform:    src.Waveform.String(),
		Type:        src.Type.String(),
		Carrier:     src.Carrier,
		Resonance:   src.Resonance,
		Amplitude:   src.Amplitude.ToPercent(),
		SourceName:  src.SourceName,
		NoiseSmooth: src.NoiseSmooth,
		Effect:      convertEffect(src.Effect),
		Line:        src.String(),
	}
}

func convertEffect(src t.Effect) effect {
	return effect{
		Type:      src.Type.String(),
		Value:     src.Value,
		Intensity: src.Intensity.ToPercent(),
	}
}

func sequenceTimeline(periods []t.Period) []timelineEntry {
	if len(periods) == 0 {
		return []timelineEntry{}
	}

	result := make([]timelineEntry, 0, len(periods))
	for i := range periods {
		line := fmt.Sprintf(
			"%s %s %s %d",
			periods[i].TimeString(),
			periods[i].PresetName,
			periods[i].Transition.String(),
			periods[i].Steps,
		)

		result = append(result, timelineEntry{
			Timestamp:  periods[i].TimeString(),
			PresetName: periods[i].PresetName,
			Transition: periods[i].Transition.String(),
			Steps:      periods[i].Steps,
			Line:       line,
		})
	}

	return result
}
