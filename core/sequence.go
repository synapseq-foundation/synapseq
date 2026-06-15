// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"fmt"
	"maps"
	"path/filepath"

	"github.com/synapseq-foundation/synapseq/v4/internal/dump"
	r "github.com/synapseq-foundation/synapseq/v4/internal/resource"
	seq "github.com/synapseq-foundation/synapseq/v4/internal/sequence"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// LoadFile loads the sequence from the input file.
func (ac *AppContext) LoadFile(path string) (*LoadedContext, error) {
	rawContent, err := r.GetFile(path, t.FormatText)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	absInputFile, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	sequence, err := seq.LoadTextSequence(rawContent, absInputFile, filepath.Dir(absInputFile))
	if err != nil {
		return nil, err
	}

	return &LoadedContext{
		appCtx:   ac,
		sequence: sequence,
	}, nil
}

// LoadContent loads the sequence from raw text content.
func (ac *AppContext) LoadContent(content string) (*LoadedContext, error) {
	sequence, err := seq.LoadTextSequence([]byte(content), "", "")
	if err != nil {
		return nil, err
	}

	return &LoadedContext{
		appCtx:   ac,
		sequence: sequence,
	}, nil
}

// Comments returns a defensive copy of sequence comments.
func (lc *LoadedContext) Comments() []string {
	if lc.sequence == nil || len(lc.sequence.Comments) == 0 {
		return []string{}
	}

	comments := make([]string, len(lc.sequence.Comments))
	copy(comments, lc.sequence.Comments)

	return comments
}

// SampleRate returns the sample rate from the loaded sequence options.
func (lc *LoadedContext) SampleRate() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.SampleRate
}

// Extends returns a defensive copy of extends list.
func (lc *LoadedContext) Extends() []string {
	if lc.sequence == nil || lc.sequence.Options == nil || len(lc.sequence.Options.Extends) == 0 {
		return []string{}
	}

	extends := make([]string, len(lc.sequence.Options.Extends))
	copy(extends, lc.sequence.Options.Extends)

	return extends
}

// Volume returns the volume from the loaded sequence options.
func (lc *LoadedContext) Volume() int {
	if lc.sequence == nil || lc.sequence.Options == nil {
		return 0
	}

	return lc.sequence.Options.Volume
}

// Ambiance returns a defensive copy of ambiance map.
func (lc *LoadedContext) Ambiance() map[string]string {
	if lc.sequence == nil || lc.sequence.Options == nil || len(lc.sequence.Options.Ambiance) == 0 {
		return map[string]string{}
	}

	ambiance := make(map[string]string, len(lc.sequence.Options.Ambiance))
	maps.Copy(ambiance, lc.sequence.Options.Ambiance)

	return ambiance
}

// Music returns a defensive copy of music map.
func (lc *LoadedContext) Music() map[string]string {
	if lc.sequence == nil || lc.sequence.Options == nil || len(lc.sequence.Options.Music) == 0 {
		return map[string]string{}
	}

	music := make(map[string]string, len(lc.sequence.Options.Music))
	maps.Copy(music, lc.sequence.Options.Music)

	return music
}

// Presets returns a defensive copy of presets map.
func (lc *LoadedContext) Presets() map[string][]string {
	if lc.sequence == nil || len(lc.sequence.Presets) == 0 {
		return map[string][]string{}
	}

	presets := make(map[string][]string, len(lc.sequence.Presets))
	for _, p := range lc.sequence.Presets {
		for _, tr := range p.Track {
			if tr.Type == t.TrackOff || tr.Type == t.TrackSilence {
				continue
			}
			pName := p.String()
			presets[pName] = append(presets[pName], tr.String())
		}
	}

	return presets
}

// Timeline returns a timeline of the sequence as a map of preset name to timeline string.
func (lc *LoadedContext) Timeline() map[string]string {
	if lc.sequence == nil || len(lc.sequence.Periods) == 0 {
		return nil
	}

	timeline := make(map[string]string, len(lc.sequence.Periods))
	for _, p := range lc.sequence.Periods {
		ln := fmt.Sprintf("%s %s %d", p.PresetName, p.Transition.String(), p.Steps)
		timeline[p.TimeString()] = ln
	}

	return timeline
}

// RawContent returns a defensive copy of raw content.
func (lc *LoadedContext) RawContent() []byte {
	if lc.sequence == nil || len(lc.sequence.RawContent) == 0 {
		return []byte{}
	}

	raw := make([]byte, len(lc.sequence.RawContent))
	copy(raw, lc.sequence.RawContent)

	return raw
}

// JSON returns an indented JSON representation of the loaded sequence.
func (lc *LoadedContext) JSON() ([]byte, error) {
	if lc.sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}

	return dump.JSON(lc.sequence)
}
