// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dump

import (
	"encoding/json"
	"fmt"
	"sort"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type document struct {
	Comments []string           `json:"comments"`
	Options  options            `json:"options"`
	Presets  map[string][]track `json:"presets"`
	Timeline []timelineEntry    `json:"timeline"`
}

type options struct {
	SampleRate int      `json:"sampleRate"`
	Volume     int      `json:"volume"`
	Sources    sources  `json:"sources"`
	Extends    []string `json:"extends"`
}

type sources struct {
	Ambiance []source `json:"ambiance"`
	Music    []source `json:"music"`
}

type source struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type track struct {
	Family    string  `json:"family"`
	Kind      string  `json:"kind,omitempty"`
	Waveform  string  `json:"waveform,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty"`
	Smooth    float64 `json:"smooth,omitempty"`
	Source    string  `json:"source,omitempty"`
	Loop      *bool   `json:"loop,omitempty"`
	Carrier   float64 `json:"carrier,omitempty"`
	Beat      float64 `json:"beat,omitempty"`
	Effect    effect  `json:"effect"`
}

type effect struct {
	Type      string  `json:"type"`
	Pan       float64 `json:"pan,omitempty"`
	Rate      float64 `json:"rate,omitempty"`
	Distance  float64 `json:"distance,omitempty"`
	Intensity float64 `json:"intensity,omitempty"`
}

type timelineEntry struct {
	Preset       string `json:"preset"`
	StartSeconds int    `json:"startSeconds"`
	Transition   string `json:"transition"`
	Steps        int    `json:"steps"`
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
		Sources: sources{
			Ambiance: sourceList(opts.Ambiance),
			Music:    sourceList(opts.Music),
		},
		Extends: extends,
	}
}

func sourceList(values map[string]string) []source {
	if len(values) == 0 {
		return []source{}
	}

	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]source, 0, len(names))
	for _, name := range names {
		result = append(result, source{Name: name, Path: values[name]})
	}

	return result
}

func sequencePresets(presets []t.Preset) map[string][]track {
	result := make(map[string][]track)

	for i := range presets {
		preset := &presets[i]
		if preset.String() == "silence" {
			continue
		}

		tracks := presetTracks(preset)
		if len(tracks) == 0 {
			continue
		}

		result[preset.String()] = tracks
	}

	return result
}

func presetTracks(preset *t.Preset) []track {
	tracks := make([]track, 0, len(preset.Track))

	for i := range preset.Track {
		if preset.Track[i].Type == t.TrackOff {
			continue
		}
		tracks = append(tracks, convertTrack(&preset.Track[i]))
	}

	return tracks
}

func convertTrack(src *t.Track) track {
	out := track{
		Amplitude: src.Amplitude.ToPercent(),
		Effect:    convertEffect(src.Effect),
	}

	switch src.Type {
	case t.TrackSilence:
		out.Family = t.KeywordSilence
	case t.TrackPureTone:
		out.Family = t.KeywordTone
		out.Kind = t.KeywordTone
		out.Waveform = src.Waveform.String()
		out.Carrier = src.Carrier
	case t.TrackBinauralBeat, t.TrackMonauralBeat, t.TrackIsochronicBeat:
		out.Family = t.KeywordTone
		out.Kind = src.Type.String()
		out.Waveform = src.Waveform.String()
		out.Carrier = src.Carrier
		out.Beat = src.Resonance
	case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
		out.Family = t.KeywordNoise
		out.Kind = src.Type.String()
		out.Waveform = src.Waveform.String()
		out.Smooth = src.NoiseSmooth
	case t.TrackAmbiance:
		out.Family = t.KeywordAmbiance
		out.Source = src.SourceName
	case t.TrackMusic:
		loop := false
		out.Family = t.KeywordMusic
		out.Source = src.SourceName
		out.Loop = &loop
	}

	return out
}

func convertEffect(src t.Effect) effect {
	out := effect{Type: src.Type.String()}

	switch src.Type {
	case t.EffectPan:
		out.Pan = src.Value
		out.Intensity = src.Intensity.ToPercent()
	case t.EffectModulation:
		out.Rate = src.Value
		out.Intensity = src.Intensity.ToPercent()
	case t.EffectDoppler:
		out.Distance = src.Value
		out.Intensity = src.Intensity.ToPercent()
	}

	return out
}

func sequenceTimeline(periods []t.Period) []timelineEntry {
	if len(periods) == 0 {
		return []timelineEntry{}
	}

	result := make([]timelineEntry, 0, len(periods))
	for i := range periods {
		result = append(result, timelineEntry{
			Preset:       periods[i].PresetName,
			StartSeconds: periods[i].Time / 1000,
			Transition:   periods[i].Transition.String(),
			Steps:        periods[i].Steps,
		})
	}

	return result
}
